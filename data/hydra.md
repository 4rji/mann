# Hydra Usage Example

Authorized credential brute-force testing — only use against systems you own or have written permission to test.

## SSH brute force (authorized)
`hydra -L users -P passwords ssh://$ip`

Try every `users × passwords` combination against the SSH service at `$ip`. Provide `users` and `passwords` as plain-text wordlists in the cwd.
