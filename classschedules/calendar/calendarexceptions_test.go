package calendar

import (
	"testing"
	"time"

	"github.com/byuoitav/common/log"
)

//TestCalendarExceptions .
func TestCalendarExceptions(t *testing.T) {

	//we just want to validate that we can get the calendar exceptions for a few days
	log.SetLevel("debug")

	ti, err := time.Parse("2006-01-02T05:04", "2018-12-13T09:03")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e, er := CheckExceptionsForDate(ti)
	if er != nil {
		t.Error(er)
		t.FailNow()
	}

	t.Logf("%v", e)

	ti, err = time.Parse("2006-01-02T05:04", "2018-11-22T09:03")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e, er = CheckExceptionsForDate(ti)
	if er != nil {
		t.Error(er)
		t.FailNow()
	}

	t.Logf("%v", e)
	return
}
