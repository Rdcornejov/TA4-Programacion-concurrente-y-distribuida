# Programación Concurrente y Distribuida

## TA4

### Integrantes:

- Jefferson Espinal Atencia (u201919607)
- Erick Aronés Garcilazo (u201924440)
- Ronaldo Cornejo Valencia (u201816502)

### Docente: 
Carlos Alberto Jara García

### Sección: CC65

### Ciclo: 2023-02

## Planteamiento del problema
El objetivo principal del trabajo es simular el juego de Ludo utilizando programación concurrente, canales y algoritmos distribuidos para la comunicación entre jugadores y el tablero del juego. Ludo es un juego en el que los jugadores compiten para guiar a sus fichas hasta la meta, a través de un laberinto lleno de obstáculos. 
La simulación debe ser capaz de manejar un grupo de jugadores de manera concurrente y usando algoritmo distribuido, donde la comunicación es a través de puertos y sincronización usando canales.
La simulación debe mostrar el progreso del juego en tiempo real, lo que significa que los jugadores deben recibir actualizaciones sobre el estado del juego.

## Explicación del código y el uso de los mecanismos de paralelización y sincronización utilizados.
### start_game.go
start_game(): Se crea un objeto llamado GameData que almacenará la variable de número de jugadores y una lista que simulará el mapa. Luego se serializa el objeto a un archivo JSON y, posteriormente, se convierte en una cadena (string) para enviarlo a través de la conexión.

```go
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
```

main(): Se inicializa el mapa del juego y luego se solicita la cantidad de jugadores y el nodo remoto. Finalmente, estos parámetros se pasan como argumentos a la función start_game().

```go
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
```

### player.go

Se definen dos estructuras: "Ficha" y "Lanzamiento". La estructura "Ficha" representa las fichas de los jugadores, mientras que "Lanzamiento" representa los resultados de lanzar dos dados. Uno de los datos más importantes es el campo "estado" de la estructura "Ficha", que nos permite saber si una ficha está entrando en una casilla con obstáculos (1) o si ya ha estado en una zona de obstáculos (2). Además, se crea una variable para el nodo remoto, un arreglo de objetos tipo "Ficha" y una lista para almacenar el mapa.

```go
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
```

enviar(): Transformación del objeto GameData a formato JSON y luego a tipo string, con el objetivo de enviar el objeto al nodo siguiente.

```go
func enviar(gameData GameData) {
	con, _ := net.Dial("tcp", direccionRemota)
	jsonBytes, _ := json.Marshal(gameData)
	jsonStr := string(jsonBytes)
	defer con.Close()
	fmt.Fprintln(con, jsonStr)
}
```

manejador(): Almacenamiento del string entrante y transformación en un objeto de tipo GameData. En la primera ronda se inicializarán las fichas de los jugadores, en las rondas siguientes se validarán si el jugador actual ha ganado el juego. En caso de que el jugador aún no haya ganado, jugará su turno y enviará el objeto GameData al nodo siguiente. En caso contrario, se imprimirá un mensaje de victoria del jugador ganador.

```go
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
```

lanzarDados(): Genera un lanzamiento de dados aleatorio y devuelve un objeto de tipo "Lanzamiento".

```go
func lanzarDados() Lanzamiento {
	valor := rand.Intn(2)
	tiro := Lanzamiento{
		dadoA:   rand.Intn(6) + 1,
		dadoB:   rand.Intn(6) + 1,
		avanzar: valor == 1,
	}
	return tiro
}
```

pierdeTurno(): Verifica si un jugador pierde su turno debido a un obstáculo en el tablero.

```go
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
```

turnoJugador(): Representa el turno de un jugador. Utiliza canales para coordinar los movimientos de las fichas y realiza cálculos para avanzar las fichas en función de los resultados de los dados. También verifica si alguna ficha ha llegado a la meta.

