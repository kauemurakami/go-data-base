package product

import (
	j "encoding/json"
	http_error "godb/core/errors"
	db_config "godb/db/config"
	"godb/packages/categories"
	"io"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID         uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CategoryID uuid.UUID           `gorm:"type:uuid;not null" json:"category_id"`
	Category   categories.Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	Name       string              `json:"name"`
	Value      uint                `json:"value"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao receber dados, favor, verifique novamente"))
		return //stop execution
	}
	var product *Product

	if err = j.Unmarshal(body, &product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao converter produto"))
		return
	}

	db_config.DB.Create(&product)

	// Carregar os detalhes da categoria associada ao produto
	if err := db_config.DB.Preload("Category").First(&product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao carregar detalhes da categoria"))
		return
	}

	u_json, err := j.Marshal(*product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao criar produto, verifique os dados e tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products *[]Product
	if err := db_config.DB.Preload("Category").Find(&products).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao recuperar produtos, tente novamente"))
		return
	}

	u_json, err := j.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao recuperar usuário, tente novamente"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u_json)
}

func DeleteProductById(w http.ResponseWriter, r *http.Request) {

}
