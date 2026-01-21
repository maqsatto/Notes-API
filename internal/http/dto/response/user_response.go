package response

import "time"

// UserResponse represents a single user in responses (without sensitive data)
type UserResponse struct {
	ID        uint64     `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AuthResponse represents the authentication response with token
type AuthResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
	TotalPages int             `json:"total_pages"`
}

// AvailabilityResponse represents the response for email/username availability checks
type AvailabilityResponse struct {
	Available bool `json:"available"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// StatsResponse represents user statistics
type StatsResponse struct {
	TotalUsers uint64 `json:"total_users"`
}

// Helper Methods


// // CalculateTotalPages calculates the total number of pages
// func (r *UserListResponse) CalculateTotalPages() {
// 	if r.Limit > 0 {
// 		r.TotalPages = int((r.Total + int64(r.Limit) - 1) / int64(r.Limit))
// 	}
// }

// // SetDefaults sets default values for pagination
// func (r *ListUsersRequest) SetDefaults() {
// 	if r.Limit == 0 {
// 		r.Limit = 10
// 	}
// 	if r.Offset < 0 {
// 		r.Offset = 0
// 	}
// }
