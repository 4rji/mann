# Sed

Stream editor for filtering and transforming text. Process files line by line without loading entire file into memory.

## Replace text in file
`sed 's/old/new/' file.txt`

Basic substitution. Replace first occurrence on each line. Use `g` flag for all occurrences: `sed 's/old/new/g' file.txt`

## Replace and save output
`sed -i 's/old/new/g' file.txt`

Modify file in-place. Creates backup with `-i.bak` flag: `sed -i.bak 's/old/new/g' file.txt`

## Delete lines matching pattern
`sed '/pattern/d' file.txt`

Remove every line containing pattern. Use negation: `sed '/pattern/!d'` keeps only matching lines.

## Print specific line numbers
`sed -n '5p' file.txt`

Print line 5. Range: `sed -n '10,20p' file.txt` prints lines 10-20. Use with `-e` for multiple ranges.

## Delete specific line
`sed '5d' file.txt`

Remove line 5. Range: `sed '5,10d' file.txt` deletes lines 5-10.

## Insert or append text
`sed '5i\new line' file.txt`

Insert before line 5. Use `a\` to append after: `sed '5a\new line' file.txt`

## Extract lines between patterns
`sed -n '/start/,/end/p' file.txt`

Print lines from start to end pattern inclusive. Remove `-n` to print all with matches highlighted.

## Case-insensitive substitution
`sed 's/old/new/i' file.txt`

Use `I` flag for case-insensitive. Combine with `g` for global: `sed 's/old/new/gi' file.txt`

## Substitute in specific line range
`sed '10,20s/old/new/g' file.txt`

Replace old with new only in lines 10-20. Works with patterns too: `sed '/start/,/end/s/old/new/g'`
