package pkg

import (
	"fmt"
)

//provides a typed error interface for the storage interface methods to return
//methods returning this interface should return a pointer so clients can check against nil
type StorageError struct {
	Code       int    //this code will follow expected http status code conventions
	Message    string //this string should be appropriate for end user messages
	WrapsError error  //the error that prompted this error, for instance a mysql error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

//provides an abstraction for the database layer
type Storage interface {
	Create(resource Resource, record Record) (*Result, *StorageError) //Creates an entry in the data store, intended for use with PUT
	Read(resource Resource, match *Record, offset int, limit int) (*Result, *StorageError) //Reads from the data store, intended for use with GET
	Update(resource Resource, record Record) (*Result, *StorageError) //Updates a record in the data store
	Delete(resource Resource, record Record) (*Result, *StorageError) //Deletes a record in the data store
}

//a resource represents the table or document within the database
type Resource struct {
	Identifier string //the identifier of the resource, for instance a mysql table name
}

//a record in the database
type Record map[string]interface{} //the data held in the resource

// a collection of records
type Records []Record

//Result abstracts the database result
type Result struct {
	Data    Records //if the db returns data it will be reflected here
	Created int64     //if the db inserts data it will be reflected here
	Updated int64     //if the db updates data it will be reflected here
	Deleted int64     //if the db deletes data it will be reflected here
}

//our storage factory
func NewStorage() (Storage, *StorageError) {
	if c := Config(); c.DB == "MYSQL" {
		s := MySqlStorage{c.ConnectionString}
		return &s, nil
	}
	return nil, &StorageError{}
}
