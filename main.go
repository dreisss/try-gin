package main

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type Todo struct {
	ID   uuid.UUID
	Done bool
	Text string `json:"text"`
}

func main() {
	var todos []Todo

	GIN := gin.Default()

	GIN.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})

	GIN.POST("/todos", func(c *gin.Context) {
		newTodoID, err := uuid.NewV7()

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var newTodo Todo
		err = c.BindJSON(&newTodo)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "empty text"})
			return
		}

		newTodo.ID = newTodoID

		todos = append(todos, newTodo)

		c.JSON(http.StatusOK, newTodo)
	})

	GIN.DELETE("/todos/:id", func(c *gin.Context) {
		requiredTodoID, err := uuid.FromString(c.Param("id"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
			return
		}

		requiredTodoIndex := slices.IndexFunc(todos, func(todo Todo) bool { return todo.ID == requiredTodoID })

		if requiredTodoIndex == -1 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}

		requiredTodo := todos[requiredTodoIndex]

		todos = slices.Delete(todos, requiredTodoIndex, requiredTodoIndex+1)

		c.JSON(http.StatusOK, requiredTodo)
	})

	GIN.PUT("/todos/:id", func(c *gin.Context) {
		requiredTodoID, err := uuid.FromString(c.Param("id"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
			return
		}

		requiredTodoIndex := slices.IndexFunc(todos, func(todo Todo) bool { return todo.ID == requiredTodoID })

		if requiredTodoIndex == -1 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
			return
		}

		requiredTodo := &todos[requiredTodoIndex]

		(*requiredTodo).Done = !(*requiredTodo).Done

		c.JSON(http.StatusOK, requiredTodo)
	})

	GIN.Run()
}
