package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type OrderTypeRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.OrderType, int, error)
	Create(companyID string, orderType models.OrderTypeInput) (models.OrderType, error)
	Update(orderType models.OrderTypeInput, id string) (models.OrderType, error)
	FindByID(id string) (models.OrderType, error)
	Delete(id string) error
}

type orderTypeRepository struct {
	db *sql.DB
}

func NewOrderTypeRepository(db *sql.DB) OrderTypeRepository {
	return &orderTypeRepository{db: db}
}

func (r *orderTypeRepository) FindAll(companyID string, params models.PaginationParams) ([]models.OrderType, int, error) {
	baseQuery := " FROM order_types WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"name": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	query := "SELECT id, company_id, name, is_active_price_adjustment, increase_type, decrease_type, increase_value, decrease_value, price_increase, price_decrease" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, params.SortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orderTypes = []models.OrderType{}
	for rows.Next() {
		var orderType models.OrderType
		if err := rows.Scan(&orderType.ID, &orderType.CompanyID, &orderType.Name, &orderType.IsActivePriceAdjustment, &orderType.IncreaseType, &orderType.DecreaseType, &orderType.IncreaseValue, &orderType.DecreaseValue, &orderType.PriceIncrease, &orderType.PriceDecrease); err != nil {
			return nil, 0, err
		}
		orderTypes = append(orderTypes, orderType)
	}

	return orderTypes, total, nil
}

func (r *orderTypeRepository) Create(companyID string, orderType models.OrderTypeInput) (models.OrderType, error) {
	query := "INSERT INTO order_types (company_id, name, is_active_price_adjustment, increase_type, decrease_type, increase_value, decrease_value, price_increase, price_decrease, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, company_id, name, is_active_price_adjustment, increase_type, decrease_type, increase_value, decrease_value, price_increase, price_decrease, is_active, created_at, updated_at"
	args := []interface{}{companyID, orderType.Name, orderType.IsActivePriceAdjustment, orderType.IncreaseType, orderType.DecreaseType, orderType.IncreaseValue, orderType.DecreaseValue, orderType.PriceIncrease, orderType.PriceDecrease, orderType.IsActive}

	var orderTypeResult models.OrderType
	err := r.db.QueryRow(query, args...).Scan(&orderTypeResult.ID, &orderTypeResult.CompanyID, &orderTypeResult.Name, &orderTypeResult.IsActivePriceAdjustment, &orderTypeResult.IncreaseType, &orderTypeResult.DecreaseType, &orderTypeResult.IncreaseValue, &orderTypeResult.DecreaseValue, &orderTypeResult.PriceIncrease, &orderTypeResult.PriceDecrease, &orderTypeResult.IsActive, &orderTypeResult.CreatedAt, &orderTypeResult.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return models.OrderType{}, err
	}
	return orderTypeResult, nil
}

func (r *orderTypeRepository) Update(orderType models.OrderTypeInput, id string) (models.OrderType, error) {
	query := "UPDATE order_types SET name = $1, is_active_price_adjustment = $2, increase_type = $3, decrease_type = $4, increase_value = $5, decrease_value = $6, price_increase = $7, price_decrease = $8, is_active = $9 WHERE id = $10 RETURNING id, company_id, name, is_active_price_adjustment, increase_type, decrease_type, increase_value, decrease_value, price_increase, price_decrease, is_active, created_at, updated_at"
	args := []interface{}{orderType.Name, orderType.IsActivePriceAdjustment, orderType.IncreaseType, orderType.DecreaseType, orderType.IncreaseValue, orderType.DecreaseValue, orderType.PriceIncrease, orderType.PriceDecrease, orderType.IsActive, id}

	var orderTypeResult models.OrderType
	err := r.db.QueryRow(query, args...).Scan(&orderTypeResult.ID, &orderTypeResult.CompanyID, &orderTypeResult.Name, &orderTypeResult.IsActivePriceAdjustment, &orderTypeResult.IncreaseType, &orderTypeResult.DecreaseType, &orderTypeResult.IncreaseValue, &orderTypeResult.DecreaseValue, &orderTypeResult.PriceIncrease, &orderTypeResult.PriceDecrease, &orderTypeResult.IsActive, &orderTypeResult.CreatedAt, &orderTypeResult.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return models.OrderType{}, err
	}
	fmt.Println(orderTypeResult)
	return orderTypeResult, nil
}

func (r *orderTypeRepository) FindByID(id string) (models.OrderType, error) {
	query := "SELECT id, company_id, name, is_active_price_adjustment, increase_type, decrease_type, increase_value, decrease_value, price_increase, price_decrease, is_active, created_at, updated_at FROM order_types WHERE id = $1"
	args := []interface{}{id}

	var orderTypeResult models.OrderType
	err := r.db.QueryRow(query, args...).Scan(&orderTypeResult.ID, &orderTypeResult.CompanyID, &orderTypeResult.Name, &orderTypeResult.IsActivePriceAdjustment, &orderTypeResult.IncreaseType, &orderTypeResult.DecreaseType, &orderTypeResult.IncreaseValue, &orderTypeResult.DecreaseValue, &orderTypeResult.PriceIncrease, &orderTypeResult.PriceDecrease, &orderTypeResult.IsActive, &orderTypeResult.CreatedAt, &orderTypeResult.UpdatedAt)
	if err != nil {
		return models.OrderType{}, err
	}
	return orderTypeResult, nil
}

func (r *orderTypeRepository) Delete(id string) error {
	query := "DELETE FROM order_types WHERE id = $1"
	args := []interface{}{id}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
