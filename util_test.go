package util

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestCatch(t *testing.T) {
	err := catch(errors.New("this is a new error"))
	if err == nil {
		t.Fail()
	}
}

func TestCurrentDir(t *testing.T) {
	dir, err := getCurrDir()
	if err != nil {
		t.Fail()
	}

	if len(dir) <= 0 {
		t.Fail()
	}

	log.Println("Current Directory:", dir)
}

func TestLogger(t *testing.T) {
	logger("This is a normal Log", "Hello there")
	logger("✔ %d:%s ", 1, "This is a formatted Log")
}

func TestTimeTracker(t *testing.T) {
	timeTrack(time.Now(), "Timecheck:")
	time.Sleep(time.Second * 5)
}
