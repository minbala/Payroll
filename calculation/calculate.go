package calculation

import (
	"payroll/models"

	"github.com/beego/beego/v2/client/orm"
)

func getOverTime(start string, end string, employeeId int) (totalHour int, err error) {
	o := orm.NewOrm()
	r := o.Raw("SELECT SUM(TIMESTAMPDIFF(SECOND,starting_time,ending_time)) AS DIFFERENCE FROM overtime WHERE employee_id=? starting_time BETWEEN ? AND ?;", employeeId, start, end)
	err = r.QueryRow(&totalHour)
	return totalHour, err
}

func getLeaveTime(start string, end string, employeeId int) (totalLeaveDays int, err error) {
	o := orm.NewOrm()
	r := o.Raw("SELECT SUM(TIMESTAMPDIFF(DAY,starting_date,ending_date)) AS DIFFERENCE FROM leave WHERE employee_id=? starting_date BETWEEN ? AND ?", employeeId, start, end)
	err = r.QueryRow(&totalLeaveDays)
	return totalLeaveDays, err

}

func calculateSalary(start string, end string, employeeId int) (totalAmount int, err error) {
	o := orm.NewOrm()
	employee := &models.EmployeesProfile{Id: employeeId}
	if err = o.Read(employee); err == nil {

	}
}
