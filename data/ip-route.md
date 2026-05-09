# IP Route Deletion

Remove default routes from the kernel routing table. Useful when forcing traffic through a specific gateway or VPN.

## Delete default route
`sudo ip route del default via 192.168.88.1 dev wlp0s20f3`

Replace `192.168.88.1` with the gateway IP and `wlp0s20f3` with the interface name. Use `ip route` to list current routes.
