# Nmap for Vulnerability and Safety

Nmap script combinations focused on safe vulnerability detection. Replace `IP` with the target host.

## Vuln and safe
`nmap -p[puerto] IP --script='vuln and safe' -sV`

Run only NSE scripts that belong to BOTH the `vuln` AND `safe` categories — non-disruptive vulnerability check with version detection.

## Examples
```
nmap -p80 192.10.0.88 --script='vuln and safe' -sV
```
