package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Sort       string `json:"sort"`
	TotalPages int    `json:"total_pages"`
}

func Paginate(c *gin.Context) (int, int, string, string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))           //Default 1
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10")) //Default 10
	sortField := c.DefaultQuery("sort_by", "id")                   // Default sorting by ID
	sortOrder := c.DefaultQuery("sort_order", "asc")               // Default sort order

	return page, pageSize, sortField, sortOrder
}
