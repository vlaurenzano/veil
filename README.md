# VEIL

"A web application is a veil around a database" - Abraham Lincoln 

Veil is a proof of concept, automatic json backend for mysql databases. It works by translating urls to database calls, implementing CRUD functionality on the table named by the url. 

## Installation

```got get github.com/vlaurenzano/veil```

```go get github.com/go-sql-driver/mysql```

## Create

```
curl -i -X PUT -H "Content-Type:application/json" http://localhost:8080/test_resource -d '{"test_field_1":"123", "test_field_2":"123"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:02:08 GMT
Content-Length: 22

{"message":"success"}
```

## Read

```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource           
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:05:02 GMT
Content-Length: 53

[{"id":1,"test_field_1":"123","test_field_2":"123"}]
```
## Update
```
curl -i -X POST -H "Content-Type:application/json" http://localhost:8080/test_resource/1 -d '{"test_field_1":"321", "test_field_2":"321"}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:06:23 GMT
Content-Length: 22

{"message":"success"}
```

## Delete
```
curl -i -X DELETE -H "Content-Type:application/json" http://localhost:8080/test_resource/1                                            
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:06:58 GMT
Content-Length: 22

{"message":"success"}
```
