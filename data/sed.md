# Sed Commands

Stream editor for filtering and transforming text. Process files line by line without loading the entire file into memory.

## Replace first match on each line
`sed 's/old/new/' file.txt`

Replace the first occurrence of `old` with `new` on every line.

## Replace every match on each line
`sed 's/old/new/g' file.txt`

Use the `g` flag to replace all matches on each line, not just the first one.

## Edit a file in place
`sed -i 's/old/new/g' file.txt`

Modify the file directly instead of printing the changed output to the terminal.

## Edit in place with backup
`sed -i.bak 's/old/new/g' file.txt`

Modify the file directly and keep the original content in a `.bak` backup file.

## Delete matching lines
`sed '/pattern/d' file.txt`

Remove every line that matches `pattern`.

## Keep only matching lines
`sed '/pattern/!d' file.txt`

Delete every line that does not match `pattern`, leaving only matching lines.

## Print a line range
`sed -n '10,20p' file.txt`

Print only lines 10 through 20. The `-n` flag prevents sed from printing every line by default.

## Delete one line
`sed '5d' file.txt`

Print the file without line 5.

## Delete a line range
`sed '5,10d' file.txt`

Print the file without lines 5 through 10.

## Insert text before a line
`sed '1i\text' file.txt`

Insert `text` before line 1.

## Append text after a line
`sed '1a\text' file.txt`

Append `text` after line 1.

## Print between two patterns
`sed -n '/start/,/end/p' file.txt`

Print the section from the first `start` match through the next `end` match.

## Replace ignoring case
`sed 's/old/new/gi' file.txt`

Replace all case variations of `old`, such as `old`, `Old`, or `OLD`.

## Replace only inside a section
`sed '/start/,/end/s/old/new/g' file.txt`

Replace text only between the lines that match `start` and `end`.

## Examples
```
sed 's/example.com/api.example.com/g' nginx.conf
sed -i.bak 's/^DEBUG=false/DEBUG=true/' .env
sed '/^#/d' app.conf
sed -n '120,150p' server.log
sed '5d' users.csv
sed '1i\# Managed by automation' app.conf
sed -n '/BEGIN CERTIFICATE/,/END CERTIFICATE/p' bundle.pem
sed 's/error/WARNING/gi' app.log
sed '/server {/,/}/s/listen 80/listen 8080/g' nginx.conf
```
