package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin" // Framework Gin para criar APIs HTTP
)

func getHostname() string {
	// Obtém o nome pela variável de ambiente HOSTNAME
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Erro ao obter o nome do host: %v", err)
	}
	// Se o nome do host não estiver definido, usa o nome do container
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}
	return hostname
}

// func getLocation(c *gin.Context) []string {
// 	return
// }

func main() {
	// Semente para gerar posições diferentes
	rand.Seed(time.Now().UnixNano())

	serverName := getHostname() // Obtém o nome do servidor
	fmt.Println("🚀 Servidor:", serverName)

	x := rand.Intn(1000) // coordenada X entre 0 e 999
	y := rand.Intn(1000) // coordenada Y entre 0 e 999

	router := gin.Default() // Cria um novo roteador Gin
	router.GET("/car/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"x": x, "y": y}) // Retorna a posição do servidor para os outros servidores
	})

	router.Run("localhost:8080") // Inicia o servidor HTTP na porta 8080

	opts := mqtt.NewClientOptions().
		AddBroker("tcp://mosquitto:1883").
		SetClientID("server-" + getHostname()) // Usa o nome do container como parte do ClientID
	client := mqtt.NewClient(opts) // Cria um novo cliente MQTT

	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}
	fmt.Println("✅ Servidor A conectado ao broker MQTT!")

	defer client.Disconnect(250) // Desconecta do broker MQTT após o término

	// Inscreve-se no tópico "car/position" para receber mensagens do carro
	// Quando uma mensagem é recebida, ela é processada na função de callback
	if token := client.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		position := string(msg.Payload())
		fmt.Println("📥 Posição recebida do carro:", position)
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
	}
}

// func sendReservationRequest(url string, position string) bool {
// 	jsonData := []byte(fmt.Sprintf(`{"car_position":"%s"}`, position))
// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		log.Println("Erro ao enviar requisição para", url, "->", err)
// 		return false
// 	}
// 	defer resp.Body.Close()
// 	return resp.StatusCode == http.StatusOK
// }
