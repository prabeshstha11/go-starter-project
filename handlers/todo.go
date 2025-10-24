package handlers

import (
	"net/http"

	"github.com/prabeshstha11/gotodo/models"

	"github.com/prabeshstha11/gotodo/db"

	"github.com/gin-gonic/gin"
)

// Create Todo
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.DB.Prepare("INSERT INTO todo (item, isCompleted) VALUES (?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, _ := stmt.Exec(todo.Item, models.BoolToInt(todo.IsCompleted))
	id, _ := result.LastInsertId()
	todo.ID = int(id)

	c.JSON(200, gin.H{"message": "Todo created", "todo": todo})
}

// Get All Todos
func GetTodos(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, item, isCompleted FROM todo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		var isCompleted int
		rows.Scan(&todo.ID, &todo.Item, &isCompleted)
		todo.IsCompleted = isCompleted == 1
		todos = append(todos, todo)
	}

	c.JSON(200, gin.H{"todos": todos})
}

// Get Todo by ID
func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	var isCompleted int
	err := db.DB.QueryRow("SELECT id, item, isCompleted FROM todo WHERE id = ?", id).
		Scan(&todo.ID, &todo.Item, &isCompleted)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	todo.IsCompleted = isCompleted == 1
	c.JSON(200, gin.H{"todo": todo})
}

// Update Todo
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var current models.Todo
	var isCompleted int
	err := db.DB.QueryRow("SELECT id, item, isCompleted FROM todo WHERE id = ?", id).
		Scan(&current.ID, &current.Item, &isCompleted)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	current.IsCompleted = isCompleted == 1

	if todo.Item != "" {
		current.Item = todo.Item
	}
	current.IsCompleted = todo.IsCompleted || current.IsCompleted

	stmt, _ := db.DB.Prepare("UPDATE todo SET item = ?, isCompleted = ? WHERE id = ?")
	defer stmt.Close()
	stmt.Exec(current.Item, models.BoolToInt(current.IsCompleted), id)

	c.JSON(200, gin.H{"message": "Todo updated", "todo": current})
}

// Delete Todo
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	stmt, _ := db.DB.Prepare("DELETE FROM todo WHERE id = ?")
	defer stmt.Close()
	res, _ := stmt.Exec(id)

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Todo deleted"})
}
