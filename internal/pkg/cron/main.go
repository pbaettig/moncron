package cron

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

var (
	CronFileLocations map[string]bool
)

func init() {
	user, _ := user.Current()

	// Path Glob, contains username
	CronFileLocations = map[string]bool{
		fmt.Sprintf("/tmp/var/spool/cron/crontabs/%s", user.Username): false,
		"/tmp/etc/cron.d/*":      true,
		"/tmp/etc/cron.weekly/*": true,
		"/tmp/etc/cron.daily/*":  true,
		"/tmp/etc/cron.hourly/*": true,
	}

	fmt.Printf("%#v\n", splitWhitespace("10 */4 * * *      pascal	sleep 10 && true sleep "))
}

func splitWhitespace(line string) []string {
	fields := make([]string, 0)
	fieldStart := false
	fieldValue := new(strings.Builder)
	for i, c := range line {
		lastChar := i == len(line)-1
		if lastChar && !unicode.IsSpace(c) {
			fmt.Printf("encounterd last character at index %d (%s): \"%+v\"\n", i, strconv.QuoteRune(c), fieldValue.String())
			fieldValue.WriteRune(c)
			fields = append(fields, fieldValue.String())
			break
		}
		// current character is whitespace...
		if unicode.IsSpace(c) {

			// ...but a field has previously been started
			// so we finalize it and add it to `fields`
			if fieldStart {

				fieldStart = false

				fields = append(fields, fieldValue.String())
				fieldValue.Reset()
				continue
			} else {
				continue
			}

		}

		// current character is not whitespace and no field has been started yet
		// so we start one
		if !fieldStart && !unicode.IsSpace(c) {
			fieldStart = true
			fieldValue.WriteRune(c)
			continue
		}

		fieldValue.WriteRune(c)
	}
	return fields
}

type cronJob struct {
	Location    string
	schedule    string
	hasUsername string
	username    string
	command     string
	environment map[string]string
}

func FromFile(path string, hasUsername bool) []cronJob {
	fd, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	jobs := make([]cronJob, 0)

	scanner := bufio.NewScanner(fd)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lineSplit := splitWhitespace(scanner.Text())
		if len(lineSplit) > 1 {

			fmt.Printf("%#v\n", lineSplit)
		}
	}

	return jobs

}

func FindAll() []cronJob {
	// jobs := make([]cronJob, 0)
	for globPath, hasUsername := range CronFileLocations {
		files, err := filepath.Glob(globPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			_ = FromFile(file, hasUsername)
		}

		// jobs = append(jobs, cronJob{
		// 	Location: path,
		// })
	}
	return []cronJob{}
}

func FindJob(args []string) {

}

// "10 */4 * * *      pascal	sleep 10 && true\n"
func parseCronJobLine(line string) {
	fields := make([]string, 0)
	fieldCounter := 0
	fieldStart := false
	fieldValue := new(strings.Builder)
	for i, c := range line {
		fmt.Println(i, string(c))

		// Field starts if:
		// - current character is not whitespace
		// - no field has been started
		if !fieldStart && !unicode.IsSpace(c) {
			fmt.Printf("field #%d started at index %d (%s)\n", fieldCounter, i, strconv.QuoteRune(c))
			fieldStart = true
		}

		// Field ends if:
		// - a field has been started
		// - the current character is whitespace OR
		// -
		if fieldStart && (unicode.IsSpace(c) || i == len(line)-1) {
			fmt.Printf("field #%d ended at index %d (%s): \"%+v\"\n", fieldCounter, i, strconv.QuoteRune(c), fieldValue.String())
			fieldStart = false

			if !unicode.IsSpace(c) {
				fieldValue.WriteRune(c)
			}
			fields = append(fields, fieldValue.String())
			fieldValue.Reset()
			fieldCounter++
		}

		// if fieldStart {
		// 	fieldValue.WriteRune(c)
		// }

	}
	for i, f := range fields {
		fmt.Printf("Field #%d: \"%s\"\n", i, f)
	}
}

// "10 */4 * * *      pascal	sleep 10 && true sleep 3"
