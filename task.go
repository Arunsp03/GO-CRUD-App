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

func (taskDb *TaskDb) AddTask(task Task, userName string) error {
	if userName == "" {
		return fmt.Errorf("UserName is empty")

	}
	if task.TaskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	if task.TaskName == "" {
		return fmt.Errorf("TaskName is empty")

	}

	_, err := taskDb.db.Exec("insert into tasks (taskId,taskName,owner) values (?,?,?)", task.TaskId, task.TaskName, userName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Unable to add task")
	}
	return nil

}
func (taskDb *TaskDb) GetTasks(userName string) ([]Task, error) {
	if userName == "" {
		return nil, fmt.Errorf("UserName is empty")

	}
	query, err := taskDb.db.Query("select taskId as TaskId,taskName as TaskName,CreatedDate as CreatedDate from tasks where owner=?", userName)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Unable to fetch tasks")
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
func (taskDb *TaskDb) EditTask(task Task, userName string) error {
	if task.TaskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	if task.TaskName == "" {
		return fmt.Errorf("TaskName is empty")

	}
	if userName == "" {
		return fmt.Errorf("UserName is empty")

	}
	_, err := taskDb.db.Exec("update tasks set taskName=? where taskId=? and owner=?", task.TaskName, task.TaskId, userName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Unable to update task with taskId %s", task.TaskId)
	}
	return nil

}
func (taskDb *TaskDb) DeleteTask(taskId string, userName string) error {
	if taskId == "" {
		return fmt.Errorf("TaskId is empty")

	}
	if userName == "" {
		return fmt.Errorf("UserName is empty")

	}
	res, err := taskDb.db.Exec("delete from tasks where taskId=? and owner=?", taskId, userName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Unable to delete task with taskId %s", taskId)
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
