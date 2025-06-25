package main
import(
	vuelo "tp2/vuelo"
)
// helper: desregistra un vuelo de todos los índices secundarios
func (sv *SistemaVuelos) eliminarIndices(v vuelo.Vuelo) {
    // porCódigo
    sv.porCodigo.Borrar(v.ObtenerCodigo())

    // porPrioridad
    borrarDeIndiceLista(sv.porPrioridad, v.ObtenerPrioridad(), func(codigo string) bool {
        return codigo == v.ObtenerCodigo()
    })

    // porConexion (no es lista, se maneja directo)
    connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
    if sv.porConexion.Pertenece(connKey) {
        abbConn := sv.porConexion.Obtener(connKey)
        abbConn.Borrar(FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()})
        if abbConn.Cantidad() == 0 {
            sv.porConexion.Borrar(connKey)
        }
    }
}

