package use_case

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/wayla99/go_gorm.git/src/entity/staff"
)

type Staff struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func (s Staff) toEntity() staff.Staff {
	return staff.Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}

func (uc UseCase) enrichStaff(s staff.Staff) (Staff, error) {
	var err error
	return Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}, err
}

func FilterStaff(id string) []string {
	return []string{fmt.Sprintf("staff_id:eq:%s", id)}
}

func (uc UseCase) CreateStaff(ctx context.Context, s Staff) (Staff, error) {
	sf := s.toEntity()
	if err := sf.Validate(); err != nil {
		return Staff{}, err
	}

	err := uc.staffRepository.Create(ctx, &sf)
	if err != nil {
		return Staff{}, err
	}
	return uc.enrichStaff(sf)
}

func (uc UseCase) GetStaffs(ctx context.Context, input *List) ([]Staff, int, error) {
	sliceStaff, total, err := uc.staffRepository.List(ctx, input, &[]*staff.Staff{})
	if err != nil {
		return nil, 0, err
	}
	sf := make([]*staff.Staff, len(sliceStaff))
	for i, i2 := range sliceStaff {
		sf[i] = i2.(*staff.Staff)
	}

	return lo.Map(sf, func(item *staff.Staff, index int) Staff {
		s, err := uc.enrichStaff(*item)
		if err != nil {
			return Staff{}
		}

		return s
	}), total, err
}

func (uc UseCase) GetStaffById(ctx context.Context, id string) (Staff, error) {
	var s staff.Staff
	err := uc.staffRepository.Read(ctx, FilterStaff(id), &s)
	if err != nil {
		return Staff{}, err
	}
	return uc.enrichStaff(s)
}

func (uc UseCase) UpdateStaffById(ctx context.Context, id string, s Staff) (Staff, error) {
	sf := s.toEntity()
	if err := sf.Validate(); err != nil {
		return Staff{}, err
	}
	err := uc.staffRepository.Update(ctx, FilterStaff(id), &sf)
	if err != nil {
		return Staff{}, err
	}
	return uc.enrichStaff(sf)
}

func (uc UseCase) DeleteStaffById(ctx context.Context, id string) error {
	err := uc.staffRepository.SoftDelete(ctx, FilterStaff(id))
	if err != nil {
		return err
	}

	return nil
}
