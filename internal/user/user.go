package user

import (
	"fmt"
	"regexp"
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         string    `json:"role" gorm:"size:20;default:user;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DTOs for requests/responses
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Validation methods
func (u User) Validate() error {
	if err := u.ValidateUsername(); err != nil {
		return err
	}
	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidateRole(); err != nil {
		return err
	}
	return nil
}

func (u User) ValidateUsername() error {
	if len(u.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(u.Username) > 50 {
		return fmt.Errorf("username must be at most 50 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(u.Username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func (u User) ValidateEmail() error {
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(u.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (u User) ValidateRole() error {
	validRoles := []string{"user", "admin", "moderator"}
	for _, role := range validRoles {
		if u.Role == role {
			return nil
		}
	}
	return fmt.Errorf("invalid role: %s", u.Role)
}

// Business logic methods
func (u User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u User) IsModerator() bool {
	return u.Role == "moderator" || u.Role == "admin"
}

func (u User) CanModifyUser(targetUser User) bool {
	return u.IsAdmin() || u.ID == targetUser.ID
}
