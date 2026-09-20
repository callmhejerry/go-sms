package constants

type AdmissionStatus string

const (
	AdmissionPending     AdmissionStatus = "pending"
	AdmissionUnderReview AdmissionStatus = "under_review"
	AdmissionAccepted    AdmissionStatus = "accepted"
	AdmissionRejected    AdmissionStatus = "reject"
)
