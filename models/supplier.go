package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Supplier struct {
	gorm.Model
	SupplierName     string `gorm:"size:255;not null" json:"supplier_name"`
	SupplierUsername string `gorm:"size:255;not null;unique" json:"supplier_username"`
	SupplierEmail    string `gorm:"size:255;not null;unique" json:"supplier_email"`
	SupplierPhone    string `gorm:"size:255;not null" json:"supplier_phone"`
	SupplierPassword string `gorm:"size:255;not null" json:"supplier_password"`
	SupplierAddress  string `json:"supplier_address"`
}

func CreateSupplier(db *gorm.DB, Supplier *Supplier) (err error) {
	err = db.Create(Supplier).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *Supplier) GetSuppliersPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Supplier, int, error) {
	var customers []Supplier
	var count int64

	// Count total records
	if err := db.Model(&Supplier{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := (int(count) + pageSize - 1) / pageSize

	// Paginate and sort
	if err := db.Order(sortField + " " + sortOrder).Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, totalPages, nil
}

func GetSupplierByEmail(db *gorm.DB, Supplier *Supplier, email string) (err error) {
	err = db.Where("supplier_email = ?", email).First(Supplier).Error
	if err != nil {
		return err
	}
	return nil
}

// get Supplier by id
func GetSupplierById(db *gorm.DB, Supplier *Supplier, id int) (err error) {
	err = db.Where("id = ?", id).First(Supplier).Error
	if err != nil {
		return err
	}
	return nil
}

// update Supplier
func UpdateSupplier(db *gorm.DB, Supplier *Supplier, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Supplier).Error
	if err != nil {
		return err
	}
	return nil
}

// delete Supplier
func DeleteSupplier(db *gorm.DB, Supplier *Supplier, id int) (err error) {
	db.Where("id = ?", id).Delete(Supplier)
	return nil
}

// encrypt pass
func (c *Supplier) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.SupplierPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.SupplierPassword = string(hashedPassword)
	return nil
}

// validate password
func (c *Supplier) VerifyPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.SupplierPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
