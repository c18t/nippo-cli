package model

import "time"

type Calender struct {
	YearMonth CalenderYearMonth
	Weeks     [][7]CalenderDay
}

type CalenderYearMonth struct {
	Year  int
	Month time.Month
}

type CalenderDay struct {
	HasContent bool
	Date       NippoDate
}
