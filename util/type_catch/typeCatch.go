package typeCatch

import (
	"encoding/json"
	"errors"

	"github.com/jinzhu/copier"
)

// StructToString 把传入的 struct 转为 string 类型
func StructToString(data any) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// ConvertFunc 是一个用于自定义字段转换的函数类型
type ConvertFunc func(src, dst interface{}) error

// CopyToDTO 把传入的 struct 转为 DTO 类型
func CopyToDTO[SRC, TRG any](src *SRC, convertFuncs ...ConvertFunc) (*TRG, error) {
	var dst TRG
	if err := copier.Copy(&dst, src); err != nil {
		return nil, err
	}

	// 应用所有传入的转换函数
	for _, convert := range convertFuncs {
		if err := convert(src, &dst); err != nil {
			return nil, err
		}
	}

	return &dst, nil
}

func CopySliceToDTO[SRC, TRG any](src []*SRC, convertFuncs ...ConvertFunc) ([]*TRG, error) {
	var dst []*TRG
	err := copier.Copy(&dst, &src)
	if err != nil {
		return nil, err
	}

	// 对切片中的每个元素应用转换函数
	for i := range dst {
		for _, convert := range convertFuncs {
			if err := convert(src[i], dst[i]); err != nil {
				return nil, err
			}
		}
	}

	return dst, nil
}

// ExtractType 从一个 interface{} 切片中提取指定类型的元素
func ExtractType[T any](uc []any) (T, error) {
	var zero T
	for _, a := range uc {
		if v, ok := a.(T); ok {
			return v, nil
		}
	}
	return zero, errors.New("invalid context struct")
}
