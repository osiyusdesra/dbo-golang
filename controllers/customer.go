package controllers

import (
	"be-dbo-golang/database"
	"be-dbo-golang/models"
	"be-dbo-golang/utils/auth"
	"be-dbo-golang/utils/pagination"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	Db *gorm.DB
}

func NewCustomer() *CustomerRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Customer{})
	return &CustomerRepo{Db: db}
}

type CustomerResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type CustomerRegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string
}

func (repository *CustomerRepo) CustomerRegister(c *gin.Context) {

	var input CustomerRegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	u.CustomerName = input.Name
	u.CustomerUsername = input.Username
	u.CustomerEmail = input.Email
	u.CustomerPhone = input.Phone
	u.CustomerPassword = input.Password
	u.CustomerAddress = input.Address

	if err := u.HashPassword(u.CustomerPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := models.CreateCustomer(repository.Db, &u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Customer registration success", "CustomerId": u.ID, "email": u.CustomerEmail, "username": u.CustomerUsername})

}

type CustomerLoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (repository *CustomerRepo) CustomerLogin(c *gin.Context) {

	var input CustomerLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	// check if email exists and password is correct
	err := models.GetCustomerByEmail(repository.Db, &u, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	credentialError := u.VerifyPassword(input.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	// Create JWT Token for Authorization
	tokenString, err := auth.GenerateToken(int(u.ID), u.CustomerEmail, u.CustomerUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func (repository *CustomerRepo) CustomerLoggedIn(c *gin.Context) {

	id, _, err := auth.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	if err := models.GetCustomerById(repository.Db, &u, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Customer not found!"})
		c.Abort()
		return
	}

	var response CustomerResponse
	response.ID = u.ID
	response.Name = u.CustomerName
	response.Username = u.CustomerUsername
	response.Email = u.CustomerEmail
	response.Phone = u.CustomerPhone
	response.Address = u.CustomerAddress

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

type CustomerUpdateInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (repository *CustomerRepo) UpdateCustomer(c *gin.Context) {
	var input CustomerUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	u := models.Customer{}

	err := models.GetCustomerById(repository.Db, &u, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Name != "" {
		u.CustomerName = input.Name
	}

	if input.Username != "" {
		u.CustomerUsername = input.Username
	}

	if input.Email != "" {
		u.CustomerEmail = input.Email
	}

	if input.Phone != "" {
		u.CustomerPhone = input.Phone
	}

	if input.Password != "" {
		u.CustomerPassword = input.Password

		if err := u.HashPassword(u.CustomerPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}

	if input.Address != "" {
		u.CustomerAddress = input.Address
	}

	err = models.UpdateCustomer(repository.Db, &u, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var response CustomerResponse
	response.ID = u.ID
	response.Name = u.CustomerName
	response.Username = u.CustomerUsername
	response.Email = u.CustomerEmail
	response.Phone = u.CustomerPhone
	response.Address = u.CustomerAddress

	c.JSON(http.StatusOK, response)
}

func (repository *CustomerRepo) GetCustomersData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var u models.Customer

	Customers, totalPages, err := u.GetCustomersPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Customer not found!"})
		c.Abort()
		return
	}
	fmt.Println(Customers)
	var responses []CustomerResponse
	for _, Customer := range Customers {
		response := CustomerResponse{
			ID:       Customer.ID,
			Name:     Customer.CustomerName,
			Username: Customer.CustomerUsername,
			Email:    Customer.CustomerEmail,
			Phone:    Customer.CustomerPhone,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       responses,
		"totalPages": totalPages,
	})
}

func (repository *CustomerRepo) DeleteCustomer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	u := models.Customer{}

	err := models.DeleteCustomer(repository.Db, &u, id)

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
