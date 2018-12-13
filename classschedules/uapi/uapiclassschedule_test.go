package uapi

import (
	"testing"
	"time"
)

func TestGetClassSchedule(t *testing.T) {

	ti, err := time.Parse("2006-01-02T05:04", "2018-12-13T09:03")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	schedules, er := getClassSchedule("CTB-410", ti)
	if er != nil {
		t.Error(er)
		t.FailNow()
	}

	t.Log(schedules)

}
