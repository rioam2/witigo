cmake_minimum_required(VERSION 3.15)

# Project information
project(all_types_example VERSION 1.0.0 LANGUAGES CXX)

# Specify the C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Add source files and their targets
add_subdirectory(src)
