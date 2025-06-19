from cargar import construir_grafo

grafo, aeropuertos = construir_grafo("aerodata/aeropuertos.csv", "aerodata/vuelos.csv")

print("Aeropuertos cargados:")
for aeropuerto in grafo.vertice():
    print(f"{aeropuerto.codigo} - {aeropuerto.ciudad}")

print("\nVuelos cargados:")
for aeropuerto in grafo.vertice():
    for destino, vuelo in grafo.adyacentes(aeropuerto):
        print(f"{vuelo.origen.codigo} -> {vuelo.destino.codigo} | tiempo: {vuelo.tiempo}, costo: {vuelo.costo}, freq: {vuelo.frecuencia}")
