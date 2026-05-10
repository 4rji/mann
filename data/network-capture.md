# Network Capturing Tools

Packet capture and pcap analysis. Use only on networks and systems you own or are authorized to inspect.

## tcpdump live capture
`tcpdump -i INTERFACE -nn`

Capture live traffic on `INTERFACE` without resolving hostnames or service names. The `-nn` flag keeps output fast and literal.

## tcpdump write full pcap
`tcpdump -i INTERFACE -nn -s 0 -w captura.cap`

Capture full packets from `INTERFACE` and save them to `captura.cap` for later analysis.

## tcpdump capture UDP traffic
`tcpdump -i INTERFACE -nn udp`

Show only UDP packets, useful for DNS, DHCP, VoIP, and other connectionless traffic.

## tcpdump capture a specific port
`tcpdump -i INTERFACE -nn port 443`

Capture traffic where either the source or destination port is `443`.

## tcpdump capture UDP on port 53
`tcpdump -i INTERFACE -nn udp port 53`

Capture DNS-style UDP traffic on port `53`.

## tcpdump capture one host
`tcpdump -i INTERFACE -nn host IP`

Capture packets where `IP` is either the source or destination host.

## tcpdump capture source host
`tcpdump -i INTERFACE -nn src host IP`

Capture only packets sent by `IP`.

## tcpdump capture destination host
`tcpdump -i INTERFACE -nn dst host IP`

Capture only packets going to `IP`.

## tcpdump capture limited packets
`tcpdump -i INTERFACE -nn -c 100 -w captura.cap`

Stop after 100 packets and save the capture to `captura.cap`.

## tcpdump read pcap
`tcpdump -r captura.cap -nn`

Read a saved capture file without resolving hostnames or service names.

## tcpdump show ASCII payload
`tcpdump -i INTERFACE -nn -A port 80`

Print packet payload as ASCII for traffic on port `80`. Useful for plain-text protocols and lab debugging.

## tcpdump capture TCP SYN packets
`tcpdump -i INTERFACE -nn 'tcp[tcpflags] & tcp-syn != 0'`

Show TCP packets with the SYN flag set, useful for spotting new connection attempts.

## tshark read
`tshark -r captura.cap 2>/dev/null`

Replay a pcap with `tshark` (CLI Wireshark) and silence its info messages.

## Examples
```
tcpdump -i wlan0 -nn
tcpdump -i wlan0 -nn -s 0 -w captura.cap
tcpdump -i wlan0 -nn udp
tcpdump -i wlan0 -nn port 443
tcpdump -i wlan0 -nn udp port 53
tcpdump -i wlan0 -nn host 192.168.1.10
tcpdump -i wlan0 -nn src host 192.168.1.10
tcpdump -i wlan0 -nn dst host 192.168.1.10
tcpdump -i wlan0 -nn -c 100 -w captura.cap
tcpdump -r captura.cap -nn
tcpdump -i wlan0 -nn -A port 80
tcpdump -i wlan0 -nn 'tcp[tcpflags] & tcp-syn != 0'
tshark -r captura.cap 2>/dev/null
```
