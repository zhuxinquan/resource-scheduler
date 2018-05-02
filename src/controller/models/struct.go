package models

//CGroup中各个子系统的指标
type SubSystemMetric struct {
	SubSystem string            `json:"subSystem"`
	Metric    map[string]string `json:"metric"`
}

//执行EXEC的请求体
type CGExecReq struct {
	Path            string            `json:"path"`
	SubSystemMetric []SubSystemMetric `json:"subSystemMetric"`
	Cmd             string            `json:"cmd"`
	User            string            `json:"user"`
}

type SetMetricData struct {
	SubSystemMetrics []SubSystemMetric `json:"subSystemMetric"`
	Path             string            `json:"path"`
}
