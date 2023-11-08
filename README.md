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
