package repository

import "sqlmaster/internal/models"

type SQLQueryManager struct {
	db DatabaseRepo
}

func NewSQLQueryManager(db DatabaseRepo) SQLQueryManager {
	return SQLQueryManager{db: db}
}
func (m *SQLQueryManager) All() ([]*models.SQLQuery, error) {
	return m.db.SQLQuery_AllReports()
}

func (m *SQLQueryManager) Get(id int) (*models.SQLQuery, error) {
	return m.db.SQLQuery_ReportByID(id)
}

func (m *SQLQueryManager) AddNewQuery(sqlquery models.SQLQuery) error {
	return m.db.SQLQuery_AddNewQuery(sqlquery)
}
