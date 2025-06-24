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
	porPrioridad Dict.DiccionarioOrdenado[PrioridadClave, vuelo.Vuelo]
}

func CrearSistemaDeVuelos() *SistemaVuelos {
	return &SistemaVuelos{
		porCodigo:    Dict.CrearHash[string, vuelo.Vuelo](),
		porFechaAsc:  Dict.CrearABB[time.Time, Lista.Lista[vuelo.Vuelo]](compararTimeAsc),
		porFechaDesc: Dict.CrearABB[time.Time, Lista.Lista[vuelo.Vuelo]](compararTimeDesc),
		porConexion:  Dict.CrearHash[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
		porPrioridad: Dict.CrearABB[PrioridadClave, vuelo.Vuelo](cmpPrioridadClave),
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
		// árbol ascendente, lista de códigos ascendente
		resultados = sv.obtenerTablero(sv.porFechaAsc, fechaDesde, fechaHasta, K, true)
	} else {
		// árbol descendente, lista de códigos descendente
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
    // Parseo de fechas (puedes añadir chequeo de errores si lo deseas)
    fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
    fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)

    // 1) Recolectar fechas donde hay vuelos a eliminar (con sus vuelos).
    var fechasAEliminar []time.Time
    sv.porFechaAsc.IterarRango(&fechaDesde, &fechaHasta, func(fecha time.Time, lista Lista.Lista[vuelo.Vuelo]) bool {
        iter := lista.Iterador()
        for iter.HaySiguiente() {
            v := iter.VerActual()
            // 1.1) Imprimir info del vuelo
            fmt.Println(v.MostrarInfo())
            // 1.2) Borrar de índice porCódigo
            sv.porCodigo.Borrar(v.ObtenerCodigo())
            // 1.3) Borrar de índice porPrioridad
            prioridadClave := PrioridadClave{
                Prioridad: v.ObtenerPrioridad(),
                Codigo:    v.ObtenerCodigo(),
            }
            sv.porPrioridad.Borrar(prioridadClave)
            // 1.4) Borrar de índice porConexion
            connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
            if sv.porConexion.Pertenece(connKey) {
                abbConn := sv.porConexion.Obtener(connKey)
                fechaClave := FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()}
                abbConn.Borrar(fechaClave)
                // Si ya no hay vuelos para esa ruta, eliminar la entrada completa
                if abbConn.Cantidad() == 0 {
                    sv.porConexion.Borrar(connKey)
                }
            }
            iter.Siguiente()
        }
        // Marca esta fecha para eliminar todo el nodo del ABB
        fechasAEliminar = append(fechasAEliminar, fecha)
        return true
    })

    // 2) Borrar cada nodo de fecha de ambos ABB:
    for _, f := range fechasAEliminar {
        sv.porFechaAsc.Borrar(f)
        sv.porFechaDesc.Borrar(f)
    }

    // 3) Confirmación de fin de comando
    fmt.Println("OK")
}



func (sv *SistemaVuelos) prioridad_vuelos(k int) {
	maxKey := PrioridadClave{Prioridad: math.MaxInt, Codigo: _MIN_CODE_ITERADORES}
	minKey := PrioridadClave{Prioridad: math.MinInt, Codigo: _MAX_CODE_ITERADORES}
	conteo := 0

	sv.porPrioridad.IterarRango(&maxKey, &minKey, func(cl PrioridadClave, v vuelo.Vuelo) bool {
		fmt.Printf("%d - %s\n", cl.Prioridad, cl.Codigo)
		conteo++
		return conteo < k
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
    // Arrancamos justo *después* de 'fecha',
    // usando el código máximo para saltar también los de la misma fecha
    claveDesde := FechaClave{Fecha: fecha, Codigo: _MAX_CODE_ITERADORES}

    encontrado := false
    abb.IterarRango(&claveDesde, nil, func(clave FechaClave, v vuelo.Vuelo) bool {
        fmt.Println(v.MostrarInfo())
        encontrado = true
        return false // dejamos de iterar tras encontrar el primero
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
