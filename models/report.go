package models

import (
	"time"

	"github.com/iotdog/daily-report/configs"
	"github.com/leesper/holmes"

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
