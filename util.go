package util

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/fatih/color"
	mailgun "github.com/janmir/go-mailgun"
	"github.com/mitchellh/go-ps"
)

const (
	_callDepth = 3
	_callSep   = " » "
	_logFile   = "_log.txt"
)

//IPData details about the caller
//Ip address
type IPData struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Loc     string `json:"loc"`
	Postal  string `json:"postal"`
	Org     string `json:"org"`
}

var (
	_debug          = true
	_fileLogging    = false
	_consoleLogging = true

	//internet connection check
	_internetCheckURL = "http://clients3.google.com/generate_204"

	//Network Time Protocol Servers
	_appleNTP  = "time.apple.com"
	_googleNTP = "time.google.com"

	//AppleNTP public ntp variable
	AppleNTP = _appleNTP
	//GoogleNTP public ntp variable
	GoogleNTP = _googleNTP

	//Logs to standard error, can be disabled
	//by calling DisableConsoleLogging()
	logErr *log.Logger
	//Logs to standard error
	//Disabling this logger is not
	//possible, you usually use this one for
	//important INFO logs
	logInfo *log.Logger
	//Logs to a file with filename _logFile
	//can be enabled by calling EnableFileLogging()
	logFile *log.Logger
	//Null logger
	logNull *log.Logger

	//Mail logger
	mailer mailgun.Mail

	frx        *regexp.Regexp
	fmtSpecReg = `%[0-9]?\.?[+# 0-9]?[sdfpbcqxXvVtTU]`

	//Color Sprintf functions, mostly used by
	//the logInfo logger
	rfg = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	bfg = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
	gfg = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	cfg = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	mfg = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()
	yfg = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	wfg = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
)

func init() {
	var err error
	frx, err = regexp.Compile(fmtSpecReg)
	Catch(err)

	logErr = log.New(os.Stderr, gfg("✱ "), log.Ltime|log.Lmicroseconds) //|log.Lshortfile
	logInfo = log.New(os.Stderr, "", 0)                                 //|log.Lshortfile
	logNull = log.New(ioutil.Discard, "", 0)
}

/*********************************/
/*
/*		 Error Handlers
/*
/*********************************/

//Catch try..catch errors
func Catch(err error, more ...string) {
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
				fn = strings.Replace(sp[len(sp)-1], ".go", "", -1)
				caller += fmt.Sprintf("%s:%d%s", fn, line, _callSep)
			}
		}
		caller = strings.TrimSuffix(caller, _callSep)

		//Further replace
		// msg = strings.Replace(msg, "rpc", "***", -1)

		if len(more) > 0 {
			errorMessage = fmt.Sprintf("%s %s %s", cfg(caller), wfg(msg), yfg("🛈 "+strings.Join(more, ", ")))
		} else {
			errorMessage = fmt.Sprintf("%s %s", cfg(caller), wfg(msg))
		}

		//Log to standard error
		Logger(rfg("Error"), errorMessage)

		//Exit
		//os.exit does not call defered functions
		//and so does fatal from log/fmt package
		//to be able to call deferred functions,
		//track by using recover routine
		os.Exit(1)
	}
}

//HTTPCatch try..catch errors
func HTTPCatch(res *http.Response, err error, more ...string) {
	if err == nil && res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code was %d", res.StatusCode)
	}

	Catch(err, more...)
}

//Recover ...
func Recover() {
	if err := recover(); err != nil {
		Logger(yfg("Recovery"), err)
	}
}

/*********************************/
/*
/*			 Logging Fn
/*
/*********************************/

func fmtr(strs ...interface{}) string {
	if len(strs) > 1 && frx.MatchString(strs[0].(string)) {
		return fmt.Sprintf(strs[0].(string), strs[1:]...)
	}
	return fmt.Sprintln(strs...)
}

//DisableLogging ...
func DisableLogging() {
	_debug = false
}

//DisableConsoleLogging ...
func DisableConsoleLogging() {
	_consoleLogging = false
}

//EnableFileLogging ...
func EnableFileLogging() {
	_fileLogging = true
	if _fileLogging {
		path, err := GetCurrDir()
		Catch(err)

		file := filepath.Join(path, _logFile)

		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		Catch(err)

		logFile = log.New(f, "✱ ", log.Ltime|log.Lmicroseconds) //|log.Lshortfile
	}
}

//Logger logs to standard error
func Logger(strs ...interface{}) {
	if _debug {
		//Log to file
		if _fileLogging {
			logFile.Print(fmtr(strs...))
		}

		if _consoleLogging {
			logErr.Print(fmtr(strs...))
		}
	}
}

//Red prints text in red
func Red(strs ...interface{}) {
	str := strings.TrimSpace(fmtr(strs...))
	logInfo.Print(rfg(str))
}

//Green prints text in green
func Green(strs ...interface{}) {
	str := strings.TrimSpace(fmtr(strs...))
	logInfo.Print(gfg(str))
}

//Cyan prints text in cyan
func Cyan(strs ...interface{}) {
	str := strings.TrimSpace(fmtr(strs...))
	logInfo.Print(cfg(str))
}

//Magenta prints text in magenta
func Magenta(strs ...interface{}) {
	str := strings.TrimSpace(fmtr(strs...))
	logInfo.Print(mfg(str))
}

