package vuelo

import (
	"fmt"
	"strconv"
	"time"
)

// ParsearLineaCSV convierte una línea del CSV en un vuelo.Vuelo
func ParsearLineaCSV(linea []string) (Vuelo, error) {
	if len(linea) < 10 {
		return nil, fmt.Errorf("línea CSV con menos columnas de las esperadas")
	}

	prioridad, err := strconv.Atoi(linea[5])
	if err != nil {
		return nil, fmt.Errorf("prioridad inválida: %w", err)
	}

	retraso, err := strconv.Atoi(linea[7])
	if err != nil {
		return nil, fmt.Errorf("retraso inválido: %w", err)
	}

	fecha, err := time.Parse("2006-01-02T15:04:05", linea[6])
	if err != nil {
		return nil, fmt.Errorf("fecha inválida: %w", err)
	}

	tiempoVuelo, err := strconv.Atoi(linea[8])
	if err != nil {
		return nil, fmt.Errorf("tiempo de vuelo inválido: %w", err)
	}

	cancelado := linea[9] == "1"

	return CrearVuelo(linea[0], linea[1], linea[2], linea[3], linea[4], prioridad, fecha, retraso, tiempoVuelo, cancelado), nil
}
