package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type CSRFResponse struct {
	CSRFToken string `json:"csrf_token"`
}
