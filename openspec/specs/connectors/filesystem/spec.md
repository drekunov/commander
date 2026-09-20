# filesystem Specification

## Purpose

Describes the local filesystem directory listing the connector supplies to the panel, so the panel can show each entry's name, size, and modification date and time while still telling directories from files.

## Requirements

### Requirement: Listing columns are Name, Size, Date, Time
The filesystem connector SHALL return a header row declaring the columns Name, Size, Date, and Time in that order, and one row per entry with that entry's values for those columns. The Name column SHALL be marked flexible so the panel sizes it to fill the remaining width; Size, Date, and Time SHALL keep their fixed widths. The Size value SHALL be a sortable size that displays a human-readable size, so the panel can order entries by size numerically. The connector SHALL also return a hidden IsDir attribute so the panel can detect directories without showing a directory column.

#### Scenario: Header declares the four columns
- **WHEN** the connector reads a directory containing entries
- **THEN** the header row declares Name, Size, Date, and Time in that order

#### Scenario: File entries
- **WHEN** the connector reads a regular file
- **THEN** its Name is the file name, its Size displays a human-readable size, and its Date and Time are the file's modification date (YYYY-MM-DD) and time (HH:MM:SS)

#### Scenario: Directory entries
- **WHEN** the connector reads a directory
- **THEN** its Size cell is `<DIR>` and its Date and Time are the directory's modification date and time

#### Scenario: Directory flag is hidden
- **WHEN** the connector returns an entry
- **THEN** it carries an IsDir attribute marked hidden, so the panel can detect directories without rendering a directory column

#### Scenario: Name column is flexible
- **WHEN** the header row is inspected
- **THEN** the Name attribute is marked flexible and the Size, Date, and Time attributes are not

#### Scenario: Size value is sortable
- **WHEN** two file entries' Size values are compared
- **THEN** they compare by byte size rather than by their displayed text

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
