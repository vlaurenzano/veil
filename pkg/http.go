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

type Message struct {
	Message string `json:"message"`
}


func MessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	m := Message{Message: message}
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

	resource := Resource{segments[len(segments)-1]}

	result, err := storage.Read(resource, offset, limit)

	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		ObjectResponse(w, 200, result.Data)
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
		_ = result //todo check if something was inserted
		MessageResponse(w, 201, "success")
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
		_ = result
		MessageResponse(w, 200, "success")
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
		_ = result //todo check if something was inserted
		MessageResponse(w, 200, "success")
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
