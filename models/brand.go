package models

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	BrandName string `gorm:"size:255;not null" json:"brand_name"`
	BrandCode string `gorm:"size:255;not null" json:"brand_code"`
}

func CreateBrand(db *gorm.DB, Brand *Brand) (err error) {
	err = db.Create(Brand).Error

	if err != nil {
		return err
	}

	return nil
}

func (b *Brand) GetBrandsPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Brand, int, error) {
	var brands []Brand
	var count int64

	// Count total records
	if err := db.Model(&Brand{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := (int(count) + pageSize - 1) / pageSize

	// Paginate and sort
	if err := db.Order(sortField + " " + sortOrder).Offset((page - 1) * pageSize).Limit(pageSize).Find(&brands).Error; err != nil {
		return nil, 0, err
	}

	return brands, totalPages, nil
}

// get Brand by id
func GetBrandById(db *gorm.DB, Brand *Brand, id int) (err error) {
	err = db.Where("id = ?", id).First(Brand).Error
	if err != nil {
		return err
	}
	return nil
}

// update Supplier
func UpdateBrand(db *gorm.DB, Brand *Brand, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Brand).Error
	if err != nil {
		return err
	}
	return nil
}

// delete Supplier
func DeleteBrand(db *gorm.DB, Brand *Brand, id int) (err error) {
	db.Where("id = ?", id).Delete(Brand)
	return nil
}
