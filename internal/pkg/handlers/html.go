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

	TemplateFuncs["convertBytesSI"] = func(v uint64) string {
		return humanReadableBytes(v, 1000)
	}
	TemplateFuncs["convertBytes"] = func(v uint64) string {
		return humanReadableBytes(v, 1024)
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

type count struct {
	Successful int
	Failed     int
	Denied     int
	Total      int
}

type JobRunsTableData struct {
	Title       string
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
		data.Title = fmt.Sprintf(`Runs of "%s" on "%s"`, params.others["job"], params.others["host"])
	} else if params.others["job"] != "" && params.others["before"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByNameBefore(params.others["job"], params.before, params.page, params.size)
		data.Title = fmt.Sprintf(`Runs of "%s" before %s`, params.others["job"], params.others["before"])
	} else if params.others["host"] != "" {
		data.Runs, data.Total, nextPage, err = h.GetByHost(params.others["host"], params.page, params.size)
		data.Title = fmt.Sprintf(`Runs on "%s"`, params.others["host"])

	} else if params.others["job"] != "" {
		data.Title = fmt.Sprintf(`Runs of "%s"`, params.others["job"])

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
	Run               model.JobRun
	PreviousRuns      []model.JobRun
	TotalPreviousRuns int
	PreviousURL       string
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
	data.PreviousURL = getPreviousURLFromReferer(r, fmt.Sprintf("/runs.html?job=%s", data.Run.Name))

	data.PreviousRuns, data.TotalPreviousRuns, _, _ = h.GetByNameBefore(data.Run.Name, data.Run.FinishedAt, 0, 7)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("couldn't render template: %s", err.Error())
	}
}

type indexData struct {
	JobNames   []string
	Hosts      []string
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
	data.Hosts, data.TotalHosts, _ = h.GetHostNames(0, 100)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Errorf("couldn't render template: %s", err.Error())
	}
}

type jobDetails struct {
	MinDuration    time.Duration
	MinDurationID  string
	MaxDuration    time.Duration
	MaxDurationID  string
	MinRSS         uint64
	MinRSSID       string
	MaxRSS         uint64
	MaxRSSID       string
	FailedRuns     int
	SuccessfulRuns int
	DeniedRuns     int
	TotalRuns      int
	HostNames      []string
	JobName        string
	JobRunsCount   map[string]*count
	Total          count
	PreviousURL    string
}

type HtmlJobDetailsHandler struct {
	store.JobRunStorer
}

func (h *HtmlJobDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		page int = 0
		next int = 1
		runs []model.JobRun
		err  error
		jd   jobDetails
	)

	tmpl, err := template.New("JobDetails.html.tmpl").Funcs(TemplateFuncs).ParseFS(pages.FS, "JobDetails.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	jd.MinDuration = time.Duration(math.MaxInt64)
	jd.MinRSS = math.MaxInt64
	jd.JobName = r.URL.Query().Get("job")
	jd.PreviousURL = getPreviousURLFromReferer(r, fmt.Sprintf("/runs.html?job=%s", jd.JobName))
	jd.JobRunsCount = make(map[string]*count)
	// jd.HostNames = make([]string, 0)
	// hostNames := make(map[string]struct{})
	log.Debugf("MinRss: %d, MinDuration: %s", jd.MinRSS, jd.MinDuration)

	for next > 0 {
		runs, _, next, err = h.GetByJob(jd.JobName, page, 10000)
		page = next

		if err != nil {
			log.Error(err)
			break
		}
		for _, run := range runs {
			if _, ok := jd.JobRunsCount[run.Host.Name]; !ok {
				jd.JobRunsCount[run.Host.Name] = new(count)
			}

			if run.Status == model.ProcessStartedNormally {
				if run.Result.ExitCode == 0 {
					jd.JobRunsCount[run.Host.Name].Successful++
					jd.Total.Successful++
				} else {
					jd.JobRunsCount[run.Host.Name].Failed++
					jd.Total.Failed++
				}
			} else {
				jd.JobRunsCount[run.Host.Name].Denied++
				jd.Total.Denied++

			}

			jd.JobRunsCount[run.Host.Name].Total++
			jd.Total.Total++

			if run.Result.WallTime < jd.MinDuration {
				jd.MinDuration = run.Result.WallTime
				jd.MinDurationID = run.ID
			}
			if run.Result.WallTime > jd.MaxDuration {
				jd.MaxDuration = run.Result.WallTime
				jd.MaxDurationID = run.ID
			}
			if run.Result.MaxRssBytes < jd.MinRSS {
				jd.MinRSS = run.Result.MaxRssBytes
				jd.MinRSSID = run.ID
			}
			if run.Result.MaxRssBytes > jd.MaxRSS {
				jd.MaxRSS = run.Result.MaxRssBytes
				jd.MaxRSSID = run.ID
			}

		}
	}

	jd.TotalRuns = jd.SuccessfulRuns + jd.FailedRuns + jd.DeniedRuns

	err = tmpl.Execute(w, jd)
	if err != nil {
		log.Error(err)
	}

}

type hostDetails struct {
	Host         model.Host
	JobRunsCount map[string]*count
	Total        count
	PreviousURL  string
}

type HtmlHostDetailsHandler struct {
	store.JobRunStorer
}

func (h *HtmlHostDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		page int = 0
		next int = 1
		runs []model.JobRun
		err  error
		hd   hostDetails
	)

	tmpl, err := template.New("HostDetails.html.tmpl").Funcs(TemplateFuncs).ParseFS(pages.FS, "HostDetails.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	hostName := r.URL.Query().Get("host")
	hd.PreviousURL = getPreviousURLFromReferer(r, fmt.Sprintf("/runs.html?host=%s", hostName))
	hd.JobRunsCount = make(map[string]*count)

	for next > 0 {
		runs, _, next, err = h.GetByHost(hostName, page, 1000)
		page = next
		if err != nil {
			log.Error(err)
			break
		}
		if len(runs) > 0 {
			hd.Host = runs[0].Host
		}
		log.Debugf("Runs: %d, Current page: %d, Next page: %d", len(runs), page, next)
		for _, run := range runs {
			if _, ok := hd.JobRunsCount[run.Name]; !ok {
				hd.JobRunsCount[run.Name] = new(count)
			}

			if run.Status == model.ProceessNotStarted {
				hd.JobRunsCount[run.Name].Denied++
				hd.Total.Denied++
			} else {
				if run.Result.ExitCode == 0 {
					hd.JobRunsCount[run.Name].Successful++
					hd.Total.Successful++

				} else {
					hd.JobRunsCount[run.Name].Failed++
					hd.Total.Failed++
				}
			}

			hd.JobRunsCount[run.Name].Total++
			hd.Total.Total++
		}

	}

	err = tmpl.Execute(w, hd)
	if err != nil {
		log.Error(err)
	}

}
