package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


type Car struct {
	ID            string
	X, Y          float64
	DestinationX  float64
	DestinationY  float64
	Battery  	  float64
}


// ======= tem que gerar para a betrry tambéemmmmmmm ===========

// Gera coordenadas aleatórias entre os limites especificados
func generateRandomCoordinate(minLimit, maxLimit float64) float64 {
	return minLimit + rand.Float64()*(maxLimit-minLimit)
}


func main() {
	// Semente para gerar posições diferentes
	rand.Seed(time.Now().UnixNano())

	// Áreas: Empresa A: 0–100, Empresa B: 100–200, Empresa C: 200–300
	// o eixo X quem define as áreas das empresas. O y é só um complemento
	originX := generateRandomCoordinate(0, 300)
	destinationX := generateRandomCoordinate(0, 300)
	for math.Abs(destinationX-originX) < 50 {
		destinationX = generateRandomCoordinate(0, 300) // evitar origem = destino muito próximos
	}
	originY := generateRandomCoordinate(0, 100)
	destinationY := generateRandomCoordinate(0, 100)


	// ====== MUDAR O ID PARA INT E GERAR ELE AUTOMATICAMENTE E GERAR O DESCARREGAMENTO DA BATERIA  ======= 
	car := Car {
		ID:             "C1",
		X:              originX,
		Y:              originY,
		DestinationX:   destinationX,
		DestinationY:   destinationY,
		Battery:        100,
	}

	fmt.Printf("Origem: (%.2f, %.2f) → Destino: (%.2f, %.2f)\n", originX, originY, destinationX, destinationY)
	

	// Criando um novo cliente MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt:1883").SetClientID("car")
	client := mqtt.NewClient(opts)
	
	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Carro conectado ao broker MQTT!")

	
		
	


	defer client.Disconnect(250)

	
	
}