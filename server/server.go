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

// ======== VARIÁVEIS GLOBAIS ========

// Nome do servidor
var serverName string

// Coordenadas de cobertura do servidor no formato [minX, maxX, minY, maxY]
var serverLocation []int

// Última rota de carro recebida: [carID, origX, origY, destX, destY]
var carRoute []int

// Cliente MQTT
var mqttClient mqtt.Client

// Dicionário de estações e seus status de reserva
var stationsSpots = map[string]int{
	"station1":  0,
	"station2":  0,
	"station3":  0,
	"station4":  0,
	"station5":  0,
	"station6":  0,
	"station7":  0,
	"station8":  0,
	"station9":  0,
	"station10": 0,
}

// ======== FUNÇÕES MQTT ========

/**
 * Notifica o carro da confirmação de reserva de pontos de recarga.
 *
 * @param carID ID do carro a ser notificado.
 * @param firstStation nome do primeiro ponto reservado.
 * @param secondStation nome do segundo ponto reservado.
 */
func notifyCarReservation(carID int, firstStation, secondStation string) {
	topic := fmt.Sprintf("car/%d/reservation", carID)
	message := fmt.Sprintf("Reserva confirmada: %s e %s", firstStation, secondStation)

	token := mqttClient.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Erro ao notificar carro %d: %v\n", carID, token.Error())
	} else {
		fmt.Printf("📤 Notificação enviada para carro %d no tópico %s\n", carID, topic)
	}
}

/**
 * Inicializa a conexão MQTT e se inscreve no tópico de posições dos carros.
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
 * Inscreve o servidor no tópico "car/position" e processa as mensagens recebidas.
 * Atribui postos se a origem e destino estiverem dentro da área de cobertura.
 */
func subscribeToCarPosition() {
	if token := mqttClient.Subscribe("car/position", 0, func(client mqtt.Client, msg mqtt.Message) {
		position := string(msg.Payload())
		fmt.Println("📥 Posição recebida do carro:", position)

		var origX, origY, destX, destY int
		var carID int
		fmt.Sscanf(position, "%d, %d, %d, %d, %d", &carID, &origX, &origY, &destX, &destY)
		carRoute = []int{carID, origX, origY, destX, destY}

		if functions.IsPositionInCompanyLimits(origX, origY, serverLocation) {
			fmt.Println("✅ Origem dentro da cobertura do servidor.")

			first_station := functions.StationReservation(carID, stationsSpots)
			stationsSpots[first_station] = carID
			fmt.Printf("🚗 Primeiro ponto de recarga reservado na %s\n", first_station)

			if functions.IsPositionInCompanyLimits(destX, destY, serverLocation) {
				fmt.Println("✅ Destino dentro da cobertura do servidor.")

				second_station := functions.StationReservation(carID, stationsSpots)
				stationsSpots[second_station] = carID
				fmt.Printf("🚗 Segundo ponto de recarga reservado na %s\n", second_station)

				notifyCarReservation(carID, first_station, second_station)
			} else {
				fmt.Println("🚫 Destino fora da cobertura. Encaminhando a outro servidor.")
				functions.SendPositionToServers(destX, destY, serverName)
			}
		} else {
			fmt.Println("🚫 Origem fora da área de cobertura deste servidor.")
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
	}
}

// ======== SERVIDOR HTTP ========

/**
 * Inicia o servidor HTTP com rotas de status e controle.
 */
func startHTTPServer() {
	router := gin.Default()

	// Retorna os limites da área coberta pelo servidor
	router.GET("/server/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Min_x": serverLocation[0],
			"Max_x": serverLocation[1],
			"Min_y": serverLocation[2],
			"Max_y": serverLocation[3],
		})
	})

	// Retorna a última rota processada
	router.GET("/server/route", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"origX": carRoute[1],
			"origY": carRoute[2],
			"destX": carRoute[3],
			"destY": carRoute[4],
		})
	})

	// Retorna status de todas as estações
	router.GET("/server/stations", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"stationsSpots": stationsSpots,
		})
	})

	// Retorna status de uma estação específica
	router.GET("/server/stations/:station", func(c *gin.Context) {
		station := c.Param("station")
		if spots, ok := stationsSpots[station]; ok {
			c.JSON(http.StatusOK, gin.H{
				"station": station,
				"spots":   spots,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Estação não encontrada"})
		}
	})

	// Recebe posição de outro servidor para tentativa de reserva
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
			fmt.Println("✅ Este servidor cobre a posição recebida.")

			station := functions.StationReservation(carRoute[0], stationsSpots)
			stationsSpots[station] = carRoute[0]
			fmt.Printf("🚗 Segundo ponto de recarga reservado na %s\n", station)
		} else {
			fmt.Println("🚫 Fora da cobertura deste servidor.")
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
 * Função principal.
 * Define aleatoriedade, nome e área do servidor, conecta ao MQTT e inicia o servidor HTTP.
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
