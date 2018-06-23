package pkg

import (
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"net/url"
)

func parsePath(path string) []string {
	segments := strings.Split(path, "/")
	return segments[1:]
}

type Response struct {
	Data    Records   //if the db returns data it will be reflected here
	Created int64     //if the db inserts data it will be reflected here
	Updated int64     //if the db updates data it will be reflected here
	Deleted int64     //if the db deletes data it will be reflected here
	Message string `json:"message"`
}


func MessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	m := Response{Message: message}
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(&m)
}

func ObjectResponse(w http.ResponseWriter, status int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(object)
}

//gets a parameter and casts it to an int or returns given default
//if an error occurs during casting it is returned
func intParamOrDefault(values url.Values, param string, def int) (int,  error){
	value := values.Get(param)
	if(value != ""){
		intVal, e := strconv.Atoi(value)
		if e != nil {
			return def, e
		}
		return intVal, nil
	}
	return def, nil
}


//handle a get call
//GET /resource -- returns records up to default limit at the default offset
//GET /resource?limit=x -- returns records up to given limit at the default offset
//GET /resource?offset=x&limit=y -- return records up to limit from given offset
//GET /resource/id -- gets the resource at the given id
func HandleGet(w http.ResponseWriter, r *http.Request, storage Storage) {


	segments := parsePath(r.URL.Path)

	var record Record
	var resource Resource

	if len(segments) == 2 {
		record = Record{"id": segments[1]}
		resource = Resource{segments[len(segments)-2]}
	} else {
		resource = Resource{segments[len(segments)-1]}
	}

	params := r.URL.Query()

	offset, e := intParamOrDefault(params, "offset",0)
	if e != nil {
		MessageResponse(w, 400,"improper value for 'offset'")
		return
	}

	limit, e := intParamOrDefault(params, "limit", Config().LimitDefault)
	if e != nil {
		MessageResponse(w, 400,"improper value for 'limit'")
		return
	}


	result, err := storage.Read(resource, &record, offset, limit)
	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		if len(result.Data) == 0 {
			MessageResponse(w, 404, "no records found")
		} else {
			ObjectResponse(w, 200, result)
		}

	}
}

func HandlePut(w http.ResponseWriter, r *http.Request, storage Storage) {
	segments := parsePath(r.URL.Path)
	b, _ := ioutil.ReadAll(r.Body)

	record := Record{}
	e := json.Unmarshal(b, &record)
	if e != nil {
		MessageResponse(w, 400, "payload could not be parsed")
		return
	}

	resource := Resource{segments[len(segments)-1]}
	result, err := storage.Create(resource, record)

	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		status := 200
		if result.Created != 0 {
			status = 201
			result.Message = "success"
		}
		ObjectResponse(w,status,result)
	}

}

func HandlePost(w http.ResponseWriter, r *http.Request, storage Storage) {
	segments := parsePath(r.URL.Path)
	b, _ := ioutil.ReadAll(r.Body)
	record := Record{}
	e := json.Unmarshal(b, &record)
	if e != nil {
		MessageResponse(w, 400, "payload could not be parsed")
		return
	}

	record["id"] = segments[len(segments)-1]
	resource := Resource{segments[len(segments)-2]}
	result, err := storage.Update(resource, record)
	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		ObjectResponse(w,200,result)
	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request, storage Storage) {
	segments := parsePath(r.URL.Path)
	record := Record{}
	record["id"] = segments[len(segments)-1]
	resource := Resource{segments[len(segments)-2]}
	result, err := storage.Delete(resource, record)
	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		if result.Deleted == 0 {
			MessageResponse(w, 404, "record not found")
		} else {
			ObjectResponse(w,200,result)
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request, storage Storage) {

	switch r.Method {

	case "GET":
		HandleGet(w, r, storage)

	case "PUT":
		HandlePut(w, r, storage)

	case "POST":
		HandlePost(w, r, storage)

	case "DELETE":
		HandleDelete(w, r, storage)

	case "OPTIONS":
		MessageResponse(w, 200, "")
	default:
		MessageResponse(w, 400, "Unsupported method")

	}

}
