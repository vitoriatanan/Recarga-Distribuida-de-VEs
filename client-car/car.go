package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	brokerURL = "tcp://mosquitto:1883"
	clientID  = "carroA"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	client := connectMQTT(brokerURL, clientID)
	defer disconnectMQTT(client)

	subscribe(client, "car/recarga", defaultMessageHandler)
	startCarLoop(client)
}

func connectMQTT(broker, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetDefaultPublishHandler(defaultMessageHandler)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Erro ao conectar ao broker: %v", token.Error())
	}

	fmt.Println("✅ Conectado ao broker MQTT!")
	return client
}

func disconnectMQTT(client mqtt.Client) {
	client.Disconnect(250)
	fmt.Println("🚪 Cliente desconectado.")
}

func publish(client mqtt.Client, topic string, payload string) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("⚠️ Erro ao publicar no tópico %s: %v", topic, err)
	}
	fmt.Printf("📄 Publicado no tópico %s\n", topic)
}

func subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Fatalf("❌ Erro ao se inscrever no tópico %s: %v", topic, err)
	}
	fmt.Printf("📡 Inscrito no tópico %s\n", topic)
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("📩 Mensagem recebida [%s]: %s\n", msg.Topic(), msg.Payload())
}

func generatePosition() (int, int, int, int) {

	originX := rand.Intn(1000)      // coordenada X entre 0 e 999
	originY := rand.Intn(1000)      // coordenada Y entre 0 e 999
	destinationX := rand.Intn(1000) // coordenada X entre 0 e 999
	destinationY := rand.Intn(1000) // coordenada Y entre 0 e 999

	return originX, originY, destinationX, destinationY
}

func startCarLoop(client mqtt.Client) {
	for {
		origX, origY, destX, destY := generatePosition()
		route := fmt.Sprintf("%d, %d, %d, %d", origX, origY, destX, destY)

		// Envia para ambos os tópicos
		publish(client, "car/position", route)
		fmt.Println("📤 Enviado para car/position:", route)

		publish(client, "car/recarga", route)
		fmt.Println("📤 Enviado para car/recarga:", route)

		time.Sleep(2 * time.Second)
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// )

// func main() {
// 	// Semente para gerar posições diferentes
// 	rand.Seed(time.Now().UnixNano())

// 	// Criando um novo cliente MQTT
// 	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883").SetClientID("car")
// 	client := mqtt.NewClient(opts)

// 	// Conecta ao broker MQTT
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		log.Fatalf("Erro de conexão MQTT: %v", token.Error())
// 	}

// 	fmt.Println("✅ Carro conectado ao broker MQTT!")

// 	defer client.Disconnect(250)

// 	// Gera coordenadas aleatórias e envia para o tópico car/position
// 	for {

// 		// Gera coordenadas aleatórias de inicio e destino da rota
// 		x_begin := rand.Intn(1000) // coordenada X entre 0 e 999
// 		y_begin := rand.Intn(1000) // coordenada Y entre 0 e 999

// 		x_end := rand.Intn(1000) // coordenada X entre 0 e 999
// 		y_end := rand.Intn(1000) // coordenada Y entre 0 e 999

// 		position := fmt.Sprintf("Carro - saída: (%d, %d) chegada: (%d, %d)", x_begin, y_begin, x_end, y_end) // cria a string com as coordenadas

// 		//car_route := route_generator()      // gera uma rota aleatória para o carro
// 		//fmt.Println(" - Rota: ", car_route) // adiciona a rota ao vetor de posições

// 		token := client.Publish("car/position", 0, false, position) // publica a posição no tópico car/position

// 		token.Wait()

// 		fmt.Println("📤 Enviado:", position)
// 		time.Sleep(2 * time.Second)
// 	}
// }
