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
		porCodigo:   Hash.CrearHash[string, vuelo.Vuelo](),
		porFecha:    Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave),
		porConexion: Hash.CrearHash[string, Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
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
	fmt.Println("OK")
	return nil
}

func (sv *SistemaVuelos) ver_tablero(K int, modo, desde, hasta string) {
	if K <= 0 {
		fmt.Println("Error: K inválido")
		fmt.Println("OK")
		return
	}

	if modo != "asc" && modo != "desc" {
		fmt.Println("Error: modo inválido")
		fmt.Println("OK")
		return
	}

	fechaDesde, err1 := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, err2 := time.Parse("2006-01-02T15:04:05", hasta)
	if err1 != nil || err2 != nil || fechaHasta.Before(fechaDesde) {
		fmt.Println("Error: rango de fechas inválido")
		fmt.Println("OK")
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

func (sv *SistemaVuelos) info_vuelo(codigoVuelo string) {
	if !sv.porCodigo.Pertenece(codigoVuelo) {
		fmt.Printf("No se encontró vuelo con código %s\n", codigoVuelo)
	} else {
		vuelo := sv.porCodigo.Obtener(codigoVuelo)
		fmt.Println(vuelo.MostrarInfo())
	}
	fmt.Println("OK")
}

func (sv *SistemaVuelos) borrar(desde, hasta string) {
	fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)

	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: "999999999999"}

	clavesABorrar := make([]FechaClave, 0)
	vuelosABorrar := make([]vuelo.Vuelo, 0)

	sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool {
		fmt.Println(v.MostrarInfo())
		clavesABorrar = append(clavesABorrar, clave)
		vuelosABorrar = append(vuelosABorrar, v)
		return true
	})

	for i, clave := range clavesABorrar {
		v := vuelosABorrar[i]

		sv.porFecha.Borrar(clave)
		sv.porCodigo.Borrar(clave.Codigo)

		claveConexion := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
		if sv.porConexion.Pertenece(claveConexion) {
			abb := sv.porConexion.Obtener(claveConexion)
			abb.Borrar(clave)
			if abb.Cantidad() == 0 {
				sv.porConexion.Borrar(claveConexion)
			}
		}
	}
	fmt.Println("OK")
}

func (sv *SistemaVuelos) prioridad_vuelos(k int) {
	comparador := func(v1, v2 vuelo.Vuelo) int {
		p1, p2 := v1.ObtenerPrioridad(), v2.ObtenerPrioridad()
		if p1 != p2 {
			return p2 - p1
		}

		if v1.ObtenerCodigo() < v2.ObtenerCodigo() {
			return -1
		} else if v1.ObtenerCodigo() > v2.ObtenerCodigo() {
			return 1
		}
		return 0
	}

	h := Heap.CrearHeap(comparador)
	iter := sv.porCodigo.Iterador()
	for iter.HaySiguiente() {
		_, vuelo := iter.VerActual()
		h.Encolar(vuelo)
		iter.Siguiente()
	}

	for i := 0; i < k && !h.EstaVacia(); i++ {
		v := h.Desencolar()
		fmt.Printf("%d - %s\n", v.ObtenerPrioridad(), v.ObtenerCodigo())
	}
	fmt.Println("OK")
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
	fmt.Println("OK")
}
