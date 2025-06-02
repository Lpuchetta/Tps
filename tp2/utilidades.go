package tp2

import(
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	vuelo "tdas/tp2/vuelo"
)

// FechaClave se usa para ordenar en el ABB
// Se ordena por fecha, y si coinciden, por el codigo
type FechaClave struct{
	Fecha time.Time
	Codigo	string
}

// comparadorFechaClave compara primero por fecha y ordena por esta misma, si no por Codigo
func comparadorFechaClave(a, b FechaClave) int{
	if a.Fecha.Before(b.Fecha){
		return -1
	}
	if a.Fecha.After(b.Fecha){
		return 1
	}

	if a.Codigo < b.Codigo{
		return -1
	}else if a.Codigo > b.Codigo{
		return -1
	}
	return 0
}

// compradorPrioridad compara primero por prioridad, en caso de empate compara por codigo.
func comparadorPrioridad(v1, v2 vuelo.Vuelo) int{
	diff := v1.ObtenerPrioridad() - v2.ObtenerPrioridad()
	if diff != 0{
		return diff
	}

	if v1.ObtenerCodigo() < v2.ObtenerCodigo(){
		return 1
	} else if v1.ObtenerCodigo() > v2.ObtenerCodigo(){
		return -1
	}
	return 0
}


func (sv *SistemaVuelos) agregarUnVuelo(v vuelo.Vuelo){
	codigo := v.ObtenerCodigo()

	if sv.porCodigo.Pertenece(codigo){
		anterior := sv.porCodigo.Obtener(codigo)
		claveAnt := FechaClave{Fecha: anterior.ObtenerFecha(), Codigo: codigo}
		sv.porFecha.Borrar(claveAnt)
	}

	sv.porCodigo.Guardar(codigo, v)

	claveNueva := FechaClave{ Fecha: v.ObtenerFecha(), Codigo: codigo}
	sv.porFecha.Guardar(claveNueva, v)

	sv.porPrioridad.Encolar(v)
}

