package pkg

import (
	"net/http"
	"strings"
	_ "fmt"
	_ "encoding/json"
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

func respond(w http.ResponseWriter, status int, message string) {
	m := Message{Message: message}
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(&m)
}

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	storage, err := NewStorage()

	if err != nil {
		respond(w, 400, "an error occurred")
	}

	switch r.Method {

	case "GET":

		segments := parsePath(r.URL.Path)
		params := r.URL.Query()
		test := params.Get("limit")
		_ = test

		resource := Resource{segments[len(segments)-1]}

		result, err := storage.Read(resource)

		if err != nil {
			respond(w, err.Code, err.Message)
			return
		} else {
			w.WriteHeader(200)
			enc := json.NewEncoder(w)
			enc.Encode(result.Data)
		}

	case "PUT":
		segments := parsePath(r.URL.Path)
		b, _ := ioutil.ReadAll(r.Body)
		r := Record{}
		e := json.Unmarshal(b, &r)
		if e != nil {
			respond(w, 400, "payload could not be parsed")
			return
		}
		resource := Resource{segments[len(segments)-1]}
		result, err := storage.Create(resource, r)

		if err != nil {
			respond(w, err.Code, err.Message)
			return
		}

		_ = result //todo check if something was inserted
		w.WriteHeader(201)
		m := Message{"success"}
		enc := json.NewEncoder(w)
		enc.Encode(m)

	case "POST":
		segments := parsePath(r.URL.Path)
		b, _ := ioutil.ReadAll(r.Body)
		r := Record{}
		e := json.Unmarshal(b, &r)
		if e != nil {
			respond(w, 400, "payload could not be parsed")
			return
		}

		r["id"] = segments[len(segments)-1]
		resource := Resource{segments[len(segments)-2]}
		result, err := storage.Update(resource, r)
		if err != nil {
			respond(w, err.Code, err.Message)
			return
		}

		_ = result //todo check if something was inserted
		w.WriteHeader(200)
		m := Message{"success"}
		enc := json.NewEncoder(w)
		enc.Encode(m)

	case "DELETE":
		segments := parsePath(r.URL.Path)
		r := Record{}
		r["id"] = segments[len(segments)-1]
		resource := Resource{segments[len(segments)-2]}
		result, err := storage.Delete(resource, r)
		if err != nil {
			respond(w, err.Code, err.Message)
			return
		}

		_ = result //todo check if something was inserted
		w.WriteHeader(200)
		m := Message{"success"}
		enc := json.NewEncoder(w)
		enc.Encode(m)

	case "OPTIONS":
		respond(w,200,"")
	default:
		respond(w,400,"Unsupported method")

	}

}
