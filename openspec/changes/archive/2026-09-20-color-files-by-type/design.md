## Context

See `proposal.md` — Why. The relevant current state and constraints:

- `internal/ui/widgets/panel/panel.go` renders listings through `github.com/charmbracelet/bubbles/table`. The table renders each cell as `styles.Cell.Render(widthStyle.Render(runewidth.Truncate(value, width, "…")))`, then wraps the cursor row with `styles.Selected` and the header with `styles.Header`. `View()` finally wraps the whole table in `m.styles.TableStyle` (the `.table` background, `#000080`).
- `runewidth.Truncate` is not ANSI-aware. A styled cell value inserted into a table row is truncated mid-escape and corrupts both the escape sequence and the column alignment. Per-cell coloring therefore cannot be done by embedding ANSI in row values.
- The panel already reads hidden attributes by name (`IsDir`) and owns the synthetic `..` parent row. Hidden attributes are excluded from columns by the data contract.
- The theme is injected as a `config.Styles` value; each CSS class is loaded into its own field in `LoadStyles`. `charmbracelet/x/ansi` is already a direct dependency.

## Goals / Non-Goals

**Goals:**
- Color only the Name cell by entry kind, from injected theme styles.
- Keep column alignment, widths, sorting, navigation, cursor, and scroll unchanged.
- Keep the focused cursor row uniform and the blurred panel fully colored.
- Add file-kind classification without changing the `app.Connector` interface.

**Non-Goals:**
- Coloring the Size, Date, or Time cells, or the header.
- A user-configurable kind-to-extension mapping; the category sets are fixed in code, only colors are themed.
- Reworking or replacing the bubbles table widget.
- Changing `IsDir` or directory navigation.

## Decisions

**D1: Classify in the filesystem connector, carry a hidden `Kind` attribute.**
- Chosen: the connector computes one kind token per entry and adds a hidden `Attribute{AttrName: "Kind", Hidden: true}`. The panel maps the token to an injected style.
- Why: executable and symlink kinds require filesystem metadata (`os.DirEntry` type/mode), which the panel does not have. Hidden attributes are the established channel for non-column data (`IsDir`).
- Alternatives considered:
  - *Derive the kind in the panel from the name extension and `IsDir`.* Cannot distinguish executables or symlinks, and puts filesystem-specific rules in the generic widget.
  - *Add fields to `app.Attribute`.* Changes the shared data model for one connector; a hidden attribute needs no interface change.

**D2: Fixed kind set and precedence.**
- Chosen tokens: `directory`, `symlink`, `executable`, `image`, `archive`, `source`, `config`. Precedence: directory (including a symlink that resolves to a directory) → symlink → executable → extension category → empty.
- Why: matches the classic commander idiom and the existing `isDirEntry` behavior, which already follows symlinks to directories for navigation. A symlink-to-directory stays browsable and reads as a directory.
- Alternatives considered:
  - *Symlink-first precedence.* Would color a browsable symlinked directory as a link, inconsistent with its navigation behavior and the `IsDir` flag.
  - *Per-extension CSS classes.* Unbounded class surface; the app's CSS parser indexes selectors by exact string, so this would need a dynamic class registry.

**D3: Extension categories are defined in code, case-insensitive.**
- Chosen: each category owns a fixed set of suffixes (for example image: `png jpg jpeg gif bmp svg webp ico tif tiff`; archive: `zip tar gz bz2 xz zst 7z rar tgz`; source: `go c h cc cpp hpp rs py js mjs ts tsx jsx java kt rb php sh lua sql`; config: `json yaml yml toml ini cfg conf env`). Matching is case-insensitive on the name suffix.
- Why: keeps the behavior deterministic and testable while the theme still controls the colors.
- Alternatives considered:
  - *Using MIME detection.* Heavier, does not classify by intent, and would read file contents.

