include(FetchContent)

# Fetch the WASI toolchain from Github
set(FETCHCONTENT_FULLY_DISCONNECTED_OLD ${FETCHCONTENT_FULLY_DISCONNECTED})
set(FETCHCONTENT_FULLY_DISCONNECTED OFF)
FetchContent_Declare(
  wasi_sdk_toolchain
  SOURCE_DIR "${CMAKE_BINARY_DIR}/_deps/wasi-sdk"
  GIT_REPOSITORY https://github.com/rioam2/wasi-sdk-toolchain.git
  GIT_TAG 0e33447291d854dfafd4f27797294c905f5f1f93
)
FetchContent_MakeAvailable(wasi_sdk_toolchain)
set(FETCHCONTENT_FULLY_DISCONNECTED ${FETCHCONTENT_FULLY_DISCONNECTED_OLD})

# Source toolchain file(s)
include("${wasi_sdk_toolchain_SOURCE_DIR}/wasi-sdk.toolchain.cmake")

# Initialize a specific version of the WASI toolchain
initialize_wasi_toolchain(
  WIT_BINDGEN_TAG "v0.29.0"
  WASMTIME_TAG "v23.0.1"
  WASM_TOOLS_TAG "v1.215.0"
  WASI_SDK_TAG "wasi-sdk-25"
  TARGET_TRIPLET "wasm32-wasip1"
)
