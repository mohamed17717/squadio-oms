package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FirstName string    `gorm:"type:text;not null" json:"first_name" validate:"required"`
	LastName  string    `gorm:"type:text;not null" json:"last_name" validate:"required"`
	Email     string    `gorm:"type:text;not null;unique" json:"email" validate:"required,email"`
	Phone     string    `gorm:"type:text;not null" json:"phone" validate:"required"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
}
