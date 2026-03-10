package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goravel/framework/support/carbon"
)

const Associations = clause.Associations

// Model is the base model for all models in the application.
type Model struct {
	Timestamps
	ID uint `gorm:"primaryKey" json:"id"`
}

// SoftDeletes is used to add soft delete functionality to a model.
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// Timestamps is used to add created_at and updated_at timestamps to a model.
type Timestamps struct {
	CreatedAt *carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt *carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
