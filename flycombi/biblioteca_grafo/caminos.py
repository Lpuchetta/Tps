from collections import deque
import heapq

def bfs_escalas(grafo, origen):

    visitados = set([origen])
    padres   = {origen: None}
    niveles  = {origen: 0}
    q = deque([origen])

    while q:
        v = q.popleft()
        for w, vuelo in grafo.adyacentes(v):
            if w not in visitados:
                visitados.add(w)
                padres[w]  = v
                niveles[w] = niveles[v] + 1
                q.append(w)

    return niveles, padres 


def _dijkstra_generico(grafo, origen, destino, peso_func):
    dist   = {v: float('inf') for v in grafo.vertice()}
    padres = {origen: None}
    dist[origen] = 0

    heap = [(0, origen)]
    while heap:
        d_v, v = heapq.heappop(heap)
        if d_v > dist[v]:
            continue
        if v == destino:
            break
        for w, vuelo in grafo.adyacentes(v):
            d_nuevo = d_v + peso_func(v, w, vuelo)
            if d_nuevo < dist[w]:
                dist[w]   = d_nuevo
                padres[w] = v
                heapq.heappush(heap, (d_nuevo, w))

    return dist, padres

def dijkstra_costo(grafo, origen, destino):

    return _dijkstra_generico(
        grafo, origen, destino,
        lambda u, v, vuelo: vuelo.costo
    )

def dijkstra_tiempo(grafo, origen, destino):

    return _dijkstra_generico(
        grafo, origen, destino,
        lambda u, v, vuelo: vuelo.tiempo
    )


