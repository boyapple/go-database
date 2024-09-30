package condition

import "gorm.io/gorm"

func NewPageCondition(offset, limit int) Condition {
	p := &Page{
		Offset: offset,
		Limit:  limit,
	}
	if p.Offset <= 0 {
		p.Offset = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p
}

// Page 分页条件
type Page struct {
	Offset int
	Limit  int
}

func (p *Page) Compile() (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.Offset - 1) * p.Limit).Limit(p.Limit)
	}, nil
}
