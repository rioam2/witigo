<img align="right" src="./docs/witigo-logo.svg" alt="witigo-logo" width="175"/>

# [witigo](https://github.com/rioam2/witigo) &middot; [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rioam2/witigo/blob/main/LICENSE) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/rioam2/witigo?tab=readme-ov-file#contributing) ![GitHub branch check runs](https://img.shields.io/github/check-runs/rioam2/witigo/main) [![codecov](https://codecov.io/github/rioam2/witigo/graph/badge.svg?token=9GE8H9XXGF)](https://codecov.io/github/rioam2/witigo)



Command line tool and library for generating host bindings in golang for WebAssembly (WASM) Components. This is useful for building applications that consume WebAssembly modules that leverage high-level component types defined in WIT (WebAssembly Interface Types).

---

## Getting Started

To generate host bindings for a WebAssembly Component, build the project locally or download a prebuilt binary from the releases page. To build the project locally, ensure you have go1.24 or higher installed, as well as [Taskfile](https://taskfile.dev/):

```sh
# Build the project. Output binary in `./bin/witigo`
task build
```

Then, you can generate a new golang package from a WebAssembly Component by running the following command:

```sh
./bin/witigo generate <path_to_wasm_component> <output_directory>
```

The generated package will include all necessary bindings and types to interact with the WebAssembly Component. By default, Wazero is used to provide a WebAssembly Runtime. The output structure will look something like:

```txt
<output_directory>/
├── example_component_core.wasm
└── example_component.go
```

---

## Features and Roadmap

- [x] Required utilities for generating host bindings
  - [x] `AlignTo(ptr, alignment)` - Aligns a pointer to the specified alignment.
  - [x] `AlignmentOf(type)` - Returns the alignment of a given type.
  - [x] `SizeOf(type)` - Returns the size of an element of a given type.
- [ ] Lowering (writing) and lifting (reading) of interface types
  - [x] `Read(type)` - Lifts a type to its host representation.
    - [x] `s8`, `s16`, `s32`, `s64`
    - [x] `u8`, `u16`, `u32`, `u64`
    - [x] `f32`, `f64`
    - [x] `bool`
    - [x] `string`
    - [x] `list`
    - [x] `record`
    - [x] `option`
    - [ ] `variant`
    - [ ] `result`
    - [ ] `tuple`
    - [ ] `flags`
    - [ ] `enum`
  - [x] `Write(type)` - Lowers a type to its WebAssembly representation.
    - [x] `s8`, `s16`, `s32`, `s64`
    - [x] `u8`, `u16`, `u32`, `u64`
    - [x] `f32`, `f64`
    - [x] `bool`
    - [x] `string`
    - [x] `list`
    - [x] `record`
    - [x] `option`
    - [ ] `variant`
    - [ ] `result`
    - [ ] `tuple`
    - [ ] `flags`
    - [ ] `enum`
- [ ] Host binding code generation
  - [x] Generate type definitions for interface types
  - [x] Generate exported function bindings
  - [ ] Generate imported function bindings
  - [ ] Allow configuration of Wazero runtime on instantiation
- [ ] Devops
  - [ ] Github Workflows actions to run tests
  - [ ] Dockerfile for building and running the tool
  - [ ] Automated releases using Release Please
  - [ ] Publishing prebuilt binaries on Releases

---

### Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

---

### License

This project is licensed under the MIT License. See the LICENSE file for more details.
