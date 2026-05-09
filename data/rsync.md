# Rsync Usage

Pull a remote home directory to the local cwd.

## Pull home (default ssh)
`rsync -avz --progress a user@host:~/`

Default SSH on port 22.

## Pull home (ssh port 2244)
`rsync -avz --progress -e "ssh -p 2244" a user@host:~/`

Custom SSH port.

## Pull home (ssh identity)
`rsync -avz --progress -e "ssh -i ~/.ssh/id_ed25519" a user@host:~/`

Use a specific SSH key.
