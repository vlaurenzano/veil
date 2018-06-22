package pkg

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
	_ "reflect"
	"strings"
	"database/sql"
	"fmt"
)

func TestUrlParse(t *testing.T) {

	p := parsePath("/")
	if len(p) != 1 {
		t.Errorf("unexpected result parsing url")
	}
	p = parsePath("/lala")
	if len(p) != 1 && p[0] != "" {
		t.Errorf("unexpected result parsing url")
	}
	p = parsePath("/la/la")
	if len(p) != 2 {
		t.Errorf("unexpected result parsing url")
	}
}

func request(method string, url string, data string) (*http.Response) {
	request, err := http.NewRequest(method, url, strings.NewReader(data))
	check(err)
	response, err := http.DefaultClient.Do(request)
	check(err)
	return response
}


func check(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func initTestTable() {

	db, err := sql.Open("mysql", Config().ConnectionString)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS `veil_test_resource`;")
	check(err)

	_, err = db.Exec(`CREATE TABLE veil_test_resource (		
		id int(11) NOT NULL AUTO_INCREMENT,
		test_field_1 varchar(255) NOT NULL,
		test_field_2 varchar(255) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
	check(err)
}

func addXRows(x int) {
	db, err := sql.Open("mysql", Config().ConnectionString)
	defer db.Close()
	for i := 0; i < x; i++ {

		_, err = db.Exec(fmt.Sprintf("INSERT INTO veil_test_resource (test_field_1, test_field_2) VALUES ('test value %d', 'test value %d');", i, i))
		check(err)
	}
}

func testHandlerFunc(w http.ResponseWriter, r *http.Request){
	storage, err := NewStorage()
	if err != nil {
		log.Fatal(fmt.Sprint("Error connecting to database"))
	}
	Handler(w, r, storage)
}

func TestAppHandleGET(t *testing.T) {

	initTestTable()
	addXRows(2)

	ts := httptest.NewServer(http.HandlerFunc(testHandlerFunc))
	defer ts.Close()


	res := request("GET", ts.URL+"/veil_test_resource", "")

	j := loadResponseBody(res)

	if res.StatusCode != 200 || len(j.Data) != 2 {
		log.Fatal("Test App Handler GET did not return the right amount of records")
	}

	res = request("GET", ts.URL+"/veil_test_not_exist", "")

	if res.StatusCode != 404 {
		log.Fatal(fmt.Sprint("GET expected a 404, got ", res.StatusCode))
	}

	res = request("GET", ts.URL+"/veil_test_resource/1", "")

	j = loadResponseBody(res)


	if res.StatusCode != 200 || len(j.Data) != 1 {
		log.Fatal(fmt.Sprint("Tried receiving record by id, received bad response ", res.StatusCode))
	}

	res = request("GET", ts.URL+"/veil_test_resource/11", "")

	if res.StatusCode != 404 {
		log.Fatal(fmt.Sprint("Tried receiving an uknown record by id, expected 404 received ", res.StatusCode))
	}


}

func loadResponseBody(res *http.Response) Response{
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check(err)

	var j Response

	err = json.Unmarshal(body, &j)
	check(err)
	return j
}

func TestAppHandlerGETWithLimitsAndOffsets(t *testing.T){
	initTestTable()
	addXRows(5)
	ts := httptest.NewServer(http.HandlerFunc(testHandlerFunc))
	defer ts.Close()

	res := request("GET", ts.URL+"/veil_test_resource?limit=3", "")
	j := loadResponseBody(res)

	if len(j.Data) != 3 {
		log.Fatal("Test App Handler GET did not return the right amount of records")
	}

	res = request("GET", ts.URL+"/veil_test_resource?offset=4&limit=3", "")
	j = loadResponseBody(res)

	if len(j.Data) != 1 {
		log.Fatal("Test App Handler GET did not return the right amount of records")
	}

}


func TestAppHandlePUT(t *testing.T) {

	initTestTable()
	addXRows(2)

	ts := httptest.NewServer(http.HandlerFunc(testHandlerFunc))
	defer ts.Close()

	res := request("PUT", ts.URL+"/veil_test_resource", "{\"test_field_1\":\"fgfg\", \"test_field_2\":\"fgfg\"}")

	if res.StatusCode != 201 {
		log.Fatal(fmt.Sprint("PUT expected a 201, got ", res.StatusCode))
	}

	res = request("PUT", ts.URL+"/veil_test_not_exist", "{\"test_field_1\":\"t\"}")
	if res.StatusCode != 404 {
		log.Fatal(fmt.Sprint("GET expected a 404, got ", res.StatusCode))
	}

	res = request("PUT", ts.URL+"/veil_test_resource", "{\"test_field_1\":\"t\"}")

	if res.StatusCode != 400 {
		log.Fatal(fmt.Sprint("GET expected a 400, got ", res.StatusCode))
	}
}

func TestAppHandlePOST(t *testing.T) {

	initTestTable()
	addXRows(2)

	ts := httptest.NewServer(http.HandlerFunc(testHandlerFunc))
	defer ts.Close()


	res := request("POST", ts.URL+"/veil_test_resource/1", "{\"test_field_1\":\"123\", \"test_field_2\":\"123\"}")

	if res.StatusCode != 200 {
		log.Fatal(fmt.Sprint("POST expected a 200, got ", res.StatusCode))
	}

	res = request("GET", ts.URL+"/veil_test_resource", "")

	j := loadResponseBody(res)

	if j.Data[0]["test_field_1"] != "123" {
		log.Fatal("Record not properly updated")
	}


	res = request("POST", ts.URL+"/veil_test_not_exist/1", "{\"test_field_1\":\"t\"}")
	if res.StatusCode != 404 {
		log.Fatal(fmt.Sprint("GET expected a 404, got ", res.StatusCode))
	}

	//res = request("POST", ts.URL+"/veil_test_resource/3", "{\"test_field_1\":\"t\"}")
	//if res.StatusCode != 404 {
	//	log.Fatal(fmt.Sprint("GET expected a 400, got ", res.StatusCode))
	//}
}


func TestAppHandleDELETE(t *testing.T) {

	initTestTable()
	addXRows(2)

	ts := httptest.NewServer(http.HandlerFunc(testHandlerFunc))
	defer ts.Close()


	res := request("DELETE", ts.URL+"/veil_test_resource/1", "{\"test_field_1\":\"123\", \"test_field_2\":\"123\"}")

	if res.StatusCode != 200 {
		log.Fatal(fmt.Sprint("POST expected a 200, got ", res.StatusCode))
	}



	//body, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//
	//check(err)
	//
	//var j Records
	//
	//err = json.Unmarshal(body, &j)
	//check(err)
	//
	//if j[0]["test_field_1"] != "123" {
	//	log.Fatal("Record not properly updated")
	//}
	//
	//
	//res = request("POST", ts.URL+"/veil_test_not_exist/1", "{\"test_field_1\":\"t\"}")
	//if res.StatusCode != 404 {
	//	log.Fatal(fmt.Sprint("GET expected a 404, got ", res.StatusCode))
	//}

}



