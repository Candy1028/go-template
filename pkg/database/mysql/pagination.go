package mysql

import (
	"gorm.io/gorm"
)

type PageOption struct {
	PageSize int
	Page     int
	Order    []string
	Preload  []string
}

func Pagination[T any](model *T, option PageOption, fn func(tx *gorm.DB) *gorm.DB) (list []T, total int64, err error) {
	query := fn(DB.Model(model))
	for _, key := range option.Preload {
		query = query.Preload(key)
	}
	if option.Order == nil {
		option.Order = []string{
			"id desc",
		}
	}
	for _, key := range option.Order {
		query = query.Order(key)
	}
	if err := query.Count(&total).Offset((option.Page - 1) * option.PageSize).
		Limit(option.PageSize).Find(&list).
		Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
