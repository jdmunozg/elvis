package logs

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/cgalvisleon/elvis/strs"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

func log(kind string, color string, args ...any) string {
	kind = strings.ToUpper(kind)
	message := fmt.Sprint(args...)
	now := time.Now().Format("2006/01/02 15:04:05")
	var result string

	switch color {
	case "Reset":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + message + Reset
	case "Red":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Red + message + Reset
	case "Green":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Green + message + Reset
	case "Yellow":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Yellow + message + Reset
	case "Blue":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Blue + message + Reset
	case "Purple":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Purple + message + Reset
	case "Cyan":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Cyan + message + Reset
	case "Gray":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Gray + message + Reset
	case "White":
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + White + message + Reset
	default:
		result = now + Purple + strs.Format(" [%s]: ", kind) + Reset + Green + message + Reset
	}

	println(result)

	return result
}

func Log(kind string, args ...any) {
	log(kind, "", args...)
}

func Logf(kind string, format string, args ...any) {
	message := strs.Format(format, args...)
	log(kind, "", message)
}

func Traces(kind, color string, err error) ([]string, error) {
	var n int = 1
	var traces []string = []string{err.Error()}

	log(kind, color, err.Error())

	for {
		pc, file, line, more := runtime.Caller(n)
		if !more {
			break
		}
		n++
		function := runtime.FuncForPC(pc)
		name := function.Name()
		list := strings.Split(name, ".")
		if len(list) > 0 {
			name = list[len(list)-1]
		}
		if !slices.Contains([]string{"ErrorM", "ErrorF"}, name) {
			trace := strs.Format("%s:%d func:%s", file, line, name)
			traces = append(traces, trace)
			log("TRACE", color, trace)
		}
	}

	return traces, err
}

func Error(err error) error {
	_, err = Traces("Error", "red", err)

	return err
}

func Errorm(message string) error {
	err := errors.New(message)
	return Error(err)
}

func Errorf(format string, args ...any) error {
	message := strs.Format(format, args...)
	err := errors.New(message)
	return Error(err)
}

func Info(v ...any) {
	log("Info", "Blue", v...)
}

func Infof(format string, args ...any) {
	message := strs.Format(format, args...)
	log("Info", "Blue", message)
}

func Fatal(v ...any) {
	log("Fatal", "Red", v...)
	os.Exit(1)
}

func Panic(v ...any) {
	log("Panic", "Red", v...)
	os.Exit(1)
}

func Ping() {
	log("PING", "")
}

func Pong() {
	log("PONG", "")
}

func Debug(v ...any) {
	log("Debug", "Cyan", v...)
}

func Debugf(format string, args ...any) {
	message := strs.Format(format, args...)
	log("Debug", "Cyan", message)
}
