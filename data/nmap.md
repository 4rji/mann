# Nmap Commands

Network discovery and port scanning. Replace `IP` with a target host or `IP/24` with a CIDR range.

## List hostnames
`nmap -sL (hostnames)`

List target hostnames without sending probes — passive enumeration based on DNS lookups.

## Vuln scripts
`nmap --script vuln IP`

Run the full `vuln` NSE category — broad vulnerability scan.

## Malware scripts
`nmap --script malware IP`

Run `malware` NSE scripts to look for known indicators of compromise.

## Source port 53
`nmap --source-port 53 IP/24`

Spoof the source port to 53 (DNS) — sometimes bypasses naive firewall rules.

## Decoys
`nmap -D RND:10 IP/24`

Send the scan alongside 10 randomly-generated decoy source IPs to obscure the real attacker.
