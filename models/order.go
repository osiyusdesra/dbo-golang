package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderCustomerId  int     `gorm:"size:255;not null" json:"order_customer_id"`
	OrderSupplierId  int     `gorm:"size:255;not null" json:"order_supplier_id"`
	OrderProductId   int     `gorm:"size:255;not null" json:"order_product_id"`
	OrderQty         int     `gorm:"size:255;not null" json:"order_qty"`
	OrderTotalAmount float64 `gorm:"size:255;not null" json:"order_total_amount"`
	OrderIsPaid      int8    `gorm:"size:255;not null" json:"order_is_paid"`
}

func CreateOrder(db *gorm.DB, Order *Order) (err error) {
	err = db.Create(Order).Error

	if err != nil {
		return err
	}

	return nil
}

func (o *Order) GetOrdersPaginate(db *gorm.DB, page, pageSize int, sortField, sortOrder string) ([]Order, int, error) {
	var orders []Order
	var count int64

	// Count total records
	if err := db.Model(&Order{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate total pages
	totalPages := (int(count) + pageSize - 1) / pageSize

	// Paginate and sort
	if err := db.Order(sortField + " " + sortOrder).Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, totalPages, nil
}

// get Order by id
func GetOrderById(db *gorm.DB, Order *Order, id int) (err error) {
	err = db.Where("id = ?", id).First(Order).Error
	if err != nil {
		return err
	}
	return nil
}

// update Supplier
func UpdateOrder(db *gorm.DB, Order *Order, id int) (err error) {
	err = db.Where("id = ?", id).Updates(Order).Error
	if err != nil {
		return err
	}
	return nil
}

// delete Supplier
func DeleteOrder(db *gorm.DB, Order *Order, id int) (err error) {
	db.Where("id = ?", id).Delete(Order)
	return nil
}
