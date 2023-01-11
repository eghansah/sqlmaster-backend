package main

import (
	"fmt"
	"log"
	"net/http"
	"sqlmaster/internal/models"
	"sqlmaster/internal/repository"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type application struct {
	Host        string
	Port        int
	DatabaseDSN string
	DB          repository.DatabaseRepo
	Logger      *zap.SugaredLogger
	Models      struct {
		Datasources repository.DatasourcesManager
		SQLQueries  repository.SQLQueryManager
		SQLRequests repository.SQLRequestsManager
	}
}

func (app *application) Run() {
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", app.Host, app.Port), app.Routes())
}

func (app *application) AddNewReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		req := models.SQLQuery{}
		err := app.readJSON(w, r, &req)
		if err != nil {
			fmt.Println("Error Parsing Request", err)
			app.writeJSON(w, http.StatusBadRequest, Response{
				Error:   true,
				Message: "bad request",
				Data:    err,
			})
			return
		}

		fmt.Printf("\n\nRequest Received => \n%+v\n", req)
		err = app.Models.SQLQueries.AddNewQuery(req)
		if err != nil {
			fmt.Println("Error Parsing Request", err)
			app.writeJSON(w, http.StatusBadRequest, Response{
				Error:   true,
				Message: "could not save new report",
				Data:    err,
			})
			return
		}

		resp := Response{
			Error:   false,
			Message: "Query/Report successfully saved",
			Data:    "",
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) AllReports() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		reports, err := app.Models.SQLQueries.All()
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		resp := Response{
			Error:   false,
			Message: "",
			Data:    reports,
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) GetReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		log.Printf("GetReport called with id: %s", chi.URLParam(r, "id"))
		reportID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		report, err := app.Models.SQLQueries.Get(reportID)
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		//Hide query
		report.Query = ""

		resp := Response{
			Error:   false,
			Message: "",
			Data:    report,
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) NewReportRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		repRequest := models.ReportRequest{}
		err := app.readJSON(w, r, &repRequest)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		err = app.Models.SQLRequests.AddNewRequest(repRequest)
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		resp := Response{
			Error:   false,
			Message: "",
			Data:    "",
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) GetRequestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		requests, err := app.Models.SQLRequests.All()
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		// rs, _ := json.MarshalIndent(requests, "", "\t")
		// log.Printf("%s\n", rs)
		resp := Response{
			Error:   false,
			Message: "",
			Data:    requests,
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) GetDataSourcesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		datasources, err := app.Models.Datasources.All()
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		resp := Response{
			Error:   false,
			Message: "",
			Data:    datasources,
		}
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *application) AddNewDatasourceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Error   bool        `json:"error"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}

		ds := models.Datasource{}
		err := app.readJSON(w, r, &ds)
		if err != nil {
			log.Printf("Could not read input: %s", err)
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		err = app.Models.Datasources.AddNewDatasource(ds)
		if err != nil {
			log.Printf("Could not add datasource: %s", err)
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		resp := Response{
			Error:   false,
			Message: "Datasource saved successfully",
		}
		app.writeJSON(w, http.StatusAccepted, resp)
	}
}
