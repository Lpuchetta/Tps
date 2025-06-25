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
		resultados = sv.obtenerTablero(sv.porFechaAsc, fechaDesde, fechaHasta, K)
	} else {
		resultados = sv.obtenerTablero(sv.porFechaDesc, fechaHasta, fechaDesde, K)
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
            vuelo := it.VerActual()
            fmt.Println(vuelo.MostrarInfo())
            sv.eliminarIndices(vuelo)
        }
        fechasAEliminar = append(fechasAEliminar, fecha)
        return true
    })

    for _, fechas := range fechasAEliminar {
        sv.porFechaAsc.Borrar(fechas)
        sv.porFechaDesc.Borrar(fechas)
    }

    fmt.Println("OK")
}

func (sv *SistemaVuelos) prioridad_vuelos(k int) {
    prioridadMinima := math.MinInt
    conteo := 0
	sv.porPrioridad.IterarRango(nil, &prioridadMinima, func(prioridad int, lista Lista.Lista[string]) bool {
        iterador := lista.Iterador()
        for iterador.HaySiguiente() && conteo < k {
            codigo := iterador.VerActual()
            fmt.Printf("%d - %s\n", prioridad, codigo)
            conteo++
            iterador.Siguiente()
        }
        return conteo < k
    })
	fmt.Println("OK")
}


func (sv *SistemaVuelos) siguiente_vuelo(origen, destino, fechaStr string) {
    claveConexion := origen + "-" + destino
    fecha, _ := time.Parse("2006-01-02T15:04:05", fechaStr)

    if !sv.porConexion.Pertenece(claveConexion) {
        fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fechaStr)
        fmt.Println("OK")
        return
    }
    abb := sv.porConexion.Obtener(claveConexion)
    claveDesde := FechaClave{Fecha: fecha, Codigo: _MAX_CODE_ITERADORES}

    encontrado := false
    abb.IterarRango(&claveDesde, nil, func(clave FechaClave, vuelo vuelo.Vuelo) bool {
        fmt.Println(vuelo.MostrarInfo())
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


func (sv *SistemaVuelos) obtenerTablero(arbol Dict.DiccionarioOrdenado[time.Time, Lista.Lista[vuelo.Vuelo]], desde, hasta time.Time, K int) []string {
	resultados := make([]string, 0, K)

	arbol.IterarRango(&desde, &hasta, func(fecha time.Time, lista Lista.Lista[vuelo.Vuelo]) bool {
		iterador := lista.Iterador()
		for iterador.HaySiguiente() && len(resultados) < K {
			v := iterador.VerActual()
			resultados = append(resultados, fmt.Sprintf(
				"%s - %s",
				v.ObtenerFecha().Format("2006-01-02T15:04:05"),
				v.ObtenerCodigo(),
			))
			iterador.Siguiente()
		}
		return len(resultados) < K
	})
	return resultados
}
