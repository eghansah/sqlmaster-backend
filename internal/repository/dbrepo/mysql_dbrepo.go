package dbrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sqlmaster/internal/models"
	"strings"
	"time"
)

type MysqlDBRepo struct {
	DB      *sql.DB
	Timeout time.Duration
}

func (m *MysqlDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *MysqlDBRepo) SQLQuery_AllReports() ([]*models.SQLQuery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		select 
			id, created_at, title, description
		from sql_query
		where enabled = true
		order by title asc;
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []*models.SQLQuery{}
	for rows.Next() {
		var report models.SQLQuery
		err := rows.Scan(
			&report.ID,
			&report.CreatedAt,
			&report.Title,
			&report.Description,
		)
		if err != nil {
			return nil, err
		}

		reports = append(reports, &report)
	}

	return reports, nil
}

func (m *MysqlDBRepo) SQLQuery_ReportByID(id int) (*models.SQLQuery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		select 
			id, created_at, title, description
		from sql_query
		where enabled = true and id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var report models.SQLQuery
	err := row.Scan(
		&report.ID,
		&report.CreatedAt,
		&report.Title,
		&report.Description,
	)
	if err != nil {
		return nil, err
	}

	//Getting params
	paramsQuery := `
		select 
			name, label, datatype, multi_value
		from sql_query_param
		where sql_query = ?
	`
	rows, err := m.DB.QueryContext(ctx, paramsQuery, report.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	params := []models.SQLQueryParam{}
	for rows.Next() {
		var param models.SQLQueryParam
		err := rows.Scan(
			&param.Name,
			&param.Label,
			&param.DataType,
			&param.Multi,
		)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	report.Params = params

	return &report, nil
}

func (m *MysqlDBRepo) SQLQuery_AddNewQuery(sqlQuery models.SQLQuery) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		insert into sql_query(id, created_at, title, description, query, enabled, datasource)
		values (0, current_timestamp, ?, ?, ?, true, ?);
	`

	paramQuery := `
		insert into sql_query_param(id, created_at, datatype, name, label, multi_value, sql_query)
		values (0, current_timestamp, ?, ?, ?, false, ?)
	`

	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	row, err := tx.Exec(query, sqlQuery.Title, sqlQuery.Description, sqlQuery.Query, sqlQuery.DatasourceID)
	if err != nil {
		tx.Rollback()
		return err
	}

	reportID, err := row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, params := range sqlQuery.Params {
		if strings.TrimSpace(params.DataType) != "" && strings.TrimSpace(params.Name) != "" &&
			strings.TrimSpace(params.Label) != "" {
			_, err := tx.Exec(paramQuery, strings.TrimSpace(params.DataType), strings.TrimSpace(params.Name),
				strings.TrimSpace(params.Label), reportID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *MysqlDBRepo) ReportRequest_AllRequests() ([]*models.ReportRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		select a.id, a.created_at, b.title, status, comment
		from report_request a
        left outer join sql_query b on b.id = a.sql_query_id
		order by a.id desc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*models.ReportRequest{}
	for rows.Next() {
		var request models.ReportRequest
		err := rows.Scan(
			&request.ID,
			&request.CreatedAt,
			&request.ReportTitle,
			&request.Status,
			&request.Comment,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, &request)
	}

	return requests, nil
}

func (m *MysqlDBRepo) ReportRequest_AddNewRequest(request models.ReportRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	insertQuery := `
		insert into report_request(id, sql_query_id, query, params)		
			select 0, id, query, '{}'
			from sql_query
			where id = ?
		;
	`

	updateQuery := `
			update report_request set params = ?
			where id = ?
	`

	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	reqParams, err := json.Marshal(request.Params)
	if err != nil {
		log.Printf("Error occured while converting params to json string: %s\n", err)
	}
	res, err := tx.Exec(insertQuery, request.SQLQueryID)
	if err != nil {
		return err
	}

	repReqID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateQuery, reqParams, repReqID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *MysqlDBRepo) Datasources_AddNewDatasource(ds models.Datasource) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		insert into datasources (id, created_at, type, name, dsn, enabled)
		values
		(0, now(), ?, ?, ?, true);
	`

	_, err := m.DB.ExecContext(ctx, query, ds.Type, ds.Name, ds.DSN)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlDBRepo) Datasources_All() ([]models.Datasource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	query := `
		select id, name from datasources 
		where enabled=true
		order by name asc
	`

	datasources := []models.Datasource{}
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ds models.Datasource
		err = rows.Scan(&ds.ID, &ds.Name)
		if err != nil {
			return nil, err
		}

		datasources = append(datasources, ds)
	}

	return datasources, nil
}
