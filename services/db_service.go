package services

import (
	"gowes/db"
)

type TableInfo struct {
	SchemaName string `json:"schema_name"`
	TableName  string `json:"table_name"`
	TableType  string `json:"table_type"`
}

type ColumnInfo struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable string `json:"is_nullable"`
	ColumnDefault *string `json:"column_default,omitempty"`
}

func ListTables() ([]TableInfo, error) {
	query := `
		SELECT 
			table_schema,
			table_name,
			table_type
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
		ORDER BY table_schema, table_name;
	`
	
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.SchemaName, &t.TableName, &t.TableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func GetTableColumns(schema, table string) ([]ColumnInfo, error) {
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
	
	rows, err := db.DB.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		if err := rows.Scan(&c.ColumnName, &c.DataType, &c.IsNullable, &c.ColumnDefault); err != nil {
			return nil, err
		}
		columns = append(columns, c)
	}
	return columns, rows.Err()
}