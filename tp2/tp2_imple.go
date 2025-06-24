package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"math"
	Dict "tdas/diccionario"
	"time"
	vuelo "tp2/vuelo"
)

const (
	_MIN_CODE_ITERADORES = ""
	_MAX_CODE_ITERADORES = "~"
)

type SistemaVuelos struct {
	porCodigo   Dict.Diccionario[string, vuelo.Vuelo]
	porFechaAsc    Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]
	porFechaDesc   Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]
	porConexion Dict.Diccionario[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]]
	porPrioridad Dict.DiccionarioOrdenado[PrioridadClave, vuelo.Vuelo]
}

func CrearSistemaDeVuelos() *SistemaVuelos {
	return &SistemaVuelos{
		porCodigo:   Dict.CrearHash[string, vuelo.Vuelo](),
		porFechaAsc:    Dict.CrearABB[FechaClave, vuelo.Vuelo](cmpFechaClaveAsc),
		porFechaDesc:   Dict.CrearABB[FechaClave, vuelo.Vuelo](cmpFechaClaveDesc),
		porConexion: Dict.CrearHash[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
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
	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: _MIN_CODE_ITERADORES}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: _MAX_CODE_ITERADORES}

	var resultados []string
	if modo == "asc"{
		resultados = sv.obtenerTablero(sv.porFechaAsc, claveDesde, claveHasta, K)
	}else{
		resultados = sv.obtenerTablero(sv.porFechaDesc, claveHasta, claveDesde, K)
	}
	for _, linea := range resultados{
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

	claveDesde := FechaClave{Fecha: fechaDesde, Codigo: _MIN_CODE_ITERADORES}
	claveHasta := FechaClave{Fecha: fechaHasta, Codigo: _MAX_CODE_ITERADORES}

	procesar := func(cl FechaClave, v vuelo.Vuelo){
		fmt.Println(v.MostrarInfo())
		sv.porCodigo.Borrar(v.ObtenerCodigo())
		prioridadClave := PrioridadClave{
			Prioridad: v.ObtenerPrioridad(),
			Codigo: v.ObtenerCodigo(),
		}
		sv.porPrioridad.Borrar(prioridadClave)

		connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
		if sv.porConexion.Pertenece(connKey){
			abb := sv.porConexion.Obtener(connKey)
			abb.Borrar(cl)
			if abb.Cantidad() == 0{
				sv.porConexion.Borrar(connKey)
			}
		}
	}
	sv.borrarRangoEnArbol(sv.porFechaAsc, claveDesde, claveHasta, procesar)
	sv.borrarRangoEnArbol(sv.porFechaDesc, claveHasta, claveDesde, func(cl FechaClave, v vuelo.Vuelo){
		
	})

	fmt.Println("OK")
}
func (sv *SistemaVuelos) prioridad_vuelos(k int) {
	maxKey := PrioridadClave{Prioridad: math.MaxInt, Codigo: _MIN_CODE_ITERADORES}
	minKey := PrioridadClave{Prioridad: math.MinInt, Codigo: _MAX_CODE_ITERADORES}
	conteo := 0

	sv.porPrioridad.IterarRango(&maxKey, &minKey, func(cl PrioridadClave, v vuelo.Vuelo) bool{
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

	claveDesde := FechaClave{Fecha: fecha, Codigo: _MIN_CODE_ITERADORES}

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


func (sv *SistemaVuelos) obtenerTablero(arbol Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo], desde, hasta FechaClave, K int)[]string{
	resultados := make([]string, 0, K)
	arbol.IterarRango(&desde, &hasta, func(cl FechaClave, v vuelo.Vuelo) bool{
		linea := fmt.Sprintf("%s - %s",v.ObtenerFecha().Format("2006-01-02T15:04:05"), v.ObtenerCodigo())
		resultados = append(resultados, linea)
		return len(resultados) < K 
	})
	return resultados
}

func(sv *SistemaVuelos) borrarRangoEnArbol(arbol Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo], desde, hasta FechaClave, procesar func(cl FechaClave, v vuelo.Vuelo)){
	var items []struct{
		clave FechaClave
		vuelo vuelo.Vuelo
	}	
	arbol.IterarRango(&desde, &hasta, func(cl FechaClave, v vuelo.Vuelo) bool{
		items = append(items, struct{
			clave FechaClave	
			vuelo vuelo.Vuelo
		}{cl, v})
		return true
	})
	for _, it := range items{
		procesar(it.clave, it.vuelo)
		arbol.Borrar(it.clave)
	}
}
