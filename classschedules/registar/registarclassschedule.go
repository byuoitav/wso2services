package registar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wso2services/classschedules/calendar"
	"github.com/byuoitav/wso2services/wso2requests"
)

//ClassScheduleResponse .
type ClassScheduleResponse struct {
	ClassRoomService struct {
		Request struct {
			Method        string      `json:"method"`
			Resource      string      `json:"resource"`
			Attributes    interface{} `json:"attributes"`
			Status        int         `json:"status"`
			StatusMessage string      `json:"statusMessage"`
		} `json:"request"`
		Room Room `json:"response"`
	} `json:"ClassRoomService"`
}

//ClassSchedule .
type ClassSchedule struct {
	DeptName       string      `json:"dept_name"`
	CatalogNumber  string      `json:"catalog_number"`
	CatalogSuffix  interface{} `json:"catalog_suffix"`
	LabQuizSection interface{} `json:"lab_quiz_section"`
	Honors         interface{} `json:"honors"`
	ServLearning   interface{} `json:"serv_learning"`
	CreditHours    float64     `json:"credit_hours"`
	SectionType    interface{} `json:"section_type"`
	ClassTime      string      `json:"class_time"`
	Days           string      `json:"days"`
	InstructorName string      `json:"instructor_name"`
	SectionSize    int         `json:"section_size"`
	TotalEnr       int         `json:"total_enr"`
	SchedType      string      `json:"sched_type"`
	AssignTo       interface{} `json:"assign_to"`
	StartDate      interface{} `json:"start_date"`
	EndDate        interface{} `json:"end_date"`

	StartTime time.Time //the zero time with the hours/minutes of start of class set
	EndTime   time.Time //the zero time with the hours/minutes of end of class set
}

//Room .
type Room struct {
	Room              string          `json:"room"`
	Building          string          `json:"building"`
	SchedulingType    string          `json:"scheduling_type"`
	RoomDesc          string          `json:"room_desc"`
	Capacity          int             `json:"capacity"`
	EffectiveYearTerm string          `json:"effective_year_term"`
	ExpiredYearTerm   string          `json:"expired_year_term"`
	ROOMDATA          interface{}     `json:"ROOMDATA"`
	Schedules         []ClassSchedule `json:"schedules"`

	lastUpdate time.Time //when we last fetched this
}

const (
	registarRecordTTL = -24 * time.Hour // once a day
)

func init() {

	cache = map[string]map[string]Room{}
	updateTimes = map[string]time.Time{}

	classtimereg = regexp.MustCompile(datere)
}

var cache map[string]map[string]Room //yearterm -> room
var updateTimes map[string]time.Time
var classtimereg *regexp.Regexp

//GetClassScheduleForTime .
func GetClassScheduleForTime(roomname string, classtime time.Time) (ClassSchedule, *nerr.E) {

	t, err := calendar.GetYearTermForDate(classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, classtime)
	}

	rm, err := getRoomInfoForYearTerm(roomname, t)

	//figure out what class classtime falls into
	//check to see if there are any exceptions for that date
	exceptions, err := calendar.CheckExceptionsForDate(classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
	}

	if len(exceptions) > 0 {
		//we have some exeptions, check to see what kind?
		return FindScheduleWithExceptions(exceptions[0], rm, classtime)
	}
	return FindSchedule(rm, classtime, classtime.Weekday())
}

//GetClassScheduleForTimeBlock .
func GetClassScheduleForTimeBlock(roomname string, start, end time.Time) ([]ClassSchedule, *nerr.E) {
	loc, er := time.LoadLocation("America/Denver")
	if er != nil {
		log.L.Errorf("Couldn't load timezone")
		return []ClassSchedule{}, nerr.Translate(er)
	}

	t, err := calendar.GetYearTermForDate(start)
	if err != nil {
		return []ClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, start)
	}

	rm, err := getRoomInfoForYearTerm(roomname, t)

	if time.Time(t.EndDate).Before(end) {
		//TODO: we're gonna need to make multiple calls.
		log.L.Errorf("Term spanning block %v - %v", start, end)
		return []ClassSchedule{}, nerr.Create(fmt.Sprintf("Term spanning block %v - %v", start, end), "block-too-large")
	}

	//start at start
	curStart := start.In(loc)
	curEnd := endOfDay(curStart)

	toReturn := []ClassSchedule{}

	for curStart.Before(end) && !curStart.Equal(end) {
		log.L.Debugf("Checking schedules for block %v - %v", curStart.In(time.Local), curEnd.In(time.Local))

		if curEnd.After(end) {
			log.L.Debugf("Cur end %v is after block end, checking block %v - %v", curEnd.In(time.Local), curStart.In(time.Local), end)
			curEnd = end
		}

		exceptions, err := calendar.CheckExceptionsForDate(curStart)
		if err != nil {
			return []ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
		}

		if len(exceptions) > 0 {
			switch exceptions[0].Category {
			case calendar.Mondayinstruction:
				tmp, err := findSchedulesInBlock(rm, curStart, curEnd, time.Monday)
				if err != nil {
					return toReturn, err.Addf("Couldn't get classes for period %v = %v", curStart, curEnd)
				}
				toReturn = append(toReturn, tmp...)
			case calendar.Fridayinstruction:
				tmp, err := findSchedulesInBlock(rm, curStart, curEnd, time.Friday)
				if err != nil {
					return toReturn, err.Addf("Couldn't get classes for period %v = %v", curStart, curEnd)
				}
				toReturn = append(toReturn, tmp...)
			default: //holiday or no class
				return []ClassSchedule{}, nil
			}
		} else {
			//no exceptions
			tmp, err := findSchedulesInBlock(rm, curStart, curEnd, curStart.Weekday())
			if err != nil {
				return toReturn, err.Addf("Couldn't get classes for period %v = %v", curStart, curEnd)
			}

			toReturn = append(toReturn, tmp...)
		}

		curStart = startOfDay(curStart.AddDate(0, 0, 1))
		curEnd = endOfDay(curStart)

	}

	return toReturn, nil
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

