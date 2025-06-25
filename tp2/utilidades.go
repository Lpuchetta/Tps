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

func cmpFechaClaveAsc(a, b FechaClave) int {
	if cmp := compararTimeAsc(a.Fecha, b.Fecha); cmp != 0 {
		return cmp
	}
	return compararStringAsc(a.Codigo, b.Codigo)
}

func generarClaveFecha(v vuelo.Vuelo) FechaClave {
	return FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()}
}

func generarClaveConexion(v vuelo.Vuelo) string {
	return v.ObtenerOrigen() + "-" + v.ObtenerDestino()
}

// inserta v en la lista manteniendo orden lex ascendente de código
func insertarOrdenadoAsc(lista Lista.Lista[vuelo.Vuelo], vuelo vuelo.Vuelo) {
	iterador := lista.Iterador()
	for iterador.HaySiguiente() && iterador.VerActual().ObtenerCodigo() < vuelo.ObtenerCodigo() {
		iterador.Siguiente()
	}
	iterador.Insertar(vuelo)
}

// inserta v en la lista manteniendo orden lex descendente de código
func insertarOrdenadoDesc(lista Lista.Lista[vuelo.Vuelo], vuelo vuelo.Vuelo) {
	iterador := lista.Iterador()
	for iterador.HaySiguiente() && iterador.VerActual().ObtenerCodigo() > vuelo.ObtenerCodigo() {
		iterador.Siguiente()
	}
	iterador.Insertar(vuelo)
}

// inserta código en lista de strings en orden ascendente
func insertarOrdenadoAscString(lista Lista.Lista[string], codigo string) {
	iterador := lista.Iterador()
	for iterador.HaySiguiente() && iterador.VerActual() < codigo {
		iterador.Siguiente()
	}
	iterador.Insertar(codigo)
}

func (sv *SistemaVuelos) agregarUnVuelo(vueloActual vuelo.Vuelo) {
	codigo := vueloActual.ObtenerCodigo()

	// 1) Si ya existía, lo borramos completamente
	if sv.porCodigo.Pertenece(codigo) {
		viejo := sv.porCodigo.Obtener(codigo)
		sv.borrarVuelo(viejo)
	}

	// 2) Guardamos en el hash global
	sv.porCodigo.Guardar(codigo, vueloActual)

	// 3) Insertar en porFechaAsc
	fecha := vueloActual.ObtenerFecha()
	if !sv.porFechaAsc.Pertenece(fecha) {
		sv.porFechaAsc.Guardar(fecha, Lista.CrearListaEnlazada[vuelo.Vuelo]())
	}
	listaAsc := sv.porFechaAsc.Obtener(fecha)
	insertarOrdenadoAsc(listaAsc, vueloActual)

	// 4) Insertar en porFechaDesc
	if !sv.porFechaDesc.Pertenece(fecha) {
		sv.porFechaDesc.Guardar(fecha, Lista.CrearListaEnlazada[vuelo.Vuelo]())
	}
	listaDesc := sv.porFechaDesc.Obtener(fecha)
	insertarOrdenadoDesc(listaDesc, vueloActual)

	// 5) Guardar en porPrioridad
	prioridad := vueloActual.ObtenerPrioridad()
	if !sv.porPrioridad.Pertenece(prioridad) {
		sv.porPrioridad.Guardar(prioridad, Lista.CrearListaEnlazada[string]())
	}
	lista := sv.porPrioridad.Obtener(prioridad)
	insertarOrdenadoAscString(lista, codigo)

	// 6) Guardar en porConexion
	claveConexion := generarClaveConexion(vueloActual)
	if !sv.porConexion.Pertenece(claveConexion) {
		sv.porConexion.Guardar(claveConexion, Abb.CrearABB[FechaClave, vuelo.Vuelo](cmpFechaClaveAsc))
	}
	claveFecha := generarClaveFecha(vueloActual)
	sv.porConexion.Obtener(claveConexion).Guardar(claveFecha, vueloActual)
}

func (sv *SistemaVuelos) borrarVuelo(vueloActual vuelo.Vuelo) {
	fecha := vueloActual.ObtenerFecha()
	codigo := vueloActual.ObtenerCodigo()

	// 1) porFechaAsc
	if sv.porFechaAsc.Pertenece(fecha) {
		lista := sv.porFechaAsc.Obtener(fecha)
		iterador := lista.Iterador()
		for iterador.HaySiguiente() {
			if iterador.VerActual().ObtenerCodigo() == codigo {
				iterador.Borrar()
				break
			}
			iterador.Siguiente()
		}
		if lista.EstaVacia() {
			sv.porFechaAsc.Borrar(fecha)
		}
	}

	// 2) porFechaDesc
	if sv.porFechaDesc.Pertenece(fecha) {
		lista := sv.porFechaDesc.Obtener(fecha)
		iterador := lista.Iterador()
		for iterador.HaySiguiente() {
			if iterador.VerActual().ObtenerCodigo() == codigo {
				iterador.Borrar()
				break
			}
			iterador.Siguiente()
		}
		if lista.EstaVacia() {
			sv.porFechaDesc.Borrar(fecha)
		}
	}

	// 3) porPrioridad
	prioridad := vueloActual.ObtenerPrioridad()
	if sv.porPrioridad.Pertenece(prioridad) {
		lista := sv.porPrioridad.Obtener(prioridad)
		iterador := lista.Iterador()
		for iterador.HaySiguiente() {
			if iterador.VerActual() == codigo {
				iterador.Borrar()
				break
			}
			iterador.Siguiente()
		}
		if lista.EstaVacia() {
			sv.porPrioridad.Borrar(prioridad)
		}
	}

	// 4) porConexion
	claveConexion := generarClaveConexion(vueloActual)
	if sv.porConexion.Pertenece(claveConexion) {
		abb := sv.porConexion.Obtener(claveConexion)
		claveFecha := generarClaveFecha(vueloActual)
		abb.Borrar(claveFecha)
		if abb.Cantidad() == 0 {
			sv.porConexion.Borrar(claveConexion)
		}
	}

	// 5) porCodigo
	sv.porCodigo.Borrar(codigo)
}

