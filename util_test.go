package util

import (
	"log"
	"testing"
	"time"
)

var (
	token = "*********"
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

func TestMinMax(t *testing.T) {
	min := Min(0, 1, 2, 3, 4, 5, 6, 7, 7, 8, -1, 4, 3, 2, 5)
	max := Min(0, 1, 2, 3, 4, 5, 6, 7, 7, 8, -1, 4, 3, 2, 5)

	if min != -1 {
		t.Errorf("min should be %d got %d", -1, min)
	}
	if max != -1 {
		t.Errorf("max should be %d got %d", 8, max)
	}
}

func TestIPLookup(t *testing.T) {
	details := GetPublicIPDetails(token)
	Logger("My Public IP Details: %+v", details)
}
