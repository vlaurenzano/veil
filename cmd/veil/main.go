package main

import (
	"net/http"
	"log"
	"github.com/vlaurenzano/veil/pkg"
)

func main(){
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		storage, err := pkg.NewStorage()
		if err != nil {
			pkg.MessageResponse(writer, 500, "an error occurred connecting to the database")
			return
		}
		pkg.Handler(writer, request, storage)
	} )
	log.Fatal(http.ListenAndServe(":8080", nil))
}