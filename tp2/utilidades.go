package main

import (
	Abb "tdas/diccionario"
	"time"
	vuelo "tp2/vuelo"
	"strings"
)

type vueloConFecha struct {
	fecha time.Time
	codigo string
	info  string
}

// FechaClave se usa para ordenar en el ABB
// Se ordena por fecha, y si coinciden, por el codigo
type FechaClave struct {
	Fecha  time.Time
	Codigo string
}

// Para recorrer fechas en orden ascendente
func cmpFechaAsc(a, b FechaClave) int {
	if a.Fecha.Before(b.Fecha) {
		return -1
	}
	if a.Fecha.After(b.Fecha) {
		return 1
	}
	return strings.Compare(a.Codigo,b.Codigo)
}

func cmpVueloConFechaAsc(a, b vueloConFecha) int {
	if a.fecha.Before(b.fecha) {
		return -1
	}
	if a.fecha.After(b.fecha) {
		return 1
	}
	return strings.Compare(a.codigo,b.codigo)
}



// Para el ABB de porConexion (fecha + código)
func comparadorFechaClave(a, b FechaClave) int {
	if a.Fecha.Before(b.Fecha) {
		return -1
	}
	if a.Fecha.After(b.Fecha) {
		return 1
	}
	if a.Codigo < b.Codigo {
		return -1
	}
	if a.Codigo > b.Codigo {
		return 1
	}
	return 0
}


func (sv *SistemaVuelos) agregarUnVuelo(v vuelo.Vuelo) {
	codigo := v.ObtenerCodigo()

	// Si ya existía el vuelo, lo removemos de porFecha y porConexion
	if sv.porCodigo.Pertenece(codigo) {
		viejo := sv.porCodigo.Obtener(codigo)

		// Remover de porFecha
		cl := FechaClave{Fecha: viejo.ObtenerFecha(), Codigo: viejo.ObtenerCodigo()}
		sv.porFecha.Borrar(cl)

		// Remover de porConexion
		connKey := viejo.ObtenerOrigen() + "-" + viejo.ObtenerDestino()
		if sv.porConexion.Pertenece(connKey) {
			abb := sv.porConexion.Obtener(connKey)
			abb.Borrar(cl)
			if abb.Cantidad() == 0 {
				sv.porConexion.Borrar(connKey)
			}
		}
	}

	// Insertar nuevo vuelo en porCodigo
	sv.porCodigo.Guardar(codigo, v)

	// Insertar en porFecha
	cl := FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()}
	sv.porFecha.Guardar(cl, v)

	// Insertar en porConexion
	connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
	if !sv.porConexion.Pertenece(connKey) {
		sv.porConexion.Guardar(connKey, Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave))
	}
	sv.porConexion.Obtener(connKey).Guardar(cl, v)
}

