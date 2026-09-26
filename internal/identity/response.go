package identity

type UserResponse struct {
	ID        string   `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	Roles     []string `json:"roles"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
