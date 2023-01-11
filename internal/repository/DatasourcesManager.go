package repository

import "sqlmaster/internal/models"

type DatasourcesManager struct {
	db DatabaseRepo
}

func NewDatasourceManager(db DatabaseRepo) DatasourcesManager {
	return DatasourcesManager{db: db}
}
func (m *DatasourcesManager) All() ([]models.Datasource, error) {
	return m.db.Datasources_All()
}

func (m *DatasourcesManager) AddNewDatasource(request models.Datasource) error {
	return m.db.Datasources_AddNewDatasource(request)
}
