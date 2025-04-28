package main

import (
	"fmt"
	"io"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		serverName := os.Getenv("SERVER_NAME")
		c.STRING(http.StatusOK, fmt.Sprintf("pong from %s", serverName))
	})

	router.GET("/call", func(c *gin.Context) {
		otherServer := os.Getenv("SERVER_NAME")
		if otherServer == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OTHER_SERVER_URL not set"})
			return
		}

		resp, err := http.Get(otherServer + "/ping")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.String(http.StatusOK, "Responde from other server: %s", string(body))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)

	// 	// Conetando ao broker MQTT
	// 	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("server")
	// 	client := mqtt.NewClient(opts)

	// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
	// 		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	// 	}

	// 	fmt.Println("✅ Servidor conectado ao broker MQTT!")

	// 	defer client.Disconnect(250)

	// 	// Assina no tópico car/position
	// 	if token := client.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
	// 		fmt.Printf("📥 Recebido: %s\n", msg.Payload())
	// 	}); token.Wait() && token.Error() != nil {
	// 		panic(token.Error())
	// 	}

	// select {} // mantém o programa rodando
}
