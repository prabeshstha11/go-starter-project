package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func connectDb() {
	var err error
	db, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to sqlite database...")
}

func createTable() {
	table := "CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY AUTOINCREMENT, item TEXT, isCompleted INTEGER)"

	_, err := db.Exec(table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("table created...")
}

type Todo struct {
	ID          int    `json:"id"`
	Item        string `json:"item"`
	IsCompleted bool   `json:"isCompleted"`
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	connectDb()
	defer db.Close()
	createTable()

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/status", func(c *gin.Context) {
		c.String(200, "Status OK")
	})

	// create item
	router.POST("/create", func(c *gin.Context) {
		var todo Todo

		if err := c.BindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO todo (item, isCompleted) VALUES (?, ?)")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(todo.Item, boolToInt(todo.IsCompleted))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		todo.ID = int(id)

		c.JSON(200, gin.H{
			"message": "Todo created successfully",
			"todo":    todo,
		})

	})

	// get item
	router.GET("/todo", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, item, isCompleted FROM todo")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var todos []Todo

		for rows.Next() {
			var todo Todo
			var isCompletedInt int
			if err := rows.Scan(&todo.ID, &todo.Item, &isCompletedInt); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			todo.IsCompleted = isCompletedInt == 1
			todos = append(todos, todo)
		}

		c.JSON(200, gin.H{
			"todos": todos,
		})

	})

	// get item by id
	router.GET("/todo/:id", func(c *gin.Context) {
		id := c.Param("id")
		var todo Todo
		var isCompletedInt int
		err := db.QueryRow("SELECT id, item, isCompleted FROM todo WHERE id = ?", id).Scan(&todo.ID, &todo.Item, &isCompletedInt)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Invalid ID"})
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		todo.IsCompleted = isCompletedInt == 1
		c.JSON(200, gin.H{"todo": todo})
	})

	// edit
	router.PATCH("/todo/:id", func(c *gin.Context) {
		id := c.Param("id")
		var todo Todo

		// bind JSON (partial fields allowed)
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// fetch current todo first
		var current Todo
		var isCompletedInt int
		err := db.QueryRow("SELECT id, item, isCompleted FROM todo WHERE id = ?", id).
			Scan(&current.ID, &current.Item, &isCompletedInt)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		current.IsCompleted = isCompletedInt == 1

		// decide new values
		if todo.Item != "" {
			current.Item = todo.Item
		}
		// if IsCompleted was sent in JSON (not default false)
		current.IsCompleted = todo.IsCompleted || current.IsCompleted

		// update in DB
		stmt, err := db.Prepare("UPDATE todo SET item = ?, isCompleted = ? WHERE id = ?")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(current.Item, boolToInt(current.IsCompleted), id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Todo updated", "todo": current})
	})

	// delete
	router.DELETE("/todo/:id", func(c *gin.Context) {
		id := c.Param("id")

		stmt, err := db.Prepare("DELETE FROM todo WHERE id = ?")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(200, gin.H{"message": "Todo deleted successfully"})
	})

	router.Run(":8800")
}
