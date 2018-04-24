package models

//CGroup指标
type CgroupMetric struct {
	SubSystemMetric []SubSystemMetric `json:"subSystemMetric"`
}

//CGroup中各个子系统的指标
type SubSystemMetric struct {
	SubSystem string            `json:"subSystem"`
	Metric    map[string]string `json:"metric"`
}

//执行EXEC的请求体
type CGExecReq struct {
}
