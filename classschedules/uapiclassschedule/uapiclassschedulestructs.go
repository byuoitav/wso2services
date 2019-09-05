package uapiclassschedule

import "time"

//SimpleClassSchedule is what we can translate to for ease of getting what we need out of it
type SimpleClassSchedule struct {
	RoomID          string //BLDG-ROOM
	TeachingArea    string
	CourseNumber    string
	SectionNumber   string
	ScheduleType    string
	StartDateTime   time.Time
	EndDateTime     time.Time
	InstructorNames []string
}

//ClassResponse is the body we unmarshal the response into.
type ClassResponse struct {
	Metadata Metadata        `json:"metadata"`
	Values   []ClassSchedule `json:"values"`
}

//Metadata .
type Metadata struct {
	CollectionSize int `json:"collection_size"`
	PageStart      int `json:"subset_start"`
	SubsetSize     int `json:"subset_size"`
}

//ClassSchedule .
type ClassSchedule struct {
	AssignedInstructors AssignedInstructors `json:"assigned_instructors"`
	AssignedSchedules   AssignedSchedule    `json:"assigned_schedules"`

	Headers struct {
		Links struct {
		} `json:"links"`
		Metadata struct {
			CollectionSize int `json:"collection_size"`
		} `json:"metadata"`
		Values []interface{} `json:"values"`
	} `json:"headers"`
}

// AssignedInstructors .
type AssignedInstructors struct {
	Values []struct {
		EnrollmentPeriod StringValueStruct `json:"enrollment_period"`
		TeachingArea     StringValueStruct `json:"teaching_area"`
		CourseNumber     StringValueStruct `json:"course_number"`
		SectionNumber    StringValueStruct `json:"section_number"`
		InstructorType   StringValueStruct `json:"instructor_type"`
		//skipping class_identifiers
		PersonID        StringValueStruct `json:"person_id"`
		ByuID           StringValueStruct `json:"byu_id"` //description contains name
		NetID           StringValueStruct `json:"net_id"`
		EmailAddress    StringValueStruct `json:"email_address"`
		UpdatedByBYUID  StringValueStruct `json:"updated_by_byu_id"`
		DateTimeUpdated StringValueStruct `json:"date_time_updated"`
	} `json:"values"`
}

//StringValueStruct .
type StringValueStruct struct {
	Value           string `json:"value"`
	APIType         string `json:"api_type"`
	Key             bool   `json:"key"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
}

//IntValueStruct ...
type IntValueStruct struct {
	Value           int    `json:"value"`
	APIType         string `json:"api_type"`
	Key             bool   `json:"key"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
}

//BoolValueStruct ...
type BoolValueStruct struct {
	Value           bool   `json:"value"`
	APIType         string `json:"api_type"`
	Key             bool   `json:"key"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
	Domain          string `json:"domain"`
}

// AssignedSchedule .
type AssignedSchedule struct {
	Values []struct {
		EnrollmentPeriod StringValueStruct `json:"enrollment_period"`
		TeachingArea     StringValueStruct `json:"teaching_area"`
		CourseNumber     StringValueStruct `json:"course_number"`
		SectionNumber    StringValueStruct `json:"section_number"`
		ScheduleType     StringValueStruct `json:"schedule_type"`
		SequenceNumber   IntValueStruct    `json:"sequence_number"`
		//skipping class_identifiers
		CreditInstitution   StringValueStruct `json:"credit_institution"`
		Building            StringValueStruct `json:"building"`
		Room                StringValueStruct `json:"room"`
		Days                StringValueStruct `json:"days"`
		Mon                 BoolValueStruct   `json:"mon"`
		Tue                 BoolValueStruct   `json:"tue"`
		Wed                 BoolValueStruct   `json:"wed"`
		Thu                 BoolValueStruct   `json:"thu"`
		Fri                 BoolValueStruct   `json:"fri"`
		Sat                 BoolValueStruct   `json:"sat"`
		Sun                 BoolValueStruct   `json:"sun"`
		StartDate           StringValueStruct `json:"start_date"`
		EndDate             StringValueStruct `json:"end_date"`
		StartTime           StringValueStruct `json:"start_time"`
		EndTime             StringValueStruct `json:"end_time"`
		ScheduleStatus      StringValueStruct `json:"schedule_status"`
		ScheduleDescription StringValueStruct `json:"schedule_description"`
		DateTimeUpdated     StringValueStruct `json:"updated_datetime"`
		UpdatedByBYUID      StringValueStruct `json:"updated_by_byu_id"`
		AllowConflict       BoolValueStruct   `json:"allow_conflict"`
		RoomCapacity        IntValueStruct
	} `json:"values"`
}
