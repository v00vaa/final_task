package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
	INSERT INTO scheduler
		(date, title, comment, repeat)
	VALUES
		(:date, :title, :comment, :repeat)
	`
	res, err := db.Exec(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
    `
	var args []any
	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			query += `
                WHERE date = :date
            `
			args = append(args, sql.Named("date", date.Format("20060102")))
		} else {
			query += `
                WHERE title LIKE :search
                   OR comment LIKE :search
            `
			args = append(args, sql.Named("search", "%"+search+"%"))
		}
	}
	query += `
        ORDER BY date
        LIMIT :limit
    `
	args = append(args, sql.Named("limit", limit))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]*Task, 0)
	for rows.Next() {
		task := new(Task)
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := new(Task)
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=:id",
		sql.Named("id", id)).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `
	UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat
	WHERE id=:id
	`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `
	DELETE FROM scheduler WHERE id=:id
	`
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `
	UPDATE scheduler SET date=:date
	WHERE id=:id
	`
	res, err := db.Exec(query,
		sql.Named("id", id),
		sql.Named("date", next),
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
