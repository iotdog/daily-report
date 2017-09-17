package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/iotdog/daily-report/configs"
	"github.com/iotdog/json2table/j2t"
	"github.com/leesper/holmes"

	gomail "gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2/bson"
)

type DailyReport struct {
	ID           bson.ObjectId `bson:"_id"`          // 日报ID
	Worker       string        `bson:"workerID"`     // 员工信息
	MainLine     []string      `bson:"mainLine"`     // 当天的工作（主线任务）
	SubLine      []string      `bson:"subLine"`      // 当天的工作（支线任务）
	TomorrowPlan []string      `bson:"tomorrowPlan"` // 明天的计划
	CreatedAt    time.Time     `bson:"createdAt"`    // 报告时间
	UpdatedAt    time.Time     `bson:"updatedAt"`    // 更新时间
	ReportIP     string        `bson:"reportIP"`     // 报告时使用电脑的网络IP
}

func (dr *DailyReport) Find(query bson.M) error {
	return MongoCli.DB(configs.Instance().DBName).C(configs.Instance().DailyReportC).Find(query).Sort("-createdAt").One(dr)
}

func (dr *DailyReport) Update() error {
	return MongoCli.DB(configs.Instance().DBName).C(configs.Instance().DailyReportC).Update(bson.M{"_id": dr.ID}, dr)
}

func (dr *DailyReport) Insert() error {
	return MongoCli.DB(configs.Instance().DBName).C(configs.Instance().DailyReportC).Insert(dr)
}

// SubmitReport 提交报告
func SubmitReport(input SubmitReportParams, ip string) interface{} {
	// 检查输入
	worker := new(Worker)
	err := worker.Find(bson.M{"workerName": input.WorkerName})
	if err != nil {
		return &CommonResponse{
			Code: 1,
			Msg:  "员工不存在",
		}
	}

	if len(input.MainLine) == 0 {
		return &CommonResponse{
			Code: 1,
			Msg:  "未填写工作内容",
		}
	}

	if len(input.Plan) == 0 {
		weekDay := time.Now().Weekday()
		if weekDay != 5 && weekDay != 6 && weekDay != 7 {
			return &CommonResponse{
				Code: 1,
				Msg:  "未填写工作计划",
			}
		}
	}

	todayBegin := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Local)

	// 检查该员工是否已发送日报
	report := new(DailyReport)
	cond := bson.M{
		"$and": []bson.M{bson.M{
			"workerID": worker.WorkerID,
		},
			bson.M{
				"updatedAt": bson.M{
					"$gt": todayBegin,
				},
			},
		},
	}
	err = report.Find(cond)
	report.MainLine = input.MainLine
	report.SubLine = input.SubLine
	report.TomorrowPlan = input.Plan
	report.UpdatedAt = time.Now()
	report.ReportIP = ip
	if err != nil {
		holmes.Errorln(err)
		holmes.Infoln("create new report")
		report.ID = bson.NewObjectId()
		report.Worker = worker.WorkerID
		report.CreatedAt = time.Now()
		report.UpdatedAt = time.Now()
		err = report.Insert()
		if err != nil {
			holmes.Errorln("failed to insert new report: ", err)
		}
	} else {
		holmes.Infoln("update report")
		err = report.Update()
		if err != nil {
			holmes.Errorln("failed to update new report: ", err)
		}
	}

	return &CommonResponse{
		Code: 0,
		Msg:  "提交成功",
	}
}

func getTodayReports() []DailyReport {
	reports := []DailyReport{}
	todayBegin := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 0, 0, 0, time.Local)
	err := MongoCli.DB(configs.Instance().DBName).C(configs.Instance().DailyReportC).Find(
		bson.M{"updatedAt": bson.M{
			"$gt": todayBegin,
		}}).All(&reports)
	if err != nil {
		holmes.Errorln("failed to get reports: ", err)
	}
	return reports
}

