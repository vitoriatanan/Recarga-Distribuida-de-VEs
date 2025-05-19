package functions

import (
	"fmt"
	"log"
	// "math/rand"
	"net/http"
	"os"
	"strings"
	// "time"

	//mqtt "github.com/eclipse/paho.mqtt.golang"
	//"github.com/gin-gonic/gin"
)

// ======== FUNÇÕES UTILITÁRIAS ========

/**
*  Obtém o nome do host atual da máquina ou da variável de ambiente INSTANCE_NAME.
*  @param: nenhum
*  @returns:
*     - (string): nome do host ou da instância.
 */
func GetHostname() string {
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
func SetServerLocation(serverName string, serverLocation []int) []int {
	switch serverName {
	case "empresa-a":
		serverLocation = []int{0, 250}
		return serverLocation
	case "empresa-b":
		serverLocation = []int{251, 500}
		return serverLocation
	case "empresa-c":
		serverLocation = []int{501, 750}
		return serverLocation
	default:
		serverLocation = []int{-1, -1}
		return serverLocation
	}
}

/**
*  Verifica se a posição atual do carro está dentro dos limites da empresa.
*  @param:
*     - carLocation ([]int): slice contendo as posições x, y, destX, destY do carro.
*  @returns:
*     - (bool): true se a posição estiver dentro dos limites, false caso contrário.
 */
func IsPositionInCompanyLimits(origX, origY int, serverLocation []int) bool {
	//x, y := carLocation[0], carLocation[1]
	//destX, destY := carLocation[2], carLocation[3]

	if origX >= serverLocation[0] && origX <= serverLocation[1] && origY >= serverLocation[0] && origY <= serverLocation[1] {
		fmt.Println("🚗 O carro está dentro dos limites da empresa.")
		return true
		// if destX >= serverLocation[0] && destX <= serverLocation[1] && destY >= serverLocation[0] && destY <= serverLocation[1] {
		// 	return true
		// }
	}
	fmt.Println("🚗 O carro está fora dos limites da empresa.")
	return false
}

func SendPositionToServers(x, y int, serverName string) {
	servers_ip := []string{"8081", "8082", "8083"}
	servers := []string{"empresa-a", "empresa-b", "empresa-c"}

	for i := range servers {
		name := servers[i]
		ip := servers_ip[i]
		
		// Verifica se o servidor atual é o mesmo que o servidor local, ou seja, não envia a posição para si mesmo
		if name == serverName {
			continue
		}

		url := fmt.Sprintf("http://%s:%s/server/forward", name, ip)
		jsonStr := fmt.Sprintf(`{"x":%d, "y":%d}`, x, y)

		resp, err := http.Post(url, "application/json", strings.NewReader(jsonStr))
		if err != nil {
			log.Printf("❌ Erro ao enviar para %s: %v\n", name, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("📨 Posição (%d, %d) enviada para %s\n", x, y, name)
		} else {
			fmt.Printf("⚠️ Resposta do servidor %s: %s\n", name, resp.Status)
		}
	}
}

