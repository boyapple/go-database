package filter

func NewEmptyFieldFilter() Filter {
	return EmptyFieldFilter("")
}

// EmptyFieldFilter 空的字段过滤器
type EmptyFieldFilter string

// Type 类型
func (e EmptyFieldFilter) Type(_ string) (interface{}, bool) {
	return nil, false
}

// Value 值
func (e EmptyFieldFilter) Value(_ string, _ interface{}) (interface{}, bool) {
	return nil, false
}
