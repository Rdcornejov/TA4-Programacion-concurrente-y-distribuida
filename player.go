package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
)

const (
	NFICHAS = 4
)

type GameData struct {
	NumPlayers int
	GameMap    [40]int
}

type Ficha struct {
	id       int
	color    string
	posicion int
	estado   int
	meta     bool
}

type Lanzamiento struct {
	dadoA   int
	dadoB   int
	avanzar bool
}

var direccionRemota string
var fichas []Ficha
var mapa [40]int

func enviar(gameData GameData) {
	con, _ := net.Dial("tcp", direccionRemota)
	jsonBytes, _ := json.Marshal(gameData)
	jsonStr := string(jsonBytes)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)
}

func lanzarDados() Lanzamiento {
	valor := rand.Intn(2)
	tiro := Lanzamiento{
		dadoA:   rand.Intn(6) + 1,
		dadoB:   rand.Intn(6) + 1,
		avanzar: valor == 1,
	}
	return tiro
}

func initialize_player(color string) {
	for j := 0; j < NFICHAS; j++ {
		ficha := Ficha{
			id:       j + 1,
			color:    color,
			posicion: 0,
			meta:     false,
		}
		fichas = append(fichas, ficha)
	}
}

func manejador(con net.Conn, color string, chFichas []chan bool) {
	var gameData GameData
	// time.Sleep(1 * time.Second)
	defer con.Close()
	fmt.Printf("Turno del Jugador: %s\n", color)
	defer con.Close()
	br := bufio.NewReader(con)
	msg, _ := br.ReadString('\n')
	msg = strings.TrimSpace(msg)
	json.Unmarshal([]byte(msg), &gameData)
	// fmt.Printf("Llego el valor: %d\n", gameData.NumPlayers)
	if gameData.NumPlayers > 0 {
		fmt.Println("Inicializando Fichas")
		initialize_player(color)
		gameData.NumPlayers = gameData.NumPlayers - 1
		mapa = gameData.GameMap
		fmt.Println(mapa)
		fmt.Println("------------------------")
		enviar(gameData)
	} else {
		fichasCompletadas := 0
		for _, f := range fichas {
			if f.meta == true {
				fichasCompletadas++
			}
		}
		if fichasCompletadas < 4 {
			turno_jugador(chFichas[0], chFichas[1], chFichas[2], chFichas[3])
			enviar(gameData)
		} else {
			fmt.Printf("El jugador %s ha ganado el juego\n", color)
			fmt.Println(fichas)
		}
	}
}

func pierdeTurno() bool {
	for i := 0; i < 4; i++ {
		for ind, valor := range mapa {
			if valor == -1 && ind == fichas[i].posicion {
				fichas[i].estado += 1
				if fichas[i].estado > 2 {
					fichas[i].estado = 2
				}
			}
			if fichas[i].estado == 2 && valor == 0 && ind == fichas[i].posicion {
				fichas[i].estado = 0
			}
		}
	}
	for i := 0; i < 4; i++ {
		if fichas[i].estado == 1 {
			return true
		}
	}
	return false
}
