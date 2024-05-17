package sales

import (
	"encoding/json"
	j "encoding/json"
	"fmt"
	http_error "godb/core/errors"
	db_config "godb/db/config"
	"godb/packages/product"
	"godb/packages/users"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SaleProduct struct {
	gorm.Model
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null;unique" json:"id"`
	SaleID    uuid.UUID       `gorm:"type:uuid;not null;primaryKey" json:"sale_id"`
	ProductID uuid.UUID       `gorm:"type:uuid;not null;" json:"product_id"`
	Product   product.Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

type Sale struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null;unique" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;" json:"user_id"`
	ProductID uuid.UUID  `gorm:"type:uuid;not null;" json:"product_id"`
	User      users.User `gorm:"foreignKey:UserID;references:ID" json:"user"`
	// Product   product.Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Products []product.Product `gorm:"many2many:sale_products;" json:"products"`
	Value    uint              `json:"total_value"`
}

func CreateSale(w http.ResponseWriter, r *http.Request) {
	// Ler o corpo da solicitação
	// Ler o corpo da solicitação
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao receber dados, favor verifique novamente"))
		return
	}
	fmt.Println(string(body))

	// Decodificar o JSON do corpo da solicitação em um objeto Sale
	var request struct {
		UserID     uuid.UUID   `json:"user_id"`
		ProductIDs []uuid.UUID `json:"product_ids"`
		Value      uint        `json:"value"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao converter venda"))
		return
	}

	// Verificar se o ID do usuário está presente
	if request.UserID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID do usuário não fornecido"))
		return
	}

	// Verificar se os IDs dos produtos estão presentes
	if len(request.ProductIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("IDs dos produtos não fornecidos"))
		return
	}

	// Criar a venda no banco de dados
	sale := Sale{
		UserID: request.UserID,
		Value:  request.Value,
	}
	if err := db_config.DB.Create(&sale).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao criar venda"))
		return
	}

	// Adicionar produtos à venda e fazer preload
	for _, productID := range request.ProductIDs {
		saleProduct := SaleProduct{
			SaleID:    sale.ID,
			ProductID: productID,
		}
		if err := db_config.DB.Create(&saleProduct).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao associar produto à venda"))
			return
		}
	}

	// Preload dos produtos associados à venda
	if err := db_config.DB.Preload("Products").First(&sale).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao carregar detalhes dos produtos associados à venda"))
		return
	}

	// Codificar a venda como JSON e enviar como resposta
	saleJSON, err := json.Marshal(sale)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao codificar venda para JSON"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(saleJSON)

}

// ########### ADD ARRAY OF THE ProductID AND Product IN SALE
// ########### BECAUSE CAN HAVE MORE ONE PRODUCT
// func CreateSale(w http.ResponseWriter, r *http.Request) {
// 	// Ler o corpo da solicitação
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(http_error.ResponseError("Erro ao receber dados, favor verifique novamente"))
// 		return
// 	}
// 	fmt.Println(string(body))
// 	// Decodificar o JSON do corpo da solicitação em um objeto Sale
// 	var sale *Sale
// 	if err := j.Unmarshal(body, &sale); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(http_error.ResponseError("Erro ao converter venda"))
// 		return
// 	}
// 	// Verificar se o ID do produto está presente
// 	if sale.ProductID == uuid.Nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(http_error.ResponseError("ID do produto não fornecido"))
// 		return
// 	}
// 	// Verificar se o ID do usuário está presente
// 	if sale.UserID == uuid.Nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(http_error.ResponseError("ID do usuário não fornecido"))
// 		return
// 	}

// 	// Criar a venda no banco de dados
// 	if err := db_config.DB.Create(&sale).Error; err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(http_error.ResponseError("Erro ao criar venda"))
// 		return
// 	}

// 	// Carregar os detalhes do usuário e do produto associados à venda
// 	if err := db_config.DB.Preload("User").Preload("Product").Preload("Product.Category").First(&sale).Error; err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(http_error.ResponseError("Erro ao carregar detalhes do usuário e do produto"))
// 		return
// 	}

// 	// Codificar a venda como JSON e enviar como resposta
// 	saleJSON, err := j.Marshal(sale)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(http_error.ResponseError("Erro ao codificar venda para JSON"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write(saleJSON)
// }

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

func DeleteSaleById(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(mux.Vars(r)["id"])
	fmt.Println(id)
	sale := Sale{ID: id}
	result := db_config.DB.Unscoped().Where("id = ?", id).Delete(&sale)
	// result := db_config.DB.Delete(&sale)
	// result := db_config.DB.Where("id = ?", id).Delete(&sale)
	if result.Error != nil {
		// Se ocorrer um erro, retornar um status de erro interno com uma mensagem de erro
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro interno ao excluir usuário"))
		return
	}
	if result.RowsAffected == 0 {
		// Se nenhuma linha for afetada, retornar um status de não encontrado com uma mensagem adequada
		w.WriteHeader(http.StatusNotFound)
		w.Write(http_error.ResponseError("Usuário não encontrado"))
		return
	}

	// Se a exclusão for bem-sucedida, retornar true com status OK
	jsonSuccess, err := j.Marshal(map[string]bool{"success": true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(http_error.ResponseError("Erro interno ao serializar a resposta"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonSuccess)

}
