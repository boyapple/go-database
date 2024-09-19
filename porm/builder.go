package porm

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNoAssignWhereCondition = fmt.Errorf("no assign where condition")
	ErrNoAssignColumn         = fmt.Errorf("no assign column")
)

// Builder 建造sql接口
type Builder interface {
	Build(message proto.Message, opts *Options) (string, []interface{}, error)
}

func buildColumns(fdMap map[string]protoreflect.FieldDescriptor, opts *Options) ([]string, error) {
	var columns []string
	if len(opts.Fields) == 0 {
		for key := range fdMap {
			columns = append(columns, key)
		}
	} else {
		for _, field := range opts.Fields {
			if _, ok := fdMap[field]; !ok {
				return nil, ErrFieldNotExistProtobuf
			}
			columns = append(columns, field)
		}
	}
	return columns, nil
}

func buildColumnAndValue(message proto.Message, opts *Options) ([]string, []interface{}, error) {
	protoReflect := proto.MessageReflect(message)
	fdMap := FieldDescMapping(message)
	columns, err := buildColumns(fdMap, opts)
	if err != nil {
		return nil, nil, err
	}
	customField := false
	if len(opts.Fields) > 0 {
		customField = true
	}
	var fields []string
	var args []interface{}
	for _, c := range columns {
		fd, ok := fdMap[c]
		if !ok {
			return nil, nil, ErrFieldNotExistProtobuf
		}
		value := protoReflect.Get(fd)
		if !value.IsValid() {
			continue
		}
		v := value.Interface()
		// 非自定义字段默认值不进行写入
		if !customField && value.Equal(fd.Default()) {
			continue
		}
		if val, ook := opts.TimeFieldFilter.Value(string(fd.Name()), v); ook {
			v = val
		}
		fields = append(fields, c)
		args = append(args, v)
	}
	return fields, args, nil
}
