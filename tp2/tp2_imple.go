package tp2

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	Abb "tdas/diccionario"
	Hash "tdas/diccionario"
	Heap "tdas/heap"
	vuelo "tdas/tp2/vuelo"
)

// SistemaVuelos expone un TDA para:
// - guardar vuelos a partir de archivos CSV,
// - buscarlos por código,
// - listarlos ordenados por fecha,
// - y desencolar por prioridad (filtrando duplicados)
type SistemaVuelos struct {
	porCodigo    Hash.Diccionario[string, vuelo.Vuelo]
	porFecha     Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]
	porPrioridad Heap.ColaPrioridad[vuelo.Vuelo]
}

// CrearSistemaDeVuelos instancia el TDA vacío.
func CrearSistemaDeVuelos() Aeropuerto {
	return &SistemaVuelos{
		porCodigo:    Hash.CrearHash[string, vuelo.Vuelo](),
		porFecha:     Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave),
		porPrioridad: Heap.CrearHeap[vuelo.Vuelo](comparadorPrioridad),
	}
}

// Agregar_Archivo lee el CSV completo y va insertando cada línea en el TDA.
// Si un vuelo ya existía (mismo código), se reemplaza la entrada anterior.
func (sv *SistemaVuelos) Agregar_Archivo(nombreArchivo string) error {
	archivo, err := os.Open(nombreArchivo)
	if err != nil {
		return fmt.Errorf("no se pudo abrir '%s': %w", nombreArchivo, err)
	}

	defer archivo.Close()

	reader := csv.NewReader(archivo)
	reader.Comma = ','

	for {
		registro, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error leyendo CSV: %w", err)
		}

		vueloActual, err := parsearLineaCSV(registro)
		if err != nil {
			continue
		}
		sv.agregarUnVuelo(vueloActual)
	}
	return nil
}

func (sv *SistemaVuelos) Ver_Tablero(K int, modo, desde, hasta string) {
	if K <= 0 {
		fmt.Println("Error: K inválido")
		return
	}

	if modo != "asc" && modo != "desc" {
		fmt.Println("Error: modo inválido")
		return
	}

	fechaDesde, err1 := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, err2 := time.Parse("2006-01-02T15:04:05", hasta)
	if err1 != nil || err2 != nil || fechaHasta.Before(fechaDesde) {
		fmt.Println("Error: rango de fechas inválido")
		return
	}

	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: "999999999999"}

	resultados := make([]string, 0, K)

	sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool {
		linea := fmt.Sprintf("%s - %s", clave.Fecha.Format("2006-01-02T15:04:05"), clave.Codigo)
		resultados = append(resultados, linea)
		return len(resultados) < K
	})

	if modo == "desc" {
		for i := len(resultados) - 1; i >= 0; i-- {
			fmt.Println(resultados[i])
		}
	} else {
		for _, linea := range resultados {
			fmt.Println(linea)
		}
	}

	fmt.Println("OK")
}

func (sv *SistemaVuelos) Info_Vuelo(codigoVuelo string) {
	if !sv.porCodigo.Pertenece(codigoVuelo) {
		fmt.Printf("No se encontró vuelo con código %s\n", codigoVuelo)
		return
	}
	vuelo := sv.porCodigo.Obtener(codigoVuelo)
	vuelo.MostrarInfo()
	fmt.Println("OK")
}

func parsearLineaCSV(linea []string) (vuelo.Vuelo, error) {
	if len(linea) < 10 {
		return nil, fmt.Errorf("error leyendo el archivo CSV")
	}

	codigoVuelo := linea[0]
	aerolinea := linea[1]
	origen := linea[2]
	destino := linea[3]
	matricula := linea[4]
	prioridad, _ := strconv.Atoi(linea[5])
	retraso, _ := strconv.Atoi(linea[7])
	fecha, _ := time.Parse("2006-01-02T15:04:05", linea[6])
	tiempoVuelo, _ := strconv.Atoi(linea[8])
	cancelado := linea[9] == "1"
	v := vuelo.CrearVuelo(codigoVuelo, aerolinea, origen, destino, matricula, prioridad, fecha, retraso, tiempoVuelo, cancelado)
	return v, nil
}
