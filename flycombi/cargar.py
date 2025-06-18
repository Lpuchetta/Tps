import csv
from grafo.grafo import Grafo, Aeropuerto

def cargar_aeropuertos(ruta_archivo):
    """
    Lee el archivo de aeropuertos y devuelve un dict de código -> Aeropuerto.
    Formato esperado: ciudad,codigo_aeropuerto,latitud,longitud
    """
    aeropuertos = {}

    with open(ruta_archivo, newline='', encoding='utf-8') as archivo:
        lector = csv.reader(archivo)
        for fila in lector:
            if len(fila) < 4:
                continue
            ciudad, codigo, lat, lon = fila
            aeropuerto = Aeropuerto(codigo, ciudad, float(lat), float(lon))
            aeropuertos[codigo] = aeropuerto
    return aeropuertos

def cargar_vuelos(ruta_archivo, grafo, aeropuertos):
    """
    Lee el archivo de vuelos y agrega aristas al grafo.
    Formato: aeropuerto_i,aeropuerto_j,tiempo_promedio,precio,cant_vuelos
    """
    with open(ruta_archivo, newline='', encoding='utf-8') as archivo:
        lector = csv.reader(archivo)
        for fila in lector:
            if len(fila) < 5:
                continue
            origen_cod, destino_cod, tiempo, costo, frecuencia = fila

            origen = aeropuertos.get(origen_cod)
            destino = aeropuertos.get(destino_cod)

            if origen and destino:
                grafo.agregar_arista(
                    origen,
                    destino,
                    float(tiempo),
                    float(costo),
                    int(frecuencia)
                )

def construir_grafo(path_aeropuertos, path_vuelos):
    """
    Construye y devuelve un grafo dirigido con vértices y aristas cargados.
    Retorna: (grafo, dict_codigo_aeropuerto)
    """
    grafo = Grafo(dirigido=True)
    aeropuertos = cargar_aeropuertos(path_aeropuertos)

    # Agregar vértices
    for aeropuerto in aeropuertos.values():
        grafo.agregar_vertice(aeropuerto)

    # Agregar aristas desde el CSV de vuelos
    cargar_vuelos(path_vuelos, grafo, aeropuertos)

    return grafo, aeropuertos
