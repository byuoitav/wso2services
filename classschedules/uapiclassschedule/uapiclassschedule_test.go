package uapiclassschedule

import (
	"testing"
	"time"

	"github.com/byuoitav/common/log"
)

func TestGetClassSchedule(t *testing.T) {
	log.SetLevel("debug")

	ti, nerr := time.Parse("2006-01-02T15:04:05-07:00", "2019-04-16T07:59:44-06:00")
	if nerr != nil {
		t.Errorf("%v", nerr.Error())
	}

	log.L.Debugf("time: %v", ti)

	sched, err := GetSimpleClassSchedulesForRoomAndDate("B66-120", ti)
	if err != nil {
		t.Error(err)
	}

	for _, one := range sched {
		log.L.Debugf("Schedule: %+v\n", one)
	}

}
