package main

import (
	"fmt"
	"net/http"

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

	router := routers.InitRoutes()
	n := negroni.New()
	n.UseHandler(router)

	addr := fmt.Sprintf(":%d", configs.Instance().ListenPort)
	holmes.Errorln(http.ListenAndServe(addr, n))
}
