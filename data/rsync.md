# Rsync Usage

Pull a remote home directory to the local cwd. Replace `USER` and `IP` with the real username and host.

## Pull home (default ssh)
`rsync -avz --progress a USER@IP:~/`

Default SSH on port 22.

## Pull home (ssh port 2244)
`rsync -avz --progress -e "ssh -p 2244" a USER@IP:~/`

Custom SSH port.

## Pull home (ssh identity)
`rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a USER@IP:~/`

Use a specific SSH key.

## Examples
```
rsync -avz --progress a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -p 2244" a root@192.168.1.20:~/
rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a root@192.168.1.20:~/
```
