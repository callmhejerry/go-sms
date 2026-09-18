package constants

type RoleName string
type Gender string
type StudentStatus string

const (
	Owner      RoleName = "owner"
	Admin      RoleName = "admin"
	Teacher    RoleName = "teacher"
	Accountant RoleName = "accountant"
)

const (
	Male   Gender = "male"
	Female Gender = "female"
)

const (
	Active   StudentStatus = "active"
	Inactive StudentStatus = "inactive"
)
