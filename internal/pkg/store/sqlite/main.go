package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pbaettig/moncron/internal/pkg/model"
)

type DB struct {
	*sql.DB
	stmt       sq.StatementBuilderType
	addChannel chan model.JobRun
	errChannel chan error
}

type colNameValuePair struct {
	name  string
	value any
}

type where struct {
	pred string
	vals any
}

func selectFromDB[T any](db *DB, table string, cols []string, where []where, page int, size int) (records []T, total int, next int, err error) {
	// prepare a Query that gets the total amount of results
	// regardless of page size settings
	totalCol := "*"
	if len(cols) == 1 {
		totalCol = cols[0]
	}
	totalQ := db.stmt.Select(fmt.Sprintf(`COUNT(%s)`, totalCol)).From(table)
	// prepare the query that retrieves data in pages
	pageQ := db.stmt.Select(cols...).From(table).Limit(uint64(size)).Offset(uint64(page * size)).OrderByClause(`json_extract(data, "$.FinishedAt") DESC`)
	for _, w := range where {
		pageQ = pageQ.Where(w.pred, w.vals)
		totalQ = totalQ.Where(w.pred, w.vals)
	}
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()
	rows, err := totalQ.Query()
	if err != nil {
		return
	}
	defer rows.Close()

	total = 0
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return
		}
	}
	rows.Close()

	if total == 0 {
		return
	}

	rows, err = pageQ.Query()
	if err != nil {
		return
	}
	defer rows.Close()

	buf := make([]byte, 0)
	for rows.Next() {
		var r T

		if err = rows.Scan(&buf); err != nil {
			return
		}

		if err = json.Unmarshal(buf, &r); err != nil {
			return
		}

		records = append(records, r)
		clear(buf)
	}
	if (page+1)*size < total {
		next = page + 1
	}

	return
}

func (db *DB) insert(j model.JobRun) error {
	buf, err := json.Marshal(j)
	if err != nil {
		return err
	}

	idBlob, err := uuid.MustParse(j.ID).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.stmt.Insert("jobRuns").Columns("id", "jobName", "data").Values(idBlob, j.Job.Name, buf).Exec()
	return err
}

func (db *DB) Add(j model.JobRun) error {
	db.addChannel <- j

	return <-db.errChannel
}

func (db *DB) GetByHost(hostName string, page int, size int) (runs []model.JobRun, total int, next int, err error) {
	runs, total, next, err = selectFromDB[model.JobRun](
		db,
		"jobRuns",
		[]string{"data"},
		[]where{{`json_extract(data, "$.Host.Name") = ?`, hostName}}, page, size)

	return
}

func (db *DB) GetByJob(jobName string, page int, size int) (runs []model.JobRun, total int, next int, err error) {
	// runs, total, err = selectFrom[model.JobRun](db, "jobRuns", []string{"data"}, []string{`jobName`, jobName}, page, size)
	runs, total, next, err = selectFromDB[model.JobRun](
		db,
		"jobRuns",
		[]string{"data"},
		[]where{{`jobName = ?`, jobName}},
		// []string{`json_extract(data, "$.FinishedAt")`},
		page,
		size,
	)
	return
}

func (db *DB) GetByHostAndJob(jobName string, hostName string, page int, size int) (runs []model.JobRun, total int, next int, err error) {
	runs, total, next, err = selectFromDB[model.JobRun](
		db,
		"jobRuns",
		[]string{"data"},
		[]where{{`json_extract(data, "$.Host.Name") = ?`, hostName}, {`jobName = ?`, jobName}},
		// []string{`json_extract(data, "$.FinishedAt")`},
		page,
		size,
	)
	return
}

func (db *DB) GetByNameBefore(jobName string, t time.Time, page int, size int) (runs []model.JobRun, total int, next int, err error) {
	runs, total, next, err = selectFromDB[model.JobRun](
		db,
		"jobRuns",
		[]string{"data"},
		[]where{{`jobName = ?`, jobName}, {`json_extract(data, "$.FinishedAt") < ?`, t.Format(time.RFC3339Nano)}},
		// []string{`json_extract(data, "$.FinishedAt")`},
		page,
		size,
	)
	return
}

func (db *DB) Delete(id string) error {
	idBlob, err := uuid.MustParse(id).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = db.stmt.Delete("jobRuns").Where("id = ?", idBlob).Exec()
	// _, err = db.Exec(`DELETE FROM jobRuns WHERE id = ?;`, idBlob)
	return err
}

