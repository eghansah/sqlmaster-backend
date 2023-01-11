package repository

import "sqlmaster/internal/models"

type SQLRequestsManager struct {
	db DatabaseRepo
}

func NewSQLRequestsManager(db DatabaseRepo) SQLRequestsManager {
	return SQLRequestsManager{db: db}
}

func (m *SQLRequestsManager) All() ([]*models.ReportRequest, error) {
	return m.db.ReportRequest_AllRequests()
}

func (m *SQLRequestsManager) AddNewRequest(request models.ReportRequest) error {
	return m.db.ReportRequest_AddNewRequest(request)
}
