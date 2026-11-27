# Cheat catalog

Each Markdown manual in this directory supplies one section of the picker. Write section and entry headings in English and add Spanish translations above them. The picker and the static list start in Spanish. In the picker, press **Shift+E** for English or **Shift+S** to return to Spanish.

## Bilingual entry

~~~~markdown
<!-- es: Interfaces de red -->
# Network interfaces

<!-- list: Interface details | Datos de interfaz -->
<!-- es: Mostrar información de una interfaz -->
## Show interface details

`ip link show dev INTERFACE`

## Examples

```
ip link show dev eth0
```
~~~~

- Use one `#` heading for the section and a `##` heading for each command description.
- Put each Spanish translation in a single-line `<!-- es: Texto -->` comment immediately above its heading. Spanish mode displays custom headings as written when they have no translation.
- For a concise action caption, put `<!-- list: English | Español -->` immediately above the entry’s `<!-- es: ... -->` comment, or directly above its `##` heading if that comment is absent. The picker searches the displayed caption and command.
- Put the command on the first nonblank line below its description. Wrap a single line in backticks or use a triple-backtick fence for a multiline command. Add explanatory notes after the command.
- Use a single command for both languages. The picker copies the command exactly as the manual defines it.

## Preview examples

The preview below the full-width list shows the section, full entry description, and complete command, including multiline commands. Use the uppercase placeholders `IP` and `INTERFACE` where appropriate. At the end of the file, add an optional `## Examples` heading and a fenced block with one example per line. Keep this heading in English. To attach an example to an entry, match the command text and replace the placeholders with concrete values. In the example above, the preview shows `eth0` and the copied command contains `INTERFACE`.

## Add or update a manual

1. Add or edit a file under `data/` with a lowercase, hyphenated name such as `network-interfaces.md`. The picker skips `README.md`.
2. From the repository root, run `go build -o mann .` to embed the edited catalog.
3. Run `./mann`, search for the entry, and press `Shift+E` to check the English heading, then `Shift+S` to check its Spanish translation.

`mann -u` downloads published manuals into the user's catalog directory. A downloaded manual overrides the embedded file with the same name. For older downloads without translations, the picker can recover a bilingual heading when both the command and its English or Spanish heading match the embedded entry. To check a local build against its embedded manuals, point `MANN_DATA_DIR` at an empty directory. See the [update instructions](../README.md#updating-the-catalog--u).
