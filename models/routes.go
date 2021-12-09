package models

import (
	"github.com/gorilla/mux"
)

func InitAllroutes(router *mux.Router) {
	initContainerRoutes(router)
	initHostRoutes(router)
}

func initContainerRoutes(router *mux.Router) {

	router.HandleFunc("/container", AllContainers).Methods("GET")
	router.HandleFunc("/container/{containerId}", getContainerById)
	router.HandleFunc("/container-for-spec-host/{hostId}", getContainerByHostId)
	router.HandleFunc("/container", insertContainer).Methods("POST")
}

func initHostRoutes(router *mux.Router) {

	router.HandleFunc("/host", AllHosts)
	router.HandleFunc("/host/{hostId}", getHostById)
}
