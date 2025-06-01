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

// FechaClave se usa para ordenar en el ABB
// Se ordena por fecha, y si coinciden, por el codigo
type FechaClave struct{
	Fecha time.Time
	Codigo	string
}

// comparadorFechaClave compara primero por fecha y ordena por esta misma, si no por Codigo
func comparadorFechaClave(a, b FechaClave) int{
	if a.Fecha.Before(b.Fecha){
		return -1
	}
	if a.Fecha.After(b.Fecha){
		return 1
	}

	if a.Codigo < b.Codigo{
		return -1
	}else if a.Codigo > b.Codigo{
		return -1
	}
	return 0
}

// compradorPrioridad compara primero por prioridad, en caso de empate compara por codigo.
func comparadorPrioridad(v1, v2 vuelo.Vuelo) int{
	diff := v1.ObtenerPrioridad() - v2.ObtenerPrioridad()
	if diff != 0{
		return diff
	}

	if v1.ObtenerCodigo() < v2.ObtenerCodigo(){
		return 1
	} else if v1.ObtenerCodigo() > v2.ObtenerCodigo(){
		return -1
	}
	return 0
}


// SistemaVuelos expone un TDA para:
// - guardar vuelos a partir de archivos CSV,
// - buscarlos por código,
// - listarlos ordenados por fecha,
// - y desencolar por prioridad (filtrando duplicados)

type SistemaVuelos struct{
	porCodigo	Hash.Diccionario[string,vuelo.Vuelo]		
	porFecha	Abb.DiccionarioOrdenado[FechaClave, vuelo.Vuelo]
	porPrioridad	Heap.ColaPrioridad[vuelo.Vuelo]
}

// CrearSistemaDeVuelos instancia el TDA vacío.
func CrearSistemaDeVuelos() Aeropuerto{

	return &SistemaVuelos{
		porCodigo: Hash.CrearHash[string,vuelo.Vuelo](),
		porFecha: Abb.CrearABB[FechaClave, vuelo.Vuelo](comparadorFechaClave),
		porPrioridad: Heap.CrearHeap[vuelo.Vuelo](comparadorPrioridad),
	}
}

// Agregar_Archivo lee el CSV completo y va insertando cada línea en el TDA.
// Si un vuelo ya existía (mismo código), se reemplaza la entrada anterior.
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
	
	
