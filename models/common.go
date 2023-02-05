package models

import (
	"time"
)

type EmployeesDataInput struct {
	Name         string
	Email        string
	Address      string
	Gender       string
	StratingDate time.Time
	PositionId   int
	Amount       int
}

type TasksDetail struct {
	TaskName     string
	StartingDate time.Time
	EndingDate   time.Time
}

type EmployeesProfileDetail struct {
	Tasks []TasksDetail
	EmployeesProfile
}

type TasksDataInput struct {
	Tasks
	EmployeeId []int
}
