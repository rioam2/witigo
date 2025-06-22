# witigo

Command line tool and library for generating host bindings in golang for WebAssembly (WASM) Components. This is useful for building applications that consume WebAssembly modules that leverage high-level component types defined in WIT (WebAssembly Interface Types).

## Getting Started

TODO: Provide usage examples and detailed instructions on how to use the library and CLI tool.

## Features and Roadmap

- [ ] Required utilities for generating host bindings
  - [x] `align_to(ptr, alignment)` - Aligns a pointer to the specified alignment.
  - [x] `alignment(type)` - Returns the alignment of a given type.
  - [x] `elem_size(type)` - Returns the size of an element of a given type.
  - [ ] `lift(type)` - Lifts a type to its host representation.
  - [ ] `lower(type)` - Lowers a type to its WebAssembly representation.

### Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.
