<!-- es: Tmux -->
# Tmux

<!-- es: Listar las sesiones -->
## List sessions
`tmux ls`

<!-- es: Crear la sesión airsend -->
## Create the airsend session
`tmux new -s airsend`

<!-- es: Conectarse a la sesión airsend -->
## Attach to the airsend session
`tmux attach -t airsend`

<!-- list: Move pane from window 1 to 0 | Mover panel de ventana 1 a 0 -->
<!-- es: Mover un panel de la ventana 1 a la ventana 0 -->
## Join pane from window 1 into window 0
`tmux join-pane -s :1 -t :0`

<!-- list: Move pane from window 2 to 0 | Mover panel de ventana 2 a 0 -->
<!-- es: Mover un panel de la ventana 2 a la ventana 0 -->
## Join pane from window 2 into window 0
`tmux join-pane -s :2 -t :0`

<!-- list: Move pane from window 3 to 0 | Mover panel de ventana 3 a 0 -->
<!-- es: Mover un panel de la ventana 3 a la ventana 0 -->
## Join pane from window 3 into window 0
`tmux join-pane -s :3 -t :0`

<!-- es: Copiar y pegar en tmux -->
## Copy and Paste Mode
```text
Ctrl-b [       -> Enter copy mode
v              -> Start selection
y              -> Copy selection
Ctrl-b ]       -> Paste
```

## Examples
```
tmux ls
tmux new -s airsend
tmux attach -t airsend
tmux join-pane -s :1 -t :0
tmux join-pane -s :2 -t :0
tmux join-pane -s :3 -t :0
```
