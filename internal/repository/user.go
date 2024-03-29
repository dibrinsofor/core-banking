package repository

import (
	"time"

	"github.com/dibrinsofor/core-banking/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *models.Users) error {
	return u.db.Create(&user).Error
}

func (u *UserRepo) UpdateUserByID(id string, user *models.Users, action string, recipient string) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user)

	if err := u.db.Model(models.Users{}).Where("account_number = ?", id).Updates(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := u.db.Create(InsertTransaction(user, action, recipient)).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (u *UserRepo) GetUserByID(id string) (*models.Users, error) {
	var user models.Users
	db := u.db.Where("account_number = ?", id).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	return &user, db.Error
}

func (u *UserRepo) GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	db := u.db.Where("email = ?", email).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	return &user, db.Error
}

func (u *UserRepo) UpdateUsersByID(id1 string, id2 string, user1 *models.Users, user2 *models.Users, action string) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user1).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user2).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := u.db.Model(models.Users{}).Where("account_number = ?", id1).Updates(&user1).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := u.db.Model(models.Users{}).Where("account_number = ?", id2).Updates(&user2).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := u.db.Create(InsertTransaction(user1, action, id2)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func InsertTransaction(user *models.Users, action string, recipient string) *models.Transaction {
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

func (u *UserRepo) GetUserBalance(AccountNumber string) (balance float64, err error) {
	var user models.Users

	if AccountNumber == "" {
		return 0.00, nil
	}

	db := u.db.Where("account_number = ?", AccountNumber).Find(&user)
	if db.Error != nil {
		return 0.00, db.Error
	}
	return user.Balance, db.Error
}

func (u *UserRepo) Deposit(UserObj *models.Users, DepositAmount float64) (newBalance float64, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	OldBalance, err := u.GetUserBalance(UserObj.AccountNumber)
	if err != nil {
		return 0, err
	}

	if err := tx.Error; err != nil {
		return 0, err
	}

	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&UserObj).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := u.db.Model(models.Users{}).Where("account_number = ?", UserObj.AccountNumber).Update("balance", (OldBalance + DepositAmount)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := u.db.Create(InsertTransaction(UserObj, "DEPOSIT", "")).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	NewBalance, err := u.GetUserBalance(UserObj.AccountNumber)
	if err != nil {
		return 0, err
	}

	return NewBalance, tx.Commit().Error

}

func (u *UserRepo) Withdraw(UserObj *models.Users, WithdrawAmount float64) (newBalance float64, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	OldBalance, err := u.GetUserBalance(UserObj.AccountNumber)
	if err != nil {
		return 0, err
	}

	if err := tx.Error; err != nil {
		return 0, err
	}

	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&UserObj).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := u.db.Model(models.Users{}).Where("account_number = ?", UserObj.AccountNumber).Update("balance", (OldBalance - WithdrawAmount)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := u.db.Create(InsertTransaction(UserObj, "WITHDRAW", "")).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// alternative would be Gorm clauses
	NewBalance, err := u.GetUserBalance(UserObj.AccountNumber)
	if err != nil {
		return 0, err
	}

	return NewBalance, tx.Commit().Error

}

func (u *UserRepo) Transfer(user1 *models.Users, user2 *models.Users, transferAmount float64) (newBalance float64, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 0, err
	}

	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user1).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := u.db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user2).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// additional check to see if the user still has the money to perform trans
	if user1.Balance > transferAmount {
		// perform transactions
		if err := u.db.Model(models.Users{}).Where("account_number = ?", user1.AccountNumber).Update("balance", (user1.Balance - transferAmount)).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		if err := u.db.Model(models.Users{}).Where("account_number = ?", user2.AccountNumber).Update("balance", (user2.Balance + transferAmount)).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := u.db.Create(InsertTransaction(user1, "TRANSFER", user2.AccountNumber)).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	NewBalance, err := u.GetUserBalance(user1.AccountNumber)
	if err != nil {
		return 0, err
	}

	return NewBalance, tx.Commit().Error
}
