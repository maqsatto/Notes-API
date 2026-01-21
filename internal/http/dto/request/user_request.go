package request

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest represents the profile update request payload
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ChangePasswordRequest represents the password change request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ListUsersRequest represents the list users query parameters
type ListUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// CheckEmailRequest represents the email availability check request
type CheckEmailRequest struct {
	Email string `json:"email"`
}

// CheckUsernameRequest represents the username availability check request
type CheckUsernameRequest struct {
	Username string `json:"username"`
}
