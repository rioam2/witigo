# Generate C++ bindings from WebAssembly Component Interface Definition
wit_bindgen(
    INTERFACE_FILE_INPUT "${CMAKE_CURRENT_SOURCE_DIR}/component.wit"
    BINDINGS_DIR_INPUT "${CMAKE_CURRENT_BINARY_DIR}/bindings"
    GENERATED_FILES_OUTPUT bindings
)

# Add a compile target for the all_types_example WebAssembly module
add_executable(all_types_example component.cpp ${bindings})
target_link_libraries(all_types_example PRIVATE wasi_sdk_reactor_module)
set_target_properties(all_types_example PROPERTIES OUTPUT_NAME "all_types_example.wasm")
target_compile_options(all_types_example PRIVATE -Wno-narrowing)


# Add WIT bindings to interface/include directory
target_include_directories(all_types_example PUBLIC ${CMAKE_CURRENT_BINARY_DIR})

# Create a WebAssembly component from the all_types_example module
wasm_create_component(
    COMPONENT_TARGET all_types_example_component
    CORE_WASM_TARGET all_types_example
    COMPONENT_TYPE "reactor"
)