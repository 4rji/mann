<!-- es: Procesos y puertos -->
# Process and Port Investigation

Find which processes own which ports / working directories.

<!-- list: lsof by port | lsof por puerto -->
<!-- es: Buscar procesos por puerto con lsof -->
## lsof by port
`lsof -i:puerto`

List processes attached to the given TCP/UDP port.

<!-- list: pwdx working directory by PID | pwdx directorio de trabajo por PID -->
<!-- es: Consultar el directorio de trabajo por PID con pwdx -->
## pwdx by PID
`pwdx PID`

Print the current working directory of a running process by PID.
