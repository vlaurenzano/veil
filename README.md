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
curl -i -X PUT -H "Content-Type:application/json" http://localhost:8080/test_resource -d '{"test_field_1":"123", "test_field_2":"123"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:52:23 GMT
Content-Length: 70

{"Data":null,"Created":1,"Updated":0,"Deleted":0,"message":"success"}

```

### Read -- GET

Get requests support retrieval of resources of records through sensible default urls. All responses are uniformly formatted and include HATEOAS meta data if applicable.  
 
The types of retrievals supported include:
 
#### All Records

Veil will retrieve all records, up to the configuration limits, by querying for a resource without parameters. 
```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource          
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:53:05 GMT
Content-Length: 213

{
    "Data": [
            {"id":1,"test_field_1":"123","test_field_2":"123"},
            {"id":2,"test_field_1":"123","test_field_2":"123"},
            {"id":3,"test_field_1":"123","test_field_2":"123"}
    ],
    "Created":0,
    "Updated":0,
    "Deleted":0,
    "message":""
}
```

#### By ID
Veil can retrieve a singular resource by id. Note: for consistency with retrieving multiple records the Data field returns an array with 1 item. 

```
curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource/1          
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:54:24 GMT
Content-Length: 111

{
    "Data": [
            {"id":1,"test_field_1":"123","test_field_2":"123"}
     ],
    "Created":0,
    "Updated":0,
    "Deleted":0,
    "message":""
}

### More examples

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource?limit=1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 22:59:59 GMT
Content-Length: 111

{
    "Data":[
        {"id":1,"test_field_1":"123","test_field_2":"123"}
    ],
    "Created":0,
    "Updated":0,
    "Deleted":0,
    "message":""
}

curl -i -X GET -H "Content-Type:application/json" http://localhost:8080/test_resource?limit=1\&offset=1
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 22 Jun 2018 23:01:33 GMT
Content-Length: 111

{
    "Data":[
        {"id":2,"test_field_1":"123","test_field_2":"123"}
    ],
    "Created":0,
    "Updated":0,
    "Deleted":0,
    "message":""
}

```
#### Update

```
curl -i -X POST -H "Content-Type:application/json" http://localhost:8080/test_resource/1 -d '{"test_field_1":"321", "test_field_2":"321"}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:06:23 GMT
Content-Length: 22

{"message":"success"}
```

#### Delete

```
curl -i -X DELETE -H "Content-Type:application/json" http://localhost:8080/test_resource/1                                            
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 17 Jun 2018 22:06:58 GMT
Content-Length: 22

{"message":"success"}
```
