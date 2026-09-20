## ADDED Requirements

### Requirement: Entry carries a hidden file kind

The filesystem connector SHALL classify every entry into a file kind and carry that kind on the entry as a hidden attribute, alongside the existing hidden `IsDir` flag, so the panel can style an entry by kind without rendering a kind column. The kind SHALL be one of `directory`, `symlink`, `executable`, `image`, `archive`, `source`, or `config`; an entry that matches none of these SHALL carry an empty kind.

Classification SHALL apply this precedence, choosing the first that matches:

1. `directory` — the entry is a directory, including a symbolic link whose target resolves to a directory.
2. `symlink` — the entry is a symbolic link that does not resolve to a directory.
3. `executable` — the entry is a regular file with at least one execute permission bit.
4. An extension category, by the entry name's suffix compared case-insensitively: `image`, `archive`, `source`, or `config`.

#### Scenario: Directory kind

- **WHEN** the connector reads a directory entry
- **THEN** the entry carries the kind `directory` as a hidden attribute

#### Scenario: Executable kind

- **WHEN** the connector reads a regular file with an execute permission bit
- **THEN** the entry carries the kind `executable`

#### Scenario: Symlink kind

- **WHEN** the connector reads a symbolic link whose target is not a directory
- **THEN** the entry carries the kind `symlink`

#### Scenario: Extension category kind

- **WHEN** the connector reads a regular, non-executable file whose name has an extension typical of a category such as an image, archive, source, or configuration file
- **THEN** the entry carries that category's kind

#### Scenario: Unclassified entry

- **WHEN** the connector reads a regular, non-executable file that matches no category
- **THEN** the entry carries an empty kind

#### Scenario: Kind is hidden

- **WHEN** the connector returns an entry
- **THEN** its kind attribute is marked hidden, so the panel reads it without rendering a kind column
