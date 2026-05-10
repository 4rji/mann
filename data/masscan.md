# Masscan Usage Examples

Mass scanner for IPv4 ranges. Always run with `sudo` because it crafts raw packets. Replace `IP` with a target host or CIDR and `INTERFACE` with your network interface.

## Scan common ports
`sudo masscan -p22,80,445,8080,443 IP/26 --rate=10000`

Scan a /26 for the listed ports at 10k pkt/s.

## Top ports 100
`sudo masscan -p1-65535 IP/24 --top-ports 100 --rate=10000 -i INTERFACE`

Hit the top 100 ports across the /24 over the chosen interface.

## Examples
```
sudo masscan -p22,80,445,8080,443 192.168.142.156/26 --rate=10000
sudo masscan -p1-65535 10.0.4.0/24 --top-ports 100 --rate=10000 -i wlan0
```
