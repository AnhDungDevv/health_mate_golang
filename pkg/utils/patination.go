package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultSize = 10
)

type PaginationQuery struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)

	if err != nil {
		return err
	}
	q.Page = n
	return nil
}

// Set order by
func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// Get OrderBy
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

// Get OrderBy
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

func GetPaginationFromCtx(c *gin.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetPage(c.Query("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.Query("orderBy"))

	return q, nil
}

func GetTotalPage(totalCount int, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	if pageSize == 0 {
		return false
	}
	return currentPage < int(math.Ceil(float64(totalCount)/float64(pageSize)))
}
