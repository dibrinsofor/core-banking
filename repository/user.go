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

func (u *UserRepo) UpdateUsersByID(id1 string, id2 string, user1 *models.User, user2 *models.User) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := u.db.Model(models.User{}).Where("account_number = ?", id1).Updates(&user1).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := u.db.Model(models.User{}).Where("account_number = ?", id2).Updates(&user2).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
