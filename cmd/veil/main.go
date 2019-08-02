package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vlaurenzano/veil/pkg"
	"net/http"
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
	logrus.Info("Info: Starting veil server on port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}