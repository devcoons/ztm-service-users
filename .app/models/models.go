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
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	Role        int       `json:"role"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Image       []byte    `json:"image"`
	Company     string    `json:"company"`
	Email       string    `json:"email"`
	MobilePhone string    `json:"moblephone"`
	LandLine    string    `json:"landline"`
	Country     string    `json:"country"`
	Province    string    `json:"provice"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	IsEnabled   bool      `json:"is_enabled"`
	LastLogin   time.Time `json:"lastlogin"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}
