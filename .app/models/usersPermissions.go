package models

import (
	"time"

	"gorm.io/gorm"
)

type UsersPermissions struct {
	Id        int       `gorm:"primaryKey;autoIncrement;" json:"id"`
	UserId    int       `gorm:"unique;not null;column:user_id;" json:"user_id"`
	Perm01    int       `gorm:"size:1;default:0;column:p_01;default:0;" json:"p_01"`
	Perm02    int       `gorm:"size:1;default:0;column:p_02;default:0;" json:"p_02"`
	Perm03    int       `gorm:"size:1;default:0;column:p_03;default:0;" json:"p_03"`
	Perm04    int       `gorm:"size:1;default:0;column:p_04;default:0;" json:"p_04"`
	Perm05    int       `gorm:"size:1;default:0;column:p_05;default:0;" json:"p_05"`
	Perm06    int       `gorm:"size:1;default:0;column:p_06;default:0;" json:"p_06"`
	Perm07    int       `gorm:"size:1;default:0;column:p_07;default:0;" json:"p_07"`
	Perm08    int       `gorm:"size:1;default:0;column:p_08;default:0;" json:"p_08"`
	Perm09    int       `gorm:"size:1;default:0;column:p_09;default:0;" json:"p_09"`
	Perm10    int       `gorm:"size:1;default:0;column:p_10;default:0;" json:"p_10"`
	Perm11    int       `gorm:"size:1;default:0;column:p_11;default:0;" json:"p_11"`
	Perm12    int       `gorm:"size:1;default:0;column:p_12;default:0;" json:"p_12"`
	Perm13    int       `gorm:"size:1;default:0;column:p_13;default:0;" json:"p_13"`
	Perm14    int       `gorm:"size:1;default:0;column:p_14;default:0;" json:"p_14"`
	Perm15    int       `gorm:"size:1;default:0;column:p_15;default:0;" json:"p_15"`
	Perm16    int       `gorm:"size:1;default:0;column:p_16;default:0;" json:"p_16"`
	CreatedAt time.Time `gorm:"size:255;" json:"cr_at"`
	UpdatedAt time.Time `gorm:"size:255;" json:"up_at"`
}

func (u *UsersPermissions) BeforeCreate(db *gorm.DB) (err error) {

	return
}

func (u *UsersPermissions) BeforeUpdate(tx *gorm.DB) (err error) {

	return
}

func (u *UsersPermissions) Create(db *gorm.DB) bool {

	result := db.Create(u)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func (u *UsersPermissions) UpdateMapped(db *gorm.DB, mapped map[string]interface{}) bool {
	result := db.Model(u).Updates(mapped)
	return result.Error == nil
}

func UsersPermissionsGetById(db *gorm.DB, id int) *UsersPermissions {
	var r = UsersPermissions{}

	result := db.First(&r, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersPermissionsGetByUserId(db *gorm.DB, user_id int) *UsersPermissions {
	var r = UsersPermissions{}

	result := db.First(&r, "user_id = ?", user_id)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersPermissionsGetAll(db *gorm.DB) *[]UsersPermissions {
	var r = []UsersPermissions{}

	result := db.Find(&r)
	if result.Error != nil {
		return nil
	}
	return &r
}

func UsersPermissionsDeleteByUserId(db *gorm.DB, id int) bool {

	result := db.Delete(&UsersPermissions{}, "user_id = ?", id)
	if result.Error != nil {
		return false
	}
	return true
}
