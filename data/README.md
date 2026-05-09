# Cheat Catalog

Each `.md` file in this directory is one section of the picker. Add a new file to create a new section; delete a file to remove it. The `mann` binary embeds this directory at build time, so you must re-build after editing.

`README.md` (this file) is excluded by name and never shows up in the picker.

## File format

A minimal section file looks like this:

    # Section Title

    Optional intro paragraph for the section. Free text — ignored by the parser.

    ## Description of the command
    `actual command goes here`

    Optional notes about the command. Also ignored by the parser.

    ## Another command
    `second command`

## Rules

- The first `# Heading` becomes the **section name** shown in the picker.
- Each `## Heading` is one entry; the **first non-blank, non-heading line below it** is the command.
- Wrap the command in single backticks. Triple-backtick code fences are accepted for multi-line commands.
- Anything not on the command line is documentation and is ignored by the picker.
- File name doesn't matter for ordering — the picker sorts by section name then by description.

## Adding a new section

1. Drop a new `whatever.md` into this folder.
2. Re-build: `go build -o mann .` from the parent directory.
3. Run `./mann` and the new entries appear automatically.
