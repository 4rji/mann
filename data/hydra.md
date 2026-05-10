# Hydra Usage Example

Authorized credential brute-force testing — only use against systems you own or have written permission to test. Replace `IP` with the target host.

## SSH brute force (authorized)
`hydra -L users -P passwords ssh://IP`

Try every `users × passwords` combination against the SSH service at `IP`. Provide `users` and `passwords` as plain-text wordlists in the cwd.

## Examples
```
hydra -L users -P passwords ssh://192.168.1.50
```
