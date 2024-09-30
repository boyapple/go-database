package condition

import (
	"gorm.io/gorm"
)

type MultiCondition []Condition

type Condition interface {
	Compile() (func(*gorm.DB) *gorm.DB, error)
}

func (mc MultiCondition) Build() ([]func(*gorm.DB) *gorm.DB, error) {
	conditions := make([]func(*gorm.DB) *gorm.DB, 0)
	for _, i := range mc {
		condition, err := i.Compile()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}
