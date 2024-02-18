package sqlite

var (
	DropJobsTableSQL   = `DROP TABLE jobs;`
	CreateJobsTableSQL = `
		CREATE TABLE IF NOT EXISTS jobs (
			name TEXT PRIMARY KEY,
			schedule TEXT
		);`
	DropJobRunsTableSQL   = `DROP TABLE jobRuns;`
	CreateJobRunsTableSQL = `
		CREATE TABLE IF NOT EXISTS jobRuns  (
			id BLOB(16) PRIMARY KEY,
			jobName TEXT,
			data JSON,
			FOREIGN KEY(jobName) REFERENCES jobs(name)
		);
	`
)
