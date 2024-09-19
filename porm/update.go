package porm

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
)

func NewUpdateBuilder() Builder {
	return &UpdateBuilder{
		build: &strings.Builder{},
	}
}

type UpdateBuilder struct {
	build *strings.Builder
}

func (u *UpdateBuilder) Build(message proto.Message, opts *Options) (string, []interface{}, error) {
	if opts.Where == "" {
		return "", nil, ErrNoAssignWhereCondition
	}
	columns, args, err := buildColumnAndValue(message, opts)
	if err != nil {
		return "", nil, err
	}
	if len(columns) == 0 {
		return "", nil, ErrNoAssignColumn
	}
	u.build.WriteString("UPDATE ")
	u.build.WriteString(fmt.Sprintf("`%s`", opts.Table))
	u.build.WriteString(" SET ")
	for idx, c := range columns {
		u.build.WriteString(fmt.Sprintf("`%s` = ?", c))
		if idx < len(columns)-1 {
			u.build.WriteString(",")
		}
	}
	u.build.WriteString(fmt.Sprintf(" WHERE %s", opts.Where))
	args = append(args, opts.Args...)
	fmt.Printf("UpdateBuilder sql:%s args:%+v \n", u.build.String(), args)
	return u.build.String(), args, nil
}
