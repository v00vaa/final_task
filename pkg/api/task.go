package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/v00vaa/final_task/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан заголовок задачи"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10)})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search) // в параметре максимальное количество записей
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	writeJSON(w, TasksResp{
		Tasks: tasks})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан заголовок задачи"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Задача не найдена"})
		return
	}

	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)

	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error()})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error()})
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error()})
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error()})
			return
		}
	}
	writeJSON(w, map[string]string{})
}
