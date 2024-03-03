# moncron

## mon
mon is a wrapper for running scheduled jobs of any kind. It will execute whatever commandline you provide and gather metrics about the execution, which are then forwarded to a prometheus pushgateway, POSTed to a webhook URL or recorded in a local file.
### Parameters
- `-name`: a free-form name that is used to identify the run in logs / metrics
- `-quiet`: disables all extra output and only displays the executed commands stdout and stderr
- `-timeout`: max. duration to wait for the specified command. If execution exceeds the provided timeout the process will be forcefully killed


#### Targets
- `-pushgw`: URL to Prometheus Pushgateway instance. If provided metrics about the job execution are sent to it.
- `-web`: URL to a webserver accepting POST requests containing job execution details in JSON format. You can use you're own webhook or use the `server` provided by moncron. In the latter case use the path `/api/runs`.
- `-log`: path to a local file where execution details are appended to
- `-log-size`: max. file size in bytes. The program will truncate the file if it grows beyond the size specified here
- `-stdout`: print execution details to stdout

### Examples
```
./mon -name sleep -stdout -- sleep 2 | jq
INFO[0000] started                                       name=sleep
INFO[0002] command finished successfully                 exit=0 name=sleep
INFO[0002] successfully pushed results                   name=sleep target=stdout
{
  "Name": "sleep",
  "Executable": "sleep",
  "Args": [
    "2"
  ],
  "Result": {
    "User": {
      "Uid": "501",
      "Gid": "20",
      "Username": "pascal",
      "Name": "Pascal Bättig",
      "HomeDir": "/Users/pascal"
    },
    "Environment": {
      "LANG": "en_US.UTF-8",
      "USER": "pascal",
      "PWD": "/Users/pascal/git/moncron",
      "SHELL": "/bin/bash",
      ...
    },
    "WorkingDirectory": "/Users/pascal/git/moncron",
    "StartedAt": "2024-02-11T17:38:37.18226Z",
    "FinishedAt": "2024-02-11T17:38:39.190605Z",
    "ExitCode": 0,
    "Killed": false,
    "MaxRssBytes": 933888,
    "Stdout": "",
    "Stderr": "",
    "WallTime": "2.008345s",
    "UserTime": "612µs",
    "SystemTime": "1.197ms",
    "ReceivedSignal": ""
  }
}
```

## server
moncron comes with a server that stores job run results in a SQLite DB and exposes a HTML interface as well as a JSON API.

### Parameters
- `-listen`: listen address for the HTTP server (default is 0.0.0.)
- `-port`: listen port for the HTTP server (default is 8088)
- `-db`: path to the SQLite DB file used to store Job runs (default is test.db)
- `-timeout`: graceful shutdown timeout for the HTTP server (default is 15 s)