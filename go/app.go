package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
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

func generatePeople(numPeople int, numWorkers int) []map[string]interface{} {
	people := make([]map[string]interface{}, numPeople)
	var wg sync.WaitGroup
	var rwMutex sync.RWMutex

	chunkSize := (numPeople + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			start := workerId * chunkSize
			end := start + chunkSize
			if end > numPeople {
				end = numPeople
			}

			for i := start; i < end; i++ {
				person := generatePerson()
				rwMutex.Lock()
				people[i] = person
				rwMutex.Unlock()
			}
		}(w)
	}

	wg.Wait()
	return people
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/people", func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano())

		numPeople := 5000
		numWorkers := runtime.NumCPU()
		fmt.Println(numWorkers, " workers")
		people := generatePeople(numPeople, numWorkers)

		peopleJSON, err := json.MarshalIndent(people[:20], "", "    ")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error formatting JSON")
			return
		}

		c.Data(http.StatusOK, "application/json", peopleJSON)
	})
	r.Run(":3000")
}
