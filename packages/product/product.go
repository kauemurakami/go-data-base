package product

import "github.com/google/uuid"

type Product struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Code  string    `gorm:"size:4;unique"`
	Price uint
}
