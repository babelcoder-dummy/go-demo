package dto

type RegisterRequest struct {
	Name     string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required,min=8"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
	TokensResponse
}

type TokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
