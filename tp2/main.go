package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	sv := CrearSistemaDeVuelos()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Fields(linea)
		if len(partes) == 0 {
			continue
		}

		comando := partes[0]
		switch comando {
		case "agregar_archivo":
			if len(partes) < 2 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}
			archivo := partes[1]
			err := sv.agregar_archivo(archivo)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
			} else {
				fmt.Println("OK")
			}

		case "ver_tablero":
			if len(partes) < 5 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			K, err := strconv.Atoi(partes[1])
			if err != nil || K <= 0 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			modo := strings.ToLower(partes[2])
			if modo != "asc" && modo != "desc" {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			desdeStr := partes[3]
			hastaStr := partes[4]

			const layout = "2006-01-02T15:04:05"
			desde, errDesde := time.Parse(layout, desdeStr)
			hasta, errHasta := time.Parse(layout, hastaStr)

			if errDesde != nil || errHasta != nil {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			if hasta.Before(desde) {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			sv.ver_tablero(K, modo, desdeStr, hastaStr)

		case "info_vuelo":
			if len(partes) < 2 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}
			codigo := partes[1]
			sv.info_vuelo(codigo)

		case "prioridad_vuelos":
			if len(partes) < 2 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}
			K, err := strconv.Atoi(partes[1])
			if err != nil || K <= 0 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}
			sv.prioridad_vuelos(K)

		case "siguiente_vuelo":
			if len(partes) < 4 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			origen := partes[1]
			destino := partes[2]
			fechaStr := partes[3]

			const layout = "2006-01-02T15:04:05"
			_, err := time.Parse(layout, fechaStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			sv.siguiente_vuelo(origen, destino, fechaStr)

		case "borrar":
			if len(partes) < 3 {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			desdeStr := partes[1]
			hastaStr := partes[2]

			const layout = "2006-01-02T15:04:05"
			desde, errDesde := time.Parse(layout, desdeStr)
			hasta, errHasta := time.Parse(layout, hastaStr)

			if errDesde != nil || errHasta != nil {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			if hasta.Before(desde) {
				fmt.Fprintln(os.Stderr, "Error en comando", comando)
				continue
			}

			sv.borrar(desdeStr, hastaStr)

		default:
			fmt.Fprintln(os.Stderr, "Error en comando", comando)
		}
	}
}
