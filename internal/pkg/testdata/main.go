package testdata

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pbaettig/moncron/internal/pkg/model"
)

var (
	// wg    *sync.WaitGroup = new(sync.WaitGroup)
	hostNameFmt string = `srv-%03d.acme.corp`
	HostNames   []string
	Seed                 = 12345678
	JobNames    []string = []string{
		"sign-tps-reports",
		"send-daily-stats",
		"etl-production-data",
		"etl-staging-data",
		"etl-test-data",
		"update-core-service",
		"build-core-service",
		"deploy",
		"cleanup-images",
		"generate-thumbnails",
		"backup-prod-db",
		"scramble-hr-data",
		"load-staff-info",
		"execute-accounts",
		"fix_something",
		"apply-remote-config",
	}
)

func init() {
	gofakeit.Seed(Seed)

	// for i := 0; i < 10; i++ {
	// 	HostNames = append(HostNames, fmt.Sprintf(HostNameFmt, i+1))
	// }

	// for i := 0; i < 100; i++ {
	// 	gofakeit.Seed(i)
	// 	jobs = append(hosts, fmt.Sprintf("%s-%s", gofakeit.HackerVerb(), gofakeit.HackerNoun()))
	// }
}

func GenerateHostNames(n int, nameFmt string) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf(nameFmt, i+1)
	}
	return names
}

func GenerateJobRuns(n int, nHosts int) []model.JobRun {
	hosts := GenerateHostNames(nHosts, hostNameFmt)
	runs := make([]model.JobRun, n)
	for i := 0; i < n; i++ {
		jobRun := RandomJobRun()
		jobRun.Host.Name = gofakeit.RandomString(hosts)
		jobRun.Job.Name = gofakeit.RandomString(JobNames)
		runs[i] = jobRun
	}
	return runs
}

// func main() {
// 	log.SetLevel(log.DebugLevel)
// 	// gofakeit.Seed(12345678)
// 	jobRunStore, err := sqlite.NewDB("test.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := 0; i < 200; i++ {
// 		gofakeit.Seed(i)
// 		host := gofakeit.RandomString(hosts)
// 		job := gofakeit.RandomString(jobs)
// 		// jobRun := model.RandomJobRun()
// 		// if err := jobRunStore.Add(jobRun); err != nil {
// 		// 	fmt.Printf("%#v, %+v, %T\n", err, err, err)
// 		// 	log.Error(err)
// 		// }
// 		// log.Infof("inserted jobRun #%d: %s,%s (%s)", i+1, jobRun.ID, jobRun.Name, jobRun.Host.Name)
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			jobRun := model.RandomJobRun()
// 			jobRun.Host.Name = host
// 			jobRun.Job.Name = job
// 			if err := jobRunStore.Add(jobRun); err != nil {
// 				fmt.Printf("%#v, %+v, %T\n", err, err, err)
// 				log.Error(err)
// 			}
// 			log.Infof("inserted jobRun #%d: %s,%s (%s)", i+1, jobRun.ID, jobRun.Name, jobRun.Host.Name)
// 		}(i)
// 	}
// 	wg.Wait()
// }
