package utils

import (
	"fmt"
	"gowes/models"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func ParsePaginationParams(r *http.Request) models.PaginationParams {
	q := r.URL.Query()
	fmt.Println(q.Get("page"))
	page, _ := strconv.Atoi(q.Get("page"))
	fmt.Println(page)
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	sortBy := q.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := strings.ToUpper(q.Get("sort_order"))
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	search := q.Get("search")

	return models.PaginationParams{
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    search,
	}
}

func CalculateMeta(totalData int, params models.PaginationParams) models.PaginationMeta {
	totalPage := int(math.Ceil(float64(totalData) / float64(params.Limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	return models.PaginationMeta{
		CurrentPage: params.Page,
		TotalPage:   totalPage,
		TotalData:   totalData,
		Limit:       params.Limit,
	}
}
