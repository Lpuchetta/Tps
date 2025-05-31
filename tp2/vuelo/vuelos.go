package vuelo

import(
	"time"
)

type Vuelo interface{

	//ObtenerFecha devuelve la fecha del vuelo
	ObtenerFecha() time.Time

	//ObtenerPrioridad devuelve la prioridad del vuelo
	ObtenerPrioridad() int

	//EstaCancelado indica si un vuelo se ecuentra cancelado con un booleano.
	EstaCancelado() bool

	//MostrarInfo devuelve un string con los formatos correspondientes para despues mostrar por pantalla
	MostrarInfo() string

	//ObtenerRetraso devuelve el retraso que tiene ese vuelo
	ObtenerRetraso() int

	//ObtenerTiempoVuelo devuelve cuando tarda en viajar
	ObtenerTiempoVuelo() int

	//ObtenerMatricula devuelve la matricula del avion
	ObtenerMatricula() string

	//ObtenerOrigen devuelve el origen de salida del vuelo
	ObtenerOrigen() string

	//ObtenerDestino devuelve el destino del vuelo
	ObtenerDestino() string

	//ObtenerCodigo devuelve el codigo del vuelo
	ObtenerCodigo() string

	//ObtenerAeorolinea devuele la aerolinea correspondiente del avion.
	ObtenerAerolinea() string
}
