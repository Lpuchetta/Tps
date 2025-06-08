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
	porCodigo   Hash.Diccionario[string, vuelo.Vuelo]
	porFecha    Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]
	porConexion Hash.Diccionario[string, Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]]
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

		vueloActual, err := vuelo.ParsearLineaCSV(registro)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parseando línea: %v\n", err)
			continue
		}
		sv.agregarUnVuelo(vueloActual)
	}
	return nil
}

func (sv *SistemaVuelos) ver_tablero(K int, modo, desde, hasta string) {
	fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)

	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: ""}

	resultados := make([]string, 0, K)

	sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool {
		linea := fmt.Sprintf("%s - %s", v.ObtenerFecha().Format("2006-01-02T15:04:05"), v.ObtenerCodigo())
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
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: ""}

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
			return p1 - p2
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
		_, v := iter.VerActual()
		h.Encolar(v)
		iter.Siguiente()
	}

	for i := 0; i < k && !h.EstaVacia(); i++ {
		v := h.Desencolar()
		fmt.Printf("%d - %s\n", v.ObtenerPrioridad(), v.ObtenerCodigo())
	}
	fmt.Println("OK")
}

func (sv *SistemaVuelos) siguiente_vuelo(origen, destino, fecha string) {
	fechaBuscada, _ := time.Parse("2006-01-02T15:04:05", fecha)
	clave := origen + "-" + destino

	if !sv.porConexion.Pertenece(clave) {
		fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fecha)
	} else {
		abbVuelos := sv.porConexion.Obtener(clave)
		claveDesde := &FechaClave{Fecha: fechaBuscada, Codigo: ""}
		var siguiente vuelo.Vuelo

		abbVuelos.IterarRango(claveDesde, nil, func(fc FechaClave, v vuelo.Vuelo) bool {
			if siguiente == nil {
				siguiente = v
			}
			return true
		})

		if siguiente == nil {
			fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fecha)
		} else {
			fmt.Println(siguiente.MostrarInfo())
		}
	}
	fmt.Println("OK")
}
