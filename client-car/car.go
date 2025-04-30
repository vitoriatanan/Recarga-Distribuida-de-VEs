package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var all_cities = []string{"São Paulo", "Rio de Janeiro", "Belo Horizonte", "Salvador", "Curitiba"} // cidades disponíveis para os carros

/*GERA ROTAS ALEATÓRIAS PARA OS CARROS*/
func route_generator() []string {
	numberOfCities := rand.Intn(4) + 2      // número de cidades entre 2 e 5
	route := make([]string, numberOfCities) // vetor de rotas com o número de cidades

	for i := 0; i < numberOfCities; i++ {
		// Adiciona cidades aleatórias ao vetor de rotas
		randon_city := rand.Intn(len(all_cities)) // número aleatório entre 0 e o número de cidades disponíveis
		route[i] = all_cities[randon_city]        // número aleatório entre 0 e o número de cidades disponíveis
	}

	return route // retorna o vetor de rotas
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

		car_route := route_generator()        // gera uma rota aleatória para o carro
		fmt.Sprintf(" - Rota: %v", car_route) // adiciona a rota ao vetor de posições

		token := client.Publish("car/position", 0, false, position)

		token.Wait()

		fmt.Println("📤 Enviado:", position)
		time.Sleep(2 * time.Second)
	}
}
