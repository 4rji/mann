# SSH Connection Keep-Alive

Prevent SSH sessions from being killed by NAT timeouts or idle disconnects. Replace `USER` and `IP` with the real username and host.

## Keepalive flags
`ssh -o ServerAliveInterval=60 -o ServerAliveCountMax=3 USER@IP`

Send a keep-alive every 60s; close after 3 missed responses (~3 min unresponsive).

## Examples
```
ssh -o ServerAliveInterval=60 -o ServerAliveCountMax=3 root@192.168.1.20
```
