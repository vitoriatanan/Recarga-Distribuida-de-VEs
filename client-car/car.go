package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Car struct {
	ID               int
	X, Y float64
	DestinationX     float64
	DestinationY     float64
	Battery          int
}

func generateRandomCoordinate(minLimit, maxLimit float64) float64 {
	return minLimit + rand.Float64()*(maxLimit-minLimit)
}

// Gera um ID único para cada carro
func getCarID() int {
	return int(time.Now().UnixNano() % 10000)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	carID := getCarID() // Gera um ID dinâmico para o carro

	// Áreas: Empresa A: 0–100, Empresa B: 100–200, Empresa C: 200–300
	// o eixo X quem define as áreas das empresas. O y é só um complemento
	 X := generateRandomCoordinate(0, 300)
	destinationX := generateRandomCoordinate(0, 300)
	for math.Abs(destinationX- X) < 50 {
		destinationX = generateRandomCoordinate(0, 300) // evitar origem = destino muito próximos

	}
	 Y := generateRandomCoordinate(0, 100)
	destinationY := generateRandomCoordinate(0, 100)

	car := Car{
		ID:           carID,
		X:            X,
		Y:            Y,
		DestinationX: destinationX,
		DestinationY: destinationY,
		Battery:      rand.Intn(51) + 50,
	}

	fmt.Printf("Origem: (%.2f, %.2f) → Destino: (%.2f, %.2f)\n",  X,  Y, destinationX, destinationY)

	// Criando um novo cliente MQTT e conectando ao broker
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt:1883").SetClientID("car")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
	}

	fmt.Println("✅ Carro conectado ao broker MQTT!")
	defer client.Disconnect(250)

	topic := "veiculo/recarga"
	step := 5.0    // passo de movimento do carro
	for car.Battery > 0 {
		// Calcula o quanto falta para chegar no destino
		dx := car.DestinationX - car.X  //distância restante nos eixos X e Y.
		dy := car.DestinationY - car.Y
		dist := math.Hypot(dx, dy) // distância total (pitagórica)
		if dist < 5.0 {  // critério de parada. Já que pode acontecer do carro não chegar na coordenada exata
			fmt.Println("Carro chegou ao destino")
			break
		}
		dirX := dx / dist
		dirY := dy / dist
		car. X += dirX * step
		car. Y += dirY * step
		car.Battery -= rand.Intn(5) + 1 // Consumo entre 1 e 3 por passo

		fmt.Printf("[C%d] Movendo para (%.2f, %.2f) | Bateria: %d%%\n", car.ID, car.X, car.Y, car.Battery)


		// Cria uma mensagem JSON com os dados do carro
		msgBytes, _ := json.Marshal(map[string]interface{}{
			"car_id": fmt.Sprintf("C%d", car.ID),
			"x":      car.X,
			"y":      car.Y,
			"bateria": car.Battery,
		})
		// Publica a mensagem no tópico
		client.Publish(topic, 0, false, msgBytes)
		time.Sleep(2 * time.Second)
	}
}
