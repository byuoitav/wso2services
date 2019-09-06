package uapiclassschedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wso2services/classschedules/calendar"
	"github.com/byuoitav/wso2services/wso2requests"
)

//cache by yearterm -> room

var classScheduleCache map[string]map[string][]ClassSchedule
var updateTimesByRoom map[string]time.Time
var ttl = (24 * time.Hour) * -1

func init() {
	updateTimesByRoom = map[string]time.Time{}
	classScheduleCache = map[string]map[string][]ClassSchedule{}
}

//GetSimpleClassSchedulesForRoomAndTime - will use the local cache if it has been looked up before
func GetSimpleClassSchedulesForRoomAndTime(roomname string, classtime time.Time) ([]SimpleClassSchedule, *nerr.E) {
	t, err := calendar.GetYearTermForDate(classtime)

	if err != nil {
		return []SimpleClassSchedule{}, err.Addf("Couldn't get year erm for room %v and time %v", roomname, classtime)
	}

	termClassSchedules, err := GetSimpleClassSchedulesForRoomEnrollmentPeriod(roomname, strings.Replace(t.YearTermDesc, " ", "", -1))

	if err != nil {
		return []SimpleClassSchedule{}, err.Addf("Couldn't get simple class schedule for room %v and time %v", roomname, classtime)
	}

	var toReturn []SimpleClassSchedule

	for _, oneSchedule := range termClassSchedules {
		if (oneSchedule.StartDateTime.Before(classtime) || oneSchedule.StartDateTime.Equal(classtime)) &&
			(oneSchedule.EndDateTime.After(classtime) || oneSchedule.EndDateTime.Equal(classtime)) {
			toReturn = append(toReturn, oneSchedule)
		}
	}

	return toReturn, nil
}

//GetSimpleClassSchedulesForRoomEnrollmentPeriod does the translation
func GetSimpleClassSchedulesForRoomEnrollmentPeriod(roomname, enrollmentPeriod string) ([]SimpleClassSchedule, *nerr.E) {
	RawClassScheduleList, err := GetClassSchedulesForRoomEnrollmentPeriod(roomname, enrollmentPeriod)

	if err != nil {
		return []SimpleClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and enrollmentPeriod %v", roomname, enrollmentPeriod)
	}

	var toReturn []SimpleClassSchedule

	for _, oneSchedule := range RawClassScheduleList {
		//get the list of instructor names - we'll add to each date
		var instructorNames []string

		for _, instructor := range oneSchedule.AssignedInstructors.Values {
			instructorNames = append(instructorNames, instructor.ByuID.Description)
		}

		//now go through each assigned schedule and translate to dates / times
		for _, assignedSchedule := range oneSchedule.AssignedSchedules.Values {
			if assignedSchedule.Building.Value+"-"+assignedSchedule.Room.Value != roomname {
				//this API is weird and sometimes returns other rooms co-scheduled or something
				continue
			}

			startDate, err := time.Parse("2006-01-02", assignedSchedule.StartDate.Value)
			if err != nil {
				log.L.Errorf("Invalid start date when parsing schedule %v", assignedSchedule.StartDate.Value)
				continue
			}
			endDate, err := time.Parse("2006-01-02", assignedSchedule.EndDate.Value)
			if err != nil {
				log.L.Errorf("Invalid end date when parsing schedule %v", assignedSchedule.EndDate.Value)
				continue
			}
			startTime, err := time.Parse("15:04", assignedSchedule.StartTime.Value)
			if err != nil {
				log.L.Errorf("Invalid start time when parsing schedule %v", assignedSchedule.StartTime.Value)
				continue
			}
			endTime, err := time.Parse("15:04", assignedSchedule.EndTime.Value)
			if err != nil {
				log.L.Errorf("Invalid end time when parsing schedule %v", assignedSchedule.EndTime.Value)
				continue
			}

			//start at the first date, and loop through each day and see if it is part of the schedule
			//if it is, then create a SimpleClassSchedule struct for it and add to the toReturn array
			//Also translate the time to the right time zone (daylight savings.....nice....)
			curDate := startDate
			location, err := time.LoadLocation("America/Denver")
			if err != nil {
				log.L.Errorf("unable to parse America/Denver")
				continue
			}

			for curDate.Before(endDate) || curDate.Equal(endDate) {
				if (curDate.Weekday() == time.Sunday && assignedSchedule.Sun.Value) ||
					(curDate.Weekday() == time.Monday && assignedSchedule.Mon.Value) ||
					(curDate.Weekday() == time.Tuesday && assignedSchedule.Tue.Value) ||
					(curDate.Weekday() == time.Wednesday && assignedSchedule.Wed.Value) ||
					(curDate.Weekday() == time.Thursday && assignedSchedule.Thu.Value) ||
					(curDate.Weekday() == time.Friday && assignedSchedule.Fri.Value) ||
					(curDate.Weekday() == time.Saturday && assignedSchedule.Sat.Value) {
					StartDateTime := time.Date(curDate.Year(), curDate.Month(), curDate.Day(),
						startTime.Hour(), startTime.Minute(), 0, 0, location)

					EndDateTime := time.Date(curDate.Year(), curDate.Month(), curDate.Day(),
						endTime.Hour(), endTime.Minute(), 0, 0, location)

					thisNewRecord := SimpleClassSchedule{
						RoomID:          roomname,
						TeachingArea:    assignedSchedule.TeachingArea.Value,
						CourseNumber:    assignedSchedule.CourseNumber.Value,
						SectionNumber:   assignedSchedule.SectionNumber.Value,
						ScheduleType:    assignedSchedule.ScheduleType.Value,
						InstructorNames: instructorNames,
						StartDateTime:   StartDateTime,
						EndDateTime:     EndDateTime,
					}

					toReturn = append(toReturn, thisNewRecord)
				}

				curDate = curDate.Add(24 * time.Hour)
			}
		}
	}

	return toReturn, nil
}

