package util

import (
	"errors"
	"strconv"
)

func GetPage(pageStr string, pageSizeStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return 0, 0, err
	}
	if page <= 0 || pageSize <= 0 {
		return 0, 0, errors.New("page error")
	}
	return page, pageSize, nil
}
