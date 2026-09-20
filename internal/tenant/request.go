package tenant

type CreateTenantRequest struct {
	Name  string             `json:"name" validate:"required"`
	Slug  string             `json:"slug" validate:"required"`
	Owner CreateOwnerRequest `json:"owner"`
}

type CreateOwnerRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=3"`
	LastName  string `json:"last_name" validate:"required,min=3"`
}
