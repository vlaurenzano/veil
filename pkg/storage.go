package pkg

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/go-sql-driver/mysql"
)

type StorageError struct {
	Code       int    //this code will follow expected http status code conventions
	Message    string //this string should be appropriate for end user messages
	WrapsError error  //the error that prompted this error, for instance a mysql error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

type Storage interface {
	Create(resource Resource, record Record) (sql.Result, *StorageError)
	Read(resource Resource) (Records, *StorageError)
	Update(resource Resource, record Record) (sql.Result, *StorageError)
	Delete(resource Resource, record Record) (sql.Result, *StorageError)
}

type Resource struct {
	Identifier string //the identifier of the resource, for instance a mysql table name
}

type Record map[string]interface{} //the data held in the resource

type Records []Record //a list of records

type MySqlStorage struct {
	ConnectionString string
}

func NewStorage() (Storage, *StorageError) {
	if c := Config(); c.DB == "MYSQL" {
		s := MySqlStorage{c.ConnectionString}
		return &s, nil
	}
	return nil, &StorageError{}
}

func (m *MySqlStorage) Create(resource Resource, record Record) (sql.Result, *StorageError) {
	db, _ := sql.Open("mysql", m.ConnectionString)
	defer db.Close()

	var keys, sss []string
	var values []interface{}

	for k, v := range record {
		keys = append(keys, k)
		values = append(values, v)
		sss = append(sss, "?")
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", resource.Identifier, strings.Join(keys, ","), strings.Join(sss, ","))
	stmt, err := db.Prepare(sql)

	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return nil, &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			default:
				return nil, &StorageError{Code: 400, Message: "unknown error inserting resource", WrapsError: err}
			}

		default:
			return nil, &StorageError{Code: 400, Message: "unknown error inserting resource", WrapsError: err}
		}
	}

	result, err := stmt.Exec(values...)

	if err != nil {

		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return nil, &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			case 1364:
				return nil, &StorageError{Code: 400, Message: "resource does not include all required values", WrapsError: err}
			default:
				return nil, &StorageError{Code: 400, Message: "unknown error inserting resource", WrapsError: err}
			}

		default:
			return nil, &StorageError{Code: 400, Message: "unknown error inserting resource", WrapsError: err}
		}

	}

	return result, nil
}

func (m *MySqlStorage) Read(resource Resource) (Records, *StorageError) {

	sqlString := fmt.Sprintf("SELECT * FROM %s", resource.Identifier)

	db, err := sql.Open("mysql", m.ConnectionString)
	defer db.Close()

	stmt, err := db.Prepare(sqlString)

	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return nil, &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			}
		default:
			return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
		}
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
	}

	tableData := Records{}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
		}

		entry := Record{}
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if (ok) {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}

		tableData = append(tableData, entry)
	}

	return tableData, nil
}

func (m *MySqlStorage) Update(resource Resource, record Record) (sql.Result, *StorageError) {

	db, _ := sql.Open("mysql", m.ConnectionString)
	defer db.Close()

	var sss []string
	var values []interface{}
	for k, v := range record {
		if k == "id" {
			continue
		}
		values = append(values, v)
		sss = append(sss, k+"=?")
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE id=?;", resource.Identifier, strings.Join(sss, ","))
	stmt, err := db.Prepare(sql)

	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return nil, &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			}
		default:
			return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
		}
	}

	result, err := stmt.Exec(append(values, record["id"])...)

	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return nil, &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			}
		default:
			return nil, &StorageError{Code: 400, Message: "unknown error retrieving resource", WrapsError: err}
		}
	}

	return result, nil
}

func (m *MySqlStorage) Delete(resource Resource, record Record) (sql.Result, *StorageError) {

	db, _ := sql.Open("mysql", m.ConnectionString)
	defer db.Close()

	sql := fmt.Sprintf("DELETE FROM %s WHERE id=?;", resource.Identifier)
	stmt, err := db.Prepare(sql)

	if err != nil {
		panic(err.Error())
	}

	result, err := stmt.Exec(record["id"])
	if err != nil {
		panic(err.Error())
	}
	return result, nil
}
