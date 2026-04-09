package services

import (
	"database/sql"
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var (
	ErrCashierShiftUserRequired   = errors.New("user_id is required")
	ErrCashierShiftOutletRequired = errors.New("outlet_id is required")
	ErrCashierShiftAlreadyActive  = errors.New("user still has active cashier shift")
)

type CashierShiftService interface {
	StartShift(companyID string, authUserID string, input models.StartCashierShiftInput) (models.CashierShift, error)
	EndShift(companyID string, authUserID string, input models.EndCashierShiftInput) (models.CashierShift, error)
}

type cashierShiftService struct {
	repo repositories.CashierShiftRepository
}

func NewCashierShiftService(repo repositories.CashierShiftRepository) CashierShiftService {
	return &cashierShiftService{repo: repo}
}

func (s *cashierShiftService) StartShift(companyID string, authUserID string, input models.StartCashierShiftInput) (models.CashierShift, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		userID = authUserID
	}
	if strings.TrimSpace(userID) == "" {
		return models.CashierShift{}, ErrCashierShiftUserRequired
	}

	if strings.TrimSpace(input.OutletID) == "" {
		return models.CashierShift{}, ErrCashierShiftOutletRequired
	}

	if _, err := s.repo.FindActiveShiftByUser(companyID, userID); err == nil {
		return models.CashierShift{}, ErrCashierShiftAlreadyActive
	} else if !errors.Is(err, sql.ErrNoRows) {
		return models.CashierShift{}, err
	}

	startTime := input.StartTime
	if startTime.IsZero() {
		startTime = time.Now().UTC()
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "open"
	}

	now := time.Now().UTC()
	shift := models.CashierShift{
		CompanyID:    companyID,
		OutletID:     strings.TrimSpace(input.OutletID),
		UserID:       userID,
		StartTime:    startTime,
		EndTime:      startTime,
		Status:       status,
		OpeningCash:  input.OpeningCash,
		ClosingCash:  input.ClosingCash,
		ExpectedCash: input.ExpectedCash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repo.StartShift(shift)
}

func (s *cashierShiftService) EndShift(companyID string, authUserID string, input models.EndCashierShiftInput) (models.CashierShift, error) {
	userID := strings.TrimSpace(authUserID)
	if userID == "" {
		return models.CashierShift{}, ErrCashierShiftUserRequired
	}

	shift, err := s.repo.FindActiveShiftByUser(companyID, userID)
	if err != nil {
		return models.CashierShift{}, err
	}

	endTime := input.EndTime
	if endTime.IsZero() {
		endTime = time.Now().UTC()
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "closed"
	}

	shift.EndTime = endTime
	shift.Status = status
	shift.ClosingCash = input.ClosingCash
	shift.ExpectedCash = input.ExpectedCash
	shift.UpdatedAt = time.Now().UTC()

	return s.repo.EndShift(shift)
}
