package todos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id   int    `json:"id"`
	Task string `json:"task"`
}

type TodoDTO struct {
	Task string `json:"task"`
}

type TodosController struct {
	Engine string
	DB     *sql.DB
	SQL    TodosQueries
}

func (c *TodosController) CreateTodosTableIfNotExists() error {
	_, err := c.DB.Exec(c.SQL.CreateTodosTableIfNotExists())
	return err
}

func (c *TodosController) ActionList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	search := query.Get("search")

	stmt, err := c.DB.PrepareContext(ctx, c.SQL.ListTodos())
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "error in prepare statement")
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "%"+search+"%")
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "query failed")
		return
	}
	defer rows.Close()

	var todos []Todo = make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Task); err != nil {
			fmt.Println(err)
			writeServerErrorResponse(w, "scan row failed")
			return
		}

		todos = append(todos, todo)
	}

	res := make(map[string]([]Todo))
	res["todos"] = todos
	writeJsonResponse(w, res)
}

func (c *TodosController) ActionInsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto TodoDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		writeBadRequestResponse(w, "parsing body failed")
		return
	}

	if c.Engine == "pg" {
		c.insertTodoPG(ctx, w, dto)
	} else if c.Engine == "mysql" {
		c.insertTodoMySQL(ctx, w, dto)
	} else {
		log.Fatalf("Invalid controller.Engine %v", c.Engine)
	}
}

func (c *TodosController) insertTodoPG(ctx context.Context, w http.ResponseWriter, dto TodoDTO) {
	rows, err := c.DB.QueryContext(ctx, c.SQL.InsertTodo(), dto.Task)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "insert failed")
		return
	}
	defer rows.Close()

	if rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Task); err != nil {
			fmt.Println(err)
			writeServerErrorResponse(w, "scan row failed")
			return
		}

		res := make(map[string]Todo)
		res["todo"] = todo
		writeJsonResponse(w, res)
		return
	} else {
		writeServerErrorResponse(w, "no record of inserted task")
		return
	}
}

func (c *TodosController) insertTodoMySQL(ctx context.Context, w http.ResponseWriter, dto TodoDTO) {
	dbRes, err := c.DB.ExecContext(ctx, c.SQL.InsertTodo(), dto.Task)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "insert failed")
		return
	}

	lId, err := dbRes.LastInsertId()
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "cannot get insert-id")
		return
	}

	rows, err := c.DB.QueryContext(ctx, c.SQL.TodoById(), lId)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "error in query")
		return
	}
	defer rows.Close()

	if rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Task); err != nil {
			writeServerErrorResponse(w, "scan row failed")
			return
		}

		res := make(map[string]Todo)
		res["todo"] = todo
		writeJsonResponse(w, res)
		return
	} else {
		writeServerErrorResponse(w, "no record of inserted task")
		return
	}
}

func (c *TodosController) ActionUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var dto TodoDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		writeBadRequestResponse(w, "parsing body failed")
		return
	}

	var todo Todo
	rows, err := c.DB.QueryContext(ctx, c.SQL.TodoById(), id)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "error in query")
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&todo.Id, &todo.Task); err != nil {
			fmt.Println(err)
			writeServerErrorResponse(w, "scan row failed")
			return
		}
	} else {
		writeNotFoundResponse(w)
		return
	}

	_, err = c.DB.ExecContext(ctx, c.SQL.UpdateTodo(), dto.Task, todo.Id)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "update todo query failed")
		return
	}

	todo.Task = dto.Task
	res := make(map[string]Todo)
	res["todo"] = todo
	writeJsonResponse(w, res)
}

func (c *TodosController) ActionDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var todo Todo
	rows, err := c.DB.QueryContext(ctx, c.SQL.TodoById(), id)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "error in query")
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&todo.Id, &todo.Task); err != nil {
			fmt.Println(err)
			writeServerErrorResponse(w, "scan row failed")
			return
		}
	} else {
		writeNotFoundResponse(w)
		return
	}

	_, err = c.DB.ExecContext(ctx, c.SQL.DeleteTodo(), todo.Id)
	if err != nil {
		fmt.Println(err)
		writeServerErrorResponse(w, "update todo query failed")
		return
	}

	res := make(map[string]Todo)
	res["todo"] = todo
	writeJsonResponse(w, res)
}

func writeJsonResponse(w http.ResponseWriter, res any) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error: json marshall: %v", res)
	}

	w.Write(resJson)
	return
}

func writeServerErrorResponse(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")

	res := make(map[string]string)
	res["msg"] = "server error"
	res["reason"] = reason

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error: json marshall: %v", res)
	}

	w.Write(resJson)
	return
}

func writeNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")

	res := make(map[string]string)
	res["msg"] = "not found"

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error: json marshall: %v", res)
	}

	w.Write(resJson)
	return
}

func writeBadRequestResponse(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")

	res := make(map[string]string)
	res["msg"] = "bad request"
	res["reason"] = reason

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error: json marshall: %v", res)
	}

	w.Write(resJson)
	return
}
