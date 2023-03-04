package fiber_server

import (
	"github.com/wayla99/go_gorm.git/src/use_case"
)

type staffListResponse struct {
	Data  []Staff `json:"data"`
	Total uint64  `json:"total"`
}

type Staff struct {
	Id        int    `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
} //@Name Staff

func (s Staff) toUseCase() use_case.Staff {
	return use_case.Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}

func newStaff(c use_case.Staff) Staff {
	return Staff{
		Id:        c.Id,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
	}
}
