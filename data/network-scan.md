# Network Scanning

Layer-2 host discovery on the local segment. Replace `INTERFACE` with your network interface and `IP` ranges with real values.

## arp-scan range
`sudo arp-scan -I INTERFACE IP_START-IP_END`

ARP-sweep an explicit IP range on the chosen interface.

## netdiscover
`sudo netdiscover -i INTERFACE`

Passive/active ARP discovery — quick way to enumerate hosts on a network you just joined.

## Examples
```
sudo arp-scan -I wlan0 10.0.4.1-10.0.4.254
sudo netdiscover -i wlan0
```
