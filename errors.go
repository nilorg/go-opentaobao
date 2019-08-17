package opentaobao

import "errors"

var (
	ErrTypeIsNil   = errors.New("类型为Nil")
	ErrTypeUnknown = errors.New("未处理到的数据类型")
)
