package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	Abb "tdas/diccionario"
	Hash "tdas/diccionario"
	Heap "tdas/heap"
	vuelo "tp2/vuelo"
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
	porConexion  Hash.Diccionario[string, Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]]
}

// CrearSistemaDeVuelos instancia el TDA vacío.
func CrearSistemaDeVuelos() Aeropuerto {
	return &SistemaVuelos{
		porCodigo:    Hash.CrearHash[string, vuelo.Vuelo](),
		porFecha:     Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave),
		porPrioridad: Heap.CrearHeap[vuelo.Vuelo](comparadorPrioridad),
		porConexion:  Hash.CrearHash[string, Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
	}
}

// Agregar_Archivo lee el CSV completo y va insertando cada línea en el TDA.
// Si un vuelo ya existía (mismo código), se reemplaza la entrada anterior.
func (sv *SistemaVuelos) agregar_archivo(nombreArchivo string) error {
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

		vueloActual, err := vuelo.ParsearLineaCSV(registro) //En archivo vuelo/parsing_vuelo.go
		if err != nil {
			continue
		}
		sv.agregarUnVuelo(vueloActual) //En archivo utilidades.go
	}
	return nil
}

func (sv *SistemaVuelos) ver_tablero(K int, modo, desde, hasta string) {
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

}

func (sv *SistemaVuelos) info_vuelo(codigoVuelo string) {
	if !sv.porCodigo.Pertenece(codigoVuelo) {
		fmt.Printf("No se encontró vuelo con código %s\n", codigoVuelo)
		return
	}
	vuelo := sv.porCodigo.Obtener(codigoVuelo)
	fmt.Println(vuelo.MostrarInfo())
}

func (sv *SistemaVuelos) borrar(desde, hasta string) {
	fechaDesde, err1 := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, err2 := time.Parse("2006-01-02T15:04:05", hasta)
	if err1 != nil || err2 != nil || fechaHasta.Before(fechaDesde) {
		fmt.Println("Error: rango de fechas inválido")
		return
	}

	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: "999999999999"}

	clavesABorrar := make([]FechaClave, 0)

	sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool {
		fmt.Println(v.MostrarInfo())
		clavesABorrar = append(clavesABorrar, clave)
		return true
	})

	for _, clave := range clavesABorrar {
		sv.porFecha.Borrar(clave)
		sv.porCodigo.Borrar(clave.Codigo)
	}

}

func (sv *SistemaVuelos) prioridad_vuelos(k int) {
	vuelos := make([]vuelo.Vuelo, 0)
	for it := sv.porCodigo.Iterador(); it.HaySiguiente(); it.Siguiente() {
		_, v := it.VerActual()
		vuelos = append(vuelos, v)
	}

	h := Heap.CrearHeapArr(vuelos, comparadorPrioridad)

	if k > len(vuelos) {
		k = len(vuelos)
	}

	temp := make([]vuelo.Vuelo, 0, k)
	for i := 0; i < k; i++ {
		temp = append(temp, h.Desencolar())
	}

	for i := k - 1; i >= 0; i-- {
		v := temp[i]
		fmt.Printf("%d - %s\n", v.ObtenerPrioridad(), v.ObtenerCodigo())
	}

}

func (sv *SistemaVuelos) siguiente_vuelo(origen, destino, fecha string) {
	fechaBuscada, err := time.Parse("2006-01-02T15:04:05", fecha)
	if err != nil {
		fmt.Println("Error: Fecha invalida")
		return
	}

	clave := origen + "-" + destino
	if !sv.porConexion.Pertenece(clave) {
		fmt.Println("No hay vuelos para esa conexion")
		return
	}

	abbVuelos := sv.porConexion.Obtener(clave)
	encontrado := false

	claveDesde := &FechaClave{Fecha: fechaBuscada, Codigo: ""}
	abbVuelos.IterarRango(claveDesde, nil, func(fc FechaClave, v vuelo.Vuelo) bool {
		fmt.Println(v.MostrarInfo())
		encontrado = true
		return false // para cortar al primer vuelo encontrado >= fechaBuscada
	})

	if !encontrado {
		fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fecha)
	}
}
