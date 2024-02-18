package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pbaettig/moncron/internal/pkg/store"
	"github.com/pbaettig/moncron/web/static"
)

func RegisterJobRunStorerHtmlRoutes(router *mux.Router, jobStore store.JobRunStorer) {
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServerFS(static.FS)))
	router.Path("/runs.html").Handler(&HtmlJobRunsTableHandler{jobStore})
	router.Path("/run.html").Queries("id", "").Handler(&HtmlJobRunDetailsHandler{jobStore})
	router.Path("/").Handler(&HtmlIndexHandler{jobStore})
}