func getRoomInfoForYearTerm(roomname string, t calendar.YearTermDate) (Room, *nerr.E) {

	//check to see if we have it in the cache
	term, ok := cache[t.YearTerm]
	if !ok {

		rm, err := fetchClassSchedule(roomname, t.YearTerm)
		if err != nil {
			updateTimes[roomname] = time.Now()
			return Room{}, err.Addf("Couldn't get class schedule for time.")
		}
		//we need to fetch the whole thing
		term = map[string]Room{
			roomname: rm,
		}
		cache[t.YearTerm] = term
		updateTimes[roomname] = time.Now()
	}

	//check for the room
	rm, ok := term[roomname]
	if !ok {
		//check to see when the last time we checked was, if it was less than 30 minutes we don't go try again, we just assume it's a bad room
		ck, ok := updateTimes[roomname]
		if !ok || time.Now().Sub(ck) > 5*time.Minute {

			//we go get it .
			updateTimes[roomname] = time.Now()
			rm, err := fetchClassSchedule(roomname, t.YearTerm)
			if err != nil {
				return Room{}, err.Addf("Couldn't get class schedule for time.")
			}
			term[roomname] = rm
		}
	} else {
		//we need to check update times here, too
		ck, ok := updateTimes[roomname]
		if !ok || time.Now().Sub(ck) > 5*time.Minute {

			//we go get it .
			rm, err := fetchClassSchedule(roomname, t.YearTerm)
			updateTimes[roomname] = time.Now()
			if err != nil {
				return Room{}, err.Addf("Couldn't get class schedule for time.")
			}
			term[roomname] = rm
		}
	}

	return rm, nil
}

//FindScheduleWithExceptions .
func FindScheduleWithExceptions(exception calendar.Exception, rm Room, t time.Time) (ClassSchedule, *nerr.E) {

	switch exception.Category {
	case calendar.Mondayinstruction:
		return FindSchedule(rm, t, time.Monday)
	case calendar.Fridayinstruction:
		return FindSchedule(rm, t, time.Friday)
	default: //holiday or no class
		return ClassSchedule{}, nil
	}

}

