// Copyright 2025 Sergey Vinogradov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v3"

	"github.com/boogie-byte/qx/internal/db"
)

var cmd = &cli.Command{
	Name:  "qx",
	Usage: "Simple task list manager",
	Commands: []*cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize new database",
			Action: initDB,
		},
		{
			Name:    "task",
			Aliases: []string{"t"},
			Usage:   "Task operations",
			Commands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "Add new task",
					Action:  taskAdd,
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "List all tasks",
					Action:  taskList,
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Delete tasks",
					Action:  taskDelete,
					Arguments: []cli.Argument{
						&cli.Int64Args{
							Name: "ids",
							Min:  1,
							Max:  -1,
						},
					},
				},
			},
		},
	},
}

func taskAdd(ctx context.Context, c *cli.Command) error {
	if !c.Args().Present() {
		return errors.New("no title provided")
	}

	title := strings.Join(c.Args().Slice(), " ")

	conn, err := getDB()
	if err != nil {
		return err
	}

	queries := db.New(conn)

	return queries.AddTask(ctx, title)
}

func taskList(ctx context.Context, c *cli.Command) error {
	conn, err := getDB()
	if err != nil {
		return err
	}

	queries := db.New(conn)

	tasks, err := queries.ListTasks(ctx)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
	for _, task := range tasks {
		fmt.Fprintf(tw, "%d\t%s\n", task.ID, task.Title)
	}

	return tw.Flush()
}

func taskDelete(ctx context.Context, c *cli.Command) error {
	ids := c.Int64Args("ids")
	if len(ids) == 0 {
		return errors.New("no task ids provided")
	}

	conn, err := getDB()
	if err != nil {
		return err
	}

	queries := db.New(conn)

	return queries.DeleteTasks(ctx, ids)
}

func initDB(ctx context.Context, c *cli.Command) error {
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("file exists: %s", dbPath)
	}

	dsn := fmt.Sprintf("file:%s?mode=rwc", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, title TEXT NOT NULL)"); err != nil {
		return err
	}

	return nil
}

func getDB() (*sql.DB, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("file:%s?mode=rw", dbPath)

	return sql.Open("sqlite3", dsn)
}

func getDBPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(homedir, ".qx")

	if err := os.MkdirAll(dataDir, 0777); err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "qx.db"), nil
}

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
