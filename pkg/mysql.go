package pkg

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/go-sql-driver/mysql"
	"regexp"
)

type MySqlStorage struct {
	ConnectionString string
}

func (m *MySqlStorage) dbConnect() (db *sql.DB, err *StorageError) {
	db, e := sql.Open("mysql", m.ConnectionString)
	err = interpretMysqlError(e)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//interprets a mysql error and returns it as a Storage Error
func interpretMysqlError(err error) (*StorageError) {
	if err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			switch err.(*mysql.MySQLError).Number {
			case 1146:
				return &StorageError{Code: 404, Message: "resource not found", WrapsError: err}
			case 1364:
				return &StorageError{Code: 400, Message: "resource does not include all required values", WrapsError: err}
			default:
				return &StorageError{Code: 400, Message: "unknown error", WrapsError: err}
			}

		default:
			return &StorageError{Code: 400, Message: "unknown error", WrapsError: err}
		}
	}
	return nil
}

func (m *MySqlStorage) Create(resource Resource, record Record) (*Response, *StorageError) {
	db, err := m.dbConnect()
	if err != nil {
		return nil, err
	}
	var keys, sss []string
	var values []interface{}

	for k, v := range record {
		keys = append(keys, k)
		values = append(values, v)
		sss = append(sss, "?")
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", resource.Identifier, strings.Join(keys, ","), strings.Join(sss, ","))
	stmt, e := db.Prepare(sql)
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	_, e = stmt.Exec(values...)

	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	result := Response{Created: 1}
	return &result, nil
}

func (m *MySqlStorage) Read(resource Resource, match *Record, offset int, limit int) (result *Response, err *StorageError) {


	//todo: improve this validation, we can't parameterize the dynamic table name
	notOk, e := regexp.Match("[;'\"\\\\]",[]byte(resource.Identifier))
	if notOk || e != nil {
		return nil, &StorageError{Code: 404, Message: "resource not found"}
	}

	sqlString := fmt.Sprintf("SELECT * FROM %s ", resource.Identifier)

	var paramValues []interface{}
	var paramKeys []string
	if *match != nil {
		sqlString += "WHERE "
		for k, v := range *match {
			paramValues = append(paramValues, v)
			paramKeys = append(paramKeys, k + " = ?")
		}
		sqlString += strings.Join(paramKeys, " AND ")
	}

	paramValues = append(paramValues, offset)
	paramValues = append(paramValues, limit)

	sqlString += " LIMIT ?, ?"

	db, err := m.dbConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, e := db.Prepare(sqlString)
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, e := stmt.Query(paramValues...)

	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, e := rows.Columns()
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	tableData := Records{}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		e := rows.Scan(scanArgs...)
		if err = interpretMysqlError(e); err != nil {
			return nil, err
		}

		entry := Record{}
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}

		tableData = append(tableData, entry)
	}
	return &Response{Data: tableData}, nil
}

func (m *MySqlStorage) Update(resource Resource, record Record) (*Response, *StorageError) {

	db, err := m.dbConnect()
	if err != nil {
		return nil, err
	}

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

	stmt, e := db.Prepare(sql)
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	r, e := stmt.Exec(append(values, record["id"])...)
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	rows, e := r.RowsAffected()
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	return &Response{Updated:rows}, nil
}

func (m *MySqlStorage) Delete(resource Resource, record Record) (*Response, *StorageError) {

	db, err := m.dbConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := fmt.Sprintf("DELETE FROM %s WHERE id=?;", resource.Identifier)
	stmt, e := db.Prepare(sql)

	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}

	r, e := stmt.Exec(record["id"])
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}
	rows, e := r.RowsAffected()
	if err = interpretMysqlError(e); err != nil {
		return nil, err
	}
	return &Response{Deleted:rows}, nil
}
