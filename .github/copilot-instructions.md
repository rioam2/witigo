# witigo – AI Coding Agent Instructions

Concise, project-specific guidance for automated changes. Keep edits minimal, type-safe, and aligned with existing patterns.

## 1. Purpose & Big Picture
`witigo` is a Go CLI + library that generates Go host bindings for WebAssembly Components (WIT-defined). Flow:
1. Input: A component `.wasm` (component model) → `wasm-tools` (embedded) extracts canonical WIT JSON + core module.
2. Parsing: `pkg/wit` lazily wraps JSON definitions (worlds, types, functions).
3. Codegen: `pkg/codegen` converts WIT types/functions → Go typedefs + function wrappers (single `.go` file output).
4. Runtime: Generated package + extracted `<name>_core.wasm` run using Wazero; ABI marshaling in `pkg/abi`.

### 1.1 Source of truth on business logic
- `docs/CanonicalABI.md` – Detailed ABI design notes, type mappings, and rationale. This is the primary reference for understanding how WIT types map to host environment runtime types (in this case Golang types) and how the ABI functions operate. Examples are in python pseudocode, but can be used to infer and understand desired logic and patterns.

## 2. Key Directories
- `cmd/main.go` – CLI dispatch (`generate`). Keep commands simple; new commands follow same pattern.
- `pkg/codegen` – Pure string/code AST generation (gowrtr). Typename mapping lives in `generate_type.go`.
- `pkg/wit` – Thin JSON façade; avoids upfront decoding. Add fields by lazy json.RawMessage extraction.
- `pkg/abi` – Canonical ABI lifting/lowering (Read*/Write* and *Parameter* helpers) for primitives + lists/records/options/enums.
- `pkg/wasmtools` – Embedded `wasm-tools.wasm` runner using wazero; provides extraction helpers.
- `examples/*` – Source of truth for expected generated shapes. Use when changing codegen.

## 3. Build & Test Workflows
- Build: `task build` (produces `./bin/witigo`).
- Unit tests: `task test` (covers `pkg/abi` + type mapping).
- Example smoke (basic): `task basic-example`.
- All-types example (broad coverage): `task all-types-example`.
Run focused tests with `go test ./pkg/abi -run TestName`.

## 4. Type & Naming Conventions
- Primitive WIT → Go: s8→int8, u32→uint32, f64→float64, char→rune, string→string, bool→bool.
- Records → `PascalCaseNameRecord` struct.
- Enums → `PascalCaseNameEnum` underlying `uint{8|16|32|64}` chosen by `discriminantSize(len(cases))` (see `generate_type.go`). Constants: `<EnumTypename><CasePascal>`.
- Results → `OkType-ErrType` → `PascalCase + Result` struct with `Ok` / `Error` fields.
- Tuples → concatenated element type names + `Tuple` (empty → `EmptyTuple`).
- Options → Go `Option[T]` generic syntax in generated source (consumer defined? treat literally—do not rename).
- Variants (planned) / Handles follow existing stubs; mimic enum + struct combo.

## 5. ABI Patterns
- All Read/Write funcs accept `AbiOptions` (carries Memory + alloc helpers). Memory alignment via `AlignTo` before access.
- `WriteX` optionally takes `ptrHint`; if nil/0 allocates via `abiMalloc` and returns a composite `AbiFreeCallback` (use + defer free in generated call sites).
- Parameter lowering: `WriteParameterX` returns flat `[]Parameter` (Value/Size/Alignment) with zero allocations where possible.
- Enums are validated via `isEnumType` then passed through int helpers.
- Floats: canonical NaN constants enforced; compare byte patterns not `math.IsNaN`.

## 6. Adding New WIT Kinds
1. Extend `pkg` root `abi_type.go` enum if missing.
2. Add typename generation to `GenerateTypenameFromType` and (if needed) struct/typedef generator.
3. Implement Read*/Write*/WriteParameter* in `pkg/abi` (mirror existing ones; keep symmetry). Reuse existing memory + alignment patterns.
4. Update codegen emission to include new type’s typedef (if non-primitive) and adjust function parameter lowering if representation differs.
5. Add tests in analogous `*_test.go` file (see `enum_test.go`, `record_test.go`). Use minimal in-memory stub implementing `Memory` interface.
6. Add an example in `examples/all-types` if visible shape change is helpful.

## 7. Code Generation Style
- Use gowrtr builder APIs (avoid manual string concatenation except small `generator.NewRawStatementf`).
- Keep single public `GenerateFromFile` entry that: extract WIT JSON → parse → generate world[0] only.
- Maintain deterministic output: avoid map iteration without ordering; rely on WIT order as delivered.

## 8. Error Handling & Validation
- Fail fast on I/O (file abs path, existence) in CLI.
- ABI read errors bubble up with contextual pointer/size info. Do not wrap with generic messages that lose original context.
- Avoid panics except for truly impossible internal states (current code panics in type mapping when kind unknown). Follow that pattern.

## 9. Performance Considerations
- Favor zero-allocation parameter lowering (write directly into `[]Parameter`).
- Enum/Int writes create byte slice sized `SizeOf(value)`; avoid extra copies.
- Keep reflection localized; do not cache global reflect.Types unless profiling shows benefit.

## 10. When Modifying ABI Logic
- Run `task test` + regenerate both examples (`task basic-example` & `task all-types-example`) and manually diff generated output for regressions.
- Ensure `go fmt` equivalence: codegen uses `.Gofmt()` already—avoid secondary formatting passes.

## 11. Safe Changes Checklist (Before PR)
- [ ] New/changed type logic has unit test(s).
- [ ] Examples still build & run (no runtime panic).
- [ ] Generated identifiers match existing casing patterns.
- [ ] No unused exported symbols introduced.

## 12. Out of Scope
- Don’t introduce a new runtime backend without explicit design note.
- Don’t eagerly fully decode WIT JSON (maintain lazy extraction pattern).

---
Questions or unclear area? Ask for clarification—keep this file terse and current.
