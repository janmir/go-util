package util

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	_callDepth = 2
	_callSep   = "»"
)

var (
	_debug = true

	frx        *regexp.Regexp
	logErr     *log.Logger
	fmtSpecReg = `%[0-9]?\.?[+# 0-9]?[sdfpbcqxXvVtTU]`

	rfg = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	gfg = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	cfg = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
)

func init() {
	var err error
	frx, err = regexp.Compile(fmtSpecReg)
	Catch(err)

	logErr = log.New(os.Stderr, gfg("✱ "), log.Ltime|log.Lmicroseconds) //|log.Lshortfile
}

//Catch try..catch errors
func Catch(err error) error {
	if err != nil {
		//get error message
		msg := err.Error()
		caller := ""

		//get error location
		for i := _callDepth; i >= 1; i-- {
			_, fn, line, _ := runtime.Caller(i)

			if line > 0 {
				//format output message
				sp := strings.Split(fn, "/")
				fn = strings.Replace(sp[len(sp)-1], ".go", "*", -1)
				caller += fmt.Sprintf("%s@%d%s", fn, line, _callSep)
			}
		}
		caller = strings.TrimSuffix(caller, _callSep)

		msg = strings.Replace(msg, "rpc", "*", -1)
		errorMessage := fmt.Sprintf("%s %s", cfg(caller), msg)

		Logger(rfg("Error"), errorMessage)
		return errors.New(errorMessage)
	}
	return nil
}

//HTTPCatch try..catch errors
func HTTPCatch(res http.Response, err error) error {
	if err == nil && res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code was %d", res.StatusCode)
	}

	return Catch(err)
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
