package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	Dict "tdas/diccionario"
	Lista "tdas/lista"
	"time"
	vuelo "tp2/vuelo"
)

const (
	_MIN_CODE_ITERADORES = ""
	_MAX_CODE_ITERADORES = "~"
)

type SistemaVuelos struct {
	porCodigo    Dict.Diccionario[string, vuelo.Vuelo]
	porFechaAsc  Dict.DiccionarioOrdenado[time.Time, Lista.Lista[vuelo.Vuelo]]
	porFechaDesc Dict.DiccionarioOrdenado[time.Time, Lista.Lista[vuelo.Vuelo]]
	porConexion  Dict.Diccionario[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]]
	porPrioridad Dict.DiccionarioOrdenado[int, Lista.Lista[string]]
}

func CrearSistemaDeVuelos() *SistemaVuelos {
	return &SistemaVuelos{
		porCodigo:    Dict.CrearHash[string, vuelo.Vuelo](),
		porFechaAsc:  Dict.CrearABB[time.Time, Lista.Lista[vuelo.Vuelo]](compararTimeAsc),
		porFechaDesc: Dict.CrearABB[time.Time, Lista.Lista[vuelo.Vuelo]](compararTimeDesc),
		porConexion:  Dict.CrearHash[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
		porPrioridad: Dict.CrearABB[int, Lista.Lista[string]](compararIntDesc),
	}
}

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
		vueloActual, err := vuelo.ParsearLineaCSV(registro)
		sv.agregarUnVuelo(vueloActual)
	}
	return nil
}

func (sv *SistemaVuelos) ver_tablero(K int, modo, desde, hasta string) {
	fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
	fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)

	var resultados []string
	if modo == "asc" {
		resultados = sv.obtenerTablero(sv.porFechaAsc, fechaDesde, fechaHasta, K, true)
	} else {
		resultados = sv.obtenerTablero(sv.porFechaDesc, fechaHasta, fechaDesde, K, false)
	}

	for _, linea := range resultados {
		fmt.Println(linea)
	}
	fmt.Println("OK")
}

func (sv *SistemaVuelos) info_vuelo(codigoVuelo string) {
	if !sv.porCodigo.Pertenece(codigoVuelo) {
		fmt.Fprintln(os.Stderr, "Error en comando info_vuelo")
		return
	}
	vuelo := sv.porCodigo.Obtener(codigoVuelo)
	fmt.Println(vuelo.MostrarInfo())
	fmt.Println("OK")
}

func (sv *SistemaVuelos) borrar(desde, hasta string) {
    fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
    fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)

    var fechasAEliminar []time.Time
    sv.porFechaAsc.IterarRango(&fechaDesde, &fechaHasta, func(fecha time.Time, lista Lista.Lista[vuelo.Vuelo]) bool {
        for it := lista.Iterador(); it.HaySiguiente(); it.Siguiente() {
            v := it.VerActual()
            fmt.Println(v.MostrarInfo())
            sv.eliminarIndices(v)
        }
        fechasAEliminar = append(fechasAEliminar, fecha)
        return true
    })

    for _, f := range fechasAEliminar {
        sv.porFechaAsc.Borrar(f)
        sv.porFechaDesc.Borrar(f)
    }

    fmt.Println("OK")
}

func (sv *SistemaVuelos) prioridad_vuelos(k int) {
    minPrio := math.MinInt
    count := 0
	sv.porPrioridad.IterarRango(nil, &minPrio, func(prio int, lst Lista.Lista[string]) bool {
        iter := lst.Iterador()
        for iter.HaySiguiente() && count < k {
            codigo := iter.VerActual()
            fmt.Printf("%d - %s\n", prio, codigo)
            count++
            iter.Siguiente()
        }
        return count < k
    })
	fmt.Println("OK")
}


func (sv *SistemaVuelos) siguiente_vuelo(origen, destino, fechaStr string) {
    connKey := origen + "-" + destino
    fecha, _ := time.Parse("2006-01-02T15:04:05", fechaStr)

    if !sv.porConexion.Pertenece(connKey) {
        fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fechaStr)
        fmt.Println("OK")
        return
    }
    abb := sv.porConexion.Obtener(connKey)
    claveDesde := FechaClave{Fecha: fecha, Codigo: _MAX_CODE_ITERADORES}

    encontrado := false
    abb.IterarRango(&claveDesde, nil, func(clave FechaClave, v vuelo.Vuelo) bool {
        fmt.Println(v.MostrarInfo())
        encontrado = true
        return false
    })
    if !encontrado {
        fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fechaStr)
        fmt.Println("OK")
        return
    }
    fmt.Println("OK")
}


func (sv *SistemaVuelos) obtenerTablero(arbol Dict.DiccionarioOrdenado[time.Time, Lista.Lista[vuelo.Vuelo]], desde, hasta time.Time, K int, esAsc bool) []string {
	resultados := make([]string, 0, K)

	arbol.IterarRango(&desde, &hasta, func(fecha time.Time, lista Lista.Lista[vuelo.Vuelo]) bool {
		iter := lista.Iterador()
		for iter.HaySiguiente() && len(resultados) < K {
			v := iter.VerActual()
			resultados = append(resultados, fmt.Sprintf(
				"%s - %s",
				v.ObtenerFecha().Format("2006-01-02T15:04:05"),
				v.ObtenerCodigo(),
			))
			iter.Siguiente()
		}
		return len(resultados) < K
	})
	return resultados
}
