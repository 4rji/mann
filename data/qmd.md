<!-- es: Comandos de QMD -->
# QMD commands

Quick Markdown Search indexes directories of Markdown files for keyword (BM25) and vector searches. Replace `PATH`, `NAME` and `QUERY` with your own values.

<!-- es: Agregar colección -->
## Add a collection
`qmd collection add PATH`

Index a directory as a searchable collection.

<!-- es: Listar colecciones -->
## List collections
`qmd collection list`

Show all indexed collections.

<!-- es: Mostrar colección -->
## Show collection details
`qmd collection show NAME`

Show a collection’s path, document count and status.

<!-- es: Renombrar colección -->
## Rename a collection
`qmd collection rename NAME NEW_NAME`

Change the name of an existing collection.

<!-- es: Quitar colección -->
## Remove a collection
`qmd collection remove NAME`

Remove the collection from the index. Its files stay on disk.

<!-- es: Actualizar índice -->
## Update the index
`qmd update`

Reindex the collections. Use `--pull` to run `git pull` first.

<!-- es: Ver estado del índice -->
## Check index status
`qmd status`

Check the health of the index and each collection.

<!-- es: Listar archivos indexados -->
## List indexed files
`qmd ls NAME`

See which files QMD has indexed in a collection.

<!-- es: Buscar (híbrido) -->
## Run a hybrid search
`qmd query "QUERY"`

Search with query expansion and reranking. Use this for general searches.

<!-- es: Buscar por palabras clave -->
## Search by keyword
`qmd search "QUERY"`

Run a fast BM25 full-text search without an LLM.

<!-- es: Ver un documento -->
## View a document
`qmd get PATH`

Display a document with line numbers.

<!-- es: Generar embeddings -->
## Generate embeddings
`qmd embed`

Generate or refresh vector embeddings. Use `-f` to force a refresh or `-c NAME` to limit it to one collection.

<!-- es: Levantar servidor MCP -->
## Start the MCP server
`qmd mcp --http --daemon`

Run the QMD server in the background so agents and IDEs can connect over HTTP.

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
