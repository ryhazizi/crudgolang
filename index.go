package main

import "fmt"
import "html/template"
import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("template/*"))
}

type employe struct {
	Id     int
	Nama   string
	Alamat string
}

func connection() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "data"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := connection()
	selDB, err := db.Query("SELECT * FROM karyawan ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	kar := employe{}
	res := []employe{}
	for selDB.Next() {
		var id int
		var nama, alamat string
		err = selDB.Scan(&id, &nama, &alamat)
		if err != nil {
			panic(err.Error())
		}
		kar.Id = id
		kar.Nama = nama
		kar.Alamat = alamat
		res = append(res, kar)
	}
	tpl.ExecuteTemplate(w, "home.flare", res)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "insert.flare", nil)
}

func insertProcess(w http.ResponseWriter, r *http.Request) {
	db := connection()
	if r.Method == "POST" {
		nama := r.FormValue("name")
		alamat := r.FormValue("address")
		insForm, err := db.Prepare("INSERT INTO karyawan(nama, alamat) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(nama, alamat)
		http.Redirect(w, r, "/", 301)
	} else {
		tpl.ExecuteTemplate(w, "404.flare", nil)
	}
	defer db.Close()
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := connection()
	getId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM karyawan WHERE id=?", getId)
	if err != nil {
		panic(err.Error())
	}
	kar := employe{}
	for selDB.Next() {
		var id int
		var nama, alamat string
		err = selDB.Scan(&id, &nama, &alamat)
		if err != nil {
			panic(err.Error())
		}
		kar.Id = id
		kar.Nama = nama
		kar.Alamat = alamat
	}
	tpl.ExecuteTemplate(w, "edit.flare", kar)
	defer db.Close()
}

func editProcess(w http.ResponseWriter, r *http.Request) {
	db := connection()
	if r.Method == "POST" {
		id := r.FormValue("id")
		nama := r.FormValue("name")
		alamat := r.FormValue("address")
		edtForm, err := db.Prepare("UPDATE karyawan SET nama=?, alamat=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		edtForm.Exec(nama, alamat, id)
		http.Redirect(w, r, "/", 301)
	} else {
		tpl.ExecuteTemplate(w, "404.flare", nil)
	}
	defer db.Close()
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := connection()
	id := r.URL.Query().Get("id")
	delData, err := db.Prepare("DELETE FROM karyawan WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delData.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "404.flare", nil)
}

func main() {

	/*indexed page*/

	http.HandleFunc("/", Index)

	http.HandleFunc("/insert", Insert)

	http.HandleFunc("/insertProcess", insertProcess)

	http.HandleFunc("/edit", Edit)

	http.HandleFunc("/delete", Delete)

	http.HandleFunc("/editProcess", editProcess)

	/*non-indexed page*/

	http.HandleFunc("/template", ErrorPage)

	http.HandleFunc("/template/home.flare", ErrorPage)

	http.HandleFunc("/template/insert.flare", ErrorPage)

	http.HandleFunc("/template/edit.flare", ErrorPage)

	http.HandleFunc("/template/delete.flare", ErrorPage)

	http.HandleFunc("/template/404.flare", ErrorPage)

	fmt.Println("web server berjalan akses http://localhost:8080/")

	http.ListenAndServe(":8080", nil)
}