//Assumes the block starts and ends in the same day. If the block needs to span days, make separate calls for each day.
func findSchedulesInBlock(rm Room, start, end time.Time, weekday time.Weekday) ([]ClassSchedule, *nerr.E) {
	toReturn := []ClassSchedule{}
	log.L.Debugf("Finding Schedules in block %v - %v. It's a %v", start.In(time.Local), end.In(time.Local), weekday)

	for _, v := range rm.Schedules {
		wkdays := FindDaysOfWeek(v.Days)

		var ok bool
		for i := range wkdays {
			if wkdays[i] == weekday {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}

		//check to see if windows overlap at all...
		//we're on the right day of the week, check to see if we're in the right time
		dayclassStart := time.Date(start.Year(), start.Month(), start.Day(), v.StartTime.Hour(), v.StartTime.Minute(), 0, 0, v.StartTime.Location())
		dayclassStart = dayclassStart.Add(-5 * time.Minute) //for a buffer

		//we're on the right day of the week, check to see if we're in the right time
		dayclassEnd := time.Date(end.Year(), end.Month(), end.Day(), v.EndTime.Hour(), v.EndTime.Minute(), 0, 0, v.EndTime.Location())
		dayclassEnd = dayclassEnd.Add(5 * time.Minute) //for a buffer

		//check to see if start time is before block start or if class end is after start time
		if dayclassStart.Before(end) && dayclassEnd.After(start) {

			//so we can do comparisons later...
			v.EndTime = dayclassEnd
			v.StartTime = dayclassStart

			if len(toReturn) == 0 {

				//bingo.
				toReturn = append(toReturn, v)
			} else {

				toCompare := toReturn[len(toReturn)-1]
				//check to see if this section happens at the same time as another class of the same name. If yes, we just add the enrollment/size numbers.
				if toCompare.DeptName == v.DeptName && toCompare.CatalogNumber == v.CatalogNumber && toCompare.StartTime.Equal(v.StartTime) && toCompare.EndTime.Equal(v.EndTime) {
					toCompare.SectionSize = v.SectionSize + toCompare.SectionSize
					toCompare.TotalEnr = v.TotalEnr + toCompare.TotalEnr
					toReturn[len(toReturn)-1] = toCompare

				} else {
					//bingo.
					toReturn = append(toReturn, v)
				}
			}
		}
	}

	return toReturn, nil
}

//FindSchedule .
func FindSchedule(rm Room, t time.Time, weekday time.Weekday) (ClassSchedule, *nerr.E) {
	log.L.Debugf("Finding schedule for %v at %v. it's a %v", rm.Building+rm.Room, t, weekday)
	//log.L.Debugf("%v has %v classes", rm.Building+rm.Room, len(rm.Schedules))

	for _, v := range rm.Schedules {
		wkdays := FindDaysOfWeek(v.Days)

		var ok bool
		for i := range wkdays {
			if wkdays[i] == weekday {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}
		//		log.L.Debugf("Class %v occurs on %v", v.DeptName+v.CatalogNumber, weekday)

		//we're on the right day of the week, check to see if we're in the right time
		dayclassStart := time.Date(t.Year(), t.Month(), t.Day(), v.StartTime.Hour(), v.StartTime.Minute(), 0, 0, time.UTC)
		dayclassStart = dayclassStart.Add(-5 * time.Minute) //for a buffer

		//		log.L.Debugf("Class %v started at %v ", v.DeptName+v.CatalogNumber, dayclassStart)

		if t.Before(dayclassStart) {
			continue
		}
		//we're on the right day of the week, check to see if we're in the right time
		dayclassEnd := time.Date(t.Year(), t.Month(), t.Day(), v.EndTime.Hour(), v.EndTime.Minute(), 0, 0, time.UTC)
		dayclassEnd = dayclassEnd.Add(5 * time.Minute) //for a buffer

		//		log.L.Debugf("Class %v ended at %v ", v.DeptName+v.CatalogNumber, dayclassEnd)

		if t.After(dayclassEnd) {
			continue
		}

		//bingo.
		return v, nil
	}

	//didn't find a class

	return ClassSchedule{}, nil
}

//FindDaysOfWeek .
func FindDaysOfWeek(wkstr string) []time.Weekday {
	toReturn := []time.Weekday{}

	for i := 0; i < len(wkstr); i++ {
		switch wkstr[i] {
		case 'M':
			toReturn = append(toReturn, time.Monday)
		case 'T':
			if i+1 != len(wkstr) && wkstr[i+1] == 'h' {
				toReturn = append(toReturn, time.Thursday)
				i++
			} else {
				toReturn = append(toReturn, time.Tuesday)
			}
		case 'W':
			toReturn = append(toReturn, time.Wednesday)
		case 'F':
			toReturn = append(toReturn, time.Friday)
		case 'S':
			if i+1 != len(wkstr) && wkstr[i+1] == 'a' {
				toReturn = append(toReturn, time.Saturday)
				i++
			} else {
				toReturn = append(toReturn, time.Saturday)
			}
		}
	}

	return toReturn
}

const datere = `(\d{1,2}):(\d{2})([ap]) - (\d{1,2}):(\d{2})([ap])`

func fetchClassSchedule(roomname, term string) (Room, *nerr.E) {
	defaultOffset, er := time.LoadLocation("MST")
	if er != nil {
		log.L.Errorf("Couldn't load MST timezone")
		return Room{}, nerr.Translate(er)
	}

	//we figure out the building
	br := strings.Split(roomname, "-")

	var resp ClassScheduleResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu:443/domains/legacy/academic/classschedule/classroom/v1/%v/%v/%v/Schedule", term, br[0], br[1]), []byte{}, &resp)

	if err != nil {
		return resp.ClassRoomService.Room, err.Addf("Couldn't fetch class scheudle")
	}

	//for each time, we'll need to go and parse the classtime to figure out the deal
	rm := resp.ClassRoomService.Room

	for i := range rm.Schedules {
		matches := classtimereg.FindAllStringSubmatch(rm.Schedules[i].ClassTime, -1)
		if len(matches) != 1 {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		shr, err := strconv.Atoi(matches[0][1])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		smin, err := strconv.Atoi(matches[0][2])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}
		if matches[0][3] == "p" {
			shr = shr + 12
		}

		rm.Schedules[i].StartTime = time.Date(0, 0, 0, shr, smin, 0, 0, defaultOffset)

		ehr, err := strconv.Atoi(matches[0][4])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		emin, err := strconv.Atoi(matches[0][5])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		if matches[0][6] == "p" {
			ehr = ehr + 12
		}

		rm.Schedules[i].EndTime = time.Date(0, 0, 0, ehr, emin, 0, 0, defaultOffset)

		log.L.Debugf("%v-%v start %v, end %v,", rm.Schedules[i].DeptName, rm.Schedules[i].CatalogNumber, rm.Schedules[i].StartTime, rm.Schedules[i].EndTime)
	}

	return resp.ClassRoomService.Room, nil
}
