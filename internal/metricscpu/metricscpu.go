package metricscpu

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/servermaster"
	"github.com/baby-someday/isucon/pkg/util"
	"golang.org/x/crypto/ssh"
)

type process struct {
	client     *ssh.Client
	session    *ssh.Session
	vmstatFile *os.File
}

func MeasureMetrics(
	interval int,
	serverMasters []servermaster.ServerMaster,
	servers []remote.Server,
) error {
	processes := []process{}
	for _, server := range servers {
		serverMaster, err := servermaster.FindServerMaster(
			server.Name,
			serverMasters,
		)
		if err != nil {
			return util.HandleError(err)
		}

		authenticationMethod, err := remote.MakeAuthenticationMethod(serverMaster.SSH)
		if err != nil {
			return util.HandleError(err)
		}

		client, session, err := remote.NewSession(
			serverMaster.Host,
			server.Environments,
			authenticationMethod,
		)
		if err != nil {
			return util.HandleError(err)
		}

		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			return util.HandleError(err)
		}

		vmstatFilePath := path.Join(output.GetCPUMetricsDirPath(), serverMaster.Host, "vmstat")
		err = os.MkdirAll(path.Dir(vmstatFilePath), 0755)
		if err != nil {
			return util.HandleError(err)
		}

		vmstatFile, err := os.Create(vmstatFilePath)
		if err != nil {
			return util.HandleError(err)
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
