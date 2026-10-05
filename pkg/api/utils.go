package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/v00vaa/final_task/pkg/db"
)

func afterNow(taskDate, now time.Time) bool {
	nowDate := now.Format(DateFormat)
	taskDateStr := taskDate.Format(DateFormat)

	return taskDateStr > nowDate
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}
	t, err := time.Parse(DateFormat, task.Date)
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
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write response: %v", err)

	}
}
