package users

import (
	j "encoding/json"
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
