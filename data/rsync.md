<!-- es: Uso de rsync -->
# Rsync Usage

Copy a local file or directory named `a` to the remote home directory. Replace `USER` and `IP` with the real username and host.

<!-- es: SSH estándar → ~/ remoto -->
## Default SSH → remote ~/
`rsync -avz --progress a USER@IP:~/`

Default SSH on port 22.

<!-- es: Puerto 2244 → ~/ remoto -->
## Port 2244 → remote ~/
`rsync -avz --progress -e "ssh -p 2244" a USER@IP:~/`

Custom SSH port.

<!-- es: Clave SSH → ~/ remoto -->
## SSH key → remote ~/
`rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a USER@IP:~/`

Use a specific SSH key.

## Examples
```
rsync -avz --progress a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -p 2244" a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a root@192.168.1.20:~/
```
