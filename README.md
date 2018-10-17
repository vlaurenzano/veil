# VEIL

"A web application is a veil around a database" - Abraham Lincoln 

Veil is a proof of concept, automatic json backend server. It works by translating urls to database queries. 

Veil's aim is to minimize the repetative `REST query -> database -> model -> tranform -> client` development cycle. It uses straightforward conventions that are easy to learn and implement.
 
Veil isn't a singular solution, it won't replace your entire stack. It doesn't have things like authentication and TLS out of the box and is designed to sit behind a reverse proxy such as nginx.

  
## Installation
Install veil with go:

```go get github.com/vlaurenzano/veil/cmd/veil```

Then install the mysql driver:

```go get github.com/go-sql-driver/mysql```

## Requests

Veil handles CRUD via RESTFUL endpoints out of the box. 

   

### Create -- PUT

Put requests place a record in the resource determined by url. Presently PUT does not upsert, nor does it allow for inserting known ids.

```
HTTP/1.1 201 Created
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:42:24 GMT
Content-Length: 96

{"status":201,"message":"success","data":null,"created":1,"updated":0,"deleted":0,"links":null}
```

### Read -- GET

Get requests support retrieval of resources of records through sensible default urls. All responses are uniformly formatted and include HATEOAS meta data if applicable.  
 
The types of retrievals supported include:
 
#### All Records

Veil will retrieve all records, up to the configuration limits, by querying for a resource without parameters. 
```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource          
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:42:57 GMT
Content-Length: 413

{"status":200,"message":"","data":[{"id":1,"test_field_1":"123","test_field_2":"123"},{"id":2,"test_field_1":"123","test_field_2":"123"},{"id":3,"test_field_1":"123","test_field_2":"123"},{"id":4,"test_field_1":"123","test_field_2":"123"},{"id":5,"test_field_1":"123","test_field_2":"123"}],"created":0,"updated":0,"deleted":0,"links":[{"rel":"self","href":"http://localhost:8080/test_resource","method":"GET"}]}

```

#### By ID
Veil can retrieve a singular resource by id. Note: for consistency with retrieving multiple records the Data field returns an array with 1 item. 

```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource/1          
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:44:25 GMT
Content-Length: 211

{"status":200,"message":"","data":[{"id":1,"test_field_1":"123","test_field_2":"123"}],"created":0,"updated":0,"deleted":0,"links":[{"rel":"self","href":"http://localhost:8080/test_resource/1","method":"GET"}]}


### More examples

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource?limit=1
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:45:09 GMT
Content-Length: 314

{"status":200,"message":"","data":[{"id":1,"test_field_1":"123","test_field_2":"123"}],"created":0,"updated":0,"deleted":0,"links":[{"rel":"self","href":"http://localhost:8080/test_resource?limit=1","method":"GET"},{"rel":"next","href":"http://localhost:8080/test_resource?offset=1\u0026limit=1","method":"GET"}]}

HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:45:37 GMT
Content-Length: 425

{"status":200,"message":"","data":[{"id":2,"test_field_1":"123","test_field_2":"123"}],"created":0,"updated":0,"deleted":0,"links":[{"rel":"self","href":"http://localhost:8080/test_resource?limit=1\u0026offset=1","method":"GET"},{"rel":"prev","href":"http://localhost:8080/test_resource?offset=0\u0026limit=1","method":"GET"},{"rel":"next","href":"http://localhost:8080/test_resource?offset=2\u0026limit=1","method":"GET"}]}

```
#### Update

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:46:16 GMT
Content-Length: 89

{"status":200,"message":"","data":null,"created":0,"updated":1,"deleted":0,"links":null}

```

#### Delete

```
curl -i -X DELETE -H "Content-Type:application/json" http://localhost:8080/test_resource/1                                            
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Wed, 17 Oct 2018 01:46:32 GMT
Content-Length: 89

{"status":200,"message":"","data":null,"created":0,"updated":0,"deleted":1,"links":null}

```