**D4: Color the Name cell by ANSI-splicing the rendered table, not by embedding styles in cells.**
- Chosen: `View()` first renders the table, then splices a foreground-only SGR sequence around the Name column's visible span of each visible data row using `charmbracelet/x/ansi` for ANSI-aware offsets. The opening sequence carries the kind style's foreground (and bold); the closing sequence resets only those properties (for example `\x1b[39m`, `\x1b[22m`), never `\x1b[0m`. Each rendered line is matched back to its row by its plain (ANSI-stripped) text, reading only the candidate window `[clamp(cursor-height,0,cursor), clamp(cursor+height,cursor,len))` that the table builds its content from.
- Why: embedding styled strings in `table.Row` values corrupts `runewidth.Truncate`. Splicing after rendering preserves the table's ANSI structure. Resetting only the set properties keeps the panel's table background behind the rest of the line, the same lesson as the header-background fix. The Name column is the first visible column (declared flexible by the filesystem connector), so its visible span is a fixed offset — one cell of left padding — regardless of panel width. Content matching is deliberate: the bubbles table keeps an unexported viewport scroll offset, so the visible window cannot be derived from `Cursor()` and `Height()` alone (verified: at cursor 10, height 5 the first visible row is one past `cursor-height`). Matching against the candidate window recovers the mapping without relying on that hidden state; a plain-text collision only occurs when two visible rows render identically in every column, which distinct filenames make effectively impossible and which is treated as ambiguous and left uncolored.
- Alternatives considered:
  - *Embed `style.Render(name)` in the row value.* Proven to corrupt output and alignment under `runewidth.Truncate`.
  - *Render the listing with a custom row renderer instead of the bubbles table.* Removes the truncation problem but duplicates header, viewport, scroll, and cursor logic and risks the existing navigation/scroll behavior.
  - *Color by setting `styles.Cell` per row.* The table exposes a single `Cell` style for all cells, so per-row/per-cell styling is not possible.
  - *Derive the visible window from `Cursor()` and `Height()`.* Wrong when the table's hidden viewport offset is non-zero; the first visible row is off by that offset.

**D5: Suppress the kind color on the focused cursor row.**
- Chosen: when the panel is focused and the matched row is the cursor row, splice no kind color; the table's `Selected` style paints the row uniformly. When blurred, `Selected` is empty, so every row is spliced as usual.
- Why: the spec requires a uniform cursor row, and a foreground-only splice over the cursor style would recolor the highlight. Because splicing happens on each `View()`, focus and cursor movement are reflected without rebuilding rows.
- Alternatives considered:
  - *Always splice and let the cursor foreground win.* Produces a mixed highlight and weaker contrast; violates the cursor-row requirement.

**D6: Missing/unknown kind renders default.**
- Chosen: when the mapped style is the zero style (no foreground, no bold) or the kind is unknown, the panel splices nothing and the Name cell keeps the default table style. The synthetic `..` row is treated as `directory`.
- Why: implements the fallback scenarios and keeps other connectors unaffected; no kind attribute means no splicing.

**D7: CSS classes map one-to-one to `Styles` fields.**
- Chosen: add `FileDirectoryStyle`, `FileExecutableStyle`, `FileSymlinkStyle`, `FileImageStyle`, `FileArchiveStyle`, `FileSourceStyle`, `FileConfigStyle` loaded from `.file-directory`, `.file-executable`, `.file-symlink`, `.file-image`, `.file-archive`, `.file-source`, `.file-config`. The panel maps kind tokens to these fields with a switch.
- Why: follows the existing named-field convention in `config.Styles` and keeps config free of kind semantics (no `config` → `app` dependency).
- Alternatives considered:
  - *A `map[string]lipgloss.Style` field.* More extensible but diverges from the established field-per-class pattern and weakens compile-time wiring.

**D8: Embedded Norton Commander Classic colors.**
- Chosen values on the `#000080` table background, foreground only so the row background is untouched: directory `#FFFFFF` (bold); executable `#55FF55`; symlink `#55FFFF`; image `#FF55FF`; archive `#FF5555`; source `#FFFF55`; config `#AAAAAA`.
- Why: bright, mutually distinct, and legible on the classic blue panel; matches the theme's existing palette style.

## Risks / Trade-offs

- [ANSI splicing depends on the table's line layout] → The Name column's span is computed from the visible column widths, and each line is matched to its row by plain text; add panel tests that assert the SGR lands on the Name span and not elsewhere, and that the stripped layout is identical to an unstyled render.
- [Hidden viewport offset] → Visible lines are mapped by content, bounded by the candidate window the table builds from `[clamp(cursor-height,0,cursor), clamp(cursor+height,cursor,len))`; cover a scrolled listing in tests.
- [Identical-looking rows] → Two rows whose visible cells are all identical are treated as ambiguous and left uncolored rather than mis-colored; distinct filenames make this practically unreachable.
- [Bold splice interacts with the `Selected` bold] → Reset bold with `\x1b[22m` only for styles that set bold, and snapshot the opening/closing sequences from the style's actual properties.
- [Theme omits a file-kind class] → `LoadStyles` leaves the field zero, the panel splices nothing, and the entry renders default; covered by the missing-class scenario.
- [Colors assume the panel background] → Classes set foreground only in the embedded theme; a working-directory override may set a background, which the splice carries and resets after the name without clearing the inherited table background.

## Migration Plan

None — behavior-only rendering and theming change, no persisted state. Rollback is reverting the connector classification, the `Styles` fields, and the panel splice; the hidden `Kind` attribute is ignored by any consumer that does not read it.

## Open Questions

None.
