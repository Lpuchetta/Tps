package tp2
import(
	"io"
	"time"
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"fmt"
	Hash	"tdas/diccionario"
	Abb		"tdas/diccionario"
	Heap	"tdas/heap"
	vuelo	"tdas/tp2/vuelo"
)

type SistemaVuelos struct{
	porCodigo	Hash.Diccionario[string,vuelo.Vuelo]		
	porFecha	Abb.DiccionarioOrdenado[time.Time, vuelo.Vuelo]
	porPrioridad	Heap.ColaPrioridad[vuelo.Vuelo]
}

func CrearSistemaDeVuelos() Aeropuerto{
	comparador := func(v1, v2 vuelo.Vuelo) int{
		return v1.ObtenerPrioridad() - v2.ObtenerPrioridad()
	}

	return &SistemaVuelos{
		porCodigo: Hash.CrearHash[string,vuelo.Vuelo](),
		porFecha: Abb.CrearABB[time.Time, vuelo.Vuelo](func(f1, f2 time.Time) int{
			if f1.Before(f2){
				return -1
			} else if f1.After(f2){
				return 1
			}
			return 0
		}),
		porPrioridad: Heap.CrearHeap[vuelo.Vuelo](comparador),
	}
}


func (sv *SistemaVuelos) Agregar_Archivo(nombreArchivo string) error{
	archivo, err := os.Open(nombreArchivo)

	if err != nil{
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer archivo.Close()

	reader := csv.NewReader(archivo)
	reader.Comma = ','

	for {
		linea, err := reader.Read()
		if err == io.EOF{
			break
		}
		if err != nil{
			return fmt.Errorf("error leyendo el CSV: %w", err)
		}
		
		v, err := parsearLineaCSV(linea)
		if err != nil{
			fmt.Println("Linea invalida")
			continue
		}
		codigo := v.ObtenerCodigo()
		if  sv.porCodigo.Pertenece(codigo){
			vueloViejo := sv.porCodigo.Obtener(codigo)
			sv.porFecha.Borrar(vueloViejo.ObtenerFecha())
		}
		sv.porCodigo.Guardar(codigo,v)
		sv.porFecha.Guardar(v.ObtenerFecha(),v)
		sv.porPrioridad.Encolar(v)
	}
	return nil
}

func parsearLineaCSV(linea []string)(vuelo.Vuelo, error){
	if len(linea) < 10{
		return nil, fmt.Errorf("error leyendo el archivo CSV")
	}

	codigoVuelo := linea[0]
	aerolinea := linea[1]
	origen := linea[2]
	destino := linea[3]
	matricula := linea[4]
	prioridad, _ := strconv.Atoi(linea[5])
	retraso, _ := strconv.Atoi(linea[7])
	fecha, _ := time.Parse("2006-01-02T15:04:05", linea[6])
	tiempoVuelo, _ := strconv.Atoi(linea[8])
	cancelado := linea[9] == "1"
	v := vuelo.CrearVuelo(codigoVuelo,aerolinea,origen,destino,matricula,prioridad,fecha,retraso,tiempoVuelo,cancelado)
	return v, nil
}
	
	
