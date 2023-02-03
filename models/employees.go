package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Employees struct {
	Id            int              `orm:"column(id);auto"`
	Name          string           `orm:"column(name);size(45)"`
	Address       string           `orm:"column(address);size(60)"`
	Email         string           `orm:"column(email);size(45)"`
	Gender        string           `orm:"column(gender);size(45);null"`
	StratingDate  time.Time        `orm:"column(strating_date);type(datetime)"`
	Tasks         []*Tasks         `orm:"reverse(many)"`
	EmployeeTasks []*EmployeeTasks `orm:"reverse(many)"`
	Salary        *Salary          `orm:"reverse(one)"`
}

func (t *Employees) TableName() string {
	return "employees"
}

func init() {
	orm.RegisterModel(new(Employees))
}

// AddEmployees insert a new Employees into database and returns
// last inserted Id on success.
func AddEmployees(m *EmployeesDataInput) (id int64, err error) {
	o := orm.NewOrm()
	employees := &Employees{Name: m.Name, Email: m.Email, Address: m.Address, Gender: m.Gender, StratingDate: m.StratingDate}
	if id, err = o.Insert(employees); err == nil {
		salary := &Salary{EmployeeId: &Employees{Id: int(id)}, PositionId: &Position{Id: m.PositionId}, Amount: m.Amount}
		id, err = o.Insert(salary)
	}

	return
}

// GetEmployeesById retrieves Employees by Id. Returns error if
// Id doesn't exist
func GetEmployeesById(id int) (v *EmployeesProfileDetail, err error) {
	o := orm.NewOrm()

	employeeViewData := EmployeesProfile{Id: id}
	employee := &Employees{Id: id}
	var Tasks []TasksDetail

	if err = o.Read(&employeeViewData); err == nil {
		if _, err = o.LoadRelated(employee, "tasks"); err == nil {
			for _, task := range employee.Tasks {
				taskDetail := TasksDetail{TaskName: task.Taskname, StartingDate: task.StartingDate, EndingDate: task.EndingDate}
				Tasks = append(Tasks, taskDetail)
			}
			return &EmployeesProfileDetail{Tasks: Tasks, EmployeesProfile: employeeViewData}, err
		}
		fmt.Print("first")
		return nil, err

		// var salary Salary

		// if err := o.QueryTable("salary").Filter("employee_id", id).One(&salary); err == nil {
		// 	v.Salary = &salary
		// 	return v, err
		// }

	}
	fmt.Print("Second")
	return nil, err
}

// GetAllEmployees retrieves all Employees matches certain condition. Returns empty list if
// no records exist
func GetAllEmployees(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(EmployeesProfile))
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

	var l []EmployeesProfile
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

// UpdateEmployees updates Employees by Id and returns error if
// the record to be updated doesn't exist
func UpdateEmployeesById(id int, m *EmployeesDataInput) (err error) {
	o := orm.NewOrm()
	employee := Employees{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&employee); err == nil {

		if m.Address != "" {
			employee.Address = m.Address
		}
		if m.Name != "" {
			employee.Name = m.Name
		}
		if m.Email != "" {
			employee.Email = m.Email
		}
		if m.Gender != "" {
			employee.Gender = m.Gender
		}

		if m.StratingDate.IsZero() == false {
			employee.StratingDate = m.StratingDate
		}

		if _, err = o.Update(&employee); err == nil {
			salary := Salary{EmployeeId: &Employees{Id: employee.Id}}
			if err = o.Read(&salary); err == nil {
				if m.Amount != 0 {
					salary.Amount = m.Amount
				}
				if m.PositionId != 0 {
					salary.PositionId = &Position{Id: m.PositionId}
				}
				if _, err = o.Update(&salary); err == nil {
					return
				}
				return
			}
			salary = Salary{EmployeeId: &Employees{Id: employee.Id}, PositionId: &Position{Id: m.PositionId}, Amount: m.Amount}
			_, err = o.Insert(&salary)
			return err
		}
		return
	}
	return
}

// DeleteEmployees deletes Employees by Id and returns error if
// the record to be deleted doesn't exist
func DeleteEmployees(id int) (err error) {
	o := orm.NewOrm()
	v := Employees{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Employees{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func populate(o orm.Ormer, v interface{}, name ...string) error {

	for _, value := range name {

		_, err := o.LoadRelated(v, value)
		if err != nil {
			return err
		}
	}
	return nil
}
