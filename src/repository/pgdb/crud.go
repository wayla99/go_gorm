package pgdb

import (
	"context"
	"time"

	"github.com/wayla99/go_gorm.git/src/entity/staff"

	"github.com/wayla99/go_gorm.git/src/use_case"

	"github.com/golang-module/carbon/v2"
)

type SoftDelete struct {
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func (g *GoPg) Health(ctx context.Context) error {
	return g.DB.Error
}

func (g *GoPg) List(ctx context.Context, opt *use_case.List, data interface{}) (items []interface{}, total int, err error) {
	var query string
	var args []interface{}
	var optFilter []string

	if opt.Filters != nil && len(opt.Filters) > 0 {
		optFilter = opt.Filters
		query, args = g.makeFilters(opt.Filters)
	}
	skip := (opt.Offset - 1) * opt.Limit

	total, err = g.Count(ctx, optFilter)
	if err != nil {
		return nil, 0, err
	}

	db := g.DB.Table(g.TBName).Where(query, args...).Offset(skip).Limit(opt.Limit)
	if sort := g.makeSorts(opt.Sorts); sort != "" {
		db = db.Order(sort)
	}
	if err = db.Find(data).Error; err != nil {
		return nil, 0, err
	}

	items = g.interfaceToSlice(data)
	return items, total, nil
}

func (g *GoPg) Count(ctx context.Context, filters []string) (total int, err error) {
	query, args := g.makeFilters(filters)
	var cnt int64
	if err = g.DB.Table(g.TBName).Where(query, args...).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return int(cnt), nil
}

type DataStruct interface {
	staff.Staff
}

func (g *GoPg) Create(ctx context.Context, data interface{}) error {
	return g.DB.Table(g.TBName).Create(data).Error
}

func (g *GoPg) Read(ctx context.Context, filters []string, data interface{}) error {
	query, args := g.makeFilters(filters)
	sort := g.makeSorts([]string{"updated_at:desc"})
	return g.DB.Table(g.TBName).Where(query, args...).Order(sort).First(data).Error
}

func (g *GoPg) Update(ctx context.Context, filters []string, data interface{}) error {
	query, args := g.makeFilters(filters)
	return g.DB.Table(g.TBName).Where(query, args...).Updates(data).Error
}

func (g *GoPg) Delete(ctx context.Context, filters []string) error {
	query, args := g.makeFilters(filters)
	return g.DB.Table(g.TBName).Where(query, args...).Delete(nil).Error
}

func (g *GoPg) SoftDelete(ctx context.Context, filters []string) error {
	query, args := g.makeFilters(filters)
	data := &SoftDelete{
		DeletedAt: carbon.Now(carbon.Bangkok).ToStdTime(),
	}
	return g.DB.Table(g.TBName).Where(query, args...).Updates(data).Error
}
