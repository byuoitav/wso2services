package registar

import (
	"testing"
	"time"

	"github.com/byuoitav/common/log"
)

func TestDaysOfWeek(t *testing.T) {

	str := []string{
		"TTh",
		"MWF",
		"MTWThFSSu",
		"MThS",
		"M",
		"Th",
	}

	answers := [][]time.Weekday{
		[]time.Weekday{time.Tuesday, time.Thursday},
		[]time.Weekday{time.Monday, time.Wednesday, time.Friday},
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
		[]time.Weekday{time.Monday, time.Thursday, time.Saturday},
		[]time.Weekday{time.Monday},
		[]time.Weekday{time.Thursday},
	}

	for i := range str {
		t.Logf("testing %v\n", str[i])
		tmp := FindDaysOfWeek(str[i])
		t.Logf("answers %v\n", tmp)
		if len(tmp) != len(answers[i]) {
			t.Fail()
		}
		for j := range tmp {
			if tmp[j] != answers[i][j] {
				t.Fail()
			}
		}
	}
}

func TestGetSchedule(t *testing.T) {

	log.SetLevel("debug")
	classtime := time.Date(2018, 11, 20, 9, 15, 00, 00, time.UTC)
	schedule, err := GetClassScheduleForTime("CTB-410", classtime)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	log.L.Debugf("%v", schedule)

}
