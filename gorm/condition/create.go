package condition

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewConflictCondition(onlyColumns []string, updateColumns ...string) Condition {
	return &Conflict{
		OnlyColumns:   onlyColumns,
		UpdateColumns: updateColumns,
	}
}

type Conflict struct {
	OnlyColumns   []string
	UpdateColumns []string
}

func (c *Conflict) Compile() (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		var columns []clause.Column
		for _, col := range c.OnlyColumns {
			columns = append(columns, clause.Column{Name: col})
		}
		conflict := &clause.OnConflict{
			Columns: columns,
		}
		if len(c.UpdateColumns) > 0 {
			conflict.DoUpdates = clause.AssignmentColumns(c.UpdateColumns)
		} else {
			conflict.UpdateAll = true
		}
		db.Clauses(conflict)
		return db
	}, nil
}
