package model

import (
	"fmt"
	"time"
)

type Calender struct {
	YearMonth CalenderYearMonth
	Weeks     [][7]CalenderDay
}

type CalenderYearMonth struct {
	t     time.Time
	Year  int
	Month time.Month
}

type CalenderDay struct {
	HasContent bool
	Date       NippoDate
}

func NewCalender(ym CalenderYearMonth, nippoList []Nippo) (*Calender, error) {
	month := ym.Time()
	monthFirstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)
	monthLastDay := monthFirstDay.AddDate(0, 1, -1)
	lastWeekNo := (int(monthFirstDay.Weekday()) + monthLastDay.Day() - 1) / 7

	hasContentMap := make([][7]bool, 1+lastWeekNo)
	for _, nippo := range nippoList {
		if nippo.Date.Year() == month.Year() && nippo.Date.Month() == month.Month() {
			weekNo := (int(monthFirstDay.Weekday()) + nippo.Date.Day() - 1) / 7
			weekDay := nippo.Date.Weekday()
			hasContentMap[weekNo][weekDay] = true
		}
	}

	weeks := make([][7]CalenderDay, 1+lastWeekNo)
	for day := 1; day <= monthLastDay.Day(); day++ {
		date := time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, time.Local)
		weekNo := (int(monthFirstDay.Weekday()) + day - 1) / 7
		weekDay := date.Weekday()
		hasContent := hasContentMap[weekNo][weekDay]
		weeks[weekNo][weekDay] = CalenderDay{hasContent, NewNippoDate(date.Format(time.RFC3339))}
	}
	return &Calender{
		ym,
		weeks,
	}, nil
}

func NewCalenderYearMonth(fileName string) (CalenderYearMonth, error) {
	ym := CalenderYearMonth{}
	month, err := time.Parse("2006-01-02", fileName[:7]+"-01")
	if err != nil {
		return ym, err
	}
	ym.t = month
	ym.Year = month.Year()
	ym.Month = month.Month()
	return ym, nil
}

func (ym CalenderYearMonth) String() string {
	return ym.PathString()
}

func (ym CalenderYearMonth) Time() time.Time {
	return ym.t
}

func (ym CalenderYearMonth) PathString() string {
	return fmt.Sprintf("%04d%02d", ym.Year, ym.Month)
}

func (ym CalenderYearMonth) FileString() string {
	return fmt.Sprintf("%04d-%02d", ym.Year, ym.Month)
}

func (ym CalenderYearMonth) TitleString() string {
	return fmt.Sprintf("%04d/%02d", ym.Year, ym.Month)
}

func (cday CalenderDay) String() string {
	return fmt.Sprintf("%02d", cday.Date.Day())
}
