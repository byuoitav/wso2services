package uapi

//ClassResponse  is the body we unmarshal the  response into.
type ClassResponse struct {
	Metadata Metadata        `json:"metadata"`
	Values   []ClassSchedule `json:"values"`
}

//Metadata .
type Metadata struct {
	CollectionSize     int           `json:"collection_size"`
	PageStart          int           `json:"page_start"`
	PageEnd            int           `json:"page_end"`
	PageSize           int           `json:"page_size"`
	DefaultPageSize    int           `json:"default_page_size"`
	MaxPageSize        int           `json:"max_page_size"`
	FieldSetsReturned  []interface{} `json:"field_sets_returned"`
	FieldSetsAvailable []string      `json:"field_sets_available"`
	DefaultFieldSets   []string      `json:"default_field_sets"`
	ContextsAvailable  struct {
		ClassRoll           []string `json:"class_roll"`
		ClassSchedule       []string `json:"class_schedule"`
		ClassScheduleCache  []string `json:"class_schedule_cache"`
		ClassScheduleRecord []string `json:"class_schedule_record"`
		Proof               []string `json:"proof"`
		Summary             []string `json:"summary"`
	} `json:"contexts_available"`
	DefaultDb []string `json:"default_db"`
}

//ClassSchedule .
type ClassSchedule struct {
	AssignedInstructors AssignedInstructors `json:"assigned_instructors"`
	AssignedSchedules   AssignedSchedule    `json:"assigned_schedules"`
	EnrollmentCounts    EnrollmentCounts    `json:"enrollment_counts"`
	Basic               ClassScheduleBasic  `json:"basic"`

	Headers struct {
		Links struct {
		} `json:"links"`
		Metadata struct {
			CollectionSize int `json:"collection_size"`
		} `json:"metadata"`
		Values []interface{} `json:"values"`
	} `json:"headers"`
}

// ClassScheduleBasic .
type ClassScheduleBasic struct {
	YearTerm struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
		Desc    string `json:"desc"`
		ExtDesc string `json:"ext_desc"`
	} `json:"year_term"`
	CurriculumID struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"curriculum_id"`
	TitleCode struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"title_code"`
	SectionNumber struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"section_number"`
	CreditInstitution struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"credit_institution"`
	EnrollmentStatus struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"enrollment_status"`
	TeachingArea struct {
		Value           string `json:"value"`
		APIType         string `json:"api_type"`
		RelatedResource string `json:"related_resource"`
	} `json:"teaching_area"`
	CourseNumber struct {
		Value           string `json:"value"`
		APIType         string `json:"api_type"`
		RelatedResource string `json:"related_resource"`
	} `json:"course_number"`
	CourseSuffix struct {
		Value           string `json:"value"`
		APIType         string `json:"api_type"`
		RelatedResource string `json:"related_resource"`
	} `json:"course_suffix"`
	CourseTitle struct {
		Value           string `json:"value"`
		APIType         string `json:"api_type"`
		RelatedResource string `json:"related_resource"`
	} `json:"course_title"`
	SectionType struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"section_type"`
	BlockCode struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Desc    string `json:"desc"`
	} `json:"block_code"`
	ClassStatus struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"class_status"`
	RegMethod struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"reg_method"`
	LabQuizFlag struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"lab_quiz_flag"`
	FinalExamFlag struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"final_exam_flag"`
	Fee struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"fee"`
	Honors struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"honors"`
	ServiceLearning struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"service_learning"`
	ClassSize struct {
		Value   int    `json:"value"`
		APIType string `json:"api_type"`
	} `json:"class_size"`
	FixedOrVariable struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"fixed_or_variable"`
	FixedOrVariableCourse struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"fixed_or_variable_course"`
	CreditHours struct {
		Value   float64 `json:"value"`
		APIType string  `json:"api_type"`
	} `json:"credit_hours"`
	MinimumCreditHours struct {
		Value           float64 `json:"value"`
		APIType         string  `json:"api_type"`
		RelatedResource string  `json:"related_resource"`
	} `json:"minimum_credit_hours"`
	LectureHours struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"lecture_hours"`
	LabHours struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"lab_hours"`
	ClassStartDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"class_start_date"`
	ClassEndDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"class_end_date"`
	RegStartDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"reg_start_date"`
	RegEndDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"reg_end_date"`
	WithdrawDeadlineDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"withdraw_deadline_date"`
	ControlMixDate struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"control_mix_date"`
	WaitlistStatus struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"waitlist_status"`
	GradeRule struct {
		Value       string `json:"value"`
		Description string `json:"description"`
		APIType     string `json:"api_type"`
	} `json:"grade_rule"`
	GradeRollType struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"grade_roll_type"`
	GradeYearTerm struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"grade_year_term"`
	CombineLinkedRolls struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"combine_linked_rolls"`
	CombineLinkedCounts struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"combine_linked_counts"`
	LinkRelation struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"link_relation"`
	LinkToYearTerm struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"link_to_year_term"`
	LinkToCurriculumID struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"link_to_curriculum_id"`
	LinkToTitleCode struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"link_to_title_code"`
	LinkToSectionNumber struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"link_to_section_number"`
	CourseDescription struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_description"`
	CourseWhenTaught struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_when_taught"`
	CoursePrerequisites struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_prerequisites"`
	CourseRequired struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_required"`
	CourseRecommended struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_recommended"`
	CourseNote struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"course_note"`
}

// EnrollmentCounts  .
type EnrollmentCounts struct {
	Links struct {
		GetEnrollmentCounts struct {
			Rel    string `json:"rel"`
			Href   string `json:"href"`
			Method string `json:"method"`
			Title  string `json:"title"`
		} `json:"getEnrollmentCounts"`
		EnrollStudentInClass struct {
			Rel    string `json:"rel"`
			Href   string `json:"href"`
			Method string `json:"method"`
			Title  string `json:"title"`
		} `json:"enrollStudentInClass"`
	} `json:"links"`
	Metadata struct {
		Message   string   `json:"message"`
		DefaultDb []string `json:"default_db"`
	} `json:"metadata"`
	YearTerm struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
		Desc    string `json:"desc"`
		ExtDesc string `json:"ext_desc"`
	} `json:"year_term"`
	CurriculumID struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"curriculum_id"`
	TitleCode struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"title_code"`
	SectionNumber struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"section_number"`
	CreditInstitution struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
		Key     bool   `json:"key"`
	} `json:"credit_institution"`
	ClassSize struct {
		Value   int    `json:"value"`
		APIType string `json:"api_type"`
	} `json:"class_size"`
	TotalEnrolled struct {
		Value   int    `json:"value"`
		APIType string `json:"api_type"`
	} `json:"total_enrolled"`
	SeatsAvailable struct {
		Value   int    `json:"value"`
		APIType string `json:"api_type"`
	} `json:"seats_available"`
	EnvelopeOnly struct {
		Value   string `json:"value"`
		APIType string `json:"api_type"`
	} `json:"envelope_only"`
	WaitlistCount struct {
		Value   int    `json:"value"`
		APIType string `json:"api_type"`
	} `json:"waitlist_count"`
}

