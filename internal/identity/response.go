package identity

type UserResponse struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
