package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/feedhenry/negotiator/config"
	"github.com/feedhenry/negotiator/controller"
	"github.com/feedhenry/negotiator/domain/openshift"
	"github.com/feedhenry/negotiator/endpoint"
	pkgos "github.com/feedhenry/negotiator/pkg/openshift"
	"github.com/gorilla/mux"
)

// Builds the openshift service that encapsulates logic for creating the objects via the pkg/client
func buildPaasService() openshift.Service {
	host := os.Getenv("API_HOST")
	token := os.Getenv("API_TOKEN")
	clientConf := pkgos.BuildDefaultConfig(host, token)
	client, err := pkgos.NewClient(clientConf)
	if err != nil {
		log.Panic(err)
	}
	service := openshift.NewService(client)
	return service
}

func buildDeployController() controller.Deploy {
	return controller.NewDeployController(buildPaasService())
}

func buildSysHandler() SysHandler {
	return SysHandler{}
}

func buildDeployHandler() endpoint.Deploy {
	return endpoint.NewDeployHandler(logrus.StandardLogger(), buildDeployController(), &config.Conf{})
}

func buildHTTPHandler() http.Handler {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	r := mux.NewRouter().StrictSlash(true)
	//recovery middleware for any panics in the handlers
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	//add middleware for all routes
	n := negroni.New(recovery)
	// set up sys routes
	sysHandler := buildSysHandler()
	r.HandleFunc("/sys/info/health", sysHandler.Health)
	r.HandleFunc("/sys/info/ping", sysHandler.Ping)
	// set up deploy routes
	deployHandler := buildDeployHandler()
	r.HandleFunc("/deploy/cloudapp", deployHandler.Deploy)
	n.UseHandler(r)
	return n
}

func main() {
	httpHandler := buildHTTPHandler()
	port := ":3000"
	logrus.Info("starting negotiator on  port " + port)
	if err := http.ListenAndServe(port, httpHandler); err != nil {
		logrus.Fatal(err)
	}
}
