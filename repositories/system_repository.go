package repositories

import (
	"database/sql"
	"gowes/models"
)

type SystemRepository interface {
	ListTables() ([]models.TableInfo, error)
	GetTableColumns(schema, table string) ([]models.ColumnInfo, error)
}

type systemRepository struct {
	db *sql.DB
}

func NewSystemRepository(db *sql.DB) SystemRepository {
	return &systemRepository{db: db}
}

func (r *systemRepository) ListTables() ([]models.TableInfo, error) {
	query := `
		SELECT 
			table_schema,
			table_name,
			table_type
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
		ORDER BY table_schema, table_name;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.TableInfo
	for rows.Next() {
		var t models.TableInfo
		if err := rows.Scan(&t.SchemaName, &t.TableName, &t.TableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (r *systemRepository) GetTableColumns(schema, table string) ([]models.ColumnInfo, error) {
	query := `
		SELECT 
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position;
	`

	rows, err := r.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns = []models.ColumnInfo{}
	for rows.Next() {
		var c models.ColumnInfo
		if err := rows.Scan(&c.ColumnName, &c.DataType, &c.IsNullable, &c.ColumnDefault); err != nil {
			return nil, err
		}
		columns = append(columns, c)
	}
	return columns, rows.Err()
}
