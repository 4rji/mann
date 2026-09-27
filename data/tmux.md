# Tmux

## List sessions
`tmux ls`

## Create the airsend session
`tmux new -s airsend`

## Attach to the airsend session
`tmux attach -t airsend`

## Join pane from window 1 into window 0
`tmux join-pane -s :1 -t :0`

## Join pane from window 2 into window 0
`tmux join-pane -s :2 -t :0`

## Join pane from window 3 into window 0
`tmux join-pane -s :3 -t :0`

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
