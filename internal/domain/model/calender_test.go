package model

import (
	"testing"
	"time"
)

func TestNewCalenderYearMonth(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		wantYear  int
		wantMonth time.Month
		wantErr   bool
	}{
		{
			name:      "valid year-month",
			fileName:  "2024-01-15.md",
			wantYear:  2024,
			wantMonth: time.January,
			wantErr:   false,
		},
		{
			name:      "valid December",
			fileName:  "2023-12-01.md",
			wantYear:  2023,
			wantMonth: time.December,
			wantErr:   false,
		},
		{
			name:     "invalid format",
			fileName: "invalid.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ym, err := NewCalenderYearMonth(tt.fileName)
			if tt.wantErr {
				if err == nil {
					t.Error("NewCalenderYearMonth() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewCalenderYearMonth() error = %v", err)
				return
			}
			if ym.Year != tt.wantYear {
				t.Errorf("Year = %v, want %v", ym.Year, tt.wantYear)
			}
			if ym.Month != tt.wantMonth {
				t.Errorf("Month = %v, want %v", ym.Month, tt.wantMonth)
			}
		})
	}
}

func TestCalenderYearMonth_String(t *testing.T) {
	ym, _ := NewCalenderYearMonth("2024-01-15.md")

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "String",
			method:   ym.String,
			expected: "202401",
		},
		{
			name:     "PathString",
			method:   ym.PathString,
			expected: "202401",
		},
		{
			name:     "FileString",
			method:   ym.FileString,
			expected: "2024-01",
		},
		{
			name:     "TitleString",
			method:   ym.TitleString,
			expected: "2024/01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

func TestCalenderYearMonth_Time(t *testing.T) {
	ym, _ := NewCalenderYearMonth("2024-03-15.md")
	result := ym.Time()

	if result.Year() != 2024 {
		t.Errorf("Time().Year() = %v, want 2024", result.Year())
	}
	if result.Month() != time.March {
		t.Errorf("Time().Month() = %v, want March", result.Month())
	}
}

func TestNewCalender(t *testing.T) {
	ym, _ := NewCalenderYearMonth("2024-01-15.md")

	// Create some nippo entries
	nippoList := []Nippo{
		{Date: NewNippoDate("2024-01-01.md")},
		{Date: NewNippoDate("2024-01-15.md")},
		{Date: NewNippoDate("2024-01-31.md")},
	}

	cal, err := NewCalender(ym, nippoList)
	if err != nil {
		t.Fatalf("NewCalender() error = %v", err)
	}

	if cal == nil {
		t.Fatal("NewCalender() returned nil")
	}

	// Check that calendar has correct year/month
	if cal.YearMonth.Year != 2024 {
		t.Errorf("YearMonth.Year = %v, want 2024", cal.YearMonth.Year)
	}
	if cal.YearMonth.Month != time.January {
		t.Errorf("YearMonth.Month = %v, want January", cal.YearMonth.Month)
	}

	// Check that weeks are populated
	if len(cal.Weeks) == 0 {
		t.Error("Weeks should not be empty")
	}
}

func TestNewCalender_WithoutNippos(t *testing.T) {
	ym, _ := NewCalenderYearMonth("2024-02-15.md")

	cal, err := NewCalender(ym, []Nippo{})
	if err != nil {
		t.Fatalf("NewCalender() error = %v", err)
	}

	if cal == nil {
		t.Fatal("NewCalender() returned nil")
	}

	// February 2024 should have 5 weeks (Feb 1 is Thursday)
	if len(cal.Weeks) == 0 {
		t.Error("Weeks should not be empty")
	}
}

func TestNewCalender_WithDifferentMonthNippos(t *testing.T) {
	ym, _ := NewCalenderYearMonth("2024-01-15.md")

	// Include nippos from different months
	nippoList := []Nippo{
		{Date: NewNippoDate("2024-01-15.md")},
		{Date: NewNippoDate("2024-02-15.md")}, // Different month - should be ignored
	}

	cal, err := NewCalender(ym, nippoList)
	if err != nil {
		t.Fatalf("NewCalender() error = %v", err)
	}

	if cal == nil {
		t.Fatal("NewCalender() returned nil")
	}
}

func TestCalenderDay_String(t *testing.T) {
	date := NewNippoDate("2024-01-15.md")
	day := CalenderDay{
		HasContent: true,
		Date:       date,
	}

	result := day.String()
	expected := "15"
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestCalenderDay_SingleDigit(t *testing.T) {
	date := NewNippoDate("2024-01-05.md")
	day := CalenderDay{
		HasContent: false,
		Date:       date,
	}

	result := day.String()
	expected := "05"
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}
