// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package db

import (
	"context"
	"strings"
)

const addTask = `-- name: AddTask :exec
INSERT INTO tasks (
  title
) VALUES (
  ?
)
`

func (q *Queries) AddTask(ctx context.Context, title string) error {
	_, err := q.db.ExecContext(ctx, addTask, title)
	return err
}

const deleteTasks = `-- name: DeleteTasks :exec
DELETE FROM tasks WHERE id IN (/*SLICE:ids*/?)
`

func (q *Queries) DeleteTasks(ctx context.Context, ids []int64) error {
	query := deleteTasks
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const listTasks = `-- name: ListTasks :many
SELECT id, title FROM tasks
`

func (q *Queries) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(&i.ID, &i.Title); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :exec
UPDATE tasks SET title=? WHERE id=?
`

type UpdateTaskParams struct {
	Title string
	ID    int64
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.ExecContext(ctx, updateTask, arg.Title, arg.ID)
	return err
}
