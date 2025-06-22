"""Funciones para la centralidad"""
from grafo.grafo import Grafo
from biblioteca_grafo.caminos import _dijkstra_generico

def calcular_centralidad(grafo: Grafo):
    """
    Calcula la centralidad de intermediación (betweenness centrality)
    aproximada usando Dijkstra desde cada vértice.
    Devuelve una lista de tuplas (aeropuerto, score) ordenada de mayor a menor.
    Complejidad: O(A * F log A).
    """
    # Inicializamos contador global
    centralidad = {v: 0 for v in grafo.vertice()}

    # Peso para Dijkstra: costo de vuelo
    peso_costo = lambda u, v, vuelo: vuelo.costo

    for s in grafo.vertice():
        # Obtenemos distancias y padres para todos los vértices desde s
        dist, padres = _dijkstra_generico(grafo, s, None, peso_costo)

        # Inicializar contador local para este origen
        cont_local = {v: 0 for v in grafo.vertice()}

        # Filtrar vértices alcanzables (dist < inf), excluyendo el origen
        alcanzables = [v for v, d in dist.items() if d < float('inf') and v != s]
        # Ordenar de mayor a menor distancia
        alcanzables.sort(key=lambda v: dist[v], reverse=True)

        # Propagar conteos desde los más lejanos a los más cercanos
        for v in alcanzables:
            p = padres.get(v)
            if p is not None:
                # 1 + todo lo que pase por v
                cont_local[p] += 1 + cont_local[v]

        # Acumular en centralidad global (excluyendo origen)
        for v in centralidad:
            if v != s:
                centralidad[v] += cont_local[v]

    # Convertir a lista ordenada
    ranking = sorted(centralidad.items(), key=lambda x: x[1], reverse=True)
    return ranking
