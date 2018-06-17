package main

import (
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vlaurenzano/veil/pkg"
)

func main(){
	http.HandleFunc("/", pkg.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}