package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var db *sql.DB

func InitDBConnection(dbLocation string) error {
	var err error
	db, err = sql.Open("sqlite3", dbLocation)
	if err != nil {
		log.Fatal(err)
		//return err
	}

	return db.Ping()
}

func GetAllhHosts() ([]Host, error) {
	rows, err := db.Query("SELECT * FROM hosts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var hosts []Host
	for rows.Next() {
		var tempHost Host
		err := rows.Scan(&tempHost.Id, &tempHost.Uuid, &tempHost.Name, &tempHost.Ip_Address)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, tempHost)
	}

	return hosts, nil
}

func GetAllhContainers() ([]Container, error) {
	rows, err := db.Query("SELECT * FROM containers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []Container
	for rows.Next() {
		var tempContainer Container
		err := rows.Scan(&tempContainer.Id, &tempContainer.Host_Id, &tempContainer.Name, &tempContainer.Image_Name)
		if err != nil {
			return nil, err
		}
		containers = append(containers, tempContainer)
	}

	return containers, nil
}

func insertContainer(w http.ResponseWriter, r *http.Request) {
	var partialContainer PartialContainer
	err := json.NewDecoder(r.Body).Decode(&partialContainer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hosts := getHostByIdDb(partialContainer.Host_Id, w)
	if len(hosts) == 0 {
		infoMesage := InfoMessage{"Not a valid host_id , insert cancelled "}
		res, err := PrettyStruct(infoMesage)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprint(w, res)
		return

	}
	stmt, err := db.Prepare("INSERT INTO containers ( host_id , name , image_name) VALUES (?,?,?)")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	_, err = stmt.Exec(partialContainer.Host_Id, uuid.New().String(), partialContainer.Image_Name)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	printMessage(w, "the container added")

}

func getContainerByHostId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostId := vars["hostId"]
	rows, err := db.Query("SELECT containers.id , containers.host_id , containers.name , containers.image_name , hosts.name FROM containers LEFT JOIN hosts ON containers.host_id = hosts.id where  containers.host_id = ?", hostId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var containerWithHostName []ContainerWithHostName
	if !rows.Next() {
		printMessage(w, "no such container")
		return
	} else {

		var tempContainer ContainerWithHostName
		err := rows.Scan(&tempContainer.Id, &tempContainer.Host_Id, &tempContainer.Name, &tempContainer.Image_Name, &tempContainer.Host_Name)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		containerWithHostName = append(containerWithHostName, tempContainer)
	}

	res, err := PrettyStruct(containerWithHostName)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, res)
}

func printMessage(w http.ResponseWriter, message string) {
	infoMesage := InfoMessage{message}
	res, err := PrettyStruct(infoMesage)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, res)
}

func getContainerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerId := vars["containerId"]
	rows, err := db.Query("SELECT * FROM containers WHERE id = ?", containerId)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		printMessage(w, "No such container!!!")
		return
	}
	for rows.Next() {
		var tempContainer Container
		err := rows.Scan(&tempContainer.Id, &tempContainer.Host_Id, &tempContainer.Name, &tempContainer.Image_Name)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		res, err := PrettyStruct(tempContainer)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, res)
	}
}

func getHostByIdDb(hostId int, w http.ResponseWriter) []Host {
	rows, err := db.Query("SELECT * FROM hosts WHERE id = ?", hostId)
	var hosts []Host
	if err != nil {
		return hosts
	}
	defer rows.Close()
	for rows.Next() {
		var tempHost Host
		err := rows.Scan(&tempHost.Id, &tempHost.Uuid, &tempHost.Name, &tempHost.Ip_Address)
		if err != nil {
			return hosts
		}
		hosts = append(hosts, tempHost)
	}
	return hosts
}

func getHostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostId := vars["hostId"]
	id, _ := strconv.Atoi(hostId)
	hosts := getHostByIdDb(id, w)
	if len(hosts) == 0 {
		printMessage(w, "No such host!!!")
		return
	}
	res, err := PrettyStruct(hosts)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, res)
}

func AllContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := GetAllhContainers()
	if err != nil {
		fmt.Fprint(w, err)
		return

	}

	if len(containers) == 0 {
		infoMesage := InfoMessage{"Containers table is empty !!!"}
		res, err := PrettyStruct(infoMesage)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprint(w, res)
		return
	}
	res, err := PrettyStruct(containers)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, res)

}

func AllHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := GetAllhHosts()
	if err != nil {
		fmt.Fprint(w, err)
		return

	}

	if len(hosts) == 0 {
		infoMessage := InfoMessage{"hosts table is empty !!!"}
		res, err := PrettyStruct(infoMessage)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, res)
		return
	}

	res, err := PrettyStruct(hosts)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, res)
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
