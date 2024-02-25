package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CustomerName     string `gorm:"size:255;not null" json:"customer_name"`
	CustomerUsername string `gorm:"size:255;not null;unique" json:"customer_username"`
	CustomerEmail    string `gorm:"size:255;not null;unique" json:"customer_email"`
	CustomerPhone    string `gorm:"size:255;not null" json:"customer_phone"`
	CustomerPassword string `gorm:"size:255;not null" json:"customer_password"`
	CustomerAddress  string `json:"customer_address"`
}

func CreateCustomer(db *gorm.DB, Customer *Customer) (err error) {
	err = db.Create(Customer).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *Customer) GetCustomersPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Customer, int, error) {
	var customers []Customer
	var count int64

	// Count total records
	if err := db.Model(&Customer{}).Count(&count).Error; err != nil {
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

func GetCustomerByEmail(db *gorm.DB, Customer *Customer, email string) (err error) {
	err = db.Where("customer_email = ?", email).First(Customer).Error
	if err != nil {
		return err
	}
	return nil
}

// get Customer by id
func GetCustomerById(db *gorm.DB, Customer *Customer, id int) (err error) {
	err = db.Where("id = ?", id).First(Customer).Error
	if err != nil {
		return err
	}
	return nil
}

// update Customer
func UpdateCustomer(db *gorm.DB, Customer *Customer, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Customer).Error
	if err != nil {
		return err
	}
	return nil
}

// delete Customer
func DeleteCustomer(db *gorm.DB, Customer *Customer, id int) (err error) {
	db.Where("id = ?", id).Delete(Customer)
	return nil
}

// encrypt pass
func (c *Customer) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.CustomerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.CustomerPassword = string(hashedPassword)
	return nil
}

// validate password
func (c *Customer) VerifyPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.CustomerPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
