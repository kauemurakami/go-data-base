package sales

import (
	j "encoding/json"
	"fmt"
	http_error "godb/core/errors"
	db_config "godb/db/config"
	"godb/packages/product"
	"godb/packages/users"
	"io"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null;unique" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null;" json:"user_id"`
	ProductID uuid.UUID       `gorm:"type:uuid;not null;" json:"product_id"`
	User      users.User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Product   product.Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Value     uint            `json:"value"`
}

// ########### ADD ARRAY OF THE ProductID AND Product IN SALE
// ########### BECAUSE CAN HAVE MORE ONE PRODUCT
func CreateSale(w http.ResponseWriter, r *http.Request) {
	// Ler o corpo da solicitação
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("Erro ao receber dados, favor verifique novamente"))
		return
	}
	fmt.Println(string(body))
	// Decodificar o JSON do corpo da solicitação em um objeto Sale
	var sale *Sale
	if err := j.Unmarshal(body, &sale); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao converter venda"))
		return
	}
	// Verificar se o ID do produto está presente
	if sale.ProductID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("ID do produto não fornecido"))
		return
	}
	// Verificar se o ID do usuário está presente
	if sale.UserID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(http_error.ResponseError("ID do usuário não fornecido"))
		return
	}

	// Criar a venda no banco de dados
	if err := db_config.DB.Create(&sale).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao criar venda"))
		return
	}

	// Carregar os detalhes do usuário e do produto associados à venda
	if err := db_config.DB.Preload("User").Preload("Product").Preload("Product.Category").First(&sale).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao carregar detalhes do usuário e do produto"))
		return
	}

	// Codificar a venda como JSON e enviar como resposta
	saleJSON, err := j.Marshal(sale)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao codificar venda para JSON"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(saleJSON)
}

func GetSales(w http.ResponseWriter, r *http.Request) {
	// Recuperar todas as vendas do banco de dados com detalhes do usuário, produto e categoria
	var sales []Sale
	if err := db_config.DB.Preload("User").Preload("Product").Preload("Product.Category").Find(&sales).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao recuperar vendas, tente novamente"))
		return
	}

	// Codificar as vendas como JSON
	salesJSON, err := j.Marshal(sales)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro ao codificar vendas para JSON"))
		return
	}

	// Enviar as vendas como resposta
	w.WriteHeader(http.StatusOK)
	w.Write(salesJSON)
}
