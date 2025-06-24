package main
import(
	vuelo "tp2/vuelo"
)
// helper: desregistra un vuelo de todos los índices secundarios
func (sv *SistemaVuelos) eliminarIndices(v vuelo.Vuelo) {
    // 1) porCódigo
    sv.porCodigo.Borrar(v.ObtenerCodigo())

    // 2) porPrioridad (agrupado)
    prio := v.ObtenerPrioridad()
    if sv.porPrioridad.Pertenece(prio) {
        lst := sv.porPrioridad.Obtener(prio)
        it := lst.Iterador()
        for it.HaySiguiente() {
            if it.VerActual() == v.ObtenerCodigo() {
                it.Borrar()
                break
            }
            it.Siguiente()
        }
        if lst.EstaVacia() {
            sv.porPrioridad.Borrar(prio)
        } else {
            sv.porPrioridad.Guardar(prio, lst)
        }
    }

    // 3) porConexion
    connKey := v.ObtenerOrigen() + "-" + v.ObtenerDestino()
    if sv.porConexion.Pertenece(connKey) {
        abbConn := sv.porConexion.Obtener(connKey)
        abbConn.Borrar(FechaClave{Fecha: v.ObtenerFecha(), Codigo: v.ObtenerCodigo()})
        if abbConn.Cantidad() == 0 {
            sv.porConexion.Borrar(connKey)
        }
    }
}
