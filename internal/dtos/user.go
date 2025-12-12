package dtos

type RegistrationRequest struct {
	Name string `json:"name"`
	Password string `json:"password"`
	Email string `json:"email"`
	AccountType *string `json:"account_type,omitempty"`
	JoinCode *string `json:"join_code,omitempty"`
	NewSchoolName *string `json:"new_school_name,omitempty"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`	
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in_seconds"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in_seconds"`
}