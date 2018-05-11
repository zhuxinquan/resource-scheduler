package models

//CGroup中各个子系统的指标
type SubSystemMetric struct {
	SubSystem string            `json:"subSystem"`
	Metric    map[string]string `json:"metric"`
}

//执行EXEC的请求体
type CGExecReq struct {
	Path string `json:"path"`
	Cmd  string `json:"cmd"`
	User string `json:"user"`
}

type SetMetricData struct {
	SubSystemMetrics []SubSystemMetric `json:"subSystemMetric"`
	Path             string            `json:"path"`
}

type SysInfo struct {
	CpuUserUse string `json:"cpuUserUse"`
	CpuSysUse  string `json:"cpuSysUse"`
	MemTotal   int64  `json:"memTotal"`
	MemFree    int64  `json:"memFree"`
	MemUsed    int64  `json:"memUsed"`
	MemRate    string `json:"memRate"`
	SwapTotal  int64  `json:"swapTotal"`
	SwapFree   int64  `json:"swapFree"`
	SwapUsed   int64  `json:"swapUsed"`
	SwapRate   string `json:"swapRate"`
}

type GroupInfo struct {
	GroupPath  string   `json:"groupPath"`
	SubSystems []string `json:"subSystems"`
}

type ProcessInfo struct {
	Pid   string `json:"pid"`
	User  string `json:"user"`
	State string `json:"state"`
	Cpu   string `json:"cpu"`
	Mem   string `json:"mem"`
	Cmd   string `json:"cmd"`
}
