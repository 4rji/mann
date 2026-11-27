<!-- es: Archivos tar -->
# Tar Archives

Create, extract and inspect tar archives. Replace `ARCHIVE`, `DIR`, `FILE` and `DEST` with real names.

`tar` itself only bundles many files into one; the compression is a separate layer chosen with a flag. Main flags: `c` create, `x` extract, `t` list, `v` verbose (print each file), `f` archive file name (always last, right before the name), `z` gzip, `j` bzip2, `J` xz, `--zstd` zstandard.

Which compression to choose:

| Format    | Flag      | Speed      | Size      | When to use |
|-----------|-----------|------------|-----------|-------------|
| `.tar`    | (none)    | instant    | no change | Already-compressed data (jpg, mp4, zip) or as an intermediate file you will append to. |
| `.tar.gz` | `-z`      | fast       | good      | Default choice. Works on any Linux/macOS/BSD box, ideal for sharing and quick backups. |
| `.tar.bz2`| `-j`      | slow       | better    | Legacy. Smaller than gzip but slower than xz and zstd at similar sizes; use only if someone asks for it. |
| `.tar.xz` | `-J`      | very slow  | best      | Distribution/long-term storage where size matters more than time (releases, archives you upload once and download many times). |
| `.tar.zst`| `--zstd`  | very fast  | very good | Modern default when both ends have `zstd`: near-xz size at gzip speed or better. Great for big backups. |

<!-- es: Comprimir un directorio en tar.gz (rápido y compatible) -->
## Compress dir as tar.gz — fast, works everywhere (default choice)
`tar -czvf ARCHIVE.tar.gz DIR`

Gzip compression. Not the smallest file, but quick to create and every system can open it. Use it when in doubt, when sending to someone else, or for everyday backups.

<!-- es: Comprimir un directorio en tar.bz2 (formato antiguo, más lento que gzip) -->
## Compress dir as tar.bz2 — smaller than gz but slow (legacy)
`tar -cjvf ARCHIVE.tar.bz2 DIR`

Bzip2 compression. Produces ~10–15% smaller files than gzip but takes several times longer. Today xz or zstd beat it on both size and speed, so use bz2 only when a tool or person specifically requires it.

<!-- es: Comprimir un directorio en tar.xz para distribuirlo o almacenarlo a largo plazo -->
## Compress dir as tar.xz — smallest file, slowest (releases, long-term storage)
`tar -cJvf ARCHIVE.tar.xz DIR`

XZ (LZMA2) compression. Best ratio of the classic formats but slow to compress and memory hungry. Good for files you create once and download/store many times, such as software releases or archived logs. To use all CPU cores: `XZ_OPT=-T0 tar -cJvf ARCHIVE.tar.xz DIR`.

<!-- es: Comprimir un directorio en tar.zst (rápido y compacto, requiere zstd) -->
## Compress dir as tar.zst — very fast and small (modern, needs zstd)
`tar --zstd -cvf ARCHIVE.tar.zst DIR`

Zstandard compression. Close to xz size at gzip speed or faster, and decompresses extremely fast. Best for large backups and transfers when both machines have `zstd` installed (`apt install zstd`).

<!-- es: Empaquetar un directorio en un tar sin compresión -->
## Bundle dir as plain tar — no compression (media, zips, appendable)
`tar -cvf ARCHIVE.tar DIR`

Just packs files into one, without compressing. Use it when the content is already compressed (videos, images, .zip, .gz files) so compressing again only wastes CPU, or when you plan to append files later with `-r`.

<!-- es: Comprimir un directorio excluyendo archivos o carpetas (logs, node_modules, .git) -->
## Compress dir excluding files or folders (logs, node_modules, .git)
`tar -czvf ARCHIVE.tar.gz --exclude='*.log' --exclude='node_modules' DIR`

Skip paths matching a pattern. Repeat `--exclude` once per pattern. Put the `--exclude` options before `DIR`. Useful to keep caches, dependencies and logs out of a backup.

<!-- es: Crear un tar.gz protegido con contraseña mediante GPG y AES-256 -->
## Compress dir with password — encrypted tar.gz with gpg (AES256)
`tar -czf - DIR | gpg -c --cipher-algo AES256 -o ARCHIVE.tar.gz.gpg`

`tar` cannot set passwords by itself, so the archive is piped into `gpg -c` (symmetric encryption), which asks for the password twice. Nothing unencrypted is written to disk. File names and contents are both hidden. Needs `gnupg` (installed on most Linux systems). Use `-cJf -` or `--zstd -cf -` instead of `-czf -` for xz or zstd.

<!-- es: Extraer un tar.gz.gpg protegido con contraseña -->
## Extract a password-protected tar.gz.gpg
`gpg -d ARCHIVE.tar.gz.gpg | tar -xzvf -`

Asks for the password, decrypts to stdout and unpacks it in the current directory. Add `-C DEST` at the end to unpack elsewhere. To only list the contents: `gpg -d ARCHIVE.tar.gz.gpg | tar -tzf -`.

<!-- es: Comprimir y cifrar un directorio con contraseña mediante OpenSSL -->
## Compress dir with password using openssl (no gpg available)
`tar -czf - DIR | openssl enc -aes-256-cbc -pbkdf2 -salt -out ARCHIVE.tar.gz.enc`

