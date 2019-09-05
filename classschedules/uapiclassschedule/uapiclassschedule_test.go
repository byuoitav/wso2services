package uapiclassschedule

import (
	"testing"
	"time"

	"github.com/byuoitav/common/log"
)

func TestGetClassSchedule(t *testing.T) {
	log.SetLevel("debug")

	log.L.Debugf("start: %v", time.Now().Unix())

	sched, err := GetSimpleClassSchedulesForRoomEnrollmentPeriod("B66-120", "Fall2019")
	if err != nil {
		t.Error(err)
	}

	log.L.Debugf("Schedule: %v %v", len(sched), time.Now().Unix())

	log.L.Debugf("Schedule: %v %v", len(sched), time.Now().Unix())

	ti, nerr := time.ParseInLocation("2006-01-02 15:04", "2019-10-15 14:12", time.Local)
	if nerr != nil {
		t.Errorf("%v", nerr.Error())
	}

	log.L.Debugf("%v", ti)

	sched, err = GetSimpleClassSchedulesForRoomAndTime("B66-120", ti)
	if err != nil {
		t.Error(err)
	}

	log.L.Debugf("Schedule: %+v, %v", sched, time.Now().Unix())

	complex, err := GetClassSchedulesForRoomEnrollmentPeriod("B66-120", "Fall2019")
	for _, compl := range complex {
		for _, sch := range compl.AssignedSchedules.Values {
			if sch.TeachingArea.Value == "CFM" && sch.CourseNumber.Value == "217" {
				log.L.Debugf("%v %v %v %v %v %v %v", sch.StartDate.Value, sch.EndDate.Value, sch.StartTime.Value, sch.EndTime.Value, sch.Days, sch.SectionNumber, sch.ScheduleType)
			}
		}
	}
}
