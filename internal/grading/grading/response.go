package grading

import (
	"time"

	"github.com/shopspring/decimal"
)

type ComputeResultResponse struct {
	Message      string `json:"message"`
	ResultsCount int    `json:"results_count"`
}

type SubjectResultResponse struct {
	SubjectName string          `json:"subject_name"`
	SubjectCode string          `json:"subject_code"`
	TotalScore  decimal.Decimal `json:"total_score"`
	MaxTotal    decimal.Decimal `json:"max_total"`
	Percentage  decimal.Decimal `json:"percentage"`
	Grade       *string         `json:"grade"`
	Remark      *string         `json:"remark"`
}

type ReportCardStudentInfoResponse struct {
	AdmissionNumber string `json:"admission_number"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ClassName       string `json:"class"`
	ClassArmName    string `json:"arm"`
}

type ReportCardResponse struct {
	StudentInfo     ReportCardStudentInfoResponse `json:"student_info"`
	AcademicSession string                        `json:"academic_session"`
	Subjects        []SubjectResultResponse       `json:"subjects"`
	Average         decimal.Decimal               `json:"average"`
	GeneratedAt     time.Time                     `json:"generated_at"`
}
