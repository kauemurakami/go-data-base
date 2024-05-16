package users

import (
	j "encoding/json"
	"fmt"
	db_config "godb/db/config"
	"io"
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
	var user User

	if err = j.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Convert error json to struct"))
		return
	}
	u := db_config.DB.Create(&user)
	aaa := db_config.DB.Select("id", "name", "email")
	fmt.Println(u)
	fmt.Println(aaa)
	// return
}