//GetClassSchedulesForRoomEnrollmentPeriod - will use the local cache if it has been looked up before
func GetClassSchedulesForRoomEnrollmentPeriod(roomname, enrollmentPeriod string) ([]ClassSchedule, *nerr.E) {
	rmsplit := strings.Split(roomname, "-")

	//check to see if we have the class schedule cached for that term
	if termmap, ok := classScheduleCache[enrollmentPeriod]; ok {
		if time.Now().Add(ttl).Before(updateTimesByRoom[roomname]) {
			//check for the cache
			if cachedSchedules, ok := termmap[roomname]; ok {
				//check to see if it's up to date
				return cachedSchedules, nil
			}
			return []ClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")
		}
	} else {
		//nothing for this term, we need to initialze the map
		classScheduleCache[enrollmentPeriod] = map[string][]ClassSchedule{}
	}

	var classes []ClassSchedule

	var resp ClassResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v2/?subset_size=100&enrollment_periods=%v&building=%v&room=%v&contexts=class_schedule", enrollmentPeriod, rmsplit[0], rmsplit[1]), []byte{}, &resp)

	if err != nil {
		return classes, err.Addf("Couldn't fetch class scheudle")
	}

	for i := range resp.Values {
		classes = append(classes, resp.Values[i])
	}

	for resp.Metadata.PageStart+resp.Metadata.SubsetSize < resp.Metadata.CollectionSize {

		err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v2/?subset_size=100&enrollment_periods=%v&building=%v&room=%v&contexts=class_schedule&subset_start_offset=%v", enrollmentPeriod, rmsplit[0], rmsplit[1], resp.Metadata.PageStart+resp.Metadata.SubsetSize+1), []byte{}, &resp)

		if err != nil {
			return classes, err.Addf("Couldn't fetch class scheudle")
		}

		for i := range resp.Values {
			classes = append(classes, resp.Values[i])
		}
	}

	updateTimesByRoom[roomname] = time.Now()

	m := classScheduleCache[enrollmentPeriod]

	for i := range classes {
		var validAssignedSchedules []AssignedScheduleValue

		//cull out the ones that don't match the building/room
		for _, oneAssignedSchedule := range classes[i].AssignedSchedules.Values {
			if oneAssignedSchedule.Building.Value+"-"+oneAssignedSchedule.Room.Value == roomname {
				validAssignedSchedules = append(validAssignedSchedules, oneAssignedSchedule)
			}
		}

		classes[i].AssignedSchedules.Values = validAssignedSchedules

		//we go through and update the map
		rmname := fmt.Sprintf("%v-%v", classes[i].AssignedSchedules.Values[0].Building.Value, classes[i].AssignedSchedules.Values[0].Room.Value)
		if rmname != roomname {
			log.L.Fatalf("MISMATCH %v %v", rmname, roomname)
		}

		m[rmname] = append(m[rmname], classes[i])
	}
	classScheduleCache[enrollmentPeriod] = m

	if v, ok := classScheduleCache[enrollmentPeriod][roomname]; ok {
		return v, nil
	}

	return []ClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")
}
