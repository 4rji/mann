# SSH Connection Keep-Alive

Prevent SSH sessions from being killed by NAT timeouts or idle disconnects.

## Keepalive flags
`ssh -o ServerAliveInterval=60 -o ServerAliveCountMax=3 usuario@servidor`

Send a keep-alive every 60s; close after 3 missed responses (~3 min unresponsive).
