package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const secretKey = "qwerty"

type Animal struct {
	AnimalId           string `json:"id"`
	Token              string `json:"token"`
	ConservationStatus int    `json:"status"`
}

func performPUTRequest(url string, data Animal) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "Application/json")
	req.Header.Set("Authorization", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}

// Функция для генерации случайного статуса
func randomStatus() int {
	time.Sleep(5 * time.Second) // Задержка на 5 секунд
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10)
}

// Функция для отправки статуса в отдельной горутине
func SendStatus(pk string, url string) {
	// Выполнение расчётов с randomStatus
	result := randomStatus()

	// Отправка PUT-запроса к основному серверу
	data := Animal{AnimalId: pk, Token: secretKey, ConservationStatus: result}
	response, err := performPUTRequest(url, data)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}
	if response.StatusCode == http.StatusOK {
		fmt.Println("Status sent successfully for pk:", pk)
		return
	} else {
		fmt.Println("Failed to process PUT request")
		return
	}
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.POST("/set-status", func(c *gin.Context) {
		// Получение значения "pk" из запроса
		pk := c.PostForm("pk")
		requestURL := fmt.Sprintf("http://127.0.0.1:8000/animals/%s/set-status/", pk)
		// Запуск горутины для отправки статуса
		go SendStatus(pk, requestURL) // Замените на ваш реальный URL

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
