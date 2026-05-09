# Process and Port Investigation

Find which processes own which ports / working directories.

## lsof by port
`lsof -i:puerto`

List processes attached to the given TCP/UDP port.

## pwdx by PID
`pwdx PID`

Print the current working directory of a running process by PID.
