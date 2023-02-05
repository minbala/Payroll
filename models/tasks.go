package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Tasks struct {
	Id           int          `orm:"column(id);auto"`
	Taskname     string       `orm:"column(taskname);size(255)"`
	StartingDate time.Time    `orm:"column(starting_date);type(datetime)"`
	EndingDate   time.Time    `orm:"column(ending_date);type(datetime)"`
	EmployeeId   []*Employees `orm:"rel(m2m);rel_through(payroll/models.EmployeeTasks)"`
}

func (t *Tasks) TableName() string {
	return "tasks"
}

func init() {
	orm.RegisterModel(new(Tasks))
}

// AddTasks insert a new Tasks into database and returns
// last inserted Id on success.
func AddTasks(m *TasksDataInput) (taskId int64, err error) {
	o := orm.NewOrm()
	task := &Tasks{Taskname: m.Tasks.Taskname, StartingDate: m.Tasks.StartingDate, EndingDate: m.Tasks.EndingDate}
	if taskId, err = o.Insert(task); err == nil {
		for _, value := range m.EmployeeId {
			employee_task := &EmployeeTasks{EmployeeId: &Employees{Id: value}, TasksId: &Tasks{Id: int(taskId)}}
			if _, err = o.Insert(employee_task); err != nil {
				return
			}
		}
		return
	}
	return
}

// GetTasksById retrieves Tasks by Id. Returns error if
// Id doesn't exist
func GetTasksById(id int) (v *Tasks, err error) {
	o := orm.NewOrm()
	v = &Tasks{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTasks retrieves all Tasks matches certain condition. Returns empty list if
// no records exist
func GetAllTasks(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Tasks))
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

	var l []Tasks
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

// UpdateTasks updates Tasks by Id and returns error if
// the record to be updated doesn't exist
func UpdateTasksById(m *Tasks) (err error) {
	o := orm.NewOrm()
	v := Tasks{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTasks deletes Tasks by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTasks(id int) (err error) {
	o := orm.NewOrm()
	v := Tasks{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Tasks{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
