class Aeropuerto:
    """Representa un aeropuerto"""
    def __init__(self, codigo, ciudad, lat, lon):
        self.codigo = codigo
        self.ciudad = ciudad
        self.lat = lat
        self.lon = lon

    def __eq__(self, otro):
        return isinstance(otro, Aeropuerto) and self.codigo == otro.codigo  

    def __hash__(self):
        return hash(self.codigo)

    def __repr__(self):
        return f"Aeropuerto({self.codigo})"
class Vuelo:
    """Informacion de un vuelo entre aeropuertos"""

    def __init__(self, origen, destino, tiempo, costo, frecuencia):
        self.origen = origen
        self.destino = destino
        self.tiempo = tiempo
        self.costo = costo
        self.frecuencia = frecuencia

    def __repr__(self):
        return (f"Vuelo({self.origen.codigo}->{self.destino.codigo}, "
                f"tiempo={self.tiempo}, costo={self.costo}, frecuencia={self.frecuencia})")

class Grafo:

    def __init__(self, dirigido=True):
        #mapa vertice -> lista tuplas (destino, vuelo)
        self.ady = {}
        self.dirigido = dirigido

    def agregar_vertice(self, v):
        """Agrega un vertice si no existe"""
        if v not in self.ady:
            self.ady[v] = []

    def agregar_arista(self, u, v, tiempo, costo, frecuencia):
        """Agrega una arista u -> v creando el Vuelo correspondiente"""
        self.agregar_vertice(u)
        self.agregar_vertice(v)
        vuelo = Vuelo(u, v, tiempo, costo, frecuencia)
        self.ady[u].append((v,vuelo))
        if not self.dirigido:
            #en grafo no dirigido tambien agregar v->u
            vuelo_inv = Vuelo(v,u, tiempo, costo, frecuencia)
            self.ady[v].append((u,vuelo_inv))

    def adyacentes(self, v):
        """Devuelve una lista de pares (destino, Vuelo) de v"""
        return list(self.ady.get(v,[]))

    def vertice(self):
        """Devuelve todos los vertices del grafo"""
        return list(self.ady.keys())

    def eliminar_arista(self, u, v, vuelo=None):
        """Eliminamos la conexion entre u y v. Si indicas un vuelo, eliminamos solo ese"""
        # Si u tiene aristas eliminamos manualmente
        if u in self.ady:
            nueva_lista = []
            for destino, info in self.ady[u]:
                if destino == v and (vuelo is None or info == vuelo):
                    continue
                nueva_lista.append((destino,info))
            self.ady[u] = nueva_lista
        # Si el grafo es no dirigido, hacemos lo mismo de v a u
        if not self.dirigido and v in self.ady:
            nueva_lista = []
            for destino, info in self.ady[v]:
                if destino == u and (vuelo is None or info == vuelo):
                    continue
                nueva_lista.append((destino,info))
            self.ady[v] = nueva_lista

    def eliminar_vertice(self, v):
        """Elimina el vertice pasado por parametro"""
        if v not in self.ady:
            return

        for u in list(self.ady.keys()):
            if u == v:
                continue
            nueva_lista = [(destino, info) for destino, info in self.ady[u] if destino != v]
            self.ady[u] = nueva_lista

        del self.ady[v]
        
    
