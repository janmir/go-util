package util

import (
	"fmt"
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
	Loggerf("✔ %d:%s ", 1, "This is a formatted Log")
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
	Loggerf("My Public IP Details: %+v", details)
}

func TestNTP(t *testing.T) {
	ntp, err := GetNTPTime(AppleNTP)
	if err != nil {
		t.Fail()
	}
	Loggerf("Apple NTP Time: %+v", ntp)

	ntp, err = GetNTPTime(GoogleNTP)
	if err != nil {
		t.Fail()
	}
	_ = ntp
	Loggerf("Google NTP Time: %+v", ntp)
}

func TestMapDecode(t *testing.T) {
	strct := struct {
		Key1 string
		Key2 string
	}{}
	_ = strct
	str := ""
	_ = str
	okmap := map[string]string{
		"Key1": "Value01",
		"key2": "Value02",
		"Key3": "Value02",
	}
	_ = okmap
	empty := map[string]string{}
	_ = empty

	err := MapDecode(&strct, okmap)
	fmt.Printf("Error: %v, %+v\n", err, strct)
}
