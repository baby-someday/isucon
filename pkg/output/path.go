package output

import "path"

func GetDistributeOutputDirPath() string {
	return path.Join("output", "distribute")
}

func GetNginxMetricsDirPath() string {
	return path.Join("output", "metrics", "nginx")
}
