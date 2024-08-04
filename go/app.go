package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateEmail() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	const usernameLength = 7
	const domainLength = 5

	randomString := func(length int) string {
		result := make([]byte, length)
		for i := range result {
			result[i] = charset[rand.Intn(len(charset))]
		}
		return string(result)
	}

	username := randomString(usernameLength)
	domain := randomString(domainLength)

	return username + "@" + domain + ".com"
}

func generatePerson() map[string]interface{} {
	email := generateEmail()
	age := rand.Intn(100) + 1
	return map[string]interface{}{"email": email, "age": age}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	numCPU := runtime.NumCPU()

	// Imprimir el número de CPUs
	fmt.Printf("Número de CPUs disponibles: %d\n", numCPU)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/people", func(c *gin.Context) {
		people := make([]map[string]interface{}, 5000)
		for i := 0; i < 5000; i++ {
			people[i] = generatePerson()
		}

		// Format JSON with indentation
		peopleJSON, err := json.MarshalIndent(people[:20], "", "    ")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error formatting JSON")
			return
		}

		// Return pretty-printed JSON
		c.Data(http.StatusOK, "application/json", peopleJSON)
	})
	r.Run(":3000")
}
