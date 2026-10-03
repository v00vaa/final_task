package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

var (
	ErrEmptyRepeat              = errors.New("repeat rule is empty")
	ErrInvalidStartDate         = errors.New("invalid start date format")
	ErrInvalidRepeatFormat      = errors.New("invalid repeat format")
	ErrInvalidRepeatSymbol      = errors.New("invalid symbol in repeat rule")
	ErrMissingRepeatInterval    = errors.New("repeat interval is missing")
	ErrRepeatIntervalOutOfRange = errors.New("repeat interval is out of range")
	ErrInvalidWeekDay           = errors.New("invalid weekday value")
	ErrInvalidMonthDay          = errors.New("invalid month day")
	ErrInvalidMonth             = errors.New("invalid month value")
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}
	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidStartDate, dstart)
	}
	params := strings.Split(repeat, " ")
	switch params[0] {
	case "d":
		if len(params) != 2 {
			return "", fmt.Errorf("%w: %w", ErrInvalidRepeatFormat, ErrMissingRepeatInterval)
		}
		interval, err := strconv.Atoi(params[1])
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrInvalidRepeatFormat, params[1])
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("%w: %d", ErrRepeatIntervalOutOfRange, interval)
		}
		for {
			start = start.AddDate(0, 0, interval)
			if afterNow(start, now) {
				return start.Format(DateFormat), nil
			}
		}
	case "y":
		if len(params) != 1 {
			return "", fmt.Errorf("%w: %s", ErrInvalidRepeatSymbol, repeat)
		}
		for {
			start = start.AddDate(1, 0, 0)
			if afterNow(start, now) {
				return start.Format(DateFormat), nil
			}
		}
	case "w":
		if len(params) != 2 {
			return "", fmt.Errorf("%w: %s", ErrInvalidRepeatFormat, repeat)
		}
		days := strings.Split(params[1], ",")
		targetDays := make(map[int]bool)
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("%w: %s", ErrInvalidWeekDay, v)
			}
			if day < 1 || day > 7 {
				return "", fmt.Errorf("%w: %s", ErrInvalidWeekDay, v)
			}
			targetDays[day] = true
		}

		date := now
		for i := 1; i <= 7; i++ {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if targetDays[weekday] {
				return date.Format(DateFormat), nil
			}
		}
		return "", ErrInvalidRepeatFormat
	case "m":
		if len(params) < 2 || len(params) > 3 {
			return "", fmt.Errorf("%w: %s", ErrInvalidRepeatFormat, repeat)
		}
		days := strings.Split(params[1], ",")
		months := make(map[int]bool)
		if len(params) == 3 {
			for _, v := range strings.Split(params[2], ",") {
				month, err := strconv.Atoi(v)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("%w: %s", ErrInvalidMonth, v)
				}
				months[month] = true
			}
		}
		for _, v := range days {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("%w: %s", ErrInvalidMonthDay, v)
			}
			if day == 0 || day > 31 || day < -2 {
				return "", fmt.Errorf("%w: %s", ErrInvalidMonthDay, v)
			}
		}
		date := start
		for {
			year := date.Year()
			month := date.Month()
			if len(months) == 0 || months[int(month)] {
				var result time.Time
				found := false
				for _, v := range days {
					day, _ := strconv.Atoi(v)
					lastDay := time.Date(
						year,
						month+1,
						0,
						0, 0, 0, 0,
						time.Local,
					).Day()
					if day < 0 {
						day = lastDay + day + 1
					}
					if day < 1 || day > lastDay {
						continue
					}
					candidate := time.Date(
						year,
						month,
						day,
						0, 0, 0, 0,
						time.Local,
					)
					if afterNow(candidate, now) &&
						afterNow(candidate, start) {

						if !found || candidate.Before(result) {
							result = candidate
							found = true
						}
					}
				}
				if found {
					return result.Format(DateFormat), nil
				}
			}
			date = time.Date(
				year,
				month+1,
				1,
				0, 0, 0, 0,
				time.Local,
			)
		}
	}
	return "", fmt.Errorf("%w: %w: %s", ErrInvalidRepeatFormat, ErrInvalidRepeatSymbol, params[0])
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
