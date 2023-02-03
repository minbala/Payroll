// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"payroll/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/overtime",
			beego.NSInclude(
				&controllers.OvertimeController{},
			),
		),

		beego.NSNamespace("/position",
			beego.NSInclude(
				&controllers.PositionController{},
			),
		),

		beego.NSNamespace("/record",
			beego.NSInclude(
				&controllers.RecordController{},
			),
		),

		beego.NSNamespace("/salary",
			beego.NSInclude(
				&controllers.SalaryController{},
			),
		),

		beego.NSNamespace("/tasks",
			beego.NSInclude(
				&controllers.TasksController{},
			),
		),

		beego.NSNamespace("/employee_tasks",
			beego.NSInclude(
				&controllers.EmployeeTasksController{},
			),
		),

		beego.NSNamespace("/employees",
			beego.NSInclude(
				&controllers.EmployeesController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
