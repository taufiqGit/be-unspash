package services

import (
	"gowes/models"
	"gowes/repositories"
)

type SystemService interface {
	ListTables() ([]models.TableInfo, error)
	GetTableColumns(schema, table string) ([]models.ColumnInfo, error)
}

type systemService struct {
	repo repositories.SystemRepository
}

func NewSystemService(repo repositories.SystemRepository) SystemService {
	return &systemService{repo: repo}
}

func (s *systemService) ListTables() ([]models.TableInfo, error) {
	return s.repo.ListTables()
}

func (s *systemService) GetTableColumns(schema, table string) ([]models.ColumnInfo, error) {
	return s.repo.GetTableColumns(schema, table)
}
