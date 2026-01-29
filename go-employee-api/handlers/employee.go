package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"go-employee-api/models"
	"go-employee-api/responses"
	"go-employee-api/storage"
)


func VersionHandler(_ *http.Request, _ map[string]string) responses.APIResponse {
	return responses.APIResponse{
		Status:      http.StatusOK,
		Message:     "API Version 1",
		IsPlainText: true,
	}
}


func GetEmployees(r *http.Request, _ map[string]string) responses.APIResponse {
	query := r.URL.Query().Get("is_active")

	var isActive *bool
	if query == "true" {
		v := true
		isActive = &v
	} else if query == "false" {
		v := false
		isActive = &v
	}

	employees, err := storage.GetEmployees(isActive)
	if err != nil {
		return responses.APIResponse{
			Status:      http.StatusInternalServerError,
			Err:         errors.New("Error while accessing data resource."),
			IsPlainText: true,
		}
	}

	return responses.APIResponse{
		Status: http.StatusOK,
		Data:   employees,
	}
}


func GetEmployeeByName(_ *http.Request, params map[string]string) responses.APIResponse {
	name := params["name"]

	employee, err := storage.GetEmployeeByName(name)
	if err == sql.ErrNoRows {
		return responses.APIResponse{
			Status:      http.StatusNotFound,
			Err:         errors.New("Data resource not created. Use POST to create data first."),
			IsPlainText: true,
		}
	}
	if err != nil {
		return responses.APIResponse{
			Status:      http.StatusInternalServerError,
			Err:         errors.New("Error while accessing data resource."),
			IsPlainText: true,
		}
	}

	return responses.APIResponse{
		Status: http.StatusOK,
		Data:   employee,
	}
}


func CreateEmployee(r *http.Request, _ map[string]string) responses.APIResponse {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		return responses.APIResponse{
			Status:      http.StatusBadRequest,
			Err:         errors.New("Invalid JSON format."),
			IsPlainText: true,
		}
	}

	if emp.Name == "" || emp.Age == 0 || emp.Address == "" {
		return responses.APIResponse{
			Status:      http.StatusBadRequest,
			Err:         errors.New("Name, age and address are required."),
			IsPlainText: true,
		}
	}

	if emp.IsActive == nil {
		active := true
		emp.IsActive = &active
	}

	// Check if employee already exists
	_, err := storage.GetEmployeeByName(emp.Name)
	if err == nil {
		return responses.APIResponse{
			Status:      http.StatusConflict,
			Err:         errors.New("Data resource already exists. Use PUT to update data."),
			IsPlainText: true,
		}
	}

	if err := storage.CreateEmployee(emp); err != nil {
		return responses.APIResponse{
			Status:      http.StatusInternalServerError,
			Err:         errors.New("Error while accessing data resource."),
			IsPlainText: true,
		}
	}

	return responses.APIResponse{
		Status:      http.StatusOK,
		Message:     "Data resource created successfully.",
		IsPlainText: true,
	}
}

func UpdateEmployee(r *http.Request, params map[string]string) responses.APIResponse {
	name := params["name"]

	var update models.Employee
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return responses.APIResponse{
			Status:      http.StatusBadRequest,
			Err:         errors.New("Invalid JSON format."),
			IsPlainText: true,
		}
	}
	err := storage.UpdateEmployee(name, update)
	if err == sql.ErrNoRows {
		return responses.APIResponse{
			Status:      http.StatusNotFound,
			Err:         errors.New("Data resource not created. Use POST to create data first."),
			IsPlainText: true,
		}
	}
	if err != nil {
		return responses.APIResponse{
			Status:      http.StatusInternalServerError,
			Err:         errors.New("Error while accessing data resource."),
			IsPlainText: true,
		}
	}
	return responses.APIResponse{
		Status:      http.StatusOK,
		Message:     "Data resource updated successfully.",
		IsPlainText: true,
	}
}


func DeleteEmployee(_ *http.Request, params map[string]string) responses.APIResponse {
	name := params["name"]

	err := storage.DeleteEmployee(name)
	if err == sql.ErrNoRows {
		return responses.APIResponse{
			Status:      http.StatusNotFound,
			Err:         errors.New("Data resource not created. Use POST to create data first."),
			IsPlainText: true,
		}
	}
	if err != nil {
		return responses.APIResponse{
			Status:      http.StatusInternalServerError,
			Err:         errors.New("Error while accessing data resource."),
			IsPlainText: true,
		}
	}

	return responses.APIResponse{
		Status:      http.StatusOK,
		Message:     "Data resource deleted successfully.",
		IsPlainText: true,
	}
}
