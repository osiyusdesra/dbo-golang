package controllers

import (
	"be-dbo-golang/database"
	"be-dbo-golang/models"
	"be-dbo-golang/utils/pagination"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRepo struct {
	Db *gorm.DB
}

func NewOrder() *OrderRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Order{})
	return &OrderRepo{Db: db}
}

type OrderRecordInput struct {
	SupplierId int     `json:"supplier_id" binding:"required"`
	CustomerId int     `json:"customer_id" binding:"required"`
	ProductId  int     `json:"product_id" binding:"required"`
	Qty        int     `json:"quantity" binding:"required"`
	Amount     float64 `json:"total_amount" binding:"required"`
}

func (repository *OrderRepo) SaveOrderData(c *gin.Context) {

	var input OrderRecordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o := models.Order{}

	o.OrderSupplierId = input.SupplierId
	o.OrderCustomerId = input.CustomerId
	o.OrderProductId = input.ProductId
	o.OrderQty = input.Qty
	o.OrderTotalAmount = input.Amount

	err := models.CreateOrder(repository.Db, &o)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order save successfully", "id": o.ID, "total_amount": o.OrderTotalAmount})

}

func (repository *OrderRepo) GetOrderById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	o := models.Order{}

	if err := models.GetOrderById(repository.Db, &o, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Order not found!"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
}

type OrderUpdateInput struct {
	SupplierId int     `json:"supplier_id"`
	CustomerId int     `json:"customer_id"`
	ProductId  int     `json:"product_id"`
	Qty        int     `json:"quantity"`
	Amount     float64 `json:"total_amount"`
}

func (repository *OrderRepo) UpdateOrder(c *gin.Context) {
	var input OrderUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	o := models.Order{}

	err := models.GetOrderById(repository.Db, &o, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.SupplierId > 0 {
		o.OrderSupplierId = input.SupplierId
	}

	if input.CustomerId > 0 {
		o.OrderCustomerId = input.CustomerId
	}

	if input.ProductId > 0 {
		o.OrderProductId = input.ProductId
	}

	if input.Qty > 0 {
		o.OrderQty = input.Qty
	}

	if input.Amount > 0 {
		o.OrderTotalAmount = input.Amount
	}

	err = models.UpdateOrder(repository.Db, &o, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (repository *OrderRepo) GetOrdersData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var s models.Order

	Orders, totalPages, err := s.GetOrdersPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Order not found!"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       Orders,
		"totalPages": totalPages,
	})
}

func (repository *OrderRepo) DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s := models.Order{}

	err := models.DeleteOrder(repository.Db, &s, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "data deleted"})
}
