package controllers

import (
	"be-dbo-golang/database"
	"be-dbo-golang/models"
	"be-dbo-golang/utils/auth"
	"be-dbo-golang/utils/pagination"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplierRepo struct {
	Db *gorm.DB
}

func NewSupplier() *SupplierRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Supplier{})
	return &SupplierRepo{Db: db}
}

type SupplierResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type SupplierRegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string
}

func (repository *SupplierRepo) SupplierRegister(c *gin.Context) {

	var input SupplierRegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := models.Supplier{}

	s.SupplierName = input.Name
	s.SupplierUsername = input.Username
	s.SupplierEmail = input.Email
	s.SupplierPhone = input.Phone
	s.SupplierPassword = input.Password
	s.SupplierAddress = input.Address

	if err := s.HashPassword(s.SupplierPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := models.CreateSupplier(repository.Db, &s)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Supplier registration success", "id": s.ID, "email": s.SupplierEmail, "username": s.SupplierUsername})

}

type SupplierLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (repository *SupplierRepo) SupplierLogin(c *gin.Context) {

	var input SupplierLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := models.Supplier{}

	// check if email exists and password is correct
	err := models.GetSupplierByEmail(repository.Db, &s, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	credentialError := s.VerifyPassword(input.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	// Create JWT Token for Authorization
	tokenString, err := auth.GenerateToken(int(s.ID), s.SupplierEmail, s.SupplierUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func (repository *SupplierRepo) SupplierLoggedIn(c *gin.Context) {

	id, _, err := auth.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := models.Supplier{}

	if err := models.GetSupplierById(repository.Db, &s, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Supplier not found!"})
		c.Abort()
		return
	}

	var response SupplierResponse
	response.ID = s.ID
	response.Name = s.SupplierName
	response.Username = s.SupplierUsername
	response.Email = s.SupplierEmail
	response.Phone = s.SupplierPhone
	response.Address = s.SupplierAddress

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

type SupplierUpdateInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (repository *SupplierRepo) UpdateSupplier(c *gin.Context) {
	var input SupplierUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	s := models.Supplier{}

	err := models.GetSupplierById(repository.Db, &s, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Name != "" {
		s.SupplierName = input.Name
	}

	if input.Username != "" {
		s.SupplierUsername = input.Username
	}

	if input.Email != "" {
		s.SupplierEmail = input.Email
	}

	if input.Phone != "" {
		s.SupplierPhone = input.Phone
	}

	if input.Password != "" {
		s.SupplierPassword = input.Password

		if err := s.HashPassword(s.SupplierPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}

	if input.Address != "" {
		s.SupplierAddress = input.Address
	}

	err = models.UpdateSupplier(repository.Db, &s, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var response SupplierResponse
	response.ID = s.ID
	response.Name = s.SupplierName
	response.Username = s.SupplierUsername
	response.Email = s.SupplierEmail
	response.Phone = s.SupplierPhone
	response.Address = s.SupplierAddress

	c.JSON(http.StatusOK, response)
}

func (repository *SupplierRepo) GetSuppliersData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var s models.Supplier

	Suppliers, totalPages, err := s.GetSuppliersPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Supplier not found!"})
		c.Abort()
		return
	}

	var responses []SupplierResponse
	for _, Supplier := range Suppliers {
		response := SupplierResponse{
			ID:       Supplier.ID,
			Name:     Supplier.SupplierName,
			Username: Supplier.SupplierUsername,
			Email:    Supplier.SupplierEmail,
			Phone:    Supplier.SupplierPhone,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       responses,
		"totalPages": totalPages,
	})
}

func (repository *SupplierRepo) DeleteSupplier(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s := models.Supplier{}

	err := models.DeleteSupplier(repository.Db, &s, id)

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
