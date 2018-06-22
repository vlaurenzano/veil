# VEIL

"A web application is a veil around a database" - Abraham Lincoln 

Veil is a proof of concept, automatic json backend for mysql databases. It works by translating urls to database calls, implementing CRUD functionality on the table named by the url. 


## Installation

```got get github.com/vlaurenzano/veil```

```go get github.com/go-sql-driver/mysql```

## About

This application can replace 80% of the work I've ever done as an API developer.

## Create -- PUT

Put requests place a record in the resource determined by url. Presently PUT does not upsert, nor does it allow for inserting known ids.

```
curl -i -X PUT -H "Content-Type:application/json" http://localhost:8080/test_resource -d '{"test_field_1":"123", "test_field_2":"123"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:52:23 GMT
Content-Length: 70

{"Data":null,"Created":1,"Updated":0,"Deleted":0,"message":"success"}

```

## Read -- GET
Get requests support retrival of all records, limits, offsets, and retrival by id

```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource          
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:53:05 GMT
Content-Length: 213

{"Data":[{"id":1,"test_field_1":"123","test_field_2":"123"},{"id":2,"test_field_1":"123","test_field_2":"123"},{"id":3,"test_field_1":"123","test_field_2":"123"}],"Created":0,"Updated":0,"Deleted":0,"message":""}

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource/1          
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:54:24 GMT
Content-Length: 111

{"Data":[{"id":1,"test_field_1":"123","test_field_2":"123"}],"Created":0,"Updated":0,"Deleted":0,"message":""}

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource?limit=1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:59:59 GMT
Content-Length: 111

{"Data":[{"id":1,"test_field_1":"123","test_field_2":"123"}],"Created":0,"Updated":0,"Deleted":0,"message":""}

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource?limit=1\&offset=1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 23:01:33 GMT
Content-Length: 111

{"Data":[{"id":2,"test_field_1":"123","test_field_2":"123"}],"Created":0,"Updated":0,"Deleted":0,"message":""}


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
