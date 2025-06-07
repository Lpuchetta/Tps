package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
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

        comando := strings.ToLower(partes[0])

        switch comando {
        case "agregar_archivo":
            if len(partes) < 2 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            archivo := partes[1]
            if err := sv.agregar_archivo(archivo); err != nil {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            fmt.Println("OK")

        case "info_vuelo":
            if len(partes) < 2 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            sv.info_vuelo(partes[1])
            fmt.Println("OK")

        case "prioridad_vuelos":
            if len(partes) < 2 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            K, err := strconv.Atoi(partes[1])
            if err != nil || K <= 0 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            sv.prioridad_vuelos(K)
            fmt.Println("OK")

        case "ver_tablero":
            if len(partes) < 5 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            K, err := strconv.Atoi(partes[1])
            if err != nil || K <= 0 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            modo := strings.ToLower(partes[2])
            if modo != "asc" && modo != "desc" {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            desdeStr, hastaStr := partes[3], partes[4]
            const layout = "2006-01-02T15:04:05"
            desde, errDesde := time.Parse(layout, desdeStr)
            hasta, errHasta := time.Parse(layout, hastaStr)
            if errDesde != nil || errHasta != nil || hasta.Before(desde) {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            sv.ver_tablero(K, modo, desdeStr, hastaStr)
            fmt.Println("OK")

        case "borrar":
            if len(partes) < 3 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            desdeStr, hastaStr := partes[1], partes[2]
            const layout = "2006-01-02T15:04:05"
            desde, errDesde := time.Parse(layout, desdeStr)
            hasta, errHasta := time.Parse(layout, hastaStr)
            if errDesde != nil || errHasta != nil || hasta.Before(desde) {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            sv.borrar(desdeStr, hastaStr)
            fmt.Println("OK")

        case "siguiente_vuelo":
            if len(partes) < 4 {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            origen, destino, fechaStr := partes[1], partes[2], partes[3]
            const layout = "2006-01-02T15:04:05"
            if _, err := time.Parse(layout, fechaStr); err != nil {
                fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
                continue
            }
            sv.siguiente_vuelo(origen, destino, fechaStr)
            fmt.Println("OK")

        default:
            fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
        }
    }
}

