package classschedules

import (
	"sort"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wos2services/wso2requests"
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
	return s[i].StartDate.Before(s[j].StartDate)
}

//YearTermDate .
type YearTermDate struct {
	DateType     string    `json:"date_type"`
	YearTerm     string    `json:"year_term"`
	YearTermDesc string    `json:"year_term_desc"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Description  string    `json:"description"`
}

func init() {
	checkMutex = sync.Mutex{}
}

//assume that we don't need to get it that often, like maybe once a day...
var lastUpdated time.Time
var dates []YearTermDate
var checkMutex sync.Mutex
var lastAccessedTermIndex int

var interval = time.Duration(24 * time.Hour)

func getControlDates() ([]YearTermDate, *nerr.E) {

	toReturn := []YearTermDate{}
	resp := ControlDatesWSResponse{}
	//make our WSO2 call
	err := wso2requests.MakeWSO2Request("GET", "https://api.byu.edu/domains/legacy/academic/controls/controldatesws/V1/all/CLASS_DATES", []byte{}, &resp)

	if err != nil {
		return toReturn, err.Addf("Couldn't get control dates")
	}

	for i := range resp.ConrolDatesWSService.ResponseList {
		toReturn = append(toReturn, resp.ConrolDatesWSService.ResponseList[i])
	}
}

//GetYearTermForDate .
func GetYearTermForDate(date time.Time) (YearTermDate, *nerr.E) {

	checkMutex.Lock()
	//check to see if last updated is more than interval away
	if lastUpdateTime < time.Now().Sub(interval) {
		v, err := getControlDates()
		if err != nil {
			log.L.Errorf(err.Error())
			return YearTermDate{}, err.Addf("Control dates are out of date, but can't update")
		}

		//we should sort v, just to be safe
		sort.Sort(SortedYearTermDate(v))

		dates = v
		lastUpdated = time.Now()
	}
	checkMutex.Unlock()

	//check to see if it's the last on we used, so we don't have to iterate through all of them every time.
	if date.After(dates[lastAccessedTermIndex].StartDate) && date.Before(dates[lastAccessedTermIndex].EndDate) {
		return dates[lastAccessedTermIndex], nil
	}

	//go through and check to see which we fall into
	for i := range dates {
		if date.After(dates[i].StartDate) && date.Before(dates[i].EndDate) {
			lastAccessedTermIndex = i
			return dates[i], nil
		}
	}

	return YearTermDate, nerr.Create(fmt.Spritnf("no YearTerm containing date %v was found.", date), "notfound")
}
