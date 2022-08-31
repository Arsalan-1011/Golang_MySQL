package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dbm *sql.DB

type Data struct {
	ID              int    `json:"ID"`
	First_Name      string `json:"First_Name"`
	Last_Name       string `json:"Last_Name"`
	Organization_ID int    `json:"Organization_ID"`
	Deleted         int    `json:"Deleted"`
}

type Message struct {
	Msg string `json:"msg"`
}

func connect() {
	db, err := sql.Open("mysql", "arsalan:root@tcp(127.0.0.1:3306)/details?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Database Connected...")
	}
	dbm = db

}

func getData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["ID"]

	res, err := dbm.Query("SELECT `ID`, `First_Name`, `Last_Name`, `Organization_ID`, `Deleted` FROM `data` WHERE `ID`=?", id)

	if err != nil {
		msg := Message{Msg: "User Not Found"}
		json.NewEncoder(w).Encode(msg)
	}

	for res.Next() {
		var c Data
		err := res.Scan(&c.ID, &c.First_Name, &c.Last_Name, &c.Organization_ID, &c.Deleted)

		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(c)
	}
	defer res.Close()

}

func getAllData(w http.ResponseWriter, r *http.Request) {
	res, err := dbm.Query("SELECT `ID`, `First_Name`, `Last_Name`, `Organization_ID`, `Deleted` FROM `data`")

	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {
		var c Data
		err := res.Scan(&c.ID, &c.First_Name, &c.Last_Name, &c.Organization_ID, &c.Deleted)

		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(c)
	}
	defer res.Close()

}

func createData(w http.ResponseWriter, r *http.Request) {
	var c Data

	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Last_Name := c.Last_Name
	First_Name := c.First_Name
	id := c.ID
	Organization_ID := c.Organization_ID
	Deleted := c.Deleted

	res, err := dbm.Query("INSERT INTO `data`(`ID`,`First_Name`, `Last_Name`, `Organization_ID`, `Deleted`) VALUES (?,?,?,?,?)", id, First_Name, Last_Name, Organization_ID, Deleted)

	if err != nil {
		fmt.Println(res)
		msg := Message{Msg: "Not Created"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{Msg: "Created Successfully..."}

	json.NewEncoder(w).Encode(msg)
	defer res.Close()
}

func updateData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["ID"]

	var c Data

	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Deleted := c.Deleted

	res, err := dbm.Query("UPDATE `data` SET `Deleted` = ? WHERE `ID` = ?;", Deleted, id)

	if err != nil {
		log.Fatal(err)
		msg := Message{Msg: "Not Updated"}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := Message{Msg: "Updated Successfully..."}
	json.NewEncoder(w).Encode(msg)
	defer res.Close()

}

func deleteData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["ID"]

	res, err := dbm.Query("DELETE FROM `data` WHERE `ID`=?", key)
	if err != nil {
		msg := Message{Msg: "User Not Found"}
		json.NewEncoder(w).Encode(msg)
	}

	if err != nil {
		log.Fatal(err)
		msg := Message{Msg: "Not Deleted"}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := Message{Msg: "Deleted Successfully..."}
	json.NewEncoder(w).Encode(msg)
	defer res.Close()
}

func handler() {
	muxRoutes := mux.NewRouter().StrictSlash(true)
	muxRoutes.HandleFunc("/data/{ID}", getData).Methods("GET")
	muxRoutes.HandleFunc("/data/", getAllData).Methods("GET")
	muxRoutes.HandleFunc("/data/", createData).Methods("POST")
	muxRoutes.HandleFunc("/data/{ID}", updateData).Methods("PUT")
	muxRoutes.HandleFunc("/data/{ID}", deleteData).Methods("DELETE")

	fs := http.FileServer(http.Dir("./files/"))
	muxRoutes.PathPrefix("/files/").Handler(http.StripPrefix("/files/", fs))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"DELETE", "GET", "POST", "PUT"},
	})

	handler := c.Handler(muxRoutes)
	log.Println("server started on 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func main() {
	connect()
	handler()

}
