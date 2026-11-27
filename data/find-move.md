<!-- es: Mover archivos -->
# Moving Files Example

`find` recipes for relocating files across nested directories.

<!-- es: Mover los archivos .py de los subdirectorios al directorio actual -->
## Move py files up
`find . -mindepth 2 -type f -name '*.py' -exec mv {} . \;`

Pull every `.py` file from sub-directories into the current directory. Watch for filename collisions — `mv` overwrites silently.
