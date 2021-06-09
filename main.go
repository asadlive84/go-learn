package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "golangdblearn1"
)

type Student struct {
	Name      string
	Roll      int
	className int
}
type ViewData struct {
	Template *template.Template
	Student  []Student
}

func main() {
	r := mux.NewRouter()

	Psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", Psqlconn)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected!")

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello this is %s", r.URL.Path)
	})

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/books/{title}/", CreateBook).Methods("GET")
	r.HandleFunc("/books/{title}/", CreateBook).Methods("POST")
	r.HandleFunc("/books/{title}/", CreateBook).Methods("PUT")
	r.HandleFunc("/books/{title}/", CreateBook).Methods("DELETE")

	r.HandleFunc("/create-student/", StudentCreate)
	r.HandleFunc("/list/", StudentList)

	http.ListenAndServe(":8000", r)
}

func CreateBook(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(rw, "Hello this is %s", r.URL.Path)
}

func StudentList(rw http.ResponseWriter, r *http.Request) {

	Psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", Psqlconn)
	CheckError(err)

	defer db.Close()
	var students []Student

	SqlStatement := `SELECT * FROM students`

	rows, err := db.Query(SqlStatement)
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.ParseFiles("std_list.html"))
	for rows.Next() {
		var id int
		var Name string
		var roll int
		var className int
		rows.Scan(&id, &Name, &roll, &className)
		students = append(students, Student{Name, roll, className})

	}

	data := ViewData{
		Template: tmpl,
		Student:  students,
	}
	fmt.Println(data)
	err = tmpl.Execute(rw, data)
	if err != nil {
		fmt.Printf("%s,", err)
	}
}

func StudentCreate(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	Roll := r.FormValue("Roll")
	RollNumber, _ := strconv.Atoi(Roll)

	className := r.FormValue("class_name")
	StdclassName, _ := strconv.Atoi(className)

	StudentForms := Student{
		Name:      r.FormValue("name"),
		Roll:      RollNumber,
		className: StdclassName,
	}

	Psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", Psqlconn)
	CheckError(err)

	defer db.Close()
	insertData := `INSERT INTO "students"("name","roll","class_name") VALUES ($1, $2,$3)`
	_, e := db.Exec(insertData, StudentForms.Name, StudentForms.Roll, StudentForms.className)
	CheckError(e)

	//fmt.Fprintf(rw, "student: %s , roll: %s, Class: %s\n", name, roll, className)
	tmpl.Execute(rw, struct{ Success bool }{true})
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
