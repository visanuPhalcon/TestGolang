package main

import (

	"net/http"
	"fmt"
	"../model/"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"




)


var db *sql.DB
var err error


func login(res http.ResponseWriter, req *http.Request) {
	model.Login(res,req,db)
}
func add(res http.ResponseWriter,  req *http.Request){

	action:=req.FormValue("action")
	fmt.Println("action : ",action)
		model.Admin(res,req,db)


}
func delete(res http.ResponseWriter,  req *http.Request){

	action:=req.FormValue("action")
	fmt.Println("action : ",action)
	model.DeleteUser(res,req,db)


}

func sparePart(res http.ResponseWriter,  req *http.Request){


	action:=req.FormValue("action")
	fmt.Println("action : ",action)

	switch action{
	case "add" :
		model.AddSparePart(res,req,db)
	case "edit":
		fmt.Println("edit")
		model.EditSparePart(res,req,db)
	case "delete":
		id:=req.FormValue("id")
		model.DeleteSparePart(res,req,db,id)

	}

	model.GetAllSparePart(res,req,db)

}


func fixCar(res http.ResponseWriter,  req *http.Request){

	action:=req.FormValue("action")

	switch action{

	case "cookies":

	case "":


	}
	model.FixCar(res,req,db)


}


func main() {




	openDB()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", login)
	http.HandleFunc("/login", login)
	http.HandleFunc("/adminAdd",add)
	http.HandleFunc("/adminDelete",delete)
	http.HandleFunc("/spare_part",sparePart)
	http.HandleFunc("/fix",fixCar)
	http.ListenAndServe(":3000", nil)
}

func openDB()  {
	// Create an sql.DB and check for errors
	db, err = sql.Open("mysql", "root:@/autocar")
	checkErr(err)

	// sql.DB should be long lived "defer" closes it once this function ends
	//defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	checkErr(err)
}

func checkErr(err error) {
	if err !=nil{
		log.Fatal(err)
	}

}
