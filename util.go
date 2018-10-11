package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	_debug = true
	frx    *regexp.Regexp

	logErr *log.Logger
	rfg    = color.New(color.FgHiRed, color.Bold).SprintfFunc()
)

func init() {
	var err error
	frx, err = regexp.Compile(`%[+# 0]?[sqxXvVT]`)
	Catch(err)

	logErr = log.New(os.Stderr, "✱ ", log.Ltime|log.Lmicroseconds) //|log.Lshortfile
}

//Catch try..catch errors
func Catch(err error) error {
	if err != nil {
		//get error message
		msg := err.Error()

		//get error location
		_, fn, line, _ := runtime.Caller(1)

		//format output message
		sp := strings.Split(fn, "/")
		fn = strings.Replace(sp[len(sp)-1], ".go", "***", -1)
		msg = strings.Replace(msg, "rpc", "***", -1)

		errorMessage := fmt.Sprintf("%s@%d %s", fn, line, msg)
		Logger(rfg("Error"), errorMessage)

		return errors.New(errorMessage)
	}
	return nil
}

//Logger logs to standard error
func Logger(strs ...interface{}) {
	if _debug {
		//check if contain format specifier
		if len(strs) > 1 && frx.MatchString(strs[0].(string)) {
			logErr.Printf(strs[0].(string), strs[1:]...)
		} else {
			logErr.Println(strs...)
		}
	}
}

//TimeTrack dump execution time
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Logger("%s took %s", name, elapsed)
}

//GetCurrDir current directory of executable
func GetCurrDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
