package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:age`
}

var Users []User

func getUser(c *gin.Context) {
	c.JSON(http.StatusOK, Users)
}

func createUser(c *gin.Context) {

	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{
			"err": err,
			"msg": "invalid req body",
		})
		return
	}
	reqBody.ID = uuid.New().String()
	Users = append(Users, reqBody)
	c.JSON(200, gin.H{
		"err": false,
	})
}
func editUser(c *gin.Context) {
	var id = c.Param("id")
	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "Invalid msg request body",
		})
		return
	}

	for i, user := range Users {
		if user.ID == id {
			Users[i].Name = reqBody.Name
			Users[i].Age = reqBody.Age

			c.JSON(200, gin.H{
				"error": false,
			})
			return
		}
	}
	c.JSON(400, gin.H{
		"error":   true,
		"message": "Invalid user id",
	})
}
func deleteUser(c *gin.Context) {
	var id = c.Param("id")
	for i, user := range Users {
		if user.ID == id {
			Users = append(Users[i:], Users[i+1:]...)
			c.JSON(200, gin.H{
				"error": false,
			})
			return
		}
	}

	c.JSON(400, gin.H{
		"msg": "Invalid request body",
	})
}
func main() {
	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", getUser)
		userRoutes.POST("/add", createUser)
		userRoutes.PUT("/edit/:id", editUser)
		userRoutes.DELETE("/:id", deleteUser)
	}

	r.Run(":8080")
}
