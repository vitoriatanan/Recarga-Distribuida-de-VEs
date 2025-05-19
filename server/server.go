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

// ======== DICINOÁRIO DE POSTOS E SUAS RESERVA ========
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
		var origX, origY, destX, destY int
		var carID int
		fmt.Sscanf(position, "%d, %d, %d, %d, %d", &carID, &origX, &origY, &destX, &destY)
		// parts := strings.Split(position, ",")
		// if len(parts) == 5 {
		// 	carID = strings.TrimSpace(parts[0])
		// 	origX, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
		// 	origY, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
		// 	destX, _ = strconv.Atoi(strings.TrimSpace(parts[3]))
		// 	destY, _ = strconv.Atoi(strings.TrimSpace(parts[4]))
		// }

		carRoute = []int{carID, origX, origY, destX, destY}
		//fmt.Printf("%d, %d, %d, %d, %d", carID, origX, origY, destX, destY)

		// Verifica se a localização de origem do carro está nos limites do servidor
		if functions.IsPositionInCompanyLimits(origX, origY, serverLocation) {

			// Se a posição estiver dentro dos limites, envia para o tópico de recarga !!!!!!!!!!!!! FAZER ISSO

			fmt.Println("✅ Este servidor cobre a posição de origem recebida. Pode atender o carro.")

			//Procura um ponto de recarga disponível e reserva
			first_station := functions.StationReservation(carID, stationsSpots)
			stationsSpots[first_station] = carID
			fmt.Printf("🚗 Primeiro ponto de recarga reservado na %s\n", first_station)

			//Verifica se a posição de destino do carro está nos limites do servidor
			if functions.IsPositionInCompanyLimits(destX, destY, serverLocation) {
				fmt.Println("✅ Este servidor cobre a posição de destino recebida. Pode atender o carro.")

				//Procura um ponto de recarga disponível e reserva
				second_station := functions.StationReservation(carID, stationsSpots)
				stationsSpots[second_station] = carID
				fmt.Printf("🚗 Segundo ponto de recarga reservado na %s\n", second_station)

			} else {
				fmt.Println("🚫 Destino da viagem fora da área de cobertura deste servidor.")

				// Envia localização de destino do carro para os outros servidores
				functions.SendPositionToServers(destX, destY, serverName)
			}

		} else {
			fmt.Println("🚫 Origem da viagem fora da área de cobertura deste servidor.")
		}

	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao se inscrever no tópico: %v", token.Error())
	}
}

// ======== HTTP Server ========
/**
*  Inicia o servidor HTTP usando o framework Gin.
*  @param: nenhum
*  @returns: nenhum
 */
func startHTTPServer() {
	router := gin.Default()

	router.GET("/server/position", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Min_x": serverLocation[0],
			"Max_x": serverLocation[1],
			"Min_y": serverLocation[2],
			"Max_y": serverLocation[3],
		})
	})

	router.GET("/server/route", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"origX": carRoute[0],
			"origY": carRoute[1],
			"destX": carRoute[2],
			"destY": carRoute[3],
		})
	})

	router.GET("/server/stations", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"stationsSpots": stationsSpots,
		})
	})

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

			// Reservar ponto de recarga
			station := functions.StationReservation(carRoute[0], stationsSpots)

			stationsSpots[station] = carRoute[0]
			fmt.Printf("🚗 Segundo ponto de recarga reservado na %s\n", station)

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
