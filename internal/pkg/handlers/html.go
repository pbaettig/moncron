package handlers

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/pbaettig/moncron/internal/pkg/model"
	"github.com/pbaettig/moncron/internal/pkg/store"
	"github.com/pbaettig/moncron/web/pages"
	log "github.com/sirupsen/logrus"
)

var (
	TemplateFuncs template.FuncMap
)

func init() {
	TemplateFuncs = make(template.FuncMap)

	TemplateFuncs["convertBytes"] = func(v int64) string {
		if v < 1024*1024 {
			return fmt.Sprintf("%d KB", v/1024)
		}
		if v < 1024*1024*1024 {
			return fmt.Sprintf("%.2f MB", float64(v)/1024/1024)
		}

		return fmt.Sprintf("%.2f GB", float64(v)/1024/1024/1024)
	}

	TemplateFuncs["since"] = func(t time.Time) string {
		return time.Since(t).String()
	}
	TemplateFuncs["mult"] = func(a, b int) int {
		return a * b
	}
	TemplateFuncs["sub"] = func(a, b int) int {
		return a - b
	}
	TemplateFuncs["add"] = func(a, b int) int {
		return a + b
	}
	TemplateFuncs["escapeTime"] = func(v time.Time) string {
		return url.QueryEscape(v.Format(time.RFC3339Nano))
	}

}

type JobRunsTableData struct {
	Runs        []model.JobRun
	Total       int
	PageSize    int
	PageNum     int
	TotalPages  int
	NextPageNum int
	NextPageURL string
}

type HtmlJobRunsTableHandler struct {
	store.JobRunStorer
}

func (h *HtmlJobRunsTableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := JobRunsTableData{}
	tmpl, err := template.New("JobRunsTable.html.tmpl").Funcs(TemplateFuncs).ParseFS(pages.FS, "JobRunsTable.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	r.URL.Query()
	params := getQueryParams(r)
	data.PageNum = params.page + 1
	data.PageSize = params.size
	nextPage := 0

	if params.others["host"] != "" && params.others["job"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByHostAndJob(params.others["job"], params.others["host"], params.page, params.size)
	} else if params.others["job"] != "" && params.others["before"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByNameBefore(params.others["job"], params.before, params.page, params.size)
	} else if params.others["host"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByHost(params.others["host"], params.page, params.size)
	} else if params.others["job"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByJob(params.others["job"], params.page, params.size)
	}
	data.TotalPages = int(math.Ceil(float64(data.Total) / float64(params.size)))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	if nextPage > 0 {
		data.NextPageNum = nextPage
		ps := r.URL.Query()
		ps.Set("p", strconv.Itoa(nextPage))
		r.URL.RawQuery = ps.Encode()
		data.NextPageURL = r.URL.RequestURI()

	}

	tmpl.Execute(w, data)

}

type jobRunDetailsData struct {
	Run         model.JobRun
	Others      []model.JobRun
	TotalOthers int
}

type HtmlJobRunDetailsHandler struct {
	store.JobRunStorer
}

func (h *HtmlJobRunDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := jobRunDetailsData{}

	tmpl, err := template.New("JobRunDetails.html.tmpl").Funcs(TemplateFuncs).ParseFS(pages.FS, "JobRunDetails.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	data.Run, err = h.Get(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	}

	data.Others, data.TotalOthers, _, _ = h.GetByNameBefore(data.Run.Name, data.Run.FinishedAt, 0, 7)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("couldn't render template: %s", err.Error())
	}
}

type indexData struct {
	JobNames   []string
	Hosts      []model.Host
	TotalJobs  int
	TotalHosts int
}

type HtmlIndexHandler struct {
	store.JobRunStorer
}

func (h *HtmlIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := indexData{}

	tmpl, err := template.ParseFS(pages.FS, "Index.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	data.JobNames, data.TotalJobs, _ = h.GetJobNames(0, 100)
	data.Hosts, data.TotalHosts, _ = h.GetHosts(0, 100)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("couldn't render template: %s", err.Error())
	}
}
