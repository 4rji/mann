<!-- es: Detección de vulnerabilidades con Nmap -->
# Nmap for Vulnerability and Safety

Nmap script combinations focused on safe vulnerability detection. Replace `IP` with the target host.

<!-- list: Vuln and safe scripts | Scripts de tipo vuln y safe -->
<!-- es: Scripts de las categorías vuln y safe -->
## Vuln and safe
`nmap -p[puerto] IP --script='vuln and safe' -sV`

Run only NSE scripts that belong to BOTH the `vuln` AND `safe` categories — non-disruptive vulnerability check with version detection.

## Examples
```
nmap -p80 192.10.0.88 --script='vuln and safe' -sV
```
