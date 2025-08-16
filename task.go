package main

import (
	"database/sql"
	"fmt"
)

type Task struct {
	TaskName    string
	TaskId      string
	CreatedDate string
}

type TaskDb struct {
	db *sql.DB
}

func (taskDb *TaskDb) AddTask(task Task) error {
	if task.TaskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	if task.TaskName == "" {
		return fmt.Errorf("TaskName is empty")

	}

	_, err := taskDb.db.Exec("insert into tasks (taskId,taskName) values (?,?)", task.TaskId, task.TaskName)
	if err != nil {
		return err
	}
	return nil

}
func (taskDb *TaskDb) GetTasks() ([]Task, error) {
	query, err := taskDb.db.Query("select taskId as TaskId,taskName as TaskName,CreatedDate as CreatedDate from tasks")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	var tasks []Task
	for query.Next() {
		var task Task
		if err := query.Scan(&task.TaskId, &task.TaskName, &task.CreatedDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (taskDb *TaskDb) EditTask(task Task) error {
	if task.TaskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	if task.TaskName == "" {
		return fmt.Errorf("TaskName is empty")

	}
	_, err := taskDb.db.Exec("update tasks set taskName=? where taskId=?", task.TaskName, task.TaskId)
	if err != nil {
		return err
	}
	return nil

}
func (taskDb *TaskDb) DeleteTask(taskId string) error {
	if taskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	res, err := taskDb.db.Exec("delete from tasks where taskId=?", taskId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", taskId)
	}

	return nil

}