//Yellow prints text in yellow
func Yellow(strs ...interface{}) {
	str := strings.TrimSpace(fmtr(strs...))
	logInfo.Print(yfg(str))
}

//TimeTrack dump execution time
//TimeTrack dump execution time
func TimeTrack(start time.Time, name string, cb ...func(string)) {
	elapsed := time.Since(start)
	if len(cb) > 0 {
		for _, fn := range cb {
			fn(fmt.Sprintf("%s", elapsed))
		}
	} else {
		Logger("%s %s took %s", mfg("Timestamp"), rfg(name), elapsed)
	}
}

//CreateMailer creates a mailgun client
func CreateMailer(domain, api string) {
	mailer = mailgun.DefaultMailClient(domain, api)
}

//SendMail sends the mailgun emai
func SendMail(to, from, subject, msg string) error {
	if !mailer.Initialized {
		Red("Initialize mail client first using [CreateMail] function.")
	}

	mailer.Create(to, from, subject, msg)
	out, err := mailer.Send()
	if err != nil {
		return err
	}

	Logger("Mail Sent: %+v", out)
	return nil
}

/*********************************/
/*
/*			 Helper Fn
/*
/*********************************/

//IsInterfaceAPointer checks if an interface is of type pointer
func IsInterfaceAPointer(val interface{}) {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		Catch(errors.New("non-pointer passed to Unmarshal"))
	}
}

//GetNTPTime return ntp time
//Available servers are AppleNTP, GoogleNTP
func GetNTPTime(server ...string) (time.Time, error) {
	serve := _appleNTP
	if len(server) == 1 {
		serve = server[0]
	}

	return ntp.Time(serve)
}

//GetCurrDir current directory of executable
func GetCurrDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

//Localize returns the absolute path of the file
//relative to the executable
func Localize(file string) string {
	dir, err := GetCurrDir()
	Catch(err)

	return filepath.Join(dir, file)
}

//GetFiles returns all files with the extension
//if provided if not all files
func GetFiles(path string, subdir bool, ext ...string) []string {
	files := make([]string, 0)

	//get absolute path
	root, err := filepath.Abs(path)
	Catch(err)

	if subdir {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if len(ext) > 0 {
				x := filepath.Ext(path)
				for _, v := range ext {
					if strings.Contains(x, v) {
						files = append(files, path)
					}
				}
			} else {
				files = append(files, path)
			}
			return nil
		})
		Catch(err)
	} else {
		fs, err := ioutil.ReadDir(root)
		Catch(err)

		for _, file := range fs {
			path := file.Name()
			full := filepath.Join(root, path)
			if len(ext) > 0 {
				x := filepath.Ext(path)
				for _, v := range ext {
					if strings.Contains(x, v) {
						files = append(files, full)
					}
				}
			} else {
				files = append(files, full)
			}
		}
	}

	return files
}

//IsConnectedToInternet checks if you are currently connected
//to the internet, timeout is set to 1.5 seconds
func IsConnectedToInternet(path ...string) bool {
	url := _internetCheckURL
	if len(path) == 1 {
		url = path[0]
	}

	client := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	_, err := client.Get(url)

	if err != nil {
		return false
	}

	return true
}

//GetPublicIPDetails returns the callers
//Ip address details includes location, etc...
func GetPublicIPDetails(token string) IPData {
	ipd := IPData{}
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	//Source list
	ipstack := "http://api.ipstack.com/check?access_key=" + token
	ipinfo := "https://ipinfo.io?token=" + token
	_, _ = ipstack, ipinfo

	req, err := http.NewRequest("GET", ipinfo, nil)
	req.Header.Set("User-Agent", "curl")
	resp, err := client.Do(req)
	HTTPCatch(resp, err, "Unable to fetch data from url")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Catch(err)
	}

	if err := json.Unmarshal(b, &ipd); err != nil {
		Catch(err)
	}

	return ipd
}

//Debounce calls a function only once after an interval
func Debounce(interval time.Duration, input chan interface{}, fn func(arg interface{})) {
	var item interface{}
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			fn(item)
		}
	}
}

//AmIRunning checks if process with name proc running instances
func AmIRunning(proc string) int {
	count := 0
	pss, err := ps.Processes()
	Catch(err)

	filex, err := regexp.Compile("^" + proc + "(\\.exe)?$")
	Catch(err, "Uncompilable Regular Expression")

	for _, v := range pss {
		name := v.Executable()
		if filex.MatchString(name) {
			count++
		}
	}

	return count
}

//Max a > b, c ....
func Max(a int, b ...int) int {
	max := a
	for _, val := range b {
		if max < val {
			max = val
		}
	}
	return max
}

//Min a < b, c ...
func Min(a int, b ...int) int {
	min := a
	for _, val := range b {
		if min > val {
			min = val
		}
	}
	return min
}

//Rand return a random number from a to b
func Rand(a, b int) int {
	if a < b {
		seed := rand.NewSource(time.Now().UnixNano())
		randomizer := rand.New(seed)
		no := randomizer.Intn(b - a)
		return a + no
	}

	return -1
}

//Encrypt encrypts a text using a key
//plain text is variable length, key should either be
// 16, 24, or 32 bytes/length to select from
// AES-128(for 16), AES-192(for 24), or AES-256(for 32).
//returns the encrypted text
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

//Decrypt decrypts the encrypted text using the
//function above
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
