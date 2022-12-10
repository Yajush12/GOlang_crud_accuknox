package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SignupRequest struct{
	UserId		string	`json:"userid"`
	Email		string 	`json:"email"`
	Password	string	`json:"password"` 	
}
type addNoteRequest struct{
	Sid 	string	`json:"sid"`
	Note	string	`json:"note"`
}
type getAllNotes struct{
	Sid		string	`json:"sid"`
}
type LoginRequest struct{
	Email		string 	`json:"email"`
	Password	string	`json:"password"` 	
}
type User struct {
	UserId		string	`gorm:"primaryKey"`
	Email		string 	`gorm:"unique"`
	Password	string	 	 
	Sid			string	`gorm:"unique,null"`
}
type Notes struct {
	UserId		string	`gorm:"ForeignKey"`
	Nid			string	`gorm:"primaryKey"`
	Notes		string 	
}
type SidRes struct {
    Sid		string	`json:"Sid"`
}
type NotesRes struct{
	AllNotes	[]Notes	`json:"AllNotes"`
}
var DB *gorm.DB
func main() {
	DB = GetConnection()
	DB.AutoMigrate(&User{},&Notes{})
	r := mux.NewRouter()
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/addNote", addNote).Methods("POST")
	r.HandleFunc("/getNotes", getNotes).Methods("GET")
	fmt.Println("started listening on port 8000")
	fmt.Println(http.ListenAndServe(":8000",r))
}

func signup(w http.ResponseWriter, r *http.Request){

	p := SignupRequest{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sid := rand_str()
	user := User{UserId: p.UserId,Email: p.Email,Password: p.Password,Sid: sid}
	fmt.Println("details ->",user)
	result := DB.Model(&User{}).Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
		// return
	}
	fmt.Println(result)
	w.WriteHeader(http.StatusOK)
	
}

func rand_str() string {
	rand.Seed(time.Now().Unix())
    length := 6
  
    ran_str := make([]byte, length)
  
    // Generating Random string
    for i := 0; i < length; i++ {
        // ran_str[i] = ran_str(65 + rand.Intn(25))
		ran_str[i] = byte(65 + rand.Intn(25))

	}
	
    
    str := string(ran_str)
    return str;
}

func login(w http.ResponseWriter, r *http.Request){
	

	p :=  LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := User{}
	result := DB.Where(&User{Email: p.Email, Password: p.Password}).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}
	fmt.Println("details are",user)
// create a random string to create sid
	newSid := rand_str()
//update the user detail
//example db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
	DB.Model(User{}).Where("Email = ?", p.Email).Updates(User{Sid:newSid})
	w.WriteHeader(http.StatusOK)
//send response
	response := SidRes{Sid : newSid}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(response)
	

}

func addNote(w http.ResponseWriter, r *http.Request){
	//using session id get username 
	//create new entry for note 
	p :=  addNoteRequest{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := User{}
	result := DB.Model(&User{}).Where(&User{Sid: p.Sid}).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	
	fmt.Println("details are",user)
	nid := rand_str()
	note := Notes{UserId: user.UserId, Nid: nid, Notes:p.Note}

	res := DB.Model(&Notes{}).Create(&note)
	if res.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
		// return
	}

	w.WriteHeader(http.StatusOK)

}

func getNotes(w http.ResponseWriter, r *http.Request){
	p :=  getAllNotes{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//getting username 
	user := User{}
	result := DB.Model(&User{}).Where(&User{Sid: p.Sid}).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	
	fmt.Println("details are",user)
	
	//getting all the notes 
	note := []Notes{}
	res := DB.Model(&Notes{}).Where(&Notes{UserId: user.UserId}).Find(&note)
	if res.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
		// return
	}

	//struct to json parse 
	//how to parse array of struct into json 

	fmt.Println("All notes are",note)
	w.WriteHeader(http.StatusOK)

	response := NotesRes{AllNotes : note}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(response)
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Yajush"
	dbname   = "postgres"
)

func GetConnection() *gorm.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo),&gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	log.Println("DB Connection established...")

	return db
}