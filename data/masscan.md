# Masscan Usage Examples

Mass scanner for IPv4 ranges. Always run with `sudo` because it crafts raw packets.

## Scan common ports
`sudo masscan -p22,80,445,8080,443 192.168.142.156/26 --rate=10000`

Scan a /26 for the listed ports at 10k pkt/s.

## Top ports 100
`sudo masscan -p1-65535 10.0.4.0/24 --top-ports 100 --rate=10000 -i wlan0`

Hit the top 100 ports across the /24 over `wlan0`.
