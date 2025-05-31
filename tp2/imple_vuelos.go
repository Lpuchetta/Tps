package tp2
import(
	"time"
	"fmt"
)

type vuelo struct{
	codigo        string
	aerolinea     string
	origen        string
	destino       string
	matricula     string
	prioridad     int
	fecha         time.Time
	retraso       int
	tiempoVuelo   int
	cancelado     bool
}

func CrearVuelo(codigo, aerolinea, origen, destino, matricula string, prioridad int, fecha time.Time, retraso, tiempo int, cancelado bool) Vuelos{
	return &vuelo{
		codigo: codigo,
		aerolinea: aerolinea,
		origen: origen,
		destino: destino,
		matricula: matricula,
		prioridad: prioridad,
		fecha: fecha,
		retraso: retraso,
		tiempoVuelo: tiempo,
		cancelado: cancelado,
	}
}

func (v *vuelo) ObtenerFecha() time.Time{
	return v.fecha
}

func (v *vuelo) ObtenerPrioridad() int{
	return v.prioridad
}

func (v * vuelo) EstaCancelado() bool{
	return v.cancelado
}

func (v *vuelo) MostrarInfo() string{
	return fmt.Sprintf("%s %s %s %s %s %d %s %02d %d %d",
		v.codigo,
		v.aerolinea,
		v.origen,
		v.destino,
		v.matricula,
		v.prioridad,
		v.fecha.Format("2006-01-02T15:04:05"),
		v.retraso,
		v.tiempoVuelo,
		deBoolAInt(v.cancelado),
	)
}

func (v *vuelo) ObtenerRetraso() int{
	return v.retraso
}

func (v *vuelo) ObtenerTiempoVuelo() int{
	return v.tiempoVuelo
}

func (v *vuelo) ObtenerMatricula() string{
	return v.matricula
}

func (v *vuelo) ObtenerOrigen() string{
	return v.origen
}

func (v *vuelo) ObtenerDestino() string{
	return v.destino
}

func (v *vuelo) ObtenerCodigo() string{
	return v.codigo
}

func (v *vuelo) ObtenerAerolinea() string{
	return v.aerolinea
}


func deBoolAInt(b bool) int{
	if b{
		return 1
	}
	return 0
}
