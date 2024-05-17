package users

import (
	j "encoding/json"
	"errors"
	http_error "godb/core/errors"
	db_config "godb/db/config"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao receber dados, favor, verifique novamente"))
		return //stop execution
	}
	var user *User

	if err = j.Unmarshal(body, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao converter usuário"))
		return
	}

	db_config.DB.Create(&user)

	u_json, err := j.Marshal(*user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao criar usuário, verifique os dados e tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users *[]User
	db_config.DB.Find(&users)
	u_json, err := j.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(http_error.ResponseError("Erro ao recuperar usuário, tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user User
	id := mux.Vars(r)["id"]
	user_id, _ := uuid.Parse(id)
	if err := db_config.DB.First(&user, user_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Registro não encontrado
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(http_error.ResponseError("Erro desconhecido, por favor tente novamente"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Usuário não encontrado"))
		return
	}

	userJSONBytes, err := j.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao converter usuário, verifique os dados e tente novamente"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJSONBytes)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result := db_config.DB.Where("id = ?", id).Delete(&User{})
	if result.Error != nil || result.RowsAffected == 0 {
		// Se ocorrer um erro ou nenhuma linha for afetada, retornar false com status BadRequest
		jsonError, _ := j.Marshal(map[string]bool{"success": false})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonError)
		return
	}
	// Se a exclusão for bem-sucedida, retornar true com status OK
	jsonSuccess, err := j.Marshal(map[string]bool{"success": true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro interno, favor tente novamente ou contate o suporte"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonSuccess)
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	// Decodificar o JSON do corpo da requisição em um objeto User
	var updatedUser *User
	if err := j.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro nas informações fornecidas, verifique err tente novamente"))
		return
	}

	// Encontrar o usuário no banco de dados pelo ID
	var existingUser User
	if err := db_config.DB.First(&existingUser, "id = ?", updatedUser.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Retornar um erro 404 se o usuário não for encontrado
			w.WriteHeader(http.StatusNotFound)
			w.Write(http_error.ResponseError("Usuário não encontrado"))
			return
		}
		// Retornar um erro 500 se ocorrer um erro ao buscar o usuário
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao procurar usuário, tente novamente"))
		return
	}

	// Atualizar apenas os campos que foram enviados na requisição
	if err := db_config.DB.Model(&existingUser).Updates(&updatedUser).Error; err != nil {
		// Retornar um erro 500 se ocorrer um erro ao atualizar o usuário
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao atualizar usuário, tente novamente"))
		return
	}

	// Verificar se o campo Name em updatedUser é vazio e, se for, manter o valor original em existingUser
	if updatedUser.Name == "" {
		updatedUser.Name = existingUser.Name
	}
	// Verificar se o campo Email em updatedUser é vazio e, se for, manter o valor original em existingUser
	if updatedUser.Email == "" {
		updatedUser.Email = existingUser.Email
	}

	// Retornar o usuário atualizado como JSON
	updatedUserJSON, err := j.Marshal(updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro nas informações fornecidas, verifique err tente novamente"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(updatedUserJSON)
}
