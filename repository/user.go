package repository

import (
	"time"

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

func (u *UserRepo) UpdateUserByID(id string, user *models.User, action string, recipient string) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := u.db.Model(models.User{}).Where("account_number = ?", id).Updates(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := u.db.Create(InsertTransaction(user, action, recipient)).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (u *UserRepo) GetUserByID(id string) (*models.User, error) {
	var user models.User
	db := u.db.Where("account_number = ?", id).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	return &user, db.Error
}

func (u *UserRepo) UpdateUsersByID(id1 string, id2 string, user1 *models.User, user2 *models.User, action string) error {
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

	if err := u.db.Create(InsertTransaction(user1, action, id2)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func InsertTransaction(user *models.User, action string, recipient string) *models.Transaction {
	t := time.Now()
	t1 := t.Format("2006-01-02")
	trans := models.Transaction{
		AccountNumber:   user.AccountNumber,
		ActionPerformed: action,
		Recipient:       recipient,
		Balance:         user.Balance,
		CreatedDate:     t1,
	}
	return &trans
}

func (u *UserRepo) GetAllTransactions(AccountNumber string) (*[]models.Transaction, error) {
	var transaction []models.Transaction

	db := u.db.Limit(10).Where("account_number = ?", AccountNumber).Find(&transaction)
	if db.Error != nil {
		return nil, db.Error
	}
	return &transaction, db.Error
}

func (u *UserRepo) QueryTransactions(AccountNumber string, Action string, Date string) (*[]models.Transaction, error) {
	var transaction []models.Transaction

	db := u.db.Limit(10).Where(map[string]interface{}{"account_number": AccountNumber, "action_performed": Action, "created_date": Date}).Find(&transaction)
	if db.Error != nil {
		return nil, db.Error
	}
	return &transaction, db.Error
}
