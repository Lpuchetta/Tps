package tp2
type Aeropuerto interface{
	/*Agregar_archivo recbie el nombre d eun archivo .csv y lo procesa*/ 
	Agregar_Archivo(nombreArchivo string)
	
	Ver_Tablero(K int, modo, desde, hasta string )

	Info_Vuelo(codigo string)

	Prioridad_Vuelos(K int)

	Siguiente_Vuelo(origen, destino, fecha string)

	Borrar(desde, hasta string)
}
