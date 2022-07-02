package metricsnginx

import (
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/util"
)

func CopyFiles(servers []remote.Server) error {
	for _, server := range servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		if err != nil {
			return util.HandleError(err)
		}

		err = nginx.CopyLogFiles(
			output.GetNginxMetricsDirPath(),
			server.Host,
			server.Nginx.Log.Access,
			server.Nginx.Log.Error,
			server.Nginx.Log.Persistence.Access,
			server.Nginx.Log.Persistence.Error,
			authenticationMethod,
		)
		if err != nil {
			return util.HandleError(err)
		}

		err = nginx.RotateLogFile(
			server.Host,
			server.Nginx.Log.Access,
			server.Nginx.Log.Persistence.Access,
			authenticationMethod,
		)
		if err != nil {
			return util.HandleError(err)
		}

		err = nginx.RotateLogFile(
			server.Host,
			server.Nginx.Log.Error,
			server.Nginx.Log.Persistence.Error,
			authenticationMethod,
		)
		if err != nil {
			return util.HandleError(err)
		}
	}

	return nil
}
