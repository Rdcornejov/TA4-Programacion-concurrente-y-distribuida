package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	MAXOBSTACULOS = 10
)

type GameData struct {
	NumPlayers int
	GameMap    [40]int
}

func start_game(direccionRemota string, numPlayers int, gameMap [40]int) {
	gameData := GameData{
		NumPlayers: numPlayers,
		GameMap:    gameMap,
	}
	jsonBytes, _ := json.Marshal(gameData)
	jsonStr := string(jsonBytes)
	fmt.Println(jsonStr)
	con, _ := net.Dial("tcp", direccionRemota)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)
}

func initialize_game_map(tabla *[40]int, invalid_positions []int) {
	var contador int

	for contador < MAXOBSTACULOS {
		number := rand.Intn(40)
		found := false
		for _, v := range invalid_positions {
			if number == v {
				found = true
				break
			}
		}
		if !found {
			contador++
			(*tabla)[number] = -1
		}
	}
}

func main() {
	var game_map [40]int
	invalid_positions := []int{0, 39}

	initialize_game_map(&game_map, invalid_positions)

	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa el numero de jugadores: ")
	numPlayersStr, _ := br.ReadString('\n')
	numPlayersStr = strings.TrimSpace(numPlayersStr)
	numPlayers, _ := strconv.Atoi(numPlayersStr)

	fmt.Print("Ingrese el puerto del nodo remoto: ")
	puertoRemotaStr, _ := br.ReadString('\n')
	puertoRemotaStr = strings.TrimSpace(puertoRemotaStr)
	direccionRemota := fmt.Sprintf("localhost:%s", puertoRemotaStr)
	start_game(direccionRemota, numPlayers, game_map)
}
