package condition

import "gorm.io/gorm"

func NewPageCondition(offset, limit int) Condition {
	return &Page{
		Offset: offset,
		Limit:  limit,
	}
}

// Page 分页条件
type Page struct {
	Offset int
	Limit  int
}

func (p *Page) Compile() (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		offset := p.Offset
		limit := p.Limit
		if offset <= 0 {
			offset = 1
		}
		if limit <= 0 {
			limit = 10
		}
		return db.Offset((offset - 1) * limit).Limit(limit)
	}, nil
}
