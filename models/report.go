package models

import "time"

type Report struct {
	TodayReport  []string `bson:"todayReport"`  // 当天的工作，分为主线任务和支线任务
	TomorrowPlan []string `bson:"tomorrowPlan"` // 明天的计划
}

type DailyReport struct {
	Worker       string    `bson:"workerID"`     // 员工信息
	TodayReport  []string  `bson:"todayReport"`  // 当天的工作，分为主线任务和支线任务
	TomorrowPlan []string  `bson:"tomorrowPlan"` // 明天的计划
	CreatedAt    time.Time `bson:"createdAt"`    // 报告时间
	ReportIP     string    `bson:"reportIP"`     // 报告时使用电脑的网络IP
}

// SubmitReport 提交报告
func SubmitReport(input SubmitReportParams) interface{} {

	// worker := new(Worker)

	// if !nameExists {
	// 	return &CommonResponse{
	// 		Code: 1,
	// 		Msg:  "姓名错误",
	// 	}
	// }

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

	return nil
}
