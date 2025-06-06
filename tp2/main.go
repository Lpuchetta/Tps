package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "tdas/tp2" 
)
func main() {
	sv := tp2.CrearSistemaDeVuelos()

	scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("Ingrese un comando (agregar_archivo, ver_tablero, info_vuelo, prioridad_vuelos, siguiente_vuelo, borrar):")
        if !scanner.Scan() {
            break
        }
        linea := scanner.Text()
        partes := strings.Fields(linea)

        if len(partes) == 0 {
            continue
        }

        comando := strings.ToLower(partes[0])

        switch comando {
	        case "agregar_archivo":
	            if len(partes) < 2 {
	                fmt.Println("Error: falta nombre de archivo")
	                continue
	            }
	            archivo := partes[1]
	            err := sv.agregar_archivo(archivo)
	            if err != nil {
	                fmt.Println("Error agregando archivo:", err)
	            } else {
	                fmt.Println("Archivo agregado correctamente")
	            }

			case "ver_tablero":
			    if len(partes) < 5 {
			        fmt.Println("Error: ver_tablero requiere 4 argumentos: K modo desde hasta")
			        continue
			    }
			
			    K, err := strconv.Atoi(partes[1])
			    if err != nil || K <= 0 {
			        fmt.Println("Error: K debe ser un número entero mayor a 0")
			        continue
			    }
			
			    modo := strings.ToLower(partes[2])
			    if modo != "asc" && modo != "desc" {
			        fmt.Println("Error: modo debe ser 'asc' o 'desc'")
			        continue
			    }
			
			    desdeStr := partes[3]
			    hastaStr := partes[4]
			
			
			    const layout = "2006-01-02T15:04:05"
			    desde, errDesde := time.Parse(layout, desdeStr)
			    hasta, errHasta := time.Parse(layout, hastaStr)
			
			    if errDesde != nil {
			        fmt.Println("Error: formato inválido para 'desde'. Debe ser YYYY-MM-DDTHH:MM:SS")
			        continue
			    }
			    if errHasta != nil {
			        fmt.Println("Error: formato inválido para 'hasta'. Debe ser YYYY-MM-DDTHH:MM:SS")
			        continue
			    }
			
			    if hasta.Before(desde) {
			        fmt.Println("Error: 'hasta' no puede ser anterior a 'desde'")
			        continue
			    }
			
			    sv.ver_tablero(K, modo, desdeStr, hastaStr)

			case "info_vuelo":
				if len(partes) < 2 {
					fmt.Println("Error: falta el código de vuelo")
					continue
				}
				codigo := partes[1]
				sv.info_vuelo(codigo)
			
			case "prioridad_vuelos":
				if len(partes) < 2 {
					fmt.Println("Error: falta K (cantidad de vuelos)")
					continue
				}
				K, err := strconv.Atoi(partes[1])
				if err != nil || K <= 0 {
					fmt.Println("K inválido. Debe ser un número entero mayor a 0.")
					continue
				}
				sv.prioridad_vuelos(K)

			case "siguiente_vuelo":
				if len(partes) < 4 {
					fmt.Println("Error: siguiente_vuelo requiere 3 argumentos: origen destino fecha")
					continue
				}
			
				origen := partes[1]
				destino := partes[2]
				fechaStr := partes[3]
			
				const layout = "2006-01-02T15:04:05"
				_, err := time.Parse(layout, fechaStr)
				if err != nil {
					fmt.Println("Error: formato inválido de fecha. Debe ser YYYY-MM-DDTHH:MM:SS")
					continue
				}
			
				sv.siguiente_vuelo(origen, destino, fechaStr)

			case "borrar":
				if len(partes) < 3 {
					fmt.Println("Error: borrar requiere 2 argumentos: desde hasta")
					continue
				}
			
				desdeStr := partes[1]
				hastaStr := partes[2]
			
				const layout = "2006-01-02T15:04:05"
				desde, errDesde := time.Parse(layout, desdeStr)
				hasta, errHasta := time.Parse(layout, hastaStr)
			
				if errDesde != nil || errHasta != nil {
					fmt.Println("Error: formato inválido de fecha. Debe ser YYYY-MM-DDTHH:MM:SS")
					continue
				}
			
				if hasta.Before(desde) {
					fmt.Println("Error: 'hasta' no puede ser anterior a 'desde'")
					continue
				}
			
				sv.borrar(desdeStr, hastaStr)

			default:
				fmt.Println("Comando desconocido")
		}	
	}
}
