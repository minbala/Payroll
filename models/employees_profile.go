package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type EmployeesProfile struct {
	Id           int       `orm:"column(id)"`
	Name         string    `orm:"column(name);size(45)"`
	PositionName string    `orm:"column(position_name);size(60)"`
	Salary       int       `orm:"column(salary)"`
	Address      string    `orm:"column(address);size(60)"`
	Email        string    `orm:"column(email);size(45)"`
	Gender       string    `orm:"column(gender);size(45);null"`
	StratingDate time.Time `orm:"column(strating_date);type(datetime)"`
}

func (t *EmployeesProfile) TableName() string {
	return "employees_profile"
}

func init() {
	orm.RegisterModel(new(EmployeesProfile))
}
