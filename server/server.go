package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

// Informações sobre o servidor
var serverName = os.Getenv("INSTANCE_NAME")
var serverLocation []int
var carRoute []int
var mqttClient mqtt.Client

// ======== Funções utilitárias ========

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Erro ao obter o nome do host: %v", err)
	}
	if hostname == "" {
		hostname = os.Getenv("INSTANCE_NAME")
	}
	return hostname
}

func setServerLocation() {
	switch serverName {
	case "empresa-a":
		serverLocation = []int{0, 250}
	case "empresa-b":
		serverLocation = []int{251, 500}
	case "empresa-c":
		serverLocation = []int{501, 750}
	default:
		serverLocation = []int{-1, -1}
	}
}

// Função que verifica se a localização do carro está dentro do limite da empresa
func isCarInCompanyLimits(carLocation []int) bool {
	if len(carLocation) != 4 {
		return false
	}

	x, y := carLocation[0], carLocation[1]
	//destX, destY := carLocation[2], carLocation[3]

	if x >= serverLocation[0] && x <= serverLocation[1] && y >= serverLocation[0] && y <= serverLocation[1] {
		fmt.Println("🚗 O carro está dentro dos limites da empresa.")
		return true
		// if destX >= serverLocation[0] && destX <= serverLocation[1] && destY >= serverLocation[0] && destY <= serverLocation[1] {
		// 	return true
		// }
	}
	fmt.Println("🚗 O carro está fora dos limites da empresa.")
	return false
}

// ======== MQTT ========

func initMQTT() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://mosquitto:1883").
		SetClientID("server-" + getHostname())

	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}
	fmt.Printf("✅ %s conectado ao broker MQTT!\n", serverName)

	subscribeToCarPosition()
}

func subscribeToCarPosition() {
	if token := mqttClient.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		position := string(msg.Payload())
		fmt.Println("📥 Posição recebida do carro:", position)

		// Transformar a posição recebida de string para slice de inteiros
		var x, y, destX, destY int
		fmt.Sscanf(position, "%d, %d, %d, %d", &x, &y, &destX, &destY)
		carRoute = []int{x, y, destX, destY}
		isCarInCompanyLimits(carRoute)

	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
	}
}

// ======== HTTP Server ========

func startHTTPServer() {
	router := gin.Default()

	router.GET("/server/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"x": serverLocation[0],
			"y": serverLocation[1],
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// ======== Função Principal ========

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🚀 Servidor:", serverName)

	setServerLocation()
	fmt.Println("📍 Localização do servidor:", serverLocation)

	initMQTT()

	startHTTPServer()
}
