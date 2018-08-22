package pkg

import (
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"net/url"
	"fmt"
)

func parsePath(path string) []string {
	segments := strings.Split(path, "/")
	return segments[1:]
}

//our response struct is always used to return data to the client
//this keeps our api nice anc consistent
type Response struct {
	Status  int     `json:"status"`  //our api status code
	Message string  `json:"message"` //our api message
	Data    Records `json:"data"`    //if the db returns data it will be reflected here
	Created int64   `json:"created"` //if the db inserts data it will be reflected here
	Updated int64   `json:"updated"` //if the db updates data it will be reflected here
	Deleted int64   `json:"deleted"` //if the db deletes data it will be reflected here
	Links   []Link  `json:"links"`
}

//Write our response to the client
func (response *Response) Write(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response.Status = status
	enc := json.NewEncoder(w)
	enc.Encode(response)
}

type Link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method"`
}

//shortcut for returning just a message to our client
func MessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	m := Response{Message: message, Status: status}
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(&m)
}

//gets a parameter and casts it to an int or returns given default
//if an error occurs during casting it is returned
func intParamOrDefault(values url.Values, param string, def int) (int, error) {
	value := values.Get(param)
	if (value != "") {
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
		HandleGetMulti(w, r, storage)
		return
	}

	result, err := storage.Read(resource, &record, 0, 1)
	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {
		if len(result.Data) == 0 {
			MessageResponse(w, 404, "no records found")
		} else {
			result.Links = append(result.Links, Link{Rel: "self", Href: "http://" + r.Host + r.RequestURI, Method: "GET"})
			result.Write(w, 200)
		}

	}
}

func HandleGetMulti(w http.ResponseWriter, r *http.Request, storage Storage) {
	segments := parsePath(r.URL.Path)
	var record Record
	var resource Resource
	resource = Resource{segments[len(segments)-1]}
	params := r.URL.Query()

	offset, e := intParamOrDefault(params, "offset", 0)
	if e != nil {
		MessageResponse(w, 400, "improper value for 'offset'")
		return
	}

	limit, e := intParamOrDefault(params, "limit", Config().LimitDefault)
	if e != nil {
		MessageResponse(w, 400, "improper value for 'limit'")
		return
	}

	if limit < 1 {
		limit = 1
	}

	if offset < 0 {
		offset = 0
	}

	result, err := storage.Read(resource, &record, offset, limit)
	if err != nil {
		MessageResponse(w, err.Code, err.Message)
	} else {

		result.Links = append(result.Links, Link{"self", "http://" + r.Host + r.RequestURI, "GET"})

		//todo this could be smarter
		if offset > 0 {
			previousPageOffset := offset - limit
			if previousPageOffset < 0 {
				previousPageOffset = 0
			}
			link := Link{Rel: "prev", Method: "GET"}
			link.Href = fmt.Sprintf("http://%s%s?offset=%d&limit=%d", r.Host, r.URL.Path, previousPageOffset, limit)
			result.Links = append(result.Links, link)
		}

		//todo this could be smarter
		if len(result.Data) == limit {
			nextPageOffset := offset + limit
			link := Link{Rel: "next", Method: "GET"}
			link.Href = fmt.Sprintf("http://%s%s?offset=%d&limit=%d", r.Host, r.URL.Path, nextPageOffset, limit)
			result.Links = append(result.Links, link)
		}
		result.Write(w, 200)
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
		result.Write(w, status)

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
		result.Write(w, 200)
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
			result.Write(w, 200)
		}
	}
}



func checkPermission(r *http.Request) bool {
	switch r.Method {
	case "GET":
		return Config().GetPermissions["global"] == "allow"
	case "PUT":
		return Config().PutPermissions["global"] == "allow"
	case "POST":
		return Config().PostPermissions["global"] == "allow"
	case "DELETE":
		return Config().DeletePermissions["global"] == "allow"
	default:
		return true
	}
}

func Handler(w http.ResponseWriter, r *http.Request, storage Storage) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	canContinue := checkPermission(r)
	if !canContinue {
		MessageResponse(w, 401, "Permission denied")
		return
	}

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
