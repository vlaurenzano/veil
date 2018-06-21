package pkg

import (
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
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

func HandleGet(w http.ResponseWriter, r *http.Request, storage Storage) {

	segments := parsePath(r.URL.Path)

	//params := r.URL.Query()
	//_ = params.Get("limit")

	resource := Resource{segments[len(segments)-1]}

	result, err := storage.Read(resource)

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
