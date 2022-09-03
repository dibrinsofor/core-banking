package repository

import (
	"github.com/dibrinsofor/core-banking/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *models.User) error {
	return u.db.Create(&user).Error
}

func (u *UserRepo) UpdateUserByID(id string, user *models.User) error {
	return u.db.Model(models.User{}).Where("account_number = ?", id).Updates(&user).Error
}

func (u *UserRepo) GetUserByID(id string) (*models.User, error) {
	var user models.User
	db := u.db.Where("account_number = ?", id).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	return &user, db.Error
}
