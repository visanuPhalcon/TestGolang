package model

import (
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"log"
	"database/sql"
	"html/template"
)

type Users struct {
	Id int
	Name string
	Email string
	PhoneNum string
	Password string

}

var err error

func Login(res http.ResponseWriter,  req *http.Request , db *sql.DB) {


	if req.Method != "POST" {
		http.ServeFile(res, req, "static/login.html")
		return
	}


	email := req.FormValue("email")
	password := req.FormValue("password")


	var databaseUsername string
	var databasePassword string

	// Search the database for the username provided
	// If it exists grab the password for validation
	err = db.QueryRow("SELECT email, password FROM users WHERE email=?", email).Scan(&databaseUsername, &databasePassword)

	fmt.Println("email: "+databaseUsername)



	// If not then redirect to the login page
	if err != nil {
		fmt.Println("not exist email")
		http.Redirect(res, req, "/", 301)
		return
	}


	// Validate the password
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	fmt.Println(err)


	fmt.Println("databasePassword: "+databasePassword)
	fmt.Println("password: "+password)


	// If wrong password redirect to the login
	if err != nil {
		fmt.Println("wrong password")
		http.Redirect(res, req, "/", 301)
		return
	}

	// If the login succeeded then go to User.index
	http.Redirect(res, req, "/addUser",301)
	// Sending it:





	stmt, err := db.Prepare("select name from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var name string
	err = stmt.QueryRow(1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)


}
func Admin(res http.ResponseWriter,  req *http.Request , db *sql.DB){

	if req.Method != "POST" {
		log.Println("test")
		http.ServeFile(res, req, "static/admin.html")
		return
	}

	firstName := req.FormValue("firstName")
	email := req.FormValue("email")
	password := req.FormValue("password")
	phoneNumber := req.FormValue("phoneNumber")

	fmt.Println("firstName: "+firstName)
	fmt.Println("email: "+email)
	fmt.Println("password: "+password)
	fmt.Println("phonenumber: "+phoneNumber)

	var nameDatabase string
	var emailDatabase string
	var phNumberDatabase string


	// Search the database for the username provided
	// If it exists grab the password for validation
	err := db.QueryRow("SELECT name, email , phonenumber FROM users WHERE name=? OR email=? OR phonenumber=?", firstName,email,phoneNumber).Scan(&nameDatabase, &emailDatabase,phNumberDatabase)

	switch  {

	case err == sql.ErrNoRows:

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//http.Error(res, "Server error, unable to create your account.", 500)
			http.ServeFile(res, req, "static/addUserFail.html")
			return
		}

		_, err = db.Exec("INSERT INTO users(name, password ,email ,phonenumber ) VALUES(?, ? , ? , ?)", firstName, hashedPassword ,email,phoneNumber)
		if err != nil {
			//http.Error(res, "Server error, unable to create your account.", 500)
			http.ServeFile(res, req, "static/addUserFail.html")
			return
		}

		http.ServeFile(res,req,"static/addUserSuccess.html")
		return

	case err != nil:
		log.Println("addUserFail")
		http.ServeFile(res, req, "static/addUserFail.html")
		//http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		log.Println("/")
		http.Redirect(res, req, "/", 301)
	}


}
func DeleteUser(res http.ResponseWriter,  req *http.Request , db *sql.DB){


	stmt, err := db.Prepare("select id, name , email , phonenumber,password from users")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()


	//var user []model.User
	//user := make([]*model.Users,0)
	user  := []*Users{}


	for rows.Next() {

		temp := new(Users)
		err := rows.Scan(&temp.Id, &temp.Name,&temp.Email,&temp.PhoneNum,&temp.Password)

		if err != nil {
			log.Fatal(err)

		}

		user = append(user,temp)


	}


	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//for i:=range(user) {
	//	log.Println()
	//	log.Println(user[i].Id)
	//	log.Println(user[i].Name)
	//	log.Println(user[i].Email)
	//	log.Println(user[i].PhoneNum)
	//	log.Println(user[i].Password)
	//}

	t,err  := template.ParseFiles("static/delete.html")
	if err != nil {
		fmt.Println(err) // Ugly debug output
		res.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	t.Execute(res , user)












}