Same idea with `openssl`, which is on almost every machine. Always keep `-pbkdf2`: without it the password is turned into a key in a weak way. Decrypt and extract with `openssl enc -d -aes-256-cbc -pbkdf2 -in ARCHIVE.tar.gz.enc | tar -xzvf -`, using the same options as when encrypting.

<!-- es: Extraer aquí un archivo tar detectando el formato automáticamente -->
## Extract archive here — any format, auto-detected
`tar -xvf ARCHIVE.tar.gz`

Unpacks into the current directory. GNU tar detects gzip, bzip2, xz and zstd automatically, so the same command works for `.tar.gz`, `.tgz`, `.tar.bz2`, `.tar.xz` and `.tar.zst`. Files with the same name are overwritten. Check the contents first with `tar -tf` if you are unsure where it will unpack.

<!-- es: Extraer un archivo tar en otro directorio -->
## Extract archive into another directory
`tar -xvf ARCHIVE.tar.gz -C DEST`

Same as above but unpacks into `DEST` instead of the current directory. `DEST` must already exist (`mkdir -p DEST`).

<!-- es: Extraer un solo archivo del tar -->
## Extract only one file from the archive
`tar -xvf ARCHIVE.tar.gz path/inside/FILE`

Pulls out a single file (or folder) without unpacking everything. The path must match exactly what `tar -tf` shows, usually without a leading `/`.

<!-- es: Extraer el contenido omitiendo la carpeta de nivel superior -->
## Extract dropping the top folder (project-1.0/ → contents directly)
`tar -xvf ARCHIVE.tar.gz --strip-components=1 -C DEST`

Many downloads wrap everything in one folder like `project-1.0/`. This removes that first level so the files land directly in `DEST`. Use `=2` to drop two levels.

<!-- es: Restaurar una copia conservando permisos y propietarios -->
## Extract keeping original permissions and owners (restore backups)
`sudo tar -xpvf ARCHIVE.tar.gz -C DEST`

`-p` preserves the file modes stored in the archive; running as root also restores the original owners. Use it when restoring system backups such as `/etc` or `/var/www`.

<!-- es: Listar los nombres de los archivos del tar sin extraerlos -->
## List archive contents — names only, no extraction
`tar -tf ARCHIVE.tar.gz`

Shows every path inside the archive without unpacking. Use it before extracting to see the folder layout and avoid unpacking hundreds of files into the current directory.

<!-- es: Listar permisos, propietarios, tamaños y fechas del contenido -->
## List archive contents with details (permissions, owner, size, date)
`tar -tvf ARCHIVE.tar.gz`

Like `ls -l` for the archive. Useful to check sizes, owners and dates before restoring.

<!-- es: Leer un archivo del tar sin extraerlo al disco -->
## Read a file inside the archive without extracting it
`tar -xOf ARCHIVE.tar.gz path/inside/FILE`

`-O` sends the file to stdout instead of writing it to disk. Pipe it to `less`, `grep` or redirect it: `tar -xOf ARCHIVE.tar.gz etc/hosts | less`.

<!-- es: Buscar un archivo dentro del tar -->
## Search for a file inside the archive
`tar -tf ARCHIVE.tar.gz | grep FILE`

Quick way to check whether a file is in the archive and get its exact path for extracting.

<!-- es: Añadir un archivo a un tar sin comprimir -->
## Add a file to an existing tar (uncompressed only)
`tar -rvf ARCHIVE.tar FILE`

Appends `FILE` to the end of the archive. Only works with plain `.tar`; compressed archives must be extracted and recreated.

<!-- es: Comparar el tar con los archivos en disco -->
## Compare archive against files on disk (what changed)
`tar -dvf ARCHIVE.tar`

Run it from the directory where the archive was created. It reports files whose size, date, mode or content differ from the copy in the archive. Handy to verify a backup.

<!-- es: Copiar un directorio por SSH sin crear un archivo temporal -->
## Copy dir to a remote host over ssh without temp file
`tar -czf - DIR | ssh USER@IP 'tar -xzf - -C DEST'`

`-f -` writes the archive to stdout, which goes through SSH and is unpacked on the other side. Keeps permissions and symlinks, and avoids filling the disk with a temporary archive. An alternative to `rsync` or `scp -r` for many small files.

<!-- es: Comprobar la integridad de un tar.gz -->
## Check that a tar.gz is not corrupted
`gzip -t ARCHIVE.tar.gz && echo OK`

Reads the whole file and verifies the checksum. Prints `OK` if fine, an error if the file is truncated or damaged. For other formats: `xz -t`, `bzip2 -t`, `zstd -t`.

## Examples
```
tar -czvf backup.tar.gz /etc/nginx
tar -cJvf release-1.0.tar.xz project/
tar --zstd -cvf home-backup.tar.zst --exclude='.cache' /home/user
tar -cvf photos.tar ~/Pictures
tar -czf - ~/Documents | gpg -c --cipher-algo AES256 -o docs.tar.gz.gpg
gpg -d docs.tar.gz.gpg | tar -xzvf - -C /tmp/restore
tar -xvf backup.tar.gz -C /tmp/restore
tar -xvf node-v20.tar.xz --strip-components=1 -C /opt/node
tar -tvf backup.tar.gz
tar -xOf backup.tar.gz etc/nginx/nginx.conf | less
tar -czf - /srv/www | ssh root@192.168.1.20 'tar -xzf - -C /srv'
```
