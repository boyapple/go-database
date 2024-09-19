package porm

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
)

func NewSelectBuilder(buildType string) Builder {
	return &SelectBuilder{
		buildType: buildType,
		build:     &strings.Builder{},
		args:      []interface{}{},
	}
}

// SelectBuilder select语句
type SelectBuilder struct {
	buildType string
	build     *strings.Builder
	args      []interface{}
}

// Build select builder
func (sb *SelectBuilder) Build(message proto.Message, opts *Options) (string, []interface{}, error) {
	sb.writeString("SELECT ")
	if err := sb.writeColumn(message, opts); err != nil {
		return "", nil, err
	}
	sb.writeString(" FROM ")
	sb.writeString(fmt.Sprintf(" `%s` ", opts.Table))
	if opts.Join != "" {
		sb.writeString(opts.Join)
	}
	if sb.buildType == SelectOne {
		if opts.Where == "" {
			return "", nil, ErrNoAssignWhereCondition
		}
		sb.writeString(fmt.Sprintf(" WHERE %s ", opts.Where))
		sb.args = append(sb.args, opts.Args...)
		sb.writeString(" LIMIT ?,? ")
		sb.args = append(sb.args, 0, 1)
		query, args := sb.builder()
		return query, args, nil
	}
	if opts.Where != "" {
		sb.writeString(fmt.Sprintf(" WHERE %s ", opts.Where))
		sb.args = append(sb.args, opts.Args...)
	}
	if sb.buildType == SelectList {
		if len(opts.OrderBy) > 0 {
			sb.writeString(" ORDER BY ")
			for i, order := range opts.OrderBy {
				if order.Desc {
					sb.writeString(fmt.Sprintf("%s DESC", order.Key))
				} else {
					sb.writeString(fmt.Sprintf("%s ASC", order.Key))
				}
				if i < len(opts.OrderBy)-1 {
					sb.writeString(",")
				}
			}
		}
		if opts.Page != nil && opts.Page.Offset > 0 && opts.Page.Limit > 0 {
			sb.writeString(" LIMIT ?,?")
			sb.args = append(sb.args, (opts.Page.Offset-1)*opts.Page.Limit, opts.Page.Limit)
		}
	}
	query, args := sb.builder()
	return query, args, nil
}

func (sb *SelectBuilder) writeColumn(message proto.Message, opts *Options) error {
	switch sb.buildType {
	case SelectOne, SelectList:
		fdMap := FieldDescMapping(message)
		columns, err := buildColumns(fdMap, opts)
		if err != nil {
			return err
		}
		if len(columns) == 0 {
			return ErrNoAssignColumn
		}
		for i, c := range columns {
			sb.writeString(fmt.Sprintf("`%s`", c))
			if i < len(columns)-1 {
				sb.writeString(",")
			}
		}
	case SelectCount:
		sb.writeString(" count(1) ")
	}
	return nil
}

func (sb *SelectBuilder) builder() (string, []interface{}) {
	return sb.build.String(), sb.args
}

func (sb *SelectBuilder) writeString(s string) {
	sb.build.WriteString(s)
}
