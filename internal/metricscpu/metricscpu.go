package metricscpu

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"golang.org/x/crypto/ssh"
)

type process struct {
	client     *ssh.Client
	session    *ssh.Session
	vmstatFile *os.File
}

func MeasureMetrics(interval int, servers []remote.Server) error {
	processes := []process{}
	for _, server := range servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(
			server,
		)
		if err != nil {
			return err
		}

		client, session, err := remote.NewSession(
			server.Host,
			server.Environments,
			authenticationMethod,
		)
		if err != nil {
			return err
		}

		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			return err
		}

		vmstatFilePath := path.Join(output.GetCPUMetricsDirPath(), server.Host, "vmstat")
		err = os.MkdirAll(path.Dir(vmstatFilePath), 0755)
		if err != nil {
			return err
		}

		vmstatFile, err := os.Create(vmstatFilePath)
		if err != nil {
			return err
		}

		go io.Copy(vmstatFile, stdoutPipe)

		go session.Run(fmt.Sprintf("vmstat -n %d", interval))

		processes = append(processes, process{
			client:     client,
			session:    session,
			vmstatFile: vmstatFile,
		})
	}

	for {
		println("🤖    終了しますか？")
		println("👉    y/n")
		var in string
		fmt.Scan(&in)

		if in == "y" {
			break
		}
	}

	return nil
}
