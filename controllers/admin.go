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

type AdminRepo struct {
	Db *gorm.DB
}

func NewAdmin() *AdminRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.Admin{})
	return &AdminRepo{Db: db}
}

type AdminResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (repository *AdminRepo) AdminRegister(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := models.Admin{}

	a.AdminName = input.Name
	a.AdminUsername = input.Username
	a.AdminEmail = input.Email
	a.AdminPhone = input.Phone
	a.AdminPassword = input.Password

	if err := a.HashPassword(a.AdminPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := models.CreateAdmin(repository.Db, &a)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "admin registration success", "adminId": a.ID, "email": a.AdminEmail, "username": a.AdminUsername})

}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (repository *AdminRepo) AdminLogin(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := models.Admin{}

	// check if email exists and password is correct
	err := models.GetAdminByEmail(repository.Db, &a, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	credentialError := a.VerifyPassword(input.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	// Create JWT Token for Authorization
	tokenString, err := auth.GenerateToken(int(a.ID), a.AdminEmail, a.AdminUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func (repository *AdminRepo) AdminLoggedIn(c *gin.Context) {

	id, _, err := auth.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := models.Admin{}

	if err := models.GetAdminById(repository.Db, &a, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Admin not found!"})
		c.Abort()
		return
	}

	var response AdminResponse
	response.ID = a.ID
	response.Name = a.AdminName
	response.Username = a.AdminUsername
	response.Email = a.AdminEmail
	response.Phone = a.AdminPhone

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

type AdminUpdateInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (repository *AdminRepo) UpdateAdmin(c *gin.Context) {
	var input AdminUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	a := models.Admin{}

	err := models.GetAdminById(repository.Db, &a, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if input.Name != "" {
		a.AdminName = input.Name
	}

	if input.Username != "" {
		a.AdminUsername = input.Username
	}

	if input.Email != "" {
		a.AdminEmail = input.Email
	}

	if input.Phone != "" {
		a.AdminPhone = input.Phone
	}

	if input.Password != "" {
		a.AdminPassword = input.Password

		if err := a.HashPassword(a.AdminPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}

	err = models.UpdateAdmin(repository.Db, &a, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var response AdminResponse
	response.ID = a.ID
	response.Name = a.AdminName
	response.Email = a.AdminEmail
	response.Phone = a.AdminPhone

	c.JSON(http.StatusOK, response)
}

func (repository *AdminRepo) GetAdminsData(c *gin.Context) {

	page, pageSize, sortField, sortOrder := pagination.Paginate(c)

	var a models.Admin

	admins, totalPages, err := a.GetAdminsPaginate(repository.Db, page, pageSize, sortField, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Admin not found!"})
		c.Abort()
		return
	}

	var responses []AdminResponse
	for _, admin := range admins {
		response := AdminResponse{
			ID:       admin.ID,
			Name:     admin.AdminName,
			Username: admin.AdminUsername,
			Email:    admin.AdminEmail,
			Phone:    admin.AdminPhone,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       responses,
		"totalPages": totalPages,
	})
}

func (repository *AdminRepo) DeleteAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	a := models.Admin{}

	err := models.DeleteAdmin(repository.Db, &a, id)

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
