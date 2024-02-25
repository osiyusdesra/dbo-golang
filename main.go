package main

import (
	"be-dbo-golang/controllers"
	"os"

	"be-dbo-golang/utils/middlewares"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowHeaders("*")
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "DELETE"}
	router.Use(cors.New(config))

	adminRepo := controllers.NewAdmin()

	customerRepo := controllers.NewCustomer()

	supplierRepo := controllers.NewSupplier()

	productRepo := controllers.NewProduct()

	brandRepo := controllers.NewBrand()

	orderRepo := controllers.NewOrder()

	apiEndpoint := router.Group("/api/v1")
	{
		apiEndpoint.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "This is backend by Osi Yusdesra")
		})

		//ADMIN OPEN API
		apiEndpoint.POST("/admin/register", adminRepo.AdminRegister)
		apiEndpoint.POST("/admin/login", adminRepo.AdminLogin)
		apiEndpoint.GET("/admins", adminRepo.GetAdminsData)

		//CUSTOMER OPEN API
		apiEndpoint.POST("/customer/register", customerRepo.CustomerRegister)
		apiEndpoint.POST("/customer/login", customerRepo.CustomerLogin)
		apiEndpoint.GET("/customer/list", customerRepo.GetCustomersData)

		//SUPPLIER OPEN API
		apiEndpoint.POST("/supplier/register", supplierRepo.SupplierRegister)
		apiEndpoint.POST("/supplier/login", supplierRepo.SupplierLogin)
		apiEndpoint.GET("/supplier/list", supplierRepo.GetSuppliersData)

		//PRODUCT OPEN API
		apiEndpoint.GET("/product/list", productRepo.GetProductsData)
		apiEndpoint.GET("/product/data/:id", productRepo.GetProductById)

		//BRAND OPEN API
		apiEndpoint.GET("/brand/list", brandRepo.GetBrandsData)
		apiEndpoint.GET("/brand/data/:id", brandRepo.GetBrandById)

		// PRIVATE API
		secured := apiEndpoint.Group("/secured").Use(middlewares.JwtAuthMiddleware())
		{
			// ADMIN
			secured.GET("/admin/data", adminRepo.AdminLoggedIn)
			secured.PUT("admin/update-password/:id", adminRepo.UpdateAdmin)
			secured.DELETE("admin/delete/:id", adminRepo.DeleteAdmin)

			// CUSTOMER
			secured.GET("/customer/data", customerRepo.CustomerLoggedIn)
			secured.PUT("/customer/update/:id", customerRepo.UpdateCustomer)
			secured.DELETE("/customer/delete/:id", customerRepo.DeleteCustomer)

			// SUPPLIER
			secured.GET("/supplier/data", supplierRepo.SupplierLoggedIn)
			secured.PUT("/supplier/update/:id", supplierRepo.UpdateSupplier)
			secured.DELETE("/supplier/delete/:id", supplierRepo.DeleteSupplier)

			// PRODUCT
			secured.POST("/product/create", productRepo.SaveProductData)
			secured.PUT("/product/update/:id", productRepo.UpdateProduct)
			secured.DELETE("/product/delete/:id", productRepo.DeleteProduct)

			// BRAND
			secured.POST("/brand/create", brandRepo.SaveBrandData)
			secured.PUT("/brand/update/:id", brandRepo.UpdateBrand)
			secured.DELETE("/brand/delete/:id", brandRepo.DeleteBrand)

			// ORDER
			secured.GET("/order/list", orderRepo.GetOrdersData)
			secured.GET("/order/data/:id", orderRepo.GetOrderById)
			secured.POST("/order/create", orderRepo.SaveOrderData)
			secured.PUT("/order/update/:id", orderRepo.UpdateOrder)
			secured.DELETE("/order/delete/:id", orderRepo.DeleteOrder)
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}

	//RUN ON PORT 3000
	router.Run(":" + port)
}
