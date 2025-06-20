import csv
import sys
from biblioteca_grafo.caminos import dijsktra_tiempo dijkstra_costo, bfs_escalas
from biblioteca_grafo.costos import kruskal
from biblioteca_grafo.ranking import calcular_centralidad
from biblioteca_grafo.cultural import orden_topologico


def ejecutar_comando(linea: str, grafo, aeropuertos):
    # separar comando y parametros
    partes = linea.replace(',', ' ').split()
    if not partes:
        return
    comando = partes[0]

    if comando == 'camino_mas':
        # aca es camino_mas <barato|rapido> origen destino
        criterio = partes[1]
        origen = partes[2]
        destino = partes[3]
        _camino_mas(grafo, aeropuertos, criterio, origen, destino)

    elif comando == 'camino_escalas':
        # aca es camino_escalas origen destino
        origen = partes[1]
        destino = partes[2]
        _camino_escalas(grafo, aeropuertos, origen, destino)

    elif comando == 'centralidad':
        # aca es centralidad n (n cantidad de mas importantes a mostrar)
        n = int(partes[1])
        _centralidad(grafo, aeropuertos, n)

    elif comando == 'nueva_aerolinea':
        salida = partes[1]
        _nueva_aerolinea(grafo, ruta_salida)

    elif comando == 'itinerario':
        ruta = int(partes[1])
        _itinerario(grafo, aeropuertos, ruta)
    else:
        print(f'Comando desconocido: {comando}')

def _camino_mas(grafo, aeropuertos, criterio, origen_ciudad, destino_ciudad):
    # obtengo codigos de aeropuertos de cada ciudad
    origenes = [a for codigo, a in aeropuertos.items() if a.ciudad == origen_ciudad]
    destinos = [a for codigo, a in aeropuertos.items() if a.ciudad == destino_ciudad]
    mejor_camino = None
    for o in origenes:
        for d in destinos:
            if criterio == 'barato':
                dist, prev = dijkstra_costo(grafo, o, d)
            else:
                dist, prev = dijkstra_tiempo(grafo, o, d)
            camino = _reconstruir_camino(prev, o, d)
            if camino and (mejor_camino is None or dist < mejor_camino[0]):
                mejor_camino = (dist, camino)
    if mejor_camino:
        print(' -> '.join(a.codigo for a in mejor_camino[1]))

def _camino_escalas(grafo, aeropuertos, origen_ciudad, destino_ciudad):
    
    origenes = [a for codigo, a in aeropuertos.items() if a.ciudad == origen_ciudad]
    destinos = [a for codigo, a in aeropuertos.items() if a.ciudad == destino_ciudad]
    mejor = None
    for o in origenes:
        for d in destinos:
            niveles, prev = bfs_escalas(grafo, o)
            if d in niveles:
                escalas = niveles[d]
                if mejor is None or escalas < mejor[0]:
                    mejor = (escalas, o, d, prev)
    if mejor:
        _, o, d, prev = mejor
        camino = _reconstruir_camino(prev, o, d)
        print(' -> '.join(a.codigo for a in camino))

def _centralidad():
    #despues la hago, medio confusa


def _nueva_aerolinea(grafo, ruta_salida):
    arbol = kruskal(grafo)
    with open(ruta_salida, 'w', newline='', encoding='utf-8') as archivo:
        writer = csv.writer(archivo)
        for u, v, tiempo, costo, freq in arbol:
            writer.writerow([u.codigo, v.codigo, tiempo, costo, freq])
    print('OK')

def _itinerario(grafo, aeropuertos, ruta_entrada):
    ciudad = []
    restricciones = []
    with open(ruta_entrada, newline='', encoding='utf-8') as archivo:
        reader = csv.reader(archivo)
        filas = list(reader)
        ciudad = [c.strip() for c in filas[0] if c.strip()]
        for pre, post in filas[1:]:
            restricciones.append((pre.strip(), post.strip()))
    orden = ordenar_topologico(ciudades, restricciones)
    print(', '.join(orden))
    for i in range(len(orden) - 1):
        _camino_mas(grafo, aeropuertos, 'rapido', orden[i], orden[i + 1])
        
def _reconstruir_camino(prev, origen, destino):
    camino = []
    u = destino
    while u:
        camino.append(u)
        u = prev.get(u)
    if camino[-1] != origen:
        return None
    return list(reversed(camino))
