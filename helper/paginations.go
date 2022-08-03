package helper

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Pagination struct {
	Limit        int         `json:"limit"`
	Page         int         `json:"page"`
	Sort         string      `json:"sort"`
	TotalRows    int         `json:"total_rows"`
	FirstPage    string      `json:"first_page"`
	PreviousPage string      `json:"previous_page"`
	NextPage     string      `json:"next_page"`
	LastPage     string      `json:"last_page"`
	Rows         interface{} `json:"rows"`
	TotalPages   int         `json:"total_pages"`
}

func GeneratePaginationRequest(c *gin.Context) *Pagination {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sort := c.DefaultQuery("sort", "name asc")

	return &Pagination{Limit: limit, Page: page, Sort: sort}
}
