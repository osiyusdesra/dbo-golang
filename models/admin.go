package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	AdminName     string `gorm:"size:255;not null" json:"admin_name"`
	AdminUsername string `gorm:"size:255;not null;unique" json:"admin_username"`
	AdminEmail    string `gorm:"size:255;not null;unique" json:"admin_email"`
	AdminPassword string `gorm:"size:255;not null" json:"admin_password"`
	AdminPhone    string `gorm:"size:255;not null" json:"admin_phone"`
}

// Record admin data
func CreateAdmin(db *gorm.DB, Admin *Admin) (err error) {
	err = db.Create(Admin).Error

	if err != nil {
		return err
	}

	return nil
}

// Get Admin List
func (a *Admin) GetAdminsPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Admin, int, error) {
	var admins []Admin
	var count int64

	// Count total records
	if err := db.Model(&Admin{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := (int(count) + pageSize - 1) / pageSize

	// Paginate and sort
	if err := db.Order(sortField + " " + sortOrder).Offset((page - 1) * pageSize).Limit(pageSize).Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, totalPages, nil
}

// get admin by email
func GetAdminByEmail(db *gorm.DB, Admin *Admin, email string) (err error) {
	err = db.Where("admin_email = ?", email).First(Admin).Error
	if err != nil {
		return err
	}
	return nil
}

// get admin by id
func GetAdminById(db *gorm.DB, Admin *Admin, id int) (err error) {
	err = db.Where("id = ?", id).First(Admin).Error
	if err != nil {
		return err
	}
	return nil
}

// update admin
func UpdateAdmin(db *gorm.DB, Admin *Admin, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Admin).Error
	if err != nil {
		return err
	}
	return nil
}

// delete admin
func DeleteAdmin(db *gorm.DB, Admin *Admin, id int) (err error) {
	db.Where("id = ?", id).Delete(Admin)
	return nil
}

// encrypt pass
func (a *Admin) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.AdminPassword = string(hashedPassword)
	return nil
}

// validate password
func (a *Admin) VerifyPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.AdminPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
