package grading

type CreateAssessmentTypeRequest struct {
	Name     string   `json:"name" validate:"required"`
	MaxScore float64  `json:"max_score" validate:"required,min=1"`
	Weight   *float64 `json:"weight"`
	IsExam   bool     `json:"is_exam"`
}
