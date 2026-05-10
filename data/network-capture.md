# Network Capturing Tools

Packet capture and pcap analysis.

## tcpdump capture
`tcpdump -i INTERFACE -w captura.cap -v`

Capture verbose traffic from `INTERFACE` and write a pcap file.

## tshark read
`tshark -r captura.cap 2>/dev/null`

Replay a pcap with `tshark` (CLI Wireshark) and silence its info messages.

## Examples
```
tcpdump -i wlan0 -w captura.cap -v
tshark -r captura.cap 2>/dev/null
```
