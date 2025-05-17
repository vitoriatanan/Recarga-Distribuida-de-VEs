package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	functions "server/functions"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

// Informações sobre o servidor
var serverName string
var serverLocation []int
var carRoute []int
var mqttClient mqtt.Client

// ======== MQTT ========
/**
*  Inicializa a conexão MQTT, define opções de conexão e se inscreve no tópico de posições de carro.
*  @param: nenhum
*  @returns: nenhum
 */
func initMQTT() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://mosquitto:1883").
		SetClientID("server-" + functions.GetHostname())

	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}
	fmt.Printf("✅ %s conectado ao broker MQTT!\n", serverName)

	subscribeToCarPosition()
}

/**
*  Inscreve o servidor no tópico 'car/position' e define a função callback
*  para processar as mensagens recebidas, verificando se o carro está nos limites.
*  @param: nenhum
*  @returns: nenhum
 */
func subscribeToCarPosition() {
	if token := mqttClient.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		position := string(msg.Payload())
		fmt.Println("📥 Posição recebida do carro:", position)

		// Transformar a posição recebida de string para slice de inteiros
		var x, y, destX, destY int
		fmt.Sscanf(position, "%d, %d, %d, %d", &x, &y, &destX, &destY)
		carRoute = []int{x, y, destX, destY}

		// Verifica se a localização de origem do carro está nos limites do servidor
		if (functions.IsCarInCompanyLimits(carRoute, serverLocation)) {
			// Reserva posto
		} else {
			functions.SendPositionToOtherServers(x, y int)
		}

	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
	}
}

// ======== HTTP Server ========
/**
*  Inicia o servidor HTTP usando o framework Gin.
*  Cria a rota '/server/position' para retornar a posição atual do servidor.
*  @param: nenhum
*  @returns: nenhum
 */
func startHTTPServer() {
	router := gin.Default()

	router.GET("/server/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"x": serverLocation[0],
			"y": serverLocation[1],
		})
	})

	router.POST("/server/forward", func(c *gin.Context) {
		var req struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
	
		fmt.Printf("📨 Coordenada recebida: (%d, %d)\n", req.X, req.Y)
	
		if functions.IsPositionInCompanyLimits(req.X, req.Y, serverLocation) {
			fmt.Println("✅ Este servidor cobre a posição recebida. Pode atender o carro.")
			// selecionar ponto de recarga ou destino
		} else {
			fmt.Println("🚫 Fora da área de cobertura deste servidor.")
		}
	
		c.JSON(http.StatusOK, gin.H{"status": "recebido"})
	})	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// ======== FUNÇÃO PRINCIPAL ========
/**
*  Função principal do programa.
*  Inicializa a semente de números aleatórios, define o nome e localização do servidor,
*  conecta ao broker MQTT e inicia o servidor HTTP.
*  @param: nenhum
*  @returns: nenhum
 */
func main() {
	rand.Seed(time.Now().UnixNano())

	serverName = os.Getenv("INSTANCE_NAME")

	fmt.Println("🚀 Servidor:", serverName)

	serverLocation = functions.SetServerLocation(serverName, serverLocation)
	fmt.Println("📍 Localização do servidor:", serverLocation)

	initMQTT()

	startHTTPServer()
}
