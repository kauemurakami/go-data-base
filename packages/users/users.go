package users

import (
	j "encoding/json"
	"fmt"
	db_config "godb/db/config"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:id`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error with received datas"))
		return //stop execution
	}
	var user *User

	if err = j.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Convert error json to struct"))
		return
	}

	db_config.DB.Create(user)

	u_json, e := j.Marshal(user)
	if e != nil {
		log.Fatalf("Erro ao retornar usuário")
	}
	w.Write(u_json)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users *[]User
	db_config.DB.Find(&users)
	u_json, e := j.Marshal(users)
	if e != nil {
		log.Fatalf("Erro ao retornar usuários")
	}
	w.Write(u_json)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user *User
	id := r.URL.Query().Get("id")
	fmt.Println(id)

	db_config.DB.First(&user, id)
	u_json, e := j.Marshal(user)
	if e != nil {
		log.Fatalf("Erro ao retornar usuários")
	}
	w.Write(u_json)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	var user *User
	db_config.DB.First(&user, r.URL.Query())
	u_json, e := j.Marshal(user)
	if e != nil {
		log.Fatalf("Erro ao retornar usuários")
	}
	w.Write(u_json)
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	var users *User
	db_config.DB.Find(&users)
	u_json, e := j.Marshal(users)
	if e != nil {
		log.Fatalf("Erro ao retornar usuários")
	}
	w.Write(u_json)
}
