<!-- es: Uso de rsync -->
# Rsync Usage

Copy a local file or directory named `a` to the remote home directory. Replace `USER` and `IP` with the real username and host.

<!-- es: Copiar al directorio personal remoto por SSH -->
## Copy to remote home over SSH
`rsync -avz --progress a USER@IP:~/`

Default SSH on port 22.

<!-- es: Copiar al directorio personal remoto por el puerto SSH 2244 -->
## Copy to remote home over SSH port 2244
`rsync -avz --progress -e "ssh -p 2244" a USER@IP:~/`

Custom SSH port.

<!-- es: Copiar al directorio personal remoto con una clave SSH -->
## Copy to remote home with an SSH key
`rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a USER@IP:~/`

Use a specific SSH key.

## Examples
```
rsync -avz --progress a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -p 2244" a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a root@192.168.1.20:~/
```
