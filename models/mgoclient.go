package models

import (
	"github.com/iotdog/daily-report/configs"

	mgo "gopkg.in/mgo.v2"
)

var (
	MongoCli *mgo.Session
)

func InitMongoDB() error {
	var err error
	MongoCli, err = mgo.Dial(configs.Instance().MongoDBAddr)
	if err != nil {
		return err
	}
	return nil
}
