package routers

import (
	"net/http"

	"github.com/iotdog/daily-report/controllers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	baseMiddleware = negroni.New(
		negroni.HandlerFunc(controllers.CorsMiddleware),
		negroni.HandlerFunc(controllers.LoggerMiddleware),
	)
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./html")))) // for static pages
	apiV1Router := router.PathPrefix("/api/1.0").Subrouter()                                             // for apis of version 1.0
	SetWorkerRouters(apiV1Router)
	return router
}
