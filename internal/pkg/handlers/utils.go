package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type queryParams struct {
	page   int
	size   int
	before time.Time
	others map[string]string
}

func getQueryParams(r *http.Request) (params queryParams) {
	q := r.URL.Query()
	params.others = make(map[string]string)
	params.page, _ = strconv.Atoi(q.Get("p"))
	q.Del("p")

	params.size, _ = strconv.Atoi(q.Get("size"))
	q.Del("size")

	if params.size == 0 {
		params.size = 50
	}

	var err error
	params.before, err = time.Parse(time.RFC3339Nano, r.URL.Query().Get("before"))
	if err != nil {
		params.before, err = time.Parse(time.RFC3339, r.URL.Query().Get("before"))
	}

	for k, vs := range q {
		if len(vs) == 0 {
			continue
		}
		params.others[k] = vs[0]
	}

	return
}

func convertToAnySlice[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

func getPreviousURLFromReferer(r *http.Request, defaultValue string) string {
	referer := r.Header.Get("Referer")
	if referer == "" || strings.HasSuffix(referer, r.RequestURI) {
		return defaultValue
	}

	return referer
}

func humanReadableBytes(v uint64, base uint64) string {
	mega := base * base
	giga := base * base * base

	if v < mega {
		return fmt.Sprintf("%d KB", v/1000)
	}
	if v < giga {
		return fmt.Sprintf("%.2f MB", float64(v)/float64(mega))
	}

	return fmt.Sprintf("%.2f GB", float64(v)/float64(giga))
}
