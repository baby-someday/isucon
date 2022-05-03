package remote

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type fileInfos struct {
	message     string
	filename    string
	permissions string
	size        int64
}

type responseType = uint8

type response struct {
	responseType responseType
	message      string
}

func (r *response) isOk() bool {
	return r.responseType == RESPONSE_TYPE_OK
}

func (r *response) isWarning() bool {
	return r.responseType == RESPONSE_TYPE_WARNING
}

func (r *response) isError() bool {
	return r.responseType == RESPONSE_TYPE_ERROR
}

func (r *response) isFailure() bool {
	return r.isWarning() || r.isError()
}

func (r *response) parseFileInfos() (*fileInfos, error) {
	message := strings.ReplaceAll(r.message, "\n", "")
	parts := strings.Split(message, " ")
	if len(parts) < 3 {
		return nil, errors.New("unable to parse message as file infos")
	}

	size, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return &fileInfos{
		message:     r.message,
		permissions: parts[0],
		size:        int64(size),
		filename:    parts[2],
	}, nil
}

const (
	RESPONSE_TYPE_OK      responseType = 0
	RESPONSE_TYPE_WARNING responseType = 1
	RESPONSE_TYPE_ERROR   responseType = 2
)

func CopyFromLocal(ctx context.Context, host, src, dst string, authenticationMethod AuthenticationMethod) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	client, err := NewClient(host, authenticationMethod)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	errorChannel := make(chan error, 2)

	go func() {
		defer waitGroup.Done()

		writer, err := session.StdinPipe()
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer writer.Close()

		_, err = fmt.Fprintln(writer, "C0755", stat.Size(), filepath.Base(src))
		if err != nil {
			errorChannel <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errorChannel <- err
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			errorChannel <- err
			return
		}

		_, err = fmt.Fprint(writer, "\x00")
		if err != nil {
			errorChannel <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errorChannel <- err
			return
		}
	}()

	go func() {
		defer waitGroup.Done()

		err := session.Run(fmt.Sprintf("scp -qt %q", dst))
		if err != nil {
			errorChannel <- err
			return
		}
	}()

	if err = wait(ctx, &waitGroup); err != nil {
		return err
	}

	close(errorChannel)
	for err := range errorChannel {
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyFromRemote(ctx context.Context, writer io.Writer, host, remotePath string, authenticationMethod AuthenticationMethod) error {
	waitGroup := sync.WaitGroup{}
	errorChannel := make(chan error, 1)

	waitGroup.Add(1)
	go func() {
		var err error

		defer func() {
			waitGroup.Done()

			errorChannel <- err
		}()

		client, err := NewClient(host, authenticationMethod)
		if err != nil {
			errorChannel <- err
			return
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			errorChannel <- err
			return
		}
		defer session.Close()

		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			errorChannel <- err
			return
		}

		stdinPipe, err := session.StdinPipe()
		if err != nil {
			errorChannel <- err
			return
		}
		defer stdinPipe.Close()

		err = session.Start(fmt.Sprintf("%s -f %q", "scp", remotePath))
		if err != nil {
			errorChannel <- err
			return
		}

		err = ack(stdinPipe)
		if err != nil {
			errorChannel <- err
			return
		}

		response, err := parseResponse(stdoutPipe)
		if err != nil {
			errorChannel <- err
			return
		}
		if response.isFailure() {
			errorChannel <- errors.New(response.message)
			return
		}

		fileInfos, err := response.parseFileInfos()
		if err != nil {
			errorChannel <- err
			return
		}

		err = ack(stdinPipe)
		if err != nil {
			errorChannel <- err
			return
		}

		_, err = copyN(writer, stdoutPipe, fileInfos.size)
		if err != nil {
			errorChannel <- err
			return
		}

		err = ack(stdinPipe)
		if err != nil {
			errorChannel <- err
			return
		}

		err = session.Wait()
		if err != nil {
			errorChannel <- err
			return
		}
	}()

	if err := wait(ctx, &waitGroup); err != nil {
		return err
	}
	finalErr := <-errorChannel
	close(errorChannel)
	return finalErr
}

func wait(ctx context.Context, waitGroup *sync.WaitGroup) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		waitGroup.Wait()
	}()

	select {
	case <-c:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func checkResponse(reader io.Reader) error {
	response, err := parseResponse(reader)
	if err != nil {
		return err
	}

	if response.isFailure() {
		return errors.New(response.message)
	}

	return nil
}

func parseResponse(reader io.Reader) (response, error) {
	buffer := make([]uint8, 1)
	_, err := reader.Read(buffer)
	if err != nil {
		return response{}, err
	}

	responseType := buffer[0]
	message := ""
	if 0 < responseType {
		bufferedReader := bufio.NewReader(reader)
		message, err = bufferedReader.ReadString('\n')
		if err != nil {
			return response{}, err
		}
	}

	return response{responseType, message}, nil
}

func ack(writer io.Writer) error {
	var msg = []byte{0}
	n, err := writer.Write(msg)
	if err != nil {
		return err
	}
	if n < len(msg) {
		return errors.New("failed to write ack buffer")
	}
	return nil
}

func copyN(writer io.Writer, src io.Reader, size int64) (int64, error) {
	var total int64
	total = 0
	for total < size {
		n, err := io.CopyN(writer, src, size)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}
