package grading

type ComputeResultResponse struct {
	Message      string `json:"message"`
	ResultsCount int    `json:"results_count"`
}
