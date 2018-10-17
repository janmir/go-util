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
	_callDepth = 3
	_callSep   = " » "
)

var (
	_debug = true

	logErr  *log.Logger
	logInfo *log.Logger

	frx        *regexp.Regexp
	fmtSpecReg = `%[0-9]?\.?[+# 0-9]?[sdfpbcqxXvVtTU]`

	rfg = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	gfg = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	cfg = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	mfg = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()
	yfg = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
)

func init() {
	var err error
	frx, err = regexp.Compile(fmtSpecReg)
	Catch(err)

	logErr = log.New(os.Stderr, gfg("✱ "), log.Ltime|log.Lmicroseconds) //|log.Lshortfile
	logInfo = log.New(os.Stderr, "", 0)                                 //|log.Lshortfile
}

//Catch try..catch errors
func Catch(err error, more ...string) error {
	if err != nil {
		//get error message
		errorMessage := ""
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
		if len(more) > 0 {
			errorMessage = fmt.Sprintf("%s %s %s", cfg(caller), msg, yfg("🛈 "+strings.Join(more, ", ")))
		} else {
			errorMessage = fmt.Sprintf("%s %s", cfg(caller), msg)
		}

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

//Recover ...
func Recover() {
	if err := recover(); err != nil {
		Logger(yfg("Recovery"), err)
	}
}

func fmtr(strs ...interface{}) string {
	if len(strs) > 1 && frx.MatchString(strs[0].(string)) {
		return fmt.Sprintf(strs[0].(string), strs[1:]...)
	}
	return fmt.Sprintln(strs...)
}

//Logger logs to standard error
func Logger(strs ...interface{}) {
	if _debug {
		//check if contain format specifier
		logErr.Println(fmtr(strs...))
	}
}

//Red prints text in red
func Red(strs ...interface{}) {
	logInfo.Println(rfg(fmtr(strs...)))
}

//Green prints text in green
func Green(strs ...interface{}) {
	logInfo.Println(gfg(fmtr(strs...)))
}

//Cyan prints text in cyan
func Cyan(strs ...interface{}) {
	logInfo.Println(cfg(fmtr(strs...)))
}

//Magenta prints text in magenta
func Magenta(strs ...interface{}) {
	logInfo.Println(mfg(fmtr(strs...)))
}

//Yellow prints text in yellow
func Yellow(strs ...interface{}) {
	logInfo.Println(yfg(fmtr(strs...)))
}

//TimeTrack dump execution time
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Logger("%s %s took %s", mfg("Timestamp"), name, elapsed)
}

//GetCurrDir current directory of executable
func GetCurrDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

//Max a > b
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//Min a < b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
