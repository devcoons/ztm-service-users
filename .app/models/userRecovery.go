package models

import (
	"time"

	"gorm.io/gorm"
)

type UsersRecovery struct {
	Id            int       `gorm:"primaryKey;autoIncrement;" json:"id"`
	UserId        string    `gorm:"unique;not null;size:128;" json:"user_id"`
	RecoveryToken string    `gorm:"size:256;" json:"rec_token"`
	ExpireAt      time.Time `json:"ex_at"`
	CreatedAt     time.Time `json:"cr_at"`
}

func (u *UsersRecovery) BeforeCreate(db *gorm.DB) (err error) {
	u.ExpireAt = time.Now().Add(time.Hour * 2)
	return
}

func (u *UsersRecovery) Create(db *gorm.DB) bool {

	result := db.Create(u)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func UsersRecoveryGetById(db *gorm.DB, id int) *UsersRecovery {
	var r = UsersRecovery{}

	result := db.First(&r, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersRecoveryGetByTokenId(db *gorm.DB, rec_token string) *UsersRecovery {
	var r = UsersRecovery{}

	result := db.First(&r, "rec_token = ?", rec_token)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersRecoveryGetByUserId(db *gorm.DB, user_id int) *UsersRecovery {
	var r = UsersRecovery{}

	result := db.First(&r, "user_id = ?", user_id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersRecoveryDeleteById(db *gorm.DB, id int) bool {

	result := db.Delete(&UsersRecovery{}, id)
	if result.Error != nil {
		return false
	}
	return true
}
