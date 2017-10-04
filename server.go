package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iotdog/daily-report/configs"

	"github.com/iotdog/daily-report/models"

	"github.com/iotdog/daily-report/routers"

	"github.com/leesper/holmes"
	"github.com/urfave/negroni"
)

func main() {
	defer holmes.Start().Stop()

	if err := models.InitMongoDB(); err != nil {
		panic(err)
	}
	defer models.MongoCli.Close()

	go func() {
		// ticker := time.Tick(15 * time.Minute)
		ticker := time.Tick(15 * time.Second) // test
		now := time.Now()
		for {
			holmes.Infoln("running check at ", now)
			if models.DailyReportSent {
				holmes.Infoln("工作纪要已发送")
			} else {
				holmes.Infoln("工作纪要未发送")
			}
			holmes.Debugln("curr hour: ", now.Hour())
			if now.Hour() > 1 && now.Hour() < 8 {
				holmes.Infoln("reset state")
				models.DailyReportSent = false
			}
			if now.Hour() >= 23 {
				weekDay := time.Now().Weekday()
				if !models.DailyReportSent && weekDay != 6 && weekDay != 7 {
					holmes.Infoln("send daily report email")
					models.SendDailyReportMail()
				}
			}
			now = <-ticker
		}
	}()

	router := routers.InitRoutes()
	n := negroni.New()
	n.UseHandler(router)

	addr := fmt.Sprintf(":%d", configs.Instance().ListenPort)
	holmes.Errorln(http.ListenAndServe(addr, n))
}
