package output

import "path"

func GetCPUMetricsDirPath() string {
	return path.Join("output", "metrics", "cpu")
}

func GetDistributeOutputDirPath() string {
	return path.Join("output", "distribute")
}

func GetNginxMetricsDirPath() string {
	return path.Join("output", "metrics", "nginx")
}
