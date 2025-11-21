package utils

import "time"

const layout = "2006-01-02"

func StringToDate(dateStr string) (time.Time, error) {
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func DateToString(t time.Time) string {
	return t.Format(layout)
}