func getWorkers() []Worker {
	workers := []Worker{}
	err := MongoCli.DB(configs.Instance().DBName).C(configs.Instance().WorkerC).Find(bson.M{"leave": false}).Sort("workerID").All(&workers)
	if err != nil {
		holmes.Errorln("failed to get workers: ", err)
	}
	return workers
}

type JSONReport struct {
	WorkerName string `json:"姓名"`
	Department string `json:"部门"`
	Group      string `json:"小组"`
	Tasks      string `json:"工作内容"`
	Plans      string `json:"明日计划"`
	ReportTime string `json:"报告时间"`
}

func genHtmlTable(reports []DailyReport, workers []Worker) string {
	if 0 == len(reports) || 0 == len(workers) {
		holmes.Errorln("invalid input")
		return ""
	}
	jsonReports := []JSONReport{}
	for _, worker := range workers {
		jsonReport := new(JSONReport)
		jsonReport.WorkerName = worker.WorkerName
		jsonReport.Department = worker.Department
		jsonReport.Group = worker.Group
		jsonReport.ReportTime = ""
		jsonReport.Tasks = ""
		jsonReport.Plans = ""
		for _, report := range reports {
			if report.Worker == worker.WorkerID {
				jsonReport.ReportTime = report.UpdatedAt.Format("2006-01-02 15:04:05")
				tmp := append(report.MainLine, report.SubLine...)
				jsonReport.Tasks = strings.Join(tmp, "<br />")
				jsonReport.Plans = strings.Join(report.TomorrowPlan, "<br />")
				break
			}
		}
		jsonReports = append(jsonReports, *jsonReport)
	}
	tmp, err := json.Marshal(jsonReports)
	if err != nil {
		return ""
	}

	_, html := j2t.JSON2HtmlTable(string(tmp), []string{"姓名", "部门", "小组",
		"工作内容", "明日计划", "报告时间"})
	holmes.Infoln(html)
	return html
}

func sendDailyReportMail(reportTable string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "")
	m.SetHeader("To", "", "")
	m.SetHeader("Subject", fmt.Sprintf("工作纪要%d.%d.%d", time.Now().Year(), time.Now().Month(), time.Now().Day()))
	m.SetBody("text/html", reportTable)
	d := gomail.NewDialer(configs.Instance().MailBoxSMTP, configs.Instance().MailBoxPort,
		configs.Instance().MailBoxUserName, configs.Instance().MailBoxPwd)
	d.SSL = configs.Instance().MailBoxSSL

	err := d.DialAndSend(m)
	holmes.Errorln(err)

	return err
}

func GetDailyReport() interface{} {
	reportInfos := []ReportInfo{}
	reports := getTodayReports()
	if 0 == len(reports) {
		return DailyReportResp{
			&CommonResponse{
				Code: 0,
				Msg:  "请求成功",
			},
			reportInfos,
		}
	}
	workers := getWorkers()
	if 0 == len(workers) {
		return DailyReportResp{
			&CommonResponse{
				Code: 0,
				Msg:  "请求成功",
			},
			reportInfos,
		}
	}
	getAllReports := true
	for _, worker := range workers {
		reportInfo := new(ReportInfo)
		reportInfo.WorkerName = worker.WorkerName
		for _, report := range reports {
			if report.Worker == worker.WorkerID {
				reportInfo.Plans = report.TomorrowPlan
				reportInfo.ReportTime = report.UpdatedAt.String()
				reportInfo.Tasks = append(report.MainLine, report.SubLine...)
				break
			}
		}
		if 0 == len(reportInfo.Tasks) {
			getAllReports = false
		}
		reportInfos = append(reportInfos, *reportInfo)
	}

	genHtmlTable(reports, workers)
	if getAllReports { // send daily report email if all workers have submitted their reports
		sendDailyReportMail(genHtmlTable(reports, workers))
	} else {
		holmes.Infoln("somebody does not submit his report")
	}

	return DailyReportResp{
		&CommonResponse{
			Code: 0,
			Msg:  "请求成功",
		},
		reportInfos,
	}
}
