package models

type AddWorkerParams struct {
	Name  string `json:"name" description:"员工姓名"`
	Num   string `json:"number" description:"员工编号"`
	Dept  string `json:"dept" description:"部门名称"`
	Group string `json:"group" description:"小组名称"`
	Email string `json:"email" description:"员工邮箱"`
}

type SubmitReportParams struct {
	WorkerName string   `json:"name" description:"员工姓名"`
	MainLine   []string `json:"mainLine" description:"主线任务"`
	SubLine    []string `json:"subLine" description:"支线任务"`
	Plan       []string `json:"plan" description:"下一步计划"`
}
