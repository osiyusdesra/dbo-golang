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

type ProductRepo struct {
	Db *gorm.DB
}

func NewProduct() *ProductRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Product{})
	return &ProductRepo{Db: db}
}

type ProductResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	BrandId    int     `json:"brand_id"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	SupplierId int     `json:"supplier_id"`
}

type ProductRecordInput struct {
	Name       string  `json:"name" binding:"required"`
	BrandId    int     `json:"brand_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Stock      int     `json:"stock" binding:"required"`
	SupplierId int     `json:"supplier_id" binding:"required"`
}

func (repository *ProductRepo) SaveProductData(c *gin.Context) {

	var input ProductRecordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := models.Product{}

	p.ProductName = input.Name
	p.ProductBrandId = input.BrandId
	p.ProductStock = input.Stock
	p.ProductPrice = input.Price
	p.ProductSupplierId = input.SupplierId

	err := models.CreateProduct(repository.Db, &p)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product registration success", "id": p.ID, "name": p.ProductName})

}

func (repository *ProductRepo) GetProductById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	p := models.Product{}

	if err := models.GetProductById(repository.Db, &p, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Product not found!"})
		c.Abort()
		return
	}

	var response ProductResponse
	response.ID = p.ID
	response.Name = p.ProductName
	response.BrandId = p.ProductBrandId
	response.Stock = p.ProductStock
	response.Price = p.ProductPrice
	response.SupplierId = p.ProductStock

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

type ProductUpdateInput struct {
	Name       string  `json:"name"`
	BrandId    int     `json:"brand_id"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	SupplierId int     `json:"supplier_id"`
}

func (repository *ProductRepo) UpdateProduct(c *gin.Context) {
	var input ProductUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	p := models.Product{}

	err := models.GetProductById(repository.Db, &p, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Name != "" {
		p.ProductName = input.Name
	}

	if input.Price >= 0 {
		p.ProductPrice = input.Price
	}

	if input.Stock >= 0 {
		p.ProductStock = input.Stock
	}

	if input.BrandId > 0 {
		p.ProductBrandId = input.BrandId
	}

	if input.SupplierId > 0 {
		p.ProductSupplierId = input.SupplierId
	}

	err = models.UpdateProduct(repository.Db, &p, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var response ProductResponse
	response.ID = p.ID
	response.Name = p.ProductName
	response.BrandId = p.ProductBrandId
	response.Stock = p.ProductStock
	response.Price = p.ProductPrice
	response.SupplierId = p.ProductStock

	c.JSON(http.StatusOK, response)
}

func (repository *ProductRepo) GetProductsData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var s models.Product

	Products, totalPages, err := s.GetProductsPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Product not found!"})
		c.Abort()
		return
	}

	var responses []ProductResponse
	for _, Product := range Products {
		response := ProductResponse{
			ID:         Product.ID,
			Name:       Product.ProductName,
			BrandId:    Product.ProductBrandId,
			Stock:      Product.ProductStock,
			Price:      Product.ProductPrice,
			SupplierId: Product.ProductStock,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       responses,
		"totalPages": totalPages,
	})
}

func (repository *ProductRepo) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s := models.Product{}

	err := models.DeleteProduct(repository.Db, &s, id)

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
