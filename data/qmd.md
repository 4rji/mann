# QMD Commands

Quick Markdown Search — indexa carpetas de `.md` y busca sobre ellas (BM25 + vectorial). Reemplaza `PATH`, `NAME` y `QUERY` según corresponda.

## Agregar colección
`qmd collection add PATH`

Indexa una carpeta nueva como colección buscable.

## Listar colecciones
`qmd collection list`

Muestra todas las colecciones indexadas.

## Mostrar colección
`qmd collection show NAME`

Detalle de una colección: ruta, cantidad de docs, estado.

## Renombrar colección
`qmd collection rename NAME NEW_NAME`

Cambia el nombre de una colección existente.

## Quitar colección
`qmd collection remove NAME`

Elimina una colección del índice (no borra los archivos en disco).

## Actualizar índice
`qmd update`

Reindexa las colecciones (usa `--pull` para hacer `git pull` antes).

## Ver estado del índice
`qmd status`

Salud del índice y de cada colección.

## Listar archivos indexados
`qmd ls NAME`

Inspecciona qué archivos quedaron indexados dentro de una colección.

## Buscar (híbrido)
`qmd query "QUERY"`

Búsqueda híbrida con expansión y rerank — la recomendada para uso normal.

## Buscar por palabras clave
`qmd search "QUERY"`

Búsqueda full-text BM25, sin LLM, más rápida.

## Ver un documento
`qmd get PATH`

Muestra el documento con números de línea.

## Generar embeddings
`qmd embed`

Genera/refresca los embeddings vectoriales (usa `-f` para forzar, `-c NAME` para limitar a una colección).

## Levantar servidor MCP
`qmd mcp --http --daemon`

Arranca el servidor QMD vía HTTP en background (modo daemon) para que agentes/IDEs se conecten.

## Examples
```
qmd collection add ./docs
qmd collection list
qmd collection show docs
qmd collection rename docs documentation
qmd collection remove documentation
qmd update --pull
qmd status
qmd ls docs
qmd query "como funciona el auth"
qmd search "rate limit"
qmd get docs/auth.md
qmd embed -f -c docs
qmd mcp --http --daemon
```
