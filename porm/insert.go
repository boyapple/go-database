package porm

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
)

func NewInsertBuilder() Builder {
	return &InsertBuilder{
		sqlBuild: &strings.Builder{},
		valBuild: &strings.Builder{},
	}
}

type InsertBuilder struct {
	sqlBuild *strings.Builder
	valBuild *strings.Builder
}

func (i *InsertBuilder) Build(message proto.Message, opts *Options) (string, []interface{}, error) {
	columns, args, err := buildColumnAndValue(message, opts)
	if err != nil {
		return "", nil, err
	}
	if len(columns) == 0 {
		return "", nil, ErrNoAssignColumn
	}
	i.sqlBuild.WriteString("INSERT INTO ")
	i.sqlBuild.WriteString(fmt.Sprintf(" `%s` ", opts.Table))
	i.sqlBuild.WriteString(" (")
	i.valBuild.WriteString(" VALUES (")
	for idx, c := range columns {
		i.sqlBuild.WriteString(fmt.Sprintf("`%s`", c))
		i.valBuild.WriteString("?")
		if idx < len(columns)-1 {
			i.sqlBuild.WriteString(",")
			i.valBuild.WriteString(",")
		}
	}
	i.sqlBuild.WriteString(")")
	i.valBuild.WriteString(")")
	i.sqlBuild.WriteString(i.valBuild.String())
	fmt.Printf("InsertBuilder sql:%s args:%+v \n", i.sqlBuild.String(), args)
	return i.sqlBuild.String(), args, nil
}
