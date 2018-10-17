package util

import (
	"log"
	"testing"
	"time"
)

func TestCatch(t *testing.T) {
	//err := Catch(errors.New("this is a new error"))
	//if err == nil {
	//	t.Fail()
	//}
}

func TestCurrentDir(t *testing.T) {
	dir, err := GetCurrDir()
	if err != nil {
		t.Fail()
	}

	if len(dir) <= 0 {
		t.Fail()
	}

	log.Println("Current Directory:", dir)
}

func TestLogger(t *testing.T) {
	Logger("This is a normal Log", "Hello there")
	Logger("✔ %d:%s ", 1, "This is a formatted Log")
}

func TestTimeTracker(t *testing.T) {
	TimeTrack(time.Now(), "Timecheck:")
	time.Sleep(time.Second * 5)
}
