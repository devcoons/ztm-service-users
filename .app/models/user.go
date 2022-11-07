package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/hints"
)

type User struct {
	Id          int    `gorm:"primaryKey;autoIncrement;"`
	Username    string `gorm:"unique;not null;size:128;"`
	Password    string `gorm:"size:256;"`
	Role        int    `gorm:"size:2;"`
	Nonce       string `gorm:"size:6;"`
	FirstName   string `gorm:"size:64;"`
	LastName    string `gorm:"size:64;"`
	Image       []byte
	Company     string `gorm:"size:128;"`
	Email       string `gorm:"size:128;"`
	MobilePhone string `gorm:"size:128;"`
	LandLine    string `gorm:"size:128;"`
	Country     string `gorm:"size:64;"`
	Province    string `gorm:"size:64;"`
	City        string `gorm:"size:64;"`
	Address     string `gorm:"size:128;"`
	IsEnabled   bool   `gorm:"default:true"`
	LastLogin   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	if !strings.HasPrefix(u.Password, "::") {
		h := sha1.New()
		h.Write([]byte(u.Password))
		u.Password = "::" + hex.EncodeToString(h.Sum(nil))
	}
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Id == 0 {
		return errors.New("id = 0 or null not valid")
	}
	if !strings.HasPrefix(u.Password, "::") {
		h := sha1.New()
		h.Write([]byte(u.Password))
		u.Password = "::" + hex.EncodeToString(h.Sum(nil))
	}
	return
}

func (u *User) Create(db *gorm.DB) bool {
	if !strings.HasPrefix(u.Password, "::") {
		h := sha1.New()
		h.Write([]byte(u.Password))
		u.Password = "::" + hex.EncodeToString(h.Sum(nil))
	}
	result := db.Create(u)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func (u *User) UpdateMapped(db *gorm.DB, mapped map[string]interface{}) bool {
	result := db.Model(u).Updates(mapped)
	return result.Error == nil
}

func UsersGetOneByUsernamePassword(db *gorm.DB, username string, password string) *User {
	var r = User{}

	if !strings.HasPrefix(password, "::") {
		h := sha1.New()
		h.Write([]byte(password))
		password = "::" + hex.EncodeToString(h.Sum(nil))
	}

	result := db.First(&r, "username = ? AND password = ?", username, password)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersGetById(db *gorm.DB, id int) *User {
	var r = User{}

	result := db.Clauses(hints.UseIndex("id")).First(&r, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersGetAll(db *gorm.DB) *[]User {
	var r = []User{}

	result := db.Find(&r)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersDeleteById(db *gorm.DB, id int) bool {

	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return false
	}
	return true
}

func UsersGetAllWithFields(db *gorm.DB, fields []string) *[]User {
	var r = []User{}

	result := db.Select(fields).Find(&r)
	if result.Error != nil {
		return nil
	}
	return &r
}
