package condition

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OpType uint32

const (
	OpTypeNone           OpType = iota
	OpTypeLarger                // >
	OpTypeLesser                // <
	OpTypeLargeAndLesser        // <>
	OpTypeLesserOrEq            // <=
	OpTypeLargerOrEq            // >=
	OpTypeLike                  // like
	OpTypeEq                    // =
	OpTypeNotEq                 // !=

	OpTypeIn    // in
	OpTypeNotIn // not in
)

func NewIDCondition(id uint64) Condition {
	return NewWhereCondition("id", id, OpTypeEq)
}

func NewEqCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeEq)
}

func NewNotEqCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeNotEq)
}

func NewGtCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeLarger)
}

func NewGteCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeLargerOrEq)
}

func NewLtCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeLesser)
}

func NewLteCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeLesserOrEq)
}

func NewLikeCondition(column string, value interface{}) Condition {
	return NewWhereCondition(column, value, OpTypeLike)
}

func NewInCondition(column string, values []interface{}) Condition {
	return NewWhereCondition(column, values, OpTypeIn)
}

func NewNotInCondition(column string, values []interface{}) Condition {
	return NewWhereCondition(column, values, OpTypeNotIn)
}

func NewWhereCondition(column string, value interface{}, opType OpType) Condition {
	return &Where{
		OpType: opType,
		Column: column,
		Value:  value,
	}
}

type Where struct {
	OpType OpType
	Column string
	Value  interface{}
}

func (w *Where) Compile() (func(*gorm.DB) *gorm.DB, error) {
	switch w.OpType {
	case OpTypeLarger:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Gt{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeLargerOrEq:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Gte{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeLesser:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Lt{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeLesserOrEq:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Lte{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeLike:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Like{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeEq, OpTypeIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Eq{Column: w.Column, Value: w.Value})
		}, nil
	case OpTypeNotEq, OpTypeLargeAndLesser, OpTypeNotIn:
		return func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clause.Neq{Column: w.Column, Value: w.Value})
		}, nil
	default:
		return nil, fmt.Errorf("not support op type")
	}
}
