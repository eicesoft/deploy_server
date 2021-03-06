///////////////////////////////////////////////////////////
// THIS FILE IS AUTO GENERATED by gormgen, DON'T EDIT IT //
//        ANY CHANGES DONE HERE WILL BE LOST             //
///////////////////////////////////////////////////////////

package role

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Predicate string

const (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
)

func NewModel() *Role {
	return new(Role)
}

func NewQueryBuilder() *QueryBuilder {
	return new(QueryBuilder)
}

func (t *Role) Create(db *gorm.DB) (id int32, err error) {
	if err = db.Create(t).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return t.Id, nil
}

func (t *Role) Delete(db *gorm.DB) (err error) {
	if err = db.Delete(t).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

func (t *Role) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	if err = db.Model(&Role{}).Where("id = ?", t.Id).Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

type QueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func (qb *QueryBuilder) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	if err = qb.buildUpdateQuery(db).Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *QueryBuilder) buildUpdateQuery(db *gorm.DB) *gorm.DB {
	ret := db.Model(&Role{})
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	return ret
}

func (qb *QueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
	ret := db
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *QueryBuilder) Count(db *gorm.DB) (int64, error) {
	var c int64
	res := qb.buildQuery(db).Model(&Role{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *QueryBuilder) First(db *gorm.DB) (*Role, error) {
	ret := &Role{}
	res := qb.buildQuery(db).First(ret)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		ret = nil
	}
	return ret, res.Error
}

func (qb *QueryBuilder) QueryOne(db *gorm.DB) (*Role, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *QueryBuilder) QueryAll(db *gorm.DB) ([]*Role, error) {
	var ret []*Role
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) WhereId(p Predicate, value int32) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", p),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereIdIn(value []int32) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereIdNotIn(value []int32) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "NOT IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) OrderById(asc bool) *QueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "id "+order)
	return qb
}

func (qb *QueryBuilder) WhereRoleName(p Predicate, value string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_name", p),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereRoleNameIn(value []string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_name", "IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereRoleNameNotIn(value []string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_name", "NOT IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) OrderByRoleName(asc bool) *QueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "role_name "+order)
	return qb
}

func (qb *QueryBuilder) WhereRoleDesc(p Predicate, value string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_desc", p),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereRoleDescIn(value []string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_desc", "IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) WhereRoleDescNotIn(value []string) *QueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role_desc", "NOT IN"),
		value,
	})
	return qb
}

func (qb *QueryBuilder) OrderByRoleDesc(asc bool) *QueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "role_desc "+order)
	return qb
}
