package kafkaStruct

type CommonResp struct {
	Database  string            `json:"database"`
	Es        int64             `json:"es"`
	ID        int               `json:"id"`
	IsDdl     bool              `json:"isDdl"`
	MysqlType map[string]string `json:"mysqlType"`
	PkNames   []string          `json:"pkNames"`
	Sql       string            `json:"sql"`
	SqlType   map[string]int    `json:"sqlType"`
	Table     string            `json:"table"`
	Ts        int64             `json:"ts"`
	Type      string            `json:"type"`
}
type OrganizationSchool struct {
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	DeletedAt         string `json:"deleted_at,omitempty"`
	SchoolID          string `json:"school_id"`
	Name              string `json:"name"`
	ExternalID        string `json:"external_id"`
	OrgCode           string `json:"org_code"`
	Address           string `json:"address"`
	AreaName          string `json:"area_name"`
	AreaID            string `json:"area_id"`
	StreetName        string `json:"street_name"`
	ProvinceName      string `json:"province_name"`
	ProvinceID        string `json:"province_id"`
	CityName          string `json:"city_name"`
	CityID            string `json:"city_id"`
	OrganizedTypeCode string `json:"organized_type_code"`
	TeachingTypeCode  string `json:"teaching_type_code"`
}

type SchoolEvent struct {
	Data []OrganizationSchool `json:"data"`
	Old  []OrganizationSchool `json:"old"`
	CommonResp
}

type SchoolClass struct {
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	DeletedAt        string `json:"deleted_at,omitempty"`
	ClassID          string `json:"class_id"`
	Name             string `json:"name"`
	YearOfEnrollment string `json:"year_of_enrollment"`
	AcademyID        string `json:"academy_id"`
	SpecialtyID      string `json:"specialty_id"`
	GradeID          string `json:"grade_id"`
	SemesterID       string `json:"semester_id"`
	Graduated        string `json:"graduated"`
	SchoolID         string `json:"school_id"`
	ClassType        string `json:"class_type"`
	Members          string `json:"members"`
	TeacherID        string `json:"teacher_id"`
	ExternalID       string `json:"external_id"`
}
type ClassEvent struct {
	Data []SchoolClass `json:"data"`
	Old  []SchoolClass `json:"old"`
	CommonResp
}
type Person struct {
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	PersonId     string `json:"person_id"`
	Username     string `json:"username"`
	SexCode      string `json:"sex_code"`
	Phone        string `json:"phone"`
	CardTypeCode string `json:"card_type_code"`
	CardNumber   string `json:"card_number"`
	ExternalId   string `json:"external_id"`
	Status       string `json:"status"`
}
type PersonEvent struct {
	Data []Person `json:"data"`
	Old  []Person `json:"old"`
	CommonResp
}

type SchoolClassStudent struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	SchoolId  string `json:"school_id"`
	ClassId   string `json:"class_id"`
	StudentId string `json:"student_id"`
	Status    int64  `json:"status omitempty"`
	RecordId  string `json:"record_id"`
}

type SchoolClassStudentEvent struct {
	Data []SchoolClassStudent `json:"data"`
	Old  []SchoolClassStudent `json:"old"`
	CommonResp
}

type SchoolClassTeacher struct {
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	SchoolId    string `json:"school_id"`
	TeacherId   string `json:"teacher_id"`
	SubjectCode string `json:"subject_code"`
	ClassId     string `json:"class_id"`
	Status      int64  `json:"status omitempty"`
	RecordId    string `json:"record_id"`
}

type SchoolClassTeacherEvent struct {
	Data []SchoolClassTeacher `json:"data"`
	Old  []SchoolClassTeacher `json:"old"`
	CommonResp
}
