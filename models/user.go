package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // omit from JSON output
	Role     string `json:"role"`               // e.g., "admin", "user"
}
