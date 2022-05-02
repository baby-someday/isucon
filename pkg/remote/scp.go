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
	"sync"
)

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

const (
	RESPONSE_TYPE_OK      responseType = 0
	RESPONSE_TYPE_WARNING responseType = 1
	RESPONSE_TYPE_ERROR   responseType = 2
)

func Copy(ctx context.Context, host, src, dst string, authenticationMethod AuthenticationMethod) error {
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
