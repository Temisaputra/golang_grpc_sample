package entity

import (
	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
)

type Users struct {
	ID        int32  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"column:username"`
	Email     string `json:"email" gorm:"column:email"`
	Role      string `json:"role" gorm:"column:role"`
	Password  string `json:"password" gorm:"column:password"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt string `json:"deleted_at" gorm:"column:deleted_at"`
}

func (r *Users) TableName() string {
	return "users"
}

func (u *Users) ToPresenter() *presenter.UserResponse {
	return &presenter.UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Username,
		Role:  u.Role,
	}
}
