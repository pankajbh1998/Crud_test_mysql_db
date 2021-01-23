package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Employee struct {
	Id string			`json:="id"`
	Name string 		`json:="name"`
	Age int				`json:="age"`
	Gender string		`json:="gender"`
	Role string			`json:="role"`
}

func dbConnection(w http.ResponseWriter) *sql.DB {
	user:="Pankaj"
	password:="Pankaj@123"
	ip:="127.0.0.1"
	dbName:="Company"
	//db,err:=sql.Open("mysql","Pankaj:Pankaj@123@tcp(127.0.0.1)/Company")
	db,err:=sql.Open("mysql",fmt.Sprintf("%v:%v@tcp(%v)/%v",user,password,ip,dbName))
	if err !=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	//w.WriteHeader()
	return db
}

func CreateData(w http.ResponseWriter,r* http.Request){
	db:=dbConnection(w)
	defer db.Close()

	var emp Employee
	err:=json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	} else {
		insert, _ := db.Exec(fmt.Sprintf("insert into Employee(Name,Age,Gender,Role) values('%v','%v','%v','%v')", emp.Name, emp.Age, emp.Gender, emp.Role))
		num, _ := insert.LastInsertId()
		emp.Id = strconv.Itoa(int(num))
		post, _ := json.Marshal(emp)
		//w.WriteHeader(http.StatusAccepted)
		w.Write(post)
	}
}

func ReadDataAll(w http.ResponseWriter,r* http.Request){
	w.Header().Set("content-type","application/json")
	db:=dbConnection(w)
	defer db.Close()
	result,err:=db.Query("select * from Employee")
	if err !=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	} else {
		var ans []Employee
		for result.Next() {
			var emp Employee
			err := result.Scan(&emp.Id, &emp.Name, &emp.Age, &emp.Gender, &emp.Role)
			if err != nil {
				http.Error(w,err.Error(),http.StatusNotFound)
			} else {
				ans = append(ans, emp)
			}
		}
		post, err := json.Marshal(ans)
		if err !=nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
		} else {
			w.Write(post)
		}
	}
}


func ReadDataId(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("content-type","application/json")
	db:=dbConnection(w)
	defer db.Close()
	vars:=mux.Vars(r)
	id:=vars["id"]
	result:=db.QueryRow(fmt.Sprintf("select * from Employee where id = %v",id))

	var emp Employee
	err:=result.Scan(&emp.Id,&emp.Name,&emp.Age,&emp.Gender,&emp.Role)
	if err !=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	} else {
		post, err := json.Marshal(emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(post)
		}
	}
}

func checkuser(db *sql.DB,id string)bool {
	res:=db.QueryRow(fmt.Sprintf("select Id from Employee where Id= %v",id))
	var temp string
	res.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}
func UpdateData(w http.ResponseWriter,r* http.Request) {
	w.Header().Set("content-text","application/json")
	db:=dbConnection(w)
	defer db.Close()

	vars:=mux.Vars(r)
	id:=vars["id"]
	if checkuser(db,id) == false {
		http.Error(w,"ID does not exist",http.StatusNotFound)
	} else {
		var emp Employee
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			db.Exec(fmt.Sprintf("update Employee set Name='%v',Age='%v',Gender='%v',Role='%v' where Id=%v", emp.Name, emp.Age, emp.Gender, emp.Role, id))
			emp.Id = id
			post, _ := json.Marshal(emp)
			w.Write(post)
		}
	}
}

func DeleteData(w http.ResponseWriter,r* http.Request){
	db:=dbConnection(w)
	defer db.Close()

	vars:=mux.Vars(r)
	id:=vars["id"]
	delete,_:=db.Exec(fmt.Sprintf("delete from Employee where Id = '%v'",id))
	if num,_:=delete.RowsAffected();num==0 {
		http.Error(w,"Given Id does not Exist in Database",http.StatusNotAcceptable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}


func main(){
	r:=mux.NewRouter()
	r.HandleFunc("/employee",CreateData).Methods("POST")
	r.HandleFunc("/employee",ReadDataAll).Methods("GET")
	r.HandleFunc("/employee/{id}",ReadDataId).Methods("GET")
	r.HandleFunc("/employee/{id}",UpdateData).Methods("PUT")
	r.HandleFunc("/employee/{id}",DeleteData).Methods("DELETE")
	http.ListenAndServe(":8080",r)
}

