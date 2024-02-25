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

type BrandRepo struct {
	Db *gorm.DB
}

func NewBrand() *BrandRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Brand{})
	return &BrandRepo{Db: db}
}

type BrandResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type BrandRecordInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func (repository *BrandRepo) SaveBrandData(c *gin.Context) {

	var input BrandRecordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b := models.Brand{}

	b.BrandName = input.Name
	b.BrandCode = input.Code

	err := models.CreateBrand(repository.Db, &b)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product registration success", "id": b.ID, "name": b.BrandName})

}

func (repository *BrandRepo) GetBrandById(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	b := models.Brand{}

	if err := models.GetBrandById(repository.Db, &b, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Product not found!"})
		c.Abort()
		return
	}

	var response BrandResponse
	response.ID = b.ID
	response.Name = b.BrandName
	response.Code = b.BrandCode

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

type BrandUpdateInput struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (repository *BrandRepo) UpdateBrand(c *gin.Context) {
	var input BrandUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	b := models.Brand{}

	err := models.GetBrandById(repository.Db, &b, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Name != "" {
		b.BrandName = input.Name
	}

	if input.Code != "" {
		b.BrandCode = input.Code
	}

	err = models.UpdateBrand(repository.Db, &b, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var response BrandResponse
	response.ID = b.ID
	response.Name = b.BrandName
	response.Code = b.BrandCode

	c.JSON(http.StatusOK, response)
}

func (repository *BrandRepo) GetBrandsData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var b models.Brand

	Brands, totalPages, err := b.GetBrandsPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Product not found!"})
		c.Abort()
		return
	}

	var responses []BrandResponse
	for _, Brand := range Brands {
		response := BrandResponse{
			ID:   Brand.ID,
			Name: Brand.BrandName,
			Code: Brand.BrandCode,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       responses,
		"totalPages": totalPages,
	})
}

func (repository *BrandRepo) DeleteBrand(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	b := models.Brand{}

	err := models.DeleteBrand(repository.Db, &b, id)

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
