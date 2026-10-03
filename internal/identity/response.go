package identity

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

type LoginResponse struct {
	Tokens RefreshTokenResponse `json:"tokens"`
	User   UserResponse         `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
