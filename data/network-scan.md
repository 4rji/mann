# Network Scanning

Layer-2 host discovery on the local segment.

## arp-scan range
`sudo arp-scan -I wlan0 10.0.4.1-10.0.4.254`

ARP-sweep an explicit IP range on the chosen interface.

## netdiscover
`sudo netdiscover -i wlan0`

Passive/active ARP discovery — quick way to enumerate hosts on a network you just joined.
