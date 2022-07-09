package output

import "path"

func GetCPUMetricsDirPath() string {
	return path.Join("output", "metrics", "cpu")
}

func GetDistributeOutputDirPath() string {
	return path.Join("output", "distribute")
}

func GetMySQLMetricsDirPath() string {
	return path.Join("output", "metrics", "mysql")
}

func GetNginxMetricsDirPath() string {
	return path.Join("output", "metrics", "nginx")
}

func GetNginxAnalysisDirPath() string {
	return path.Join("output", "analysis", "nginx")
}
