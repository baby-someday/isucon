package metricsnginx

import (
	"path"

	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/remote"
)

func CopyFiles(network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		if err != nil {
			return err
		}
		err = nginx.CopyLogFiles(
			server.Host,
			path.Join(getOutputPath(), server.Host, "access.log"),
			server.Nginx.Log.Access,
			path.Join(getOutputPath(), server.Host, "error.log"),
			server.Nginx.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
