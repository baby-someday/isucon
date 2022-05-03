package metricsnginx

import (
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
)

func CopyFiles(network remote.Network) error {
	return nginx.CopyLogFiles(
		output.GetNginxMetricsDirPath(),
		network,
	)
}
