package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var all_cities = []string{"São Paulo", "Rio de Janeiro", "Belo Horizonte", "Salvador", "Curitiba"} // cidades disponíveis para os carros

func route_generator() []string { // ARRUMAR ESSA FUNÇAO QUE TÁ ESTRANHA, ROTAS BEM RUINS!!!
	numberOfCities := rand.Intn(4) + 2 // número de cidades entre 2 e 5
	route := make([]string, 0, numberOfCities)
	usedCities := make(map[string]bool)

	for len(route) < numberOfCities {
		randomCity := all_cities[rand.Intn(len(all_cities))]

		if !usedCities[randomCity] {
			route = append(route, randomCity)
			usedCities[randomCity] = true
		}
	}

	return route
}

func main() {
	// Semente para gerar posições diferentes
	rand.Seed(time.Now().UnixNano())

	// Criando um novo cliente MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("car")
	client := mqtt.NewClient(opts)

	// Conecta ao broker MQTT
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Carro conectado ao broker MQTT!")

	defer client.Disconnect(250)

	// Gera coordenadas aleatórias e envia para o tópico car/position
	for {
		x := rand.Intn(1000) // coordenada X entre 0 e 999
		y := rand.Intn(1000) // coordenada Y entre 0 e 999

		position := fmt.Sprintf("Carro A - posição x:%d, y:%d", x, y)

		//car_route := route_generator()      // gera uma rota aleatória para o carro
		//fmt.Println(" - Rota: ", car_route) // adiciona a rota ao vetor de posições

		token := client.Publish("car/position", 0, false, position) // publica a posição no tópico car/position

		token.Wait()

		fmt.Println("📤 Enviado:", position)
		time.Sleep(2 * time.Second)
	}
}
