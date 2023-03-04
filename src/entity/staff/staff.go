package staff

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidStaff = errors.New("invalid staff")
)

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (*Staff) TableName() string {
	return "staff"
}

type Staff struct {
	Id        int            `gorm:"column:staff_id;not null;primaryKey;" json:"staff_id"`
	FirstName string         `validate:"omitempty,min=2,max=255" gorm:"column:first_name;"`
	LastName  string         `validate:"omitempty,min=2,max=255" gorm:"column:last_name;"`
	Email     string         `validate:"omitempty,email" gorm:"column:email;unique"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index; column:deleted_at"`
}

func (s *Staff) Validate() error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidStaff, err.Error())
	}

	return nil
}