func (db *DB) Get(id string) (run model.JobRun, err error) {

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		return
	}
	idBlob, err := uuidValue.MarshalBinary()
	if err != nil {
		return
	}
	runs, _, _, err := selectFromDB[model.JobRun](
		db,
		"jobRuns",
		[]string{"data"},
		[]where{{`id = ?`, idBlob}},
		0,
		1,
	)
	if err != nil {
		return
	}
	if len(runs) == 0 {
		return run, fmt.Errorf("no run with id %s found", id)
	}

	run = runs[0]
	return
}

func (db *DB) GetHosts(page int, size int) (hosts []model.Host, total int, err error) {
	hosts, total, _, err = selectFromDB[model.Host](
		db,
		"jobRuns",
		[]string{`DISTINCT(json_extract(data, "$.Host"))`},
		[]where{},
		page,
		size,
	)
	return
}

func (db *DB) GetHostNames(page int, size int) (hosts []string, total int, err error) {
	hosts, total, _, err = selectFromDB[string](
		db,
		"jobRuns",
		[]string{`DISTINCT(json_quote(json_extract(data, "$.Host.Name")))`},
		[]where{},
		page,
		size,
	)
	return
}

func (db *DB) GetJobNames(page int, size int) (jobs []string, total int, err error) {
	jobs, total, _, err = selectFromDB[string](
		db,
		"jobRuns",
		[]string{`DISTINCT(json_quote(json_extract(data, "$.Name")))`},
		[]where{},
		page,
		size,
	)
	return
}

func (db *DB) GetStatsByHost() (stats map[string]*model.Counts, err error) {
	successQ := db.stmt.Select(`json_extract(data, "$.Host.Name") as hostname`, `COUNT(*)`).
		From(`jobRuns`).
		Where(`json_extract(data, "$.Status") = "?"`, model.ProcessStartedNormally).
		Where(`json_extract(data, "$.Result.ExitCode") = 0`).
		GroupBy(`hostname`)

	failedQ := db.stmt.Select(`json_extract(data, "$.Host.Name") as hostname`, `COUNT(*)`).
		From(`jobRuns`).
		Where(`json_extract(data, "$.Status") = "?"`, model.ProcessStartedNormally).
		Where(`json_extract(data, "$.Result.ExitCode") != 0`).
		GroupBy(`hostname`)
	deniedQ := db.stmt.Select(`json_extract(data, "$.Host.Name") as hostname`, `COUNT(*)`).
		From(`jobRuns`).
		Where(`json_extract(data, "$.Status") = "?"`, model.ProcessStartedNormally).
		Where(`json_extract(data, "$.Result.ExitCode") != 0`).
		GroupBy(`hostname`)

	type result struct {
		hostname string
		count    int
	}

	runQuery := func(q sq.SelectBuilder) (results []result, err error) {
		rows, err := q.Query()
		if err != nil {
			return
		}
		defer rows.Close()

		r := result{}
		for rows.Next() {
			if err = rows.Scan(&r); err != nil {
				return
			}

			results = append(results, r)
		}
		return
	}

	results, err := runQuery(successQ)
	if err != nil {
		return
	}
	for _, r := range results {
		if _, ok := stats[r.hostname]; !ok {
			stats[r.hostname] = new(Counts)
		}

		stats[r.hostname].Successful = r.count
		stats[r.hostname].Total += r.count
	}

	results, err = runQuery(failedQ)
	if err != nil {
		return
	}
	for _, r := range results {
		if _, ok := stats[r.hostname]; !ok {
			stats[r.hostname] = new(Counts)
		}

		stats[r.hostname].Failed = r.count
		stats[r.hostname].Total += r.count
	}

	results, err = runQuery(deniedQ)
	if err != nil {
		return
	}
	for _, r := range results {
		if _, ok := stats[r.hostname]; !ok {
			stats[r.hostname] = new(Counts)
		}

		stats[r.hostname].Denied = r.count
		stats[r.hostname].Total += r.count
	}

	return
}

func (db *DB) InitTables() error {
	_, err := db.Exec(CreateJobRunsTableSQL)
	return err
}

func (db *DB) observeAddChannel() {
	for jobRun := range db.addChannel {
		db.errChannel <- db.insert(jobRun)
	}
}

func NewDB(uri string) (*DB, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}

	d := &DB{
		DB: db,

		addChannel: make(chan model.JobRun),
		errChannel: make(chan error),
	}
	dbCache := sq.NewStmtCache(d.DB)
	d.stmt = sq.StatementBuilder.RunWith(dbCache)

	if err = d.InitTables(); err != nil {
		return nil, err
	}

	go d.observeAddChannel()

	return d, nil

}
