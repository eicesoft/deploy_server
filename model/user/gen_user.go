///////////////////////////////////////////////////////////
// THIS FILE IS AUTO GENERATED by gormgen, DON'T EDIT IT //
//        ANY CHANGES DONE HERE WILL BE LOST             //
///////////////////////////////////////////////////////////

package user

import (
	"fmt"

	"deploy_server/model"
	"deploy_server/pkg/core"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func NewModel() *User {
	return new(User)
}

func NewQueryBuilder() *userQueryBuilder {
	return new(userQueryBuilder)
}

func (t *User) Assign(src interface{}) {
	core.StructCopy(t, src)
}

func (t *User) Create(db *gorm.DB) (id int32, err error) {
	if err = db.Create(t).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return t.Id, nil
}

func (t *User) Delete(db *gorm.DB) (err error) {
	if err = db.Delete(t).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

func (t *User) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	if err = db.Model(&User{}).Where("id = ?", t.Id).Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

type userQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func (qb *userQueryBuilder) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	if err = qb.buildUpdateQuery(db).Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *userQueryBuilder) buildUpdateQuery(db *gorm.DB) *gorm.DB {
	ret := db.Model(&User{})
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	return ret
}

func (qb *userQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
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

func (qb *userQueryBuilder) Count(db *gorm.DB) (int64, error) {
	var c int64
	res := qb.buildQuery(db).Model(&User{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *userQueryBuilder) First(db *gorm.DB) (*User, error) {
	ret := &User{}
	res := qb.buildQuery(db).First(ret)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		ret = nil
	}
	return ret, res.Error
}

func (qb *userQueryBuilder) QueryOne(db *gorm.DB) (*User, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *userQueryBuilder) QueryAll(db *gorm.DB) ([]*User, error) {
	var ret []*User
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *userQueryBuilder) Limit(limit int) *userQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *userQueryBuilder) Offset(offset int) *userQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *userQueryBuilder) WhereId(p model.Predicate, value int32) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereIdIn(value []int32) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereIdNotIn(value []int32) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "NOT IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) OrderById(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "id "+order)
	return qb
}

func (qb *userQueryBuilder) WhereName(p model.Predicate, value string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "name", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereNameIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "name", "IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereNameNotIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "name", "NOT IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) OrderByName(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "name "+order)
	return qb
}

func (qb *userQueryBuilder) WherePassword(p model.Predicate, value string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WherePasswordIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", "IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WherePasswordNotIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", "NOT IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) OrderByPassword(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "password "+order)
	return qb
}

func (qb *userQueryBuilder) WhereRole(p model.Predicate, value string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role", p),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereRoleIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role", "IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) WhereRoleNotIn(value []string) *userQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "role", "NOT IN"),
		value,
	})
	return qb
}

func (qb *userQueryBuilder) OrderByRole(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "role "+order)
	return qb
}
