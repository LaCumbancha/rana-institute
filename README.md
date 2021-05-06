# Institutional Site
TP1 | 75.61 - Taller de Programación III | 1C2021 | FIUBA

## Requerimientos

### Funcionales

Se solicita un sitio web institucional que cuente con un contador de visitas en las distintas páginas, individualizado para cada una de ellas e incremental. Las secciones a desarrollar serán:

* Home
* Jobs
* About
* About/Legal

### No Funcionales

Además del correcto funcionamiento del sistema, deben tenerse en cuenta las siguientes restricciones:

* El sistema de conteo de visitas debe ser programado de forma tal que se permita su reuso en otros sitios de la empresa.
* El sistema de conteo de visitas debe ser desarrollado con tecnologías Google AppEngine.
* Se requiere mostrar cómo se comporta la aplicación bajo diferentes niveles de carga. Para esto se deberá proporcionar los resultados de pruebas de carga y stress.
* El sistema debe soportar escalamiento al tráfico recibido.
* El sistema debe mostrar alta disponibilidad hacia los clientes.
* El sistema debe ser tolerante a fallos como la caída de procesos.

## Setup

Debido a las dependencias con Google Cloud Platform, el sistema no podrá correrse de manera local sin realizarle las pertinentes modificaciones al Site (y desconectándolo del Storage). Así mismo, podrán deployarse nuevas versiones una vez que estemos autenticados con los siguientes comandos del Makefile:

```
make deploy-site version=<version>
make deploy-cache version=<version>
make deploy-storage version=<version>
```

Cada uno de estos targets permitirá deployar un único servicio por separado, pero en caso de querer deployar todos juntos se puede ejecutar:

```
make deploy version=<version>
```

Para poder conectarse directamente y a través del browser con el sistema deployado, se tiene el comando

```
make run-site
```

## Testing

Los tests de carga requieren la utilización de un framework adicional (K6), una base de datos para almacenar las métricas (InfluxDB) y un software de visualización de datos (Grafana). Estos 3 pueden ser reemplazados con Docker Compose, ya que se provee la opción de correr los tests en un ambiente dockerizado que disponibilice los resultados en la ruta `localhost:3000`. Así mismo, los tests están parametrizados a través de escenarios, por lo que para probar nuevos casos sólo debe crearse uno nuevo y almacenarlo en su ruta correspondiente (/test/scripts/scenarios/scenarioX.csv), y pasárselo al target de Makefile:

```
make run-test scenario=X
```

Debido a que esto está corriendo en Docker, para detener los containers que lo estén corriendo puede ejecutarse el comando:

```
make stop-test
```

Ha de tenerse en cuenta que una vez que se detengan los contenedores se perderá toda la información almacenada en InfluxDB por no estar persistida, por lo que no se podrán volver a ver los gráficos obtenidos.
