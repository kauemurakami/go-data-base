package categories

import (
	j "encoding/json"
	http_error "godb/core/errors"
	db_config "godb/db/config"
	"io"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null" json:"id"`
	Name string    `json:"name"`
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao receber dados, favor, verifique novamente"))
		return //stop execution
	}
	var category *Category

	if err = j.Unmarshal(body, &category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao converter categoria"))
		return
	}

	db_config.DB.Create(&category)

	u_json, err := j.Marshal(*category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao criar categoria, verifique os dados e tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	var category *[]Category
	db_config.DB.Find(&category)
	u_json, err := j.Marshal(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao recuperar usuário, tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}
