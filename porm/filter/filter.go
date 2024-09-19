package filter

// Filter 过滤器
type Filter interface {
	Type(name string) (interface{}, bool)
	Value(name string, value interface{}) (interface{}, bool)
}
