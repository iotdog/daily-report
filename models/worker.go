package models

import (
	"fmt"
	"time"

	"github.com/iotdog/daily-report/configs"

	"gopkg.in/mgo.v2/bson"
)

type Worker struct {
	WorkerID   string    `bson:"workerID"`   // 员工编号
	WorkerName string    `bson:"workerName"` // 员工姓名
	Department string    `bson:"department"` // 部门名称
	Group      string    `bson:"group"`      // 小组名称
	Email      string    `bson:"email"`      // 员工邮箱
	Leave      bool      `bson:"leave"`      // 是否离职：离职（true）、在职（false）
	CreatedAt  time.Time `bson:"createdAt"`
}

func (w *Worker) Find(query bson.M) error {
	return MongoCli.DB(configs.Instance().DBName).C(configs.Instance().WorkerC).Find(query).One(w)
}

func (w *Worker) Insert() error {
	return MongoCli.DB(configs.Instance().DBName).C(configs.Instance().WorkerC).Insert(w)
}

func AddWorker(input AddWorkerParams) interface{} {
	if input.Dept == "" || input.Name == "" || input.Num == "" {
		return &CommonResponse{
			Code: 1,
			Msg:  "信息填写不完整",
		}
	}
	worker := new(Worker)
	err := worker.Find(bson.M{"workerID": input.Num})
	if err == nil {
		return &CommonResponse{
			Code: 1,
			Msg:  fmt.Sprintf("员工已存在：%s - %s", worker.WorkerID, worker.WorkerName),
		}
	}
	worker.CreatedAt = time.Now()
	worker.Department = input.Dept
	worker.Group = input.Group
	worker.WorkerID = input.Num
	worker.WorkerName = input.Name
	worker.Email = input.Email
	worker.Leave = false
	err = worker.Insert()
	if err != nil {
		return &CommonResponse{
			Code: 1,
			Msg:  "保存员工信息失败",
		}
	}
	return &CommonResponse{
		Code: 0,
		Msg:  "添加成功",
	}
}
