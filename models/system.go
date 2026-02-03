package models

type TableInfo struct {
	SchemaName string `json:"schema_name"`
	TableName  string `json:"table_name"`
	TableType  string `json:"table_type"`
}

type ColumnInfo struct {
	ColumnName    string  `json:"column_name"`
	DataType      string  `json:"data_type"`
	IsNullable    string  `json:"is_nullable"`
	ColumnDefault *string `json:"column_default,omitempty"`
}
