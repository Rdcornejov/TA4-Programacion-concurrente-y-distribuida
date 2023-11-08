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
