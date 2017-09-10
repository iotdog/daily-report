package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iotdog/daily-report/controllers"
	"github.com/urfave/negroni"
)

func SetReportRouters(subRouter *mux.Router) {
	submitReportHandler := negroni.Wrap(http.HandlerFunc(controllers.SubmitReport))
	subRouter.Handle("/submit_report", baseMiddleware.With(submitReportHandler)).Methods(http.MethodPost, http.MethodOptions)
}
