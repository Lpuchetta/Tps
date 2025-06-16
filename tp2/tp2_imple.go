package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	Dict "tdas/diccionario"
	Heap "tdas/heap"
	"time"
	vuelo "tp2/vuelo"
)

// SistemaVuelos expone un TDA para:
// - guardar vuelos a partir de archivos CSV,
// - buscarlos por código,
// - listarlos ordenados por fecha,
// - y desencolar por prioridad (filtrando duplicados)
type SistemaVuelos struct {
	porCodigo    Dict.Diccionario[string, vuelo.Vuelo]
	porFecha  Dict.DiccionarioOrdenado[FechaClave,vuelo.Vuelo]
	porConexion  Dict.Diccionario[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]]
}

func CrearSistemaDeVuelos() *SistemaVuelos {
	return &SistemaVuelos{
		porCodigo:    Dict.CrearHash[string, vuelo.Vuelo](),
		porFecha:     Dict.CrearABB[FechaClave,vuelo.Vuelo](comparadorFechaClave),
		porConexion:  Dict.CrearHash[string, Dict.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]](),
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
    // Parseamos las fechas de los strings de entrada
    fechaDesde, _ := time.Parse("2006-01-02T15:04:05", desde)
    fechaHasta, _ := time.Parse("2006-01-02T15:04:05", hasta)
    const maxCode = "~"

    // Creamos las claves de rango (desde, hasta)
    claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
    claveHasta := FechaClave{Fecha: fechaHasta, Codigo: maxCode}

    if modo == "desc" {
        // ---- AQUÍ: creamos un min‑heap real invirtiendo cmpFechaAsc ----
        // comparadorInv := func(a, b vueloConFecha) int { return cmpFechaAsc(b, a) }
        heap := Heap.CrearHeap(func(a, b vueloConFecha) int {
            return cmpVueloConFechaAsc(b, a)
        })

        sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool{
            linea := fmt.Sprintf("%s - %s",
                v.ObtenerFecha().Format("2006-01-02T15:04:05"),
                v.ObtenerCodigo(),
            )
            vf := vueloConFecha{
                fecha: v.ObtenerFecha(),
                codigo: v.ObtenerCodigo(),
                info:  linea,
            }

            if heap.Cantidad() < K {
                heap.Encolar(vf)
            } else {
                // si vf es más reciente que la raíz (el más antiguo de los K),
                // reemplazamos:
                if cmpVueloConFechaAsc(heap.VerMax(), vf) < 0 {
                    heap.Desencolar()
                    heap.Encolar(vf)
                }
            }
            return true
        })

        // Extraemos y volcamos en orden inverso para descendente
        resultados := make([]string, 0, heap.Cantidad())
        for !heap.EstaVacia() {
            resultados = append(resultados, heap.Desencolar().info)
        }
        for i := len(resultados) - 1; i >= 0; i-- {
            fmt.Println(resultados[i])
        }

    } else {
        // Modo ascendente: iteramos y cortamos en K
        resultados := make([]string, 0, K)
        sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(clave FechaClave, v vuelo.Vuelo) bool {
            linea := fmt.Sprintf("%s - %s",
                v.ObtenerFecha().Format("2006-01-02T15:04:05"),
                v.ObtenerCodigo(),
            )
            resultados = append(resultados, linea)
            return len(resultados) < K
        })
        for _, linea := range resultados {
            fmt.Println(linea)
        }
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
   // Parsear fechas con manejo de errores
   fechaDesde, err := time.Parse("2006-01-02T15:04:05", desde)
   if err != nil {
       fmt.Fprintln(os.Stderr, "Error en comando borrar: formato de 'desde' inválido")
       return
   }
   fechaHasta, err := time.Parse("2006-01-02T15:04:05", hasta)
   if err != nil {
       fmt.Fprintln(os.Stderr, "Error en comando borrar: formato de 'hasta' inválido")
       return
   }
   if fechaHasta.Before(fechaDesde) {
       fmt.Fprintln(os.Stderr, "Error en comando borrar: 'hasta' debe ser igual o posterior a 'desde'")
       return
   }

   // Claves de rango (incluyendo todo el código)
   claveDesde := FechaClave{Fecha: fechaDesde, Codigo: ""}
   claveHasta := FechaClave{Fecha: fechaHasta, Codigo: "~"}

   // Iterar sobre porFecha (FechaClave → vuelo)
   sv.porFecha.IterarRango(&claveDesde, &claveHasta, func(cl FechaClave, v vuelo.Vuelo) bool {
       // a) imprimir y borrar de porCodigo
       fmt.Println(v.MostrarInfo())
       sv.porCodigo.Borrar(v.ObtenerCodigo())

       // b) borrar en porConexion
       connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
       if sv.porConexion.Pertenece(connKey) {
           abb := sv.porConexion.Obtener(connKey)
           abb.Borrar(cl)
           if abb.Cantidad() == 0 {
               sv.porConexion.Borrar(connKey)
           }
       }

       // c) borrar del ABB porFecha
       sv.porFecha.Borrar(cl)

       return true // continuar borrando
   })

   fmt.Println("OK")
}
func (sv *SistemaVuelos) prioridad_vuelos(k int) {
	// 1) Recolectar vuelos en un slice
	vuelos := make([]vuelo.Vuelo, 0, sv.porCodigo.Cantidad())
	iter := sv.porCodigo.Iterador()
	for iter.HaySiguiente() {
		_, v := iter.VerActual()
		vuelos = append(vuelos, v)
		iter.Siguiente()
	}

	// 2) Construir heap en O(n) usando CrearHeapArr
	comparador := func(v1, v2 vuelo.Vuelo) int {
		p1, p2 := v1.ObtenerPrioridad(), v2.ObtenerPrioridad()
		if p1 != p2 {
			if p1 > p2 {
				return 1
			}
			return -1
		}
		if v1.ObtenerCodigo() < v2.ObtenerCodigo() {
			return 1
		}
		if v1.ObtenerCodigo() > v2.ObtenerCodigo() {
			return -1
		}
		return 0
	}
	h := Heap.CrearHeapArr(vuelos, comparador)

	// 3) Extraer K vuelos en O(K log n)
	for i := 0; i < k && !h.EstaVacia(); i++ {
		v := h.Desencolar()
		fmt.Printf("%d - %s\n", v.ObtenerPrioridad(), v.ObtenerCodigo())
	}
	fmt.Println("OK")
}

func (sv *SistemaVuelos) siguiente_vuelo(origen, destino, fechaStr string) {
	connKey := origen + "-" + destino

	fecha, err := time.Parse("2006-01-02T15:04:05", fechaStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error en comando siguiente_vuelo")
		return
	}

	if !sv.porConexion.Pertenece(connKey) {
		fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fechaStr)
		fmt.Println("OK")
		return
	}

	abb := sv.porConexion.Obtener(connKey)

	claveDesde := FechaClave{Fecha: fecha, Codigo: ""}

	encontrado := false
	abb.IterarRango(&claveDesde, nil, func(clave FechaClave, v vuelo.Vuelo) bool {
		// este será el primer vuelo >= claveDesde
		fmt.Println(v.MostrarInfo())
		encontrado = true
		return false // corto la iteración, ya encontré el vuelo siguiente
	})

	if !encontrado {
		fmt.Printf("No hay vuelo registrado desde %s hacia %s desde %s\n", origen, destino, fechaStr)
		fmt.Println("OK")
		return
	}

	fmt.Println("OK")
}
