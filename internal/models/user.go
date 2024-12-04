package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:32;not null;unique" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // json:"-" 确保密码不会被序列化
	Email     string         `gorm:"size:128;not null;unique" json:"email"`
	FullName  string         `gorm:"size:128;not null" json:"full_name"`
	Bio       string         `gorm:"size:256" json:"bio"`
	Role      string         `gorm:"size:20;not null;default:'user'" json:"role"`     // admin, user
	Status    string         `gorm:"size:20;not null;default:'active'" json:"status"` // active, inactive, suspended
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// UserSettings 用户设置模型
type UserSettings struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;unique" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Theme         string         `gorm:"size:20;default:'light'" json:"theme"` // light, dark
	Language      string         `gorm:"size:10;default:'en'" json:"language"`
	Notifications bool           `gorm:"default:true" json:"notifications"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate GORM hook for handling password hashing before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate GORM hook for handling password hashing before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (User) TableName() string {
	return "users"
}

func (UserSettings) TableName() string {
	return "user_settings"
}

func (u *User) Validate() error {
	// TODO: 添加验证逻辑
	return nil
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=1,max=128"`
	Bio      string `json:"bio" binding:"max=256"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	FullName string `json:"full_name" binding:"omitempty,min=1,max=128"`
	Bio      string `json:"bio" binding:"max=256"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

type UpdateUserSettingsRequest struct {
	Theme         string `json:"theme" binding:"omitempty,oneof=light dark"`
	Language      string `json:"language" binding:"omitempty,len=2"`
	Notifications bool   `json:"notifications"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Bio       string     `json:"bio"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Bio:       u.Bio,
		Role:      u.Role,
		Status:    u.Status,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
