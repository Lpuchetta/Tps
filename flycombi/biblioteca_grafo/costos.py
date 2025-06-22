from union_find import UnionFind
from grafo.grafo import Grafo, Vuelo

def kruskal(grafo: Grafo):
    """
    Devuelve la lista de aristas del MST como
    (u, v, tiempo, costo, frecuencia). Complejidad O(F log A).
    """
    # 1. Inicializar Union-Find con todos los vértices
    uf = UnionFind(grafo.vertice())

    # 2. Recolectar y transformar las aristas
    raw_edges = obtener_aristas(grafo)  # [(u, v, vuelo), …]
    aristas = []
    for u, v, vuelo in raw_edges:
        aristas.append((u, v, vuelo.tiempo, vuelo.costo, vuelo.frecuencia))

    # 3. Ordenar por costo (índice 3)
    aristas.sort(key=lambda x: x[3])

    # 4. Construir MST
    mst = []
    for u, v, tiempo, costo, freq in aristas:
        if uf.find(u) != uf.find(v):
            uf.union(u, v)
            mst.append((u, v, tiempo, costo, freq))

    return mst


def obtener_aristas(grafo):
    aristas = []
    for u in grafo.vertice():
        for v, vuelo in grafo.adyacentes(u):
            aristas.append((u,v, vuelo))
    return aristas
