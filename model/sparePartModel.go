package model

import (
	"net/http"
	"log"
	"database/sql"
	"strconv"
	"fmt"
	"html/template"

	"encoding/json"
)

type SparePart struct {
	Id int
	Name string
	Type int64
	Amount int64
	Unit int64
	Price float64
	TypeText string
	UnitText string

}

var unit = []string{"ไม่มีหน่วย","ชิ้น","โหล","คู่","กล่อง"}
var typeSparePart = []string{"ไม่มีหน่วย","เครื่องยนต์/ครัทช์","แอร์รถยนต์","เกียร์","เครื่องกรอง"}


func GetAllSparePart(res http.ResponseWriter,  req *http.Request , db *sql.DB){

	stm , err := db.Prepare("select id ,  name  , type ,amount, unit , price from sparepart ")
	checkErr(err)
	defer stm.Close()

	rows , err := stm.Query()
	checkErr(err)
	defer rows.Close()

	sparePart := []*SparePart{}
	for rows.Next(){
		item := new (SparePart)
		err:= rows.Scan(&item.Id , &item.Name , &item.Type ,
			&item.Amount,&item.Unit,&item.Price )

		checkErr(err)
		item.TypeText= typeSparePart[item.Type]
		item.UnitText = unit[item.Unit]
		sparePart = append(sparePart,item)

	}

	if err = rows.Err(); err!=nil{
		log.Fatal(err)
	}


	//for i:=range(sparePart) {
	//	log.Println()
	//	log.Println(sparePart[i].Id)
	//	log.Println(sparePart[i].Name)
	//	log.Println(sparePart[i].TypeText)
	//	log.Println(sparePart[i].UnitText)
	//	log.Println(sparePart[i].Amount)
	//	log.Println(sparePart[i].Price)
	//}


	t,err  := template.ParseFiles("static/SparePart.html")
	if err != nil {
		fmt.Println(err) // Ugly debug output
		res.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	t.Execute(res , sparePart)

}
func AddSparePart(res http.ResponseWriter,  req *http.Request , db *sql.DB ){



	sparePart := new (SparePart)
	sparePart.Name = req.FormValue("name")
	sparePart.Amount,_ = strconv.ParseInt( req.FormValue("amount"),10,64)
	sparePart.Type,_ = strconv.ParseInt( req.FormValue("type"),10,64)
	sparePart.Unit,_ = strconv.ParseInt( req.FormValue("unit"),10,64)
	sparePart.Price,_ = strconv.ParseFloat(req.FormValue("price"),64)

	checkExist , err := db.Prepare("select id from sparepart where name=? limit 1  ")
	checkErr(err)
	defer checkExist.Close()

	var id int64

	err =checkExist.QueryRow(sparePart.Name).Scan(&id)




	switch {

	//insert a new row if this item is not exist
	case err == sql.ErrNoRows :
		log.Println("you can add")
		stm, err := db.Prepare("insert sparepart set name=? , amount=? , type=? , unit=? , price=? ")
		checkErr(err)
		defer stm.Close()
		insert, err := stm.Exec(sparePart.Name, sparePart.Amount, sparePart.Type, sparePart.Unit, sparePart.Price)
		checkErr(err)
		id, err =insert.LastInsertId()
		log.Println("add : id:",id)

		msg := SparePart{ Id: int(id) }
		js, _ :=json.Marshal(msg)
		res.Write(js)

	        //json.NewEncoder(res).Encode(msg)

		//res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//json,_:=json.Marshal(msg)
		//res.Write([]byte(json))




	//cant insert a new row
	case err!=nil:
		log.Println("you cant add")
		msg := SparePart{ Id: int(-1) }
		js, _ :=json.Marshal(msg)
		res.Write(js)



	default:
		log.Println("default: you cant add")
		msg := SparePart{ Id: int(-1) }
		js, _ :=json.Marshal(msg)
		res.Write(js)

	}




}
func DeleteSparePart(res http.ResponseWriter,  req *http.Request , db *sql.DB , id string){

	stm , err := db.Prepare("delete from sparepart where id=?")
	checkErr(err)

	row, err := stm.Exec(id)
	checkErr(err)

	affect , err :=  row.RowsAffected()
	checkErr(err)

	fmt.Println("delete id : ",affect)
	//db.Close()







}
func EditSparePart(res http.ResponseWriter,  req *http.Request , db *sql.DB){


	sparePart := new (SparePart)
	id ,_ := strconv.ParseInt( req.FormValue("id"),10,64)
	sparePart.Name = req.FormValue("name")
	sparePart.Amount,_ = strconv.ParseInt( req.FormValue("amount"),10,64)
	sparePart.Type,_ = strconv.ParseInt( req.FormValue("type"),10,64)
	sparePart.Unit,_ = strconv.ParseInt( req.FormValue("unit"),10,64)
	sparePart.Price,_ = strconv.ParseFloat(req.FormValue("price"),64)




	var idRow int
	checkExist , err := db.Prepare("select id from sparepart where name=? limit 1  ")
	checkErr(err)
	defer checkExist.Close()
	err =checkExist.QueryRow(sparePart.Name).Scan(&idRow)




	switch {

	//update a row if name is not exist
	case err == sql.ErrNoRows || int(id)==idRow :


		log.Println("you can update")
		stm, err := db.Prepare("update sparepart set name=? , amount=? , type=? , unit=? , price=? where id=? ")
		checkErr(err)
		defer stm.Close()
		update, err := stm.Exec(sparePart.Name, sparePart.Amount, sparePart.Type, sparePart.Unit, sparePart.Price , id)
		checkErr(err)


		row, _ := update.RowsAffected()
		log.Println("update : RowsAffected:",row)

		msg := SparePart{ Id: int(id) }
		js, _ :=json.Marshal(msg)
		res.Write(js)

	//json.NewEncoder(res).Encode(msg)

	//res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//json,_:=json.Marshal(msg)
	//res.Write([]byte(json))




	//cant insert a new row
	case err!=nil:
		log.Println("you cant update")
		msg := SparePart{ Id: int(-1) }
		js, _ :=json.Marshal(msg)
		res.Write(js)



	default:
		log.Println("you cant update")
		msg := SparePart{ Id: int(-1) }
		js, _ :=json.Marshal(msg)
		res.Write(js)

	}



}


func FixCar(res http.ResponseWriter,  req *http.Request , db *sql.DB){

	stm , err := db.Prepare("select id ,  name  , type ,amount, unit , price from sparepart ")
	checkErr(err)
	defer stm.Close()

	rows , err := stm.Query()
	checkErr(err)
	defer rows.Close()

	sparePart := []*SparePart{}
	for rows.Next(){
		item := new (SparePart)
		err:= rows.Scan(&item.Id , &item.Name , &item.Type ,
			&item.Amount,&item.Unit,&item.Price )

		checkErr(err)
		item.TypeText= typeSparePart[item.Type]
		item.UnitText = unit[item.Unit]
		sparePart = append(sparePart,item)

	}

	if err = rows.Err(); err!=nil{
		log.Fatal(err)
	}


	//for i:=range(sparePart) {
	//	log.Println()
	//	log.Println(sparePart[i].Id)
	//	log.Println(sparePart[i].Name)
	//	log.Println(sparePart[i].TypeText)
	//	log.Println(sparePart[i].UnitText)
	//	log.Println(sparePart[i].Amount)
	//	log.Println(sparePart[i].Price)
	//}


	t,err  := template.ParseFiles("static/fixCar.html")
	if err != nil {
		fmt.Println(err) // Ugly debug output
		res.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	t.Execute(res , sparePart)

}
func GetAllFixedCar(res http.ResponseWriter,  req *http.Request , db *sql.DB){

	stm , err := db.Prepare("select id ,  name  , type ,amount, unit , price from sparepart ")
	checkErr(err)
	defer stm.Close()

	rows , err := stm.Query()
	checkErr(err)
	defer rows.Close()

	sparePart := []*SparePart{}
	for rows.Next(){
		item := new (SparePart)
		err:= rows.Scan(&item.Id , &item.Name , &item.Type ,
			&item.Amount,&item.Unit,&item.Price )

		checkErr(err)
		item.TypeText= typeSparePart[item.Type]
		item.UnitText = unit[item.Unit]
		sparePart = append(sparePart,item)

	}

	if err = rows.Err(); err!=nil{
		log.Fatal(err)
	}


	//for i:=range(sparePart) {
	//	log.Println()
	//	log.Println(sparePart[i].Id)
	//	log.Println(sparePart[i].Name)
	//	log.Println(sparePart[i].TypeText)
	//	log.Println(sparePart[i].UnitText)
	//	log.Println(sparePart[i].Amount)
	//	log.Println(sparePart[i].Price)
	//}


	t,err  := template.ParseFiles("static/SparePart.html")
	if err != nil {
		fmt.Println(err) // Ugly debug output
		res.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	t.Execute(res , sparePart)

}



func checkErr(err error) {
	if err !=nil{
		log.Println(err)
	}

}