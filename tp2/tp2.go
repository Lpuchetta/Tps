package main

type Aeropuerto interface {

	//agregar_archivo() procesa de forma completa un archivo de .csv que contiene datos de vuelos.
	agregar_archivo(nombreArchivo string) error

	//ver_tablero() muestra los K vuelos ordenados por fecha de forma ascendente (asc) o descendente (desc), 
	//cuya fecha de despegue esté dentro de el intervalo <desde> <hasta> (inclusive).
	ver_tablero(K int, modo, desde, hasta string)

	//info_vuelo() muestra toda la información posible en sobre el vuelo que tiene el código pasado por parámetro.
	info_vuelo(codigo string)

	//prioridad_vuelos() muestra los códigos de los K vuelos que tienen mayor prioridad.
	prioridad_vuelos(K int)

	//siguiente_vuelo() muestra la información del vuelo (tal cual en info_vuelo) del próximo vuelo directo que conecte los aeropuertos de origen y destino, 
	//a partir de la fecha indicada (inclusive).
	// Si no hay un siguiente vuelo cargado, imprimir No hay vuelo registrado desde <aeropuerto origen> hacia <aeropuerto destino> desde <fecha> 
	//(con los valores que correspondan).
	siguiente_vuelo(origen, destino, fecha string)

	//borrar() borra todos los vuelos cuya fecha de despegue estén dentro del intervalo <desde> <hasta> (inclusive). 
	borrar(desde, hasta string)
}