```go
func turno_jugador(ficha1 chan bool, ficha2 chan bool, ficha3 chan bool, ficha4 chan bool) {
	var tiro Lanzamiento = lanzarDados()
	var ind int
	if !pierdeTurno() {
		go func() {
			if fichas[0].meta == false {
				ficha1 <- true
			}
		}()
		go func() {
			if fichas[1].meta == false {
				ficha2 <- true
			}
		}()
		go func() {
			if fichas[2].meta == false {
				ficha3 <- true
			}
		}()
		go func() {
			if fichas[3].meta == false {
				ficha4 <- true
			}
		}()

		select {
		case <-ficha1:
			fmt.Printf("(JUEGA FICHA 1)\n")
			ind = 0

		case <-ficha2:
			fmt.Printf("(JUEGA FICHA 2)\n")
			ind = 1

		case <-ficha3:
			fmt.Printf("(JUEGA FICHA 3)\n")
			ind = 2

		case <-ficha4:
			fmt.Printf("(JUEGA FICHA 4)\n")
			ind = 3
		}

		go func() {
			for {
				select {
				case <-ficha1:
				case <-ficha2:
				case <-ficha3:
				case <-ficha4:
					// Descartar elementos del canal
				default:
					// El canal está vacío
					return
				}
			}
		}()

		if tiro.avanzar {
			fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA+tiro.dadoB)
			fichas[ind].posicion += tiro.dadoA + tiro.dadoB
			if fichas[ind].posicion > 39 {
				fichas[ind].posicion = 39 - (fichas[ind].posicion - 39)
			}
		} else {
			fmt.Println("RESULTADO LANZAMIENTO: ", tiro.dadoA-tiro.dadoB)
			fichas[ind].posicion += tiro.dadoA - tiro.dadoB
			if fichas[ind].posicion < 0 {
				fichas[ind].posicion = 0
			}
		}
		fmt.Println("POSCION ACTUAL DE LA FICHA: ", fichas[ind].posicion)
		// gano?
		for i := 0; i < 4; i++ {
			if fichas[i].posicion == 39 {
				fichas[i].meta = true
			}
		}

	} else {
		fmt.Println("ESTE JUGADOR PERDIO SU TURNO")
	}
	fmt.Println("-------------------------")
}
```

main(): Se ingresan los valores del color del jugador, el puerto actual y el puerto de destino. Se crea una lista de canales para las fichas del jugador y se ejecuta de manera concurrente la función 'manejador()'.

```go
func main() {

	br := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa el color del jugador: ")
	color, _ := br.ReadString('\n')
	color = strings.TrimSpace(color)

	fmt.Print("Puerto Actual: ")
	strPuertoLocal, _ := br.ReadString('\n')
	strPuertoLocal = strings.TrimSpace(strPuertoLocal)
	direccionLocal := fmt.Sprintf("localhost:%s", strPuertoLocal)

	fmt.Print("Puerto Destino: ")
	strPuertoRemoto, _ := br.ReadString('\n')
	strPuertoRemoto = strings.TrimSpace(strPuertoRemoto)
	direccionRemota = fmt.Sprintf("localhost:%s", strPuertoRemoto)

	chFichas := make([]chan bool, NFICHAS)

	for i := range chFichas {
		chFichas[i] = make(chan bool)
	}

	ln, _ := net.Listen("tcp", direccionLocal)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejador(con, color, chFichas)
	}
}
```

### Uso de Mecanismos de Paralelización y Sincronización:
- Canales (ficha1, ficha2, ficha3, ficha4): Los canales se utilizan para coordinar las acciones de las fichas de los jugadores. En la función turno_jugador, los canales ficha1, ficha2, ficha3 y ficha4 se utilizan para notificar al servidor de la selección de una ficha, el resultado del lanzamiento de los dados y la finalización del turno del jugador.
- Goroutines: El código utiliza goroutines para ejecutar múltiples tareas en paralelo. El archivo player.go utiliza goroutines para manejar las conexiones entrantes, administrar el juego para cada jugador y enviar datos del juego al jugador conectado. Además, también nos permite representar a las fichas al momento de ser elegidas para jugar.

## Diagrama
![image](https://github.com/Rdcornejov/TA4-Programacion-concurrente-y-distribuida/assets/89090023/ef720722-415b-4333-9104-dde4bb279cb6)

## Explicación de las pruebas realizadas y pegar las imágenes de evidencia. 

- Simulación con 3 jugadores:

Para comenzar, ejecutamos el archivo player.go tres veces, enviando el color del jugador, el nodo actual y el nodo de destino. En la compilación del archivo start_game.go, enviamos la cantidad de jugadores y el nodo remoto que dará inicio al juego.

![image](https://github.com/Rdcornejov/TA4-Programacion-concurrente-y-distribuida/assets/66271146/5de4d7bf-0736-40b1-a146-2861afd31c6d)

En la primera ronda, todos los jugadores inicializan sus fichas. El archivo start_game.go envía el mapa del juego y la cantidad de jugadores para el proceso de inicialización. En la prueba, se observa que todos los jugadores recibieron el mismo mensaje, imprimiendo todos el mismo mapa.

![image](https://github.com/Rdcornejov/TA4-Programacion-concurrente-y-distribuida/assets/66271146/76acd185-7791-4366-981b-f0ad36eacfa0)

En las siguientes rondas, observamos que cada jugador comienza a mover sus fichas siguiendo algunas validaciones. Por ejemplo, todas las fichas comienzan en la posición 0, la meta es la posición 39, no existen casillas negativas, y si un jugador, al llegar a la meta, tiene un lanzamiento que, al sumarse a la posición actual, supera 39, entonces el jugador retrocede. También se observan validaciones para los obstáculos, mostrando un mensaje de que el jugador perdió su turno.

![image](https://github.com/Rdcornejov/TA4-Programacion-concurrente-y-distribuida/assets/66271146/c2ecdde0-d79a-4620-aad3-26d7f2f3731b)

