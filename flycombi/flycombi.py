import sys
from cargar import construir_grafo
from comandos import ejecutar_comando

def main():
    if len(sys.argv) != 3:
        print("Uso: ./flycombi.py aeropuertos.csv vuelos.csv")
        sys.exit(1)

    aeropuertos_csv, vuelos_csv = sys.argv[1], sys.argv[2]

    grafo, aeropuertos = construir_grafo(aeropuertos_csv, vuelos_csv)


    for linea in sys.stdin:
        linea = linea.strip()
        if not linea:
            continue
        ejecutar_comando(linea, grafo, aeropuertos)

if __name__ == "main":
    main()
