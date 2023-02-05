package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Leave struct {
	Id           int        `orm:"column(id);null"`
	Title        string     `orm:"column(title);size(45);null"`
	Reason       string     `orm:"column(reason);size(255);null"`
	StartingDate time.Time  `orm:"column(starting_date);type(datetime);null"`
	EndingDate   time.Time  `orm:"column(ending_date);type(datetime);null"`
	EmployeeId   *Employees `orm:"column(employee_id);rel(fk)"`
}

func (t *Leave) TableName() string {
	return "Leave"
}

func init() {
	orm.RegisterModel(new(Leave))
}

// AddLeave insert a new Leave into database and returns
// last inserted Id on success.
func AddLeave(m *Leave) (id int64, err error) {
	o := orm.NewOrm()
	if isTrue := m.StartingDate.Before(m.EndingDate); isTrue {
		id, err = o.Insert(m)
		return id, err
	}

	return id, errors.New("ending must come after starting date")
}

// GetLeaveById retrieves Leave by Id. Returns error if
// Id doesn't exist
func GetLeaveById(id int) (v *Leave, err error) {
	o := orm.NewOrm()
	v = &Leave{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLeave retrieves all Leave matches certain condition. Returns empty list if
// no records exist
func GetAllLeave(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Leave))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Leave
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateLeave updates Leave by Id and returns error if
// the record to be updated doesn't exist
func UpdateLeaveById(m *Leave) (err error) {
	o := orm.NewOrm()
	v := Leave{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLeave deletes Leave by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLeave(id int) (err error) {
	o := orm.NewOrm()
	v := Leave{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Leave{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
