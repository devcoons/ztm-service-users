package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	if db != nil {
		fmt.Println("[MDL] Automigrate models to database")
		fmt.Print("[MDL] - Users ")
		er := db.AutoMigrate(&User{})
		if er == nil {
			fmt.Println("[OK]")
		} else {
			fmt.Println("[ER]")
		}
		fmt.Print("[MDL] - UsersPermissions ")
		er = db.AutoMigrate(&UsersPermissions{})
		if er == nil {
			fmt.Println("[OK]")
		} else {
			fmt.Println("[ER]")
		}
		fmt.Print("[MDL] - UsersRecovery ")
		er = db.AutoMigrate(&UsersRecovery{})
		if er == nil {
			fmt.Println("[OK]")
		} else {
			fmt.Println("[ER]")
		}
	} else {
		fmt.Println("[MDL] Could not migrate models (db missing)")
	}

}

func ResetMigration(db *gorm.DB) {

	if db != nil {
		fmt.Println("[MDL] Reset models migration to database")
		fmt.Print("[MDL] - Delete all tables")
		db.Migrator().DropTable(&User{})
		db.Migrator().DropTable(&UsersPermissions{})
		db.Migrator().DropTable(&UsersRecovery{})
		AutoMigrate(db)
	}
}

type UserJsonOverview struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Role      int       `json:"role"`
	Image     []byte    `json:"image"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	LastLogin time.Time `json:"lastlogin"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type UserJson struct {
	Id          int       `json:"id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Role        int       `json:"role,omitempty"`
	FirstName   string    `json:"firstname,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	Image       []byte    `json:"image,omitempty"`
	Company     string    `json:"company,omitempty"`
	Email       string    `json:"email,omitempty"`
	MobilePhone string    `json:"moblephone,omitempty"`
	LandLine    string    `json:"landline,omitempty"`
	Country     string    `json:"country,omitempty"`
	Province    string    `json:"provice,omitempty"`
	City        string    `json:"city,omitempty"`
	Address     string    `json:"address,omitempty"`
	IsEnabled   bool      `json:"is_enabled,omitempty"`
	LastLogin   time.Time `json:"lastlogin,omitempty"`
	CreatedAt   time.Time `json:"createdat,omitempty"`
	UpdatedAt   time.Time `json:"updatedat,omitempty"`
}
