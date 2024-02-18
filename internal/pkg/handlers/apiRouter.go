package handlers

import (
	"github.com/gorilla/mux"
	"github.com/pbaettig/moncron/internal/pkg/store"
)

func RegisterJobRunStorerApiRoutes(router *mux.Router, jobStore store.JobRunStorer) {
	hostsHandler := &ApiHandler{ApiGetHosts{jobStore}}
	router.Path("/hosts").Methods("GET").Handler(hostsHandler)

	jobNamesHandler := &ApiHandler{ApiGetJobNames{jobStore}}
	router.Path("/jobs").Methods("GET").Handler(jobNamesHandler)

	runsByHostAndNameHandler := &ApiHandler{ApiGetJobRunsByHostAndName{jobStore}}
	router.Path("/runs").Methods("GET").Queries("host", "", "job", "").Handler(runsByHostAndNameHandler)

	runsByNameHandler := &ApiHandler{ApiGetJobRunsByName{jobStore}}
	router.Path("/runs").Methods("GET").Queries("job", "").Handler(runsByNameHandler)

	runsByNameBeforeHandler := &ApiHandler{ApiGetJobRunsByNameBefore{jobStore}}
	router.Path("/runs").Methods("GET").Queries("job", "", "before", "").Handler(runsByNameBeforeHandler)

	runsByHostHandler := &ApiHandler{ApiGetJobRunsByHost{jobStore}}
	router.Path("/runs").Methods("GET").Queries("host", "").Handler(runsByHostHandler)

	getJobRunHandler := &ApiHandler{ApiGetJobRunById{jobStore}}
	router.Path("/runs/{id}").Methods("GET").Handler(getJobRunHandler)

	deleteJobRunHandler := &ApiHandler{ApiDeleteJobRun{jobStore}}
	router.Path("/runs/{id}").Methods("DELETE").Handler(deleteJobRunHandler)

	addRunHandler := &ApiHandler{ApiAddJobRun{jobStore}}
	router.Path("/runs").Methods("POST").Handler(addRunHandler)

}
