package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	ErrEnvVariableNotDefined = fmt.Errorf("environment variable is not defined")
)

func getEnvBool(name string) (bool, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return false, fmt.Errorf("env variable %s is not defined", name)
	}

	return strconv.ParseBool(v)
}

func getEnvString(name string) (string, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("env variable %s is not defined", name)
	}

	return v, nil
}

func getEnvInt64(name string) (int64, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return 0, fmt.Errorf("env variable %s is not defined", name)
	}

	return strconv.ParseInt(v, 10, 64)
}

func getEnvDuration(name string) (time.Duration, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return 0, fmt.Errorf("env variable %s is not defined", name)
	}

	return time.ParseDuration(v)
}

type cmdlineArgs struct {
	Timeout        time.Duration
	JobName        string
	PushgatewayURL string
	MonServerURL   string
	Stdout         bool
	Quiet          bool
	Version        bool
	LogFile        string
	LogFileMaxSize int64
	Once           bool
	ProcessCmdline []string
}

func (c *cmdlineArgs) FromEnv() error {
	if timeout, err := getEnvDuration("MON_TIMEOUT"); err == nil {
		c.Timeout = timeout
	}
	if name, err := getEnvString("MON_JOB"); err == nil {
		c.JobName = name
	}
	if pushgw, err := getEnvString("MON_PUSHGW_URL"); err == nil {
		c.PushgatewayURL = pushgw
	}
	if server, err := getEnvString("MON_SERVER_URL"); err == nil {
		c.MonServerURL = server
	}
	if stdout, err := getEnvBool("MON_STDOUT"); err == nil {
		c.Stdout = stdout
	}
	if quiet, err := getEnvBool("MON_QUIET"); err == nil {
		c.Quiet = quiet
	}
	if logfile, err := getEnvString("MON_LOGFILE"); err == nil {
		c.LogFile = logfile
	}
	if once, err := getEnvBool("MON_ONCE"); err == nil {
		c.Once = once
	}
	if logfileSize, err := getEnvInt64("MON_LOGFILE_SIZE"); err == nil {
		c.LogFileMaxSize = logfileSize
	}

	return nil
}

func (c *cmdlineArgs) FromCmdline() error {
	timeout := flag.Duration("timeout", 0, "timeout value")
	name := flag.String("name", "", "job name")
	pushgw := flag.String("pushgw", "", "Prometheus pushgateway URL")
	server := flag.String("server", "", "mon Server URL")
	stdout := flag.Bool("stdout", false, "echo results to stdout")
	quiet := flag.Bool("quiet", false, "")
	version := flag.Bool("version", false, "")
	logfile := flag.String("log", "", "Log File path")
	logSize := flag.Int64("log-size", 10*1024*1024, "Log File Maximum Size")
	once := flag.Bool("once", false, "ensure the specified command is not already running before starting")

	flag.Parse()

	if timeout != nil && *timeout > 0 {
		c.Timeout = *timeout
	}
	if name != nil && *name != "" {
		c.JobName = *name
	}
	if pushgw != nil && *pushgw != "" {
		c.PushgatewayURL = *pushgw
	}
	if server != nil && *server != "" {
		c.MonServerURL = *server
	}
	if stdout != nil && *stdout {
		c.Stdout = *stdout
	}
	if quiet != nil && *quiet {
		c.Quiet = *quiet
	}
	if version != nil && *version {
		c.Version = *version
	}
	if logfile != nil && *logfile != "" {
		c.LogFile = *logfile
	}
	if logSize != nil && *logSize > 0 {
		c.LogFileMaxSize = *logSize
	}
	if once != nil && *once {
		c.Once = *once
	}

	c.ProcessCmdline = flag.Args()
	return nil
}

func (c *cmdlineArgs) Validate() error {
	if c.Version {
		return nil
	}

	if c.JobName == "" {
		return fmt.Errorf("-name is required")
	}

	if len(c.ProcessCmdline) == 0 {
		return fmt.Errorf("nothing to execute")
	}

	return nil
}
