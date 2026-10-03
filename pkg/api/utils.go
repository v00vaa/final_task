package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/v00vaa/final_task/pkg/db"
)

func afterNow(taskDate, now time.Time) bool {
	nowDate := now.Format("20060102")
	taskDateStr := taskDate.Format("20060102")

	return taskDateStr > nowDate
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	// если сегодня (now) больше task.Date (t)
	if afterNow(now, t) {
		// если правила повторения нет, то берём сегодняшнее число
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=UTF-8",
	)

	json.NewEncoder(w).Encode(data)
}
