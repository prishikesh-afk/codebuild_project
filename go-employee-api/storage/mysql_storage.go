package storage

import (
	"database/sql"
	"fmt"
	"os"

	"go-employee-api/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = connectDB()
	if err != nil {
		panic(err)
	}
}

func connectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	//tests the connections 
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetEmployees(isActive *bool) ([]models.Employee, error) {
	query := "SELECT id, name, age, address, is_active FROM employees"
	var args []any

	if isActive != nil {
		query += " WHERE is_active = ?"
		args = append(args, *isActive)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee

	for rows.Next() {
		var e models.Employee
		var active bool

		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Age,
			&e.Address,
			&active,
		); err != nil {
			return nil, err
		}

		e.IsActive = &active
		employees = append(employees, e)
	}

	return employees, nil
}

func GetEmployeeByName(name string) (*models.Employee, error) {
	row := db.QueryRow(
		"SELECT id, name, age, address, is_active FROM employees WHERE name = ?",
		name,
	)

	var e models.Employee
	var active bool

	if err := row.Scan(
		&e.ID,
		&e.Name,
		&e.Age,
		&e.Address,
		&active,
	); err != nil {
		return nil, err
	}

	e.IsActive = &active
	return &e, nil
}

func CreateEmployee(e models.Employee) error {
	_, err := db.Exec(
		"INSERT INTO employees (name, age, address, is_active) VALUES (?, ?, ?, ?)",
		e.Name,
		e.Age,
		e.Address,
		*e.IsActive,
	)
	return err
}

func UpdateEmployee(name string, e models.Employee) error {
	query := `
		UPDATE employees
		SET age = COALESCE(NULLIF(?, 0), age),
		    address = COALESCE(NULLIF(?, ''), address),
		    is_active = COALESCE(?, is_active)
		WHERE name = ?
	`
	result, err := db.Exec(
		query,
		e.Age,
		e.Address,
		e.IsActive,
		name,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteEmployee(name string) error {
	result, err := db.Exec(
		"DELETE FROM employees WHERE name = ?",
		name,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