// AssignedInstructors .
type AssignedInstructors struct {
	Metadata struct {
		AdvisoryMessages []interface{} `json:"advisory_messages"`
		CollectionSize   int           `json:"collection_size"`
		PageStart        int           `json:"page_start"`
		PageEnd          int           `json:"page_end"`
		PageSize         int           `json:"page_size"`
		DefaultPageSize  int           `json:"default_page_size"`
		MaxPageSize      int           `json:"max_page_size"`
	} `json:"metadata"`
	Values []struct {
		CreditInstitution struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"credit_institution"`
		YearTerm struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
			Desc    string `json:"desc"`
			ExtDesc string `json:"ext_desc"`
		} `json:"year_term"`
		CurriculumID struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"curriculum_id"`
		TitleCode struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"title_code"`
		SectionNumber struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"section_number"`
		InstructorType struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"instructor_type"`
		InstructorName struct {
			Value           string `json:"value"`
			APIType         string `json:"api_type"`
			RelatedResource string `json:"related_resource"`
			Desc            string `json:"desc"`
		} `json:"instructor_name"`
		PersonID struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"person_id"`
		ByuID struct {
			Value           string `json:"value"`
			APIType         string `json:"api_type"`
			RelatedResource string `json:"related_resource"`
			Desc            string `json:"desc"`
		} `json:"byu_id"`
		NetID struct {
			Value           string `json:"value"`
			APIType         string `json:"api_type"`
			RelatedResource string `json:"related_resource"`
			Desc            string `json:"desc"`
		} `json:"net_id"`
		EmailAddress struct {
			Value           string `json:"value"`
			APIType         string `json:"api_type"`
			RelatedResource string `json:"related_resource"`
			Desc            string `json:"desc"`
		} `json:"email_address"`
		UpdatedByID struct {
			Value         string `json:"value"`
			APIType       string `json:"api_type"`
			UpdatedByName string `json:"updated_by_name"`
		} `json:"updated_by_id"`
		DateTimeUpdated struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"date_time_updated"`
	} `json:"values"`
}

//Value .
type Value struct {
	Value   string `json:"value"`
	APIType string `json:"api_type"`
	Key     bool   `json:"key"`
}

// AssignedSchedule .
type AssignedSchedule struct {
	Metadata struct {
		CollectionSize  int `json:"collection_size"`
		PageStart       int `json:"page_start"`
		PageEnd         int `json:"page_end"`
		PageSize        int `json:"page_size"`
		DefaultPageSize int `json:"default_page_size"`
		MaxPageSize     int `json:"max_page_size"`
	} `json:"metadata"`
	Values []struct {
		ScheduleType struct {
		} `json:"schedule_type"`
		ScheduleID struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"schedule_id"`
		SequenceNumber struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Key     bool   `json:"key"`
		} `json:"sequence_number"`
		YearTerm struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
			ExtDesc string `json:"ext_desc"`
		} `json:"year_term"`
		CurriculumID struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"curriculum_id"`
		TitleCode struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"title_code"`
		SectionNumber struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"section_number"`
		Institution struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"institution"`
		Building struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"building"`
		Room struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"room"`
		DaysTaught struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"days_taught"`
		Mon struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"mon"`
		Tue struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"tue"`
		Wed struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"wed"`
		Thu struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"thu"`
		Fri struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"fri"`
		Sat struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"sat"`
		Sun struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"sun"`
		UseStartDate struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"use_start_date"`
		UseEndDate struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"use_end_date"`
		TimeTaught struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"time_taught"`
		BeginTime struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"begin_time"`
		EndTime struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"end_time"`
		Status struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"status"`
		SetUpID struct {
			Value     string      `json:"value"`
			APIType   string      `json:"api_type"`
			SetUpName interface{} `json:"set_up_name"`
		} `json:"set_up_id"`
		SetUpDescription struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"set_up_description"`
		ScheduleReason struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"schedule_reason"`
		DateTimeUpdated struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
		} `json:"date_time_updated"`
		UpdatedByID struct {
			Value         string `json:"value"`
			APIType       string `json:"api_type"`
			UpdatedByName string `json:"updated_by_name"`
		} `json:"updated_by_id"`
		AllowConflict struct {
			Value   string `json:"value"`
			APIType string `json:"api_type"`
			Desc    string `json:"desc"`
		} `json:"allow_conflict"`
	} `json:"values"`
}
