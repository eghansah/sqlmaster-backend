package models

import "time"

type SQLQueryParam struct {
	DataType string `json:"type"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	Multi    bool   `json:"multi"`
	Value    any    `json:"value,omitempty"`
}

type SQLQuery struct {
	ID           int             `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Query        string          `json:"query,omitempty"`
	DatasourceID string          `json:"datasource"`
	Params       []SQLQueryParam `json:"params"`
	Enabled      bool
	CreatedAt    time.Time
	UpdateAt     time.Time
}

type ReportRequest struct {
	ID          int    `json:"id"`
	SQLQueryID  int    `json:"sql_query_id"`
	ReportTitle string `json:"report_title,omitempty"`
	Params      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"params"`
	Status    string `json:"status,omitempty"`
	Comment   string `json:"comment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Datasource struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	DSN     string `json:"url"`
	Enabled bool
}
