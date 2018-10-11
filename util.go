package util

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func catch(err error) {
	if err != nil {
		//get error message
		msg := err.Error()

		//get error location
		_, fn, line, _ := runtime.Caller(1)

		//format output message
		sp := strings.Split(fn, "/")
		fn = strings.Replace(sp[len(sp)-1], ".go", "***", -1)
		msg = strings.Replace(msg, "rpc", "***", -1)

		log.Fatal("✖: ", fn, ":", line, " ", msg)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func getCurrDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
