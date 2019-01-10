package calendar

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wso2services/wso2requests"
)

// ControlDatesWSResponse is the response from the ControlDatesAPI
type ControlDatesWSResponse struct {
	ControlDatesWSService struct {
		Request struct {
			Method        string      `json:"method"`
			Resource      string      `json:"resource"`
			Attributes    interface{} `json:"attributes"`
			Status        int         `json:"status"`
			StatusMessage string      `json:"statusMessage"`
		} `json:"request"`
		Response struct {
			RequestCount int            `json:"request_count"`
			DateList     []YearTermDate `json:"date_list"`
		} `json:"response"`
	} `json:"ControlDatesWSService"`
}

//SortedYearTermDate .
type SortedYearTermDate []YearTermDate

//Len for sorted interface.
func (s SortedYearTermDate) Len() int {
	return len(s)
}

//Swap for sorted interface.
func (s SortedYearTermDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Less for sorted interface.
func (s SortedYearTermDate) Less(i, j int) bool {
	return time.Time(s[i].StartDate).Before(time.Time(s[j].StartDate))
}

//YearTermDate .
type YearTermDate struct {
	DateType     string          `json:"date_type"`
	YearTerm     string          `json:"year_term"`
	YearTermDesc string          `json:"year_term_desc"`
	StartDate    ControlDateTime `json:"start_date"`
	EndDate      ControlDateTime `json:"end_date"`
	Description  string          `json:"description"`
}

const (
	controlDateFormat = "20060102 15:04:05"
)

//ControlDateTime to have custom unmarshaler
type ControlDateTime time.Time

//UnmarshalJSON .
func (m *ControlDateTime) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(controlDateFormat, strings.TrimSpace(strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	)))

	if err != nil {
		return err
	}

	*m = ControlDateTime(t)

	return nil
}
func init() {
	checkMutex = sync.Mutex{}
}

//assume that we don't need to get it that often, like maybe once a day...
var classSchedulelastUpdated time.Time
var dates []YearTermDate
var classSchedulecheckMutex sync.Mutex
var lastAccessedTermIndex int

var classScheduleinterval = time.Duration(24*time.Hour) * -1

func getControlDates() ([]YearTermDate, *nerr.E) {

	toReturn := []YearTermDate{}
	resp := ControlDatesWSResponse{}
	//make our WSO2 call
	err := wso2requests.MakeWSO2Request("GET", "https://api.byu.edu/domains/legacy/academic/controls/controldatesws/v1/all/CLASS_DATES", []byte{}, &resp)

	if err != nil {
		return toReturn, err.Addf("Couldn't get control dates")
	}

	for i := range resp.ControlDatesWSService.Response.DateList {
		toReturn = append(toReturn, resp.ControlDatesWSService.Response.DateList[i])
	}

	return toReturn, nil
}

//GetYearTermForDate .
func GetYearTermForDate(date time.Time) (YearTermDate, *nerr.E) {

	classSchedulecheckMutex.Lock()
	//check to see if last updated is more than interval away
	if classSchedulelastUpdated.Before(time.Now().Add(classScheduleinterval)) {
		v, err := getControlDates()
		if err != nil {
			log.L.Errorf(err.Error())
			classSchedulecheckMutex.Unlock()
			return YearTermDate{}, err.Addf("Control dates are out of date, but can't update")
		}

		//we should sort v, just to be safe
		sort.Sort(SortedYearTermDate(v))

		dates = v
		classSchedulelastUpdated = time.Now()
	}
	classSchedulecheckMutex.Unlock()

	//check to see if it's the last on we used, so we don't have to iterate through all of them every time.
	if date.After(time.Time(dates[lastAccessedTermIndex].StartDate)) && date.Before(time.Time(dates[lastAccessedTermIndex].EndDate)) {
		return dates[lastAccessedTermIndex], nil
	}

	//go through and check to see which we fall into
	for i := range dates {
		if date.After(time.Time(dates[i].StartDate)) && date.Before(time.Time(dates[i].EndDate)) {
			lastAccessedTermIndex = i
			return dates[i], nil
		}
	}

	return YearTermDate{}, nerr.Create(fmt.Sprintf("no YearTerm containing date %v was found.", date), "notfound")
}
