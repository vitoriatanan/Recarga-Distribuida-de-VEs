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
var serverName string
var serverLocation []int
var carRoute []int
var mqttClient mqtt.Client

// ======== FUNÇÕES UTILITÁRIAS ========

/**
*  Obtém o nome do host atual da máquina ou da variável de ambiente INSTANCE_NAME.
*  @param: nenhum
*  @returns: 
*     - (string): nome do host ou da instância.
*/
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

/**
*  Define os limites da localização do servidor de acordo com seu nome.
*  Os limites são configurados em pares de inteiros, que indicam a área em um grid.
*  @param: nenhum
*  @returns: nenhum
*/
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

/**
*  Verifica se a posição atual do carro está dentro dos limites da empresa.
*  @param:
*     - carLocation ([]int): slice contendo as posições x, y, destX, destY do carro.
*  @returns: 
*     - (bool): true se a posição estiver dentro dos limites, false caso contrário.
*/
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
/**
*  Inicializa a conexão MQTT, define opções de conexão e se inscreve no tópico de posições de carro.
*  @param: nenhum
*  @returns: nenhum
*/
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
		isCarInCompanyLimits(carRoute)

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// Função para enviar mensagens entre servidores pela API
// func postLocation(serverName string, message string) {
// 	url := fmt.Sprintf("http://%s:8080/server/position", serverName)
// 	resp, err := http.Post(url, "application/json", nil)
// 	if err != nil {
// 		log.Fatalf("Erro ao enviar mensagem para %s: %v", serverName, err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("Erro ao enviar mensagem para %s: %s", serverName, resp.Status)
// 	} else {
// 		fmt.Printf("✅ Mensagem enviada para %s: %s\n", serverName, message)
// 	}
// }

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

	setServerLocation()
	fmt.Println("📍 Localização do servidor:", serverLocation)

	initMQTT()

	startHTTPServer()
}
