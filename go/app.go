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

	// Dividir el trabajo entre los trabajadores
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

func generatePeople_0(numPeople int, numWorkers int) []map[string]interface{} {
	people := make([]map[string]interface{}, numPeople)
	jobs := make(chan int, numPeople)
	results := make(chan map[string]interface{}, numPeople)

	var wg sync.WaitGroup

	// Crear trabajadores
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				results <- generatePerson()
			}
		}()
	}

	// Enviar trabajos a los trabajadores
	go func() {
		for i := 0; i < numPeople; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Recoger los resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	// Almacenar los resultados en el slice de people
	i := 0
	for result := range results {
		people[i] = result
		i++
	}

	return people
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Imprimir el número de CPUs
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/people", func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano())

		numPeople := 45000
		numWorkers := 4 //runtime.NumCPU() // Número de núcleos disponibles
		fmt.Println(numWorkers, " workers")
		people := generatePeople(numPeople, numWorkers)

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
