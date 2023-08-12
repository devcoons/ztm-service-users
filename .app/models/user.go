package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id          int       `gorm:"primaryKey;autoIncrement;" json:"id"`
	Username    string    `gorm:"unique;not null;size:128;" json:"username"`
	Password    string    `gorm:"size:256;" json:"password"`
	Role        int       `gorm:"size:2;" json:"role"`
	Nonce       string    `gorm:"size:6;" json:"nonce"`
	FirstName   string    `gorm:"size:64;" json:"fname"`
	LastName    string    `gorm:"size:64;" json:"lname"`
	Image       []byte    `json:"img"`
	Company     string    `gorm:"size:128;" json:"company"`
	Email       string    `gorm:"size:128;" json:"email"`
	MobilePhone string    `gorm:"size:128;" json:"cellphone"`
	LandLine    string    `gorm:"size:128;" json:"landline"`
	Country     string    `gorm:"size:64;" json:"country"`
	Province    string    `gorm:"size:64;" json:"province"`
	City        string    `gorm:"size:64;" json:"city"`
	Address     string    `gorm:"size:128;" json:"address"`
	IsEnabled   bool      `gorm:"default:true" json:"enabled"`
	IsAdmin     bool      `gorm:"default:false" json:"admin"`
	LastLogin   time.Time `json:"ll_at"`
	CreatedAt   time.Time `json:"cr_at"`
	UpdatedAt   time.Time `json:"up_at"`
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
		fmt.Println("?")
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

	if mapped["password"] != nil {
		if !strings.HasPrefix(mapped["password"].(string), "::") {
			h := sha1.New()
			h.Write([]byte(mapped["password"].(string)))
			mapped["password"] = "::" + hex.EncodeToString(h.Sum(nil))
		}
	}

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

	result := db.First(&r, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersGetByUsername(db *gorm.DB, uname string) *User {
	var r = User{}

	result := db.First(&r, "username = ?", uname)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersGetIdByUsername(db *gorm.DB, uname string) int {
	var r = User{}

	result := db.First(&r, "username = ?", uname)
	if result.Error != nil {
		return -1
	}
	return r.Id
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
