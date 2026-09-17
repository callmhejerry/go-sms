package tenant

type CreateTenantRequest struct {
	Name  string             `json:"name"`
	Slug  string             `json:"slug"`
	Owner CreateOwnerRequest `json:"owner"`
}

type CreateOwnerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
