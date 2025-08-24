# Generate C++ bindings from WebAssembly Component Interface Definition
wit_bindgen(
    INTERFACE_FILE_INPUT "${CMAKE_CURRENT_SOURCE_DIR}/component.wit"
    BINDINGS_DIR_INPUT "${CMAKE_CURRENT_BINARY_DIR}/bindings"
    GENERATED_FILES_OUTPUT bindings
)

# Add a compile target for the basic_example WebAssembly module
add_executable(basic_example component.cpp ${bindings})
target_link_libraries(basic_example PRIVATE wasi_sdk_reactor_module)
set_target_properties(basic_example PROPERTIES OUTPUT_NAME "basic_example.wasm")

# Add WIT bindings to interface/include directory
target_include_directories(basic_example PUBLIC ${CMAKE_CURRENT_BINARY_DIR})

# Create a WebAssembly component from the basic_example module
wasm_create_component(
    COMPONENT_TARGET basic_example_component
    CORE_WASM_TARGET basic_example
    COMPONENT_TYPE "reactor"
)