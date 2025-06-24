package main

import (
	"strings"
	Abb "tdas/diccionario"
	Lista "tdas/lista"
	"time"
	vuelo "tp2/vuelo"
)

type FechaClave struct {
	Fecha  time.Time
	Codigo string
}
type PrioridadClave struct {
	Prioridad int
	Codigo    string
}

func compararIntAsc(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
func compararIntDesc(a, b int) int {
	if a < b {
		return 1
	}
	if a > b {
		return -1
	}
	return 0
}

func compararStringAsc(a, b string) int {
	return strings.Compare(a, b)
}

func compararStringsDesc(a, b string) int {
	return -strings.Compare(a, b)
}

func compararTimeAsc(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}
func compararTimeDesc(a, b time.Time) int {
	if a.Before(b) {
		return 1
	}
	if a.After(b) {
		return -1
	}
	return 0
}

func cmpPrioridadClave(a, b PrioridadClave) int {
	if cmp := compararIntDesc(a.Prioridad, b.Prioridad); cmp != 0 {
		return cmp
	}
	return compararStringAsc(a.Codigo, b.Codigo)
}

func cmpFechaClaveAsc(a, b FechaClave) int {
	if cmp := compararTimeAsc(a.Fecha, b.Fecha); cmp != 0 {
		return cmp
	}
	return compararStringAsc(a.Codigo, b.Codigo)
}

/*func cmpFechaClaveDesc(a, b FechaClave) int {
	if cmp := compararTimeDesc(a.Fecha, b.Fecha); cmp != 0 {
		return cmp
	}
	return compararStringsDesc(a.Codigo, b.Codigo)
}*/

func generarClaveFecha(v vuelo.Vuelo) FechaClave {
	return FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()}
}
func generarClavePrioridad(v vuelo.Vuelo) PrioridadClave {
	return PrioridadClave{Prioridad: v.ObtenerPrioridad(), Codigo: v.ObtenerCodigo()}
}
func generarClaveConexion(v vuelo.Vuelo) string {
	return v.ObtenerOrigen() + "-" + v.ObtenerDestino()
}

// inserta v en la lista manteniendo orden lex ascendente de código
func insertarOrdenadoAsc(lista Lista.Lista[vuelo.Vuelo], v vuelo.Vuelo) {
	iter := lista.Iterador()
	for iter.HaySiguiente() && iter.VerActual().ObtenerCodigo() < v.ObtenerCodigo() {
		iter.Siguiente()
	}
	iter.Insertar(v)
}

// inserta v en la lista manteniendo orden lex descendente de código
func insertarOrdenadoDesc(lista Lista.Lista[vuelo.Vuelo], v vuelo.Vuelo) {
	iter := lista.Iterador()
	// avanzamos mientras el código actual sea > que el nuevo
	for iter.HaySiguiente() && iter.VerActual().ObtenerCodigo() > v.ObtenerCodigo() {
		iter.Siguiente()
	}
	iter.Insertar(v)
}

func (sv *SistemaVuelos) agregarUnVuelo(v vuelo.Vuelo) {
	codigo := v.ObtenerCodigo()

	// 1) Si ya existía, lo borramos completamente
	if sv.porCodigo.Pertenece(codigo) {
		viejo := sv.porCodigo.Obtener(codigo)
		sv.borrarVuelo(viejo)
	}

	// 2) Guardamos en el hash global
	sv.porCodigo.Guardar(codigo, v)

	// 3) Insertar en porFechaAsc
	fecha := v.ObtenerFecha()
	if !sv.porFechaAsc.Pertenece(fecha) {
		// creamos la lista vacía para ese día
		sv.porFechaAsc.Guardar(fecha, Lista.CrearListaEnlazada[vuelo.Vuelo]())
	}
	// ahora insertamos ordenado por código
	listaAsc := sv.porFechaAsc.Obtener(fecha)
	insertarOrdenadoAsc(listaAsc, v)

	// 4) Insertar en porFechaDesc (mismo tratamiento: la lista ya está
	//    ordenada por código, pero la recorrerás en sentido inverso
	//    si usas compararTimeDesc en el ABB externo)
	if !sv.porFechaDesc.Pertenece(fecha) {
		sv.porFechaDesc.Guardar(fecha, Lista.CrearListaEnlazada[vuelo.Vuelo]())
	}
	listaDesc := sv.porFechaDesc.Obtener(fecha)
	insertarOrdenadoDesc(listaDesc, v)

	// 5) El resto igual: prioridad y conexiones
	prioCl := generarClavePrioridad(v)
	sv.porPrioridad.Guardar(prioCl, v)

	connKey := generarClaveConexion(v)
	if !sv.porConexion.Pertenece(connKey) {
		sv.porConexion.Guardar(connKey, Abb.CrearABB[FechaClave, vuelo.Vuelo](cmpFechaClaveAsc))
	}
	cl := generarClaveFecha(v)
	sv.porConexion.Obtener(connKey).Guardar(cl, v)
}

func (sv *SistemaVuelos) borrarVuelo(v vuelo.Vuelo) {
	fecha := v.ObtenerFecha()
	codigo := v.ObtenerCodigo()

	// 1) porFechaAsc
	if sv.porFechaAsc.Pertenece(fecha) {
		lista := sv.porFechaAsc.Obtener(fecha)
		iter := lista.Iterador()
		for iter.HaySiguiente() {
			if iter.VerActual().ObtenerCodigo() == codigo {
				iter.Borrar()
				break
			}
			iter.Siguiente()
		}
		if lista.EstaVacia() {
			sv.porFechaAsc.Borrar(fecha)
		}
	}

	if sv.porFechaDesc.Pertenece(fecha) {
		lista := sv.porFechaDesc.Obtener(fecha)
		iter := lista.Iterador()
		for iter.HaySiguiente() {
			if iter.VerActual().ObtenerCodigo() == codigo {
				iter.Borrar()
				break
			}
			iter.Siguiente()
		}
		if lista.EstaVacia() {
			sv.porFechaDesc.Borrar(fecha)
		}
	}

	prioCl := generarClavePrioridad(v)
	sv.porPrioridad.Borrar(prioCl)

	connKey := generarClaveConexion(v)
	if sv.porConexion.Pertenece(connKey) {
		abb := sv.porConexion.Obtener(connKey)
		cl := generarClaveFecha(v)
		abb.Borrar(cl)
		if abb.Cantidad() == 0 {
			sv.porConexion.Borrar(connKey)
		}
	}

	sv.porCodigo.Borrar(codigo)
}
