package use_case

import (
	"context"
	"errors"
)

type List struct {
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Sorts   []string `json:"sorts"`
	Filters []string `json:"filters"`
}

var (
	ErrStaffNotFound = errors.New("staff not found")
)

type UseCase struct {
	staffRepository StaffRepository
}

type StaffRepository interface {
	Health(ctx context.Context) error
	Create(ctx context.Context, data interface{}) error
	List(ctx context.Context, opt *List, data interface{}) (items []interface{}, total int, err error)
	Read(ctx context.Context, filters []string, data interface{}) error
	Update(ctx context.Context, filters []string, data interface{}) error
	Delete(ctx context.Context, filters []string) error
	SoftDelete(ctx context.Context, filters []string) error
	Count(ctx context.Context, filters []string) (total int, err error)
}

func New(staffRepo StaffRepository) *UseCase {
	return &UseCase{staffRepository: staffRepo}
}
