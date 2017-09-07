package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iotdog/daily-report/controllers"
	"github.com/urfave/negroni"
)

func SetWorkerRouters(subRouter *mux.Router) {
	addWorkerHandler := negroni.Wrap(http.HandlerFunc(controllers.AddWorker))
	subRouter.Handle("/add_worker", baseMiddleware.With(addWorkerHandler)).Methods(http.MethodPost, http.MethodOptions)
}
