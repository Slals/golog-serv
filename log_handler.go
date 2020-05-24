package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

// LogLevel defines a log level
type LogLevel string

const (
	// TRACE is trace log level
	// Used to keep track of normal processes
	TRACE LogLevel = "trace"

	// DEBUG is debug log level
	// Same as trace but only used for development environment
	DEBUG LogLevel = "debug"

	// INFO is info log level
	// Used to keep track of scheduled operations
	INFO LogLevel = "info"

	// NOTICE is notice log level
	// Used to track noticable event from production environment
	NOTICE LogLevel = "notice"

	// WARN is warn log level
	// Used to track events that could lead to an error
	WARN LogLevel = "warn"

	// ERROR is error log level
	// Used to track errors which doesn't kill the client process from develpment
	// environment and / or production environment
	ERROR LogLevel = "error"

	// FATAL is fatal log level
	// Used to track fatal errors which kill the client process from development
	// environment and / or production environment
	FATAL LogLevel = "fatal"
)

// UnmarshalJSON unmashals enum value
func (s *LogLevel) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	switch j {
	case "trace":
		*s = TRACE
		return nil
	case "debug":
		*s = DEBUG
		return nil
	case "info":
		*s = INFO
		return nil
	case "notice":
		*s = NOTICE
		return nil
	case "warn":
		*s = WARN
		return nil
	case "error":
		*s = ERROR
		return nil
	case "fatal":
		*s = FATAL
		return nil
	}

	return errors.New("Invalid log level")
}

func stringToLogLevel(s string) LogLevel {
	switch s {
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "notice":
		return NOTICE
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return "unknown"
	}
}

// Log defines a log item
type Log struct {
	Level   LogLevel `json:"level"`
	Key     string   `json:"key"`
	Message string   `json:"message"`

	Timestamp string
	UserAgent string
}

func logHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)

		var data Log
		if err := decoder.Decode(&data); err != nil {
			http.Error(rw, "Failed to parse the body : `level`, `key` and `message` are required", http.StatusBadRequest)
			return
		}

		f, err := os.OpenFile(getPath()+"HEAD.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			http.Error(rw, "Couldn't start the log writer", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// timestamp;user-agent;level;key;message
		fmt.Fprintf(f, "%s;%s;%s;%s;%s\n", time.Now().Format(time.RFC1123), req.UserAgent(), data.Level, data.Key, data.Message)

		rw.WriteHeader(http.StatusNoContent)

	case http.MethodGet:
		f, err := os.OpenFile(getPath()+"HEAD.log", os.O_RDONLY, 0644)
		if err != nil {
			http.Error(rw, "Couldn't open log data", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var data []byte
		if _, err := f.Read(data); err != nil {
			http.Error(rw, "Couldn't read log data", http.StatusInternalServerError)
			return
		}

		logs := make([]Log, 0)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			rawData := strings.Split(line, ";")
			if len(rawData) == 5 {
				logs = append(logs, Log{
					Timestamp: rawData[0],
					UserAgent: rawData[1],
					Level:     stringToLogLevel(rawData[2]),
					Key:       rawData[3],
					Message:   rawData[4],
				})
			}
		}

		t := template.New("index")
		t = template.Must(t.ParseFiles("tmpl/index.html"))

		sort.Slice(logs, func(i, j int) bool {
			dateA, _ := time.Parse(time.RFC1123, logs[i].Timestamp)
			dateB, _ := time.Parse(time.RFC1123, logs[j].Timestamp)
			return dateA.Unix() > dateB.Unix()
		})

		p := struct {
			Title string
			Logs  []Log
		}{
			Title: os.Getenv("PAGE_TITLE"),
			Logs:  logs,
		}

		if err := t.ExecuteTemplate(rw, "index", p); err != nil {
			log.Print(err)
			http.Error(rw, "Couldn't display the index page", http.StatusInternalServerError)
			return
		}
	}
}

func getPath() string {
	path := os.Getenv("DEBUG_PATH")
	if path == "" {
		path = "./"
	}
	return path
}
