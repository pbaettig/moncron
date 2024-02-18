package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pbaettig/moncron/internal/pkg/model"
	"github.com/pbaettig/moncron/internal/pkg/store"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func shortResponse(success bool, msg string) any {
	return struct {
		Success bool
		Msg     string
	}{Success: success, Msg: msg}
}

func writeJsonResponse(v any, status int, w http.ResponseWriter) error {
	// w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	w.WriteHeader(status)
	return encoder.Encode(v)
}

// func getPagingParams(r *http.Request) (page int, size int, internalErr error, externalErr error) {
// 	query := r.URL.Query()

// 	p := query.Get("p")
// 	if p == "" {
// 		p = "0"
// 	}

// 	s := query.Get("size")
// 	if s == "" {
// 		s = "20"
// 	}

// 	page, internalErr = strconv.Atoi(p)
// 	if internalErr != nil {
// 		return page, size, internalErr, fmt.Errorf("cannot parse page number")
// 	}

// 	size, internalErr = strconv.Atoi(s)
// 	if internalErr != nil {
// 		return page, size, internalErr, fmt.Errorf("cannot parse size")
// 	}

// 	return page, size, nil, nil
// }

type HandlerResponse struct {
	total       int
	content     []any
	status      int
	internalErr error
	externalErr error
}

type ApiRequestResponder interface {
	Name() string
	Respond(r *http.Request) HandlerResponse
}

type ApiResponse struct {
	Total int
	Data  []any
}

type ApiHandler struct {
	ApiRequestResponder
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h.Respond(r)
	if resp.status == 0 {
		resp.status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var content any
	if resp.internalErr != nil {
		clientErr := resp.internalErr
		if resp.externalErr != nil {
			clientErr = resp.externalErr
		}

		content = struct{ Error string }{Error: clientErr.Error()}
	} else {
		content = ApiResponse{
			Total: resp.total,
			Data:  resp.content,
		}
	}

	w.WriteHeader(resp.status)
	encoder.Encode(content)
}

type ApiGetJobRunById struct {
	store.JobRunStorer
}

func (h ApiGetJobRunById) Name() string {
	return "get-run"
}

func (h ApiGetJobRunById) Respond(r *http.Request) (resp HandlerResponse) {
	var run model.JobRun
	vars := mux.Vars(r)
	id := vars["id"]

	run, resp.internalErr = h.Get(id)
	resp.content = []any{run}

	if resp.internalErr == store.ErrNotFound {
		resp.status = http.StatusNotFound
	}
	resp.total = 1
	return
}

type ApiAddJobRun struct {
	store.JobRunStorer
}

func (h ApiAddJobRun) Name() string {
	return "add-run"
}

func (h ApiAddJobRun) Respond(r *http.Request) (resp HandlerResponse) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	c := model.JobRun{}

	if resp.internalErr = decoder.Decode(&c); resp.internalErr != nil {
		resp.externalErr = fmt.Errorf("cannot decode request body")
		resp.status = http.StatusBadRequest
		return
	}

	if resp.internalErr = validate.Struct(&c); resp.internalErr != nil {
		resp.externalErr = fmt.Errorf("request body validation failed")
		resp.status = http.StatusBadRequest
		return
	}

	resp.internalErr = h.Add(c)
	resp.status = http.StatusCreated
	return
}

type ApiDeleteJobRun struct {
	store.JobRunStorer
}

func (h ApiDeleteJobRun) Name() string {
	return "delete-run"
}

func (h ApiDeleteJobRun) Respond(r *http.Request) (resp HandlerResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp.internalErr = h.Delete(id)
	if resp.internalErr != nil {
		resp.status = http.StatusInternalServerError
		if resp.internalErr == store.ErrNotFound {
			resp.externalErr = resp.internalErr
			resp.status = http.StatusNotFound
		}
		return
	}

	resp.status = http.StatusNoContent
	resp.content = []any{shortResponse(true, "deleted")}
	return
}

type ApiGetJobRunsByName struct {
	store.JobRunStorer
}

func (h ApiGetJobRunsByName) Name() string {
	return "get-runs-by-job"
}

func (h ApiGetJobRunsByName) Respond(r *http.Request) (resp HandlerResponse) {
	var runs []model.JobRun
	qp := getQueryParams(r)
	runs, resp.total, _, resp.internalErr = h.GetByJob(r.URL.Query().Get("job"), qp.page, qp.size)
	resp.content = convertToAnySlice(runs)
	return
}

type ApiGetJobRunsByNameBefore struct {
	store.JobRunStorer
}

func (h ApiGetJobRunsByNameBefore) Name() string {
	return "get-runs-by-job-before"
}

func (h ApiGetJobRunsByNameBefore) Respond(r *http.Request) (resp HandlerResponse) {
	var runs []model.JobRun
	qp := getQueryParams(r)

	runs, resp.total, _, resp.internalErr = h.GetByNameBefore(r.URL.Query().Get("job"), qp.before, qp.page, qp.size)
	resp.content = convertToAnySlice(runs)
	return
}

type ApiGetJobRunsByHost struct {
	store.JobRunStorer
}

func (h ApiGetJobRunsByHost) Name() string {
	return "get-runs-by-host"
}

func (h ApiGetJobRunsByHost) Respond(r *http.Request) (resp HandlerResponse) {
	var runs []model.JobRun
	qp := getQueryParams(r)

	runs, resp.total, _, resp.internalErr = h.GetByHost(r.URL.Query().Get("host"), qp.page, qp.size)
	resp.content = convertToAnySlice(runs)
	return
}

type ApiGetJobRunsByHostAndName struct {
	store.JobRunStorer
}

func (h ApiGetJobRunsByHostAndName) Name() string {
	return "get-runs-by-host-and-name"
}

func (h ApiGetJobRunsByHostAndName) Respond(r *http.Request) (resp HandlerResponse) {
	var runs []model.JobRun

	qp := getQueryParams(r)

	runs, resp.total, _, resp.internalErr = h.GetByHostAndJob(qp.others["job"], qp.others["host"], qp.page, qp.size)
	resp.content = convertToAnySlice(runs)
	return
}

type ApiGetJobNames struct {
	store.JobRunStorer
}

func (h ApiGetJobNames) Name() string {
	return "get-job-names"
}

func (h ApiGetJobNames) Respond(r *http.Request) (resp HandlerResponse) {
	var names []string
	qp := getQueryParams(r)

	names, resp.total, resp.internalErr = h.GetJobNames(qp.page, qp.size)
	resp.content = convertToAnySlice(names)
	return
}

type ApiGetHosts struct {
	store.JobRunStorer
}

func (h ApiGetHosts) Name() string {
	return "get-hosts"
}

func (h ApiGetHosts) Respond(r *http.Request) (resp HandlerResponse) {
	var hosts []model.Host
	qp := getQueryParams(r)

	hosts, resp.total, resp.internalErr = h.GetHosts(qp.page, qp.size)
	resp.content = convertToAnySlice(hosts)
	return
}
