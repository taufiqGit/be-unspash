package repositories

import (
	"database/sql"
	"gowes/models"
	"time"
)

type CashierShiftRepository interface {
	StartShift(shift models.CashierShift) (models.CashierShift, error)
	FindActiveShiftByUser(companyID string, userID string) (models.CashierShift, error)
	EndShift(shift models.CashierShift) (models.CashierShift, error)
}

type cashierShiftRepository struct {
	db *sql.DB
}

func NewCashierShiftRepository(db *sql.DB) CashierShiftRepository {
	return &cashierShiftRepository{db: db}
}

func (r *cashierShiftRepository) StartShift(shift models.CashierShift) (models.CashierShift, error) {
	row := r.db.QueryRow(`
		INSERT INTO cashier_shifts (
			company_id, outlet_id, user_id, start_time, end_time, status, opening_cash, closing_cash, expected_cash, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, company_id, outlet_id, user_id, start_time, end_time, status, opening_cash, closing_cash, expected_cash, created_at, updated_at
	`,
		shift.CompanyID,
		shift.OutletID,
		shift.UserID,
		shift.StartTime,
		shift.EndTime,
		shift.Status,
		shift.OpeningCash,
		shift.ClosingCash,
		shift.ExpectedCash,
		shift.CreatedAt,
		shift.UpdatedAt,
	)

	return scanCashierShiftRow(row)
}

func (r *cashierShiftRepository) FindActiveShiftByUser(companyID string, userID string) (models.CashierShift, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, outlet_id, user_id, start_time, end_time, status, opening_cash, closing_cash, expected_cash, created_at, updated_at
		FROM cashier_shifts
		WHERE company_id = $1
		  AND user_id = $2
		  AND (status ILIKE 'open' OR end_time <= start_time)
		ORDER BY start_time DESC
		LIMIT 1
	`, companyID, userID)

	shift, err := scanCashierShiftRow(row)
	if err != nil {
		return models.CashierShift{}, err
	}
	return shift, nil
}

func (r *cashierShiftRepository) EndShift(shift models.CashierShift) (models.CashierShift, error) {
	row := r.db.QueryRow(`
		UPDATE cashier_shifts
		SET end_time = $1,
		    status = $2,
		    closing_cash = $3,
		    expected_cash = $4,
		    updated_at = $5
		WHERE id = $6
		  AND company_id = $7
		RETURNING id, company_id, outlet_id, user_id, start_time, end_time, status, opening_cash, closing_cash, expected_cash, created_at, updated_at
	`, shift.EndTime, shift.Status, shift.ClosingCash, shift.ExpectedCash, shift.UpdatedAt, shift.ID, shift.CompanyID)

	return scanCashierShiftRow(row)
}

func scanCashierShiftRow(row *sql.Row) (models.CashierShift, error) {
	var shift models.CashierShift
	var endTime sql.NullTime

	err := row.Scan(
		&shift.ID,
		&shift.CompanyID,
		&shift.OutletID,
		&shift.UserID,
		&shift.StartTime,
		&endTime,
		&shift.Status,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.ExpectedCash,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	)
	if err != nil {
		return models.CashierShift{}, err
	}

	if endTime.Valid {
		shift.EndTime = endTime.Time
	} else {
		shift.EndTime = time.Time{}
	}

	return shift, nil
}
