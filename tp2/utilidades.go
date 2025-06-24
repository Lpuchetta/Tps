package main

import (
	"strings"
	Abb "tdas/diccionario"
	"time"
	vuelo "tp2/vuelo"
)

type vueloConFecha struct {
	fecha  time.Time
	codigo string
	info   string
}

// FechaClave se usa para ordenar en el ABB
// Se ordena por fecha, y si coinciden, por el codigo
type FechaClave struct {
	Fecha  time.Time
	Codigo string
}

type PrioridadClave struct {
	Prioridad int
	Codigo string
}

func cmpPrioridadClave(a, b PrioridadClave) int{
	if a.Prioridad != b.Prioridad{
		if a.Prioridad > b.Prioridad{
			return -1
		}
		return 1
	}
	if a.Codigo < b.Codigo{
		return -1
	} else if a.Codigo > b.Codigo{
		return 1
	}
	return 0
}
// Para ver_tablero.
func cmpVueloConFechaAsc(a, b vueloConFecha) int {
	if a.fecha.Before(b.fecha) {
		return -1
	}
	if a.fecha.After(b.fecha) {
		return 1
	}
	return strings.Compare(a.codigo, b.codigo)
}

func cmpFechaClaveAsc(a,b FechaClave) int{
	if a.Fecha.Before(b.Fecha){
		return -1
	}
	if a.Fecha.After(b.Fecha){
		return 1
	}
	if a.Codigo < b.Codigo{
		return -1
	}
	if a.Codigo > b.Codigo{
		return 1
	}
	return 0
}

func cmpFechaClaveDesc(a,b FechaClave) int{
	if a.Fecha.Before(b.Fecha){
		return 1
	}
	if a.Fecha.After(b.Fecha){
		return -1
	}
	if a.Codigo < b.Codigo{
		return 1
	}
	if a.Codigo > b.Codigo{
		return -1
	}
	return 0
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

func generarClaveFecha(v vuelo.Vuelo) FechaClave {
	return FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()}
}
func generarClavePrioridad(v vuelo.Vuelo) PrioridadClave {
	return PrioridadClave{Prioridad: v.ObtenerPrioridad(), Codigo: v.ObtenerCodigo()}
}
func generarClaveConexion(v vuelo.Vuelo) string {
	return v.ObtenerOrigen() + "-" + v.ObtenerCodigo()
} 


func (sv *SistemaVuelos) agregarUnVuelo(v vuelo.Vuelo) {
	codigo := v.ObtenerCodigo()

	// Si ya existía el vuelo, lo removemos de porFecha y porConexion
	if sv.porCodigo.Pertenece(codigo) {
		viejo := sv.porCodigo.Obtener(codigo)

		// Remover de porFecha
		cl := FechaClave{Fecha: viejo.ObtenerFecha(), Codigo: viejo.ObtenerCodigo()}
		sv.porFechaAsc.Borrar(cl)
		sv.porFechaDesc.Borrar(cl)

		prioCl := PrioridadClave{Prioridad: viejo.ObtenerPrioridad(), Codigo: viejo.ObtenerCodigo(),}
		sv.porPrioridad.Borrar(prioCl)

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
	sv.porFechaAsc.Guardar(cl, v)
	sv.porFechaDesc.Guardar(cl, v)

	prioCl := PrioridadClave{Prioridad: v.ObtenerPrioridad(), Codigo: v.ObtenerCodigo()}
	sv.porPrioridad.Guardar(prioCl, v)
	// Insertar en porConexion
	connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
	if !sv.porConexion.Pertenece(connKey) {
		sv.porConexion.Guardar(connKey, Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave))
	}
	sv.porConexion.Obtener(connKey).Guardar(cl, v)
}
