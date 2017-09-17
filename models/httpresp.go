package models

type CommonResponse struct {
	Code int    `json:"code" description:"response code"`
	Msg  string `json:"message" description:"description of response code"`
}

type ReportInfo struct {
	WorkerName string   `json:"workerName" description:"worker name"`
	Tasks      []string `json:"tasks" description:"descriptions of daily tasks"`
	Plans      []string `json:"plans" description:"tomorrow plans"`
	ReportTime string   `json:"reportTime" description:"report time"`
}

type DailyReportResp struct {
	*CommonResponse
	Reports []ReportInfo `json:"reports" description:"all reports"`
}
