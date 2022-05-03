package metricsnginx

import "path"

func getOutputPath() string {
	return path.Join("output", "metrics", "nginx")
}
