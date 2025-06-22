"""Para el itinerario, orden topologico"""
from grafo.grafo import Grafo
from collections import deque

def ordenar_topologico(ciudades, restricciones):
    """
    Dado una lista de ciudades y una lista de restricciones (ciudad_pre, ciudad_post),
    devuelve un orden topológico de las ciudades que respeta las restricciones.
    """
    grafo = Grafo(dirigido=True)

    # Agregamos los vértices (las ciudades a visitar)
    for ciudad in ciudades:
        grafo.agregar_vertice(ciudad)

    # Agregamos las restricciones como aristas dirigidas
    for pre, post in restricciones:
        grafo.agregar_arista(pre, post, 1, 1, 1)  # pesos dummy: no importan

    # Realizamos ordenamiento topológico por grados de entrada
    return _topologico_grados(grafo)


def _grados_entrada(grafo):
    """Devuelve un diccionario con la cantidad de aristas entrantes por vértice."""
    g_ent = {v: 0 for v in grafo.vertice()}
    for v in grafo.vertice():
        for w, _ in grafo.adyacentes(v):
            g_ent[w] += 1
    return g_ent


def _topologico_grados(grafo):
    """
    Ordenamiento topológico usando el algoritmo de Kahn (por grados de entrada).
    Retorna una lista con los vértices en orden válido.
    """
    g_ent = _grados_entrada(grafo)
    cola = deque()
    orden = []

    for v in grafo.vertice():
        if g_ent[v] == 0:
            cola.append(v)

    while cola:
        v = cola.popleft()
        orden.append(v)
        for w, _ in grafo.adyacentes(v):
            g_ent[w] -= 1
            if g_ent[w] == 0:
                cola.append(w)
    return orden
