package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductSupplierId int     `gorm:"size:255;not null" json:"product_supplier_id"`
	ProductName       string  `gorm:"size:255;not null" json:"product_name"`
	ProductBrandId    int     `gorm:"size:255;not null" json:"product_brand_id"`
	ProductStock      int     `gorm:"size:255;not null" json:"product_stock"`
	ProductPrice      float64 `gorm:"size:255;not null" json:"product_price"`
}

func CreateProduct(db *gorm.DB, Product *Product) (err error) {
	err = db.Create(Product).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) GetProductsPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Product, int, error) {
	var products []Product
	var count int64

	// Count total records
	if err := db.Model(&Product{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := (int(count) + pageSize - 1) / pageSize

	// Paginate and sort
	if err := db.Order(sortField + " " + sortOrder).Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalPages, nil
}

// get Product by id
func GetProductById(db *gorm.DB, Product *Product, id int) (err error) {
	err = db.Where("id = ?", id).First(Product).Error
	if err != nil {
		return err
	}
	return nil
}

// update Supplier
func UpdateProduct(db *gorm.DB, Product *Product, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Product).Error
	if err != nil {
		return err
	}
	return nil
}

// delete Supplier
func DeleteProduct(db *gorm.DB, Product *Product, id int) (err error) {
	db.Where("id = ?", id).Delete(Product)
	return nil
}
