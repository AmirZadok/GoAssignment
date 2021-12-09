package main

import (
	"net/http"

	"github.com/amirzadok/RestApi/models"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	router := mux.NewRouter()
	models.InitAllroutes(router)
	models.InitDBConnection("./aqua.db")
	http.ListenAndServe(":80", router)

}
