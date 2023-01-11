package repository

import (
	"database/sql"
	"sqlmaster/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	SQLQuery_AllReports() ([]*models.SQLQuery, error)
	SQLQuery_ReportByID(id int) (*models.SQLQuery, error)
	SQLQuery_AddNewQuery(sqlQuery models.SQLQuery) error

	ReportRequest_AllRequests() ([]*models.ReportRequest, error)
	ReportRequest_AddNewRequest(request models.ReportRequest) error

	Datasources_All() ([]models.Datasource, error)
	Datasources_AddNewDatasource(ds models.Datasource) error
}
