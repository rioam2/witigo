#include <cstdlib>
#include <string>
#include <vector>

#include "bindings/basic_example.c"
#include "bindings/basic_example.h"

void exports_basic_example_double(basic_example_double_operation_t* input,
                                  basic_example_double_result_t* ret) {
  // "Double" the input string
  std::string transformed_string((char*)input->double_string.ptr,
                                 input->double_string.len);
  transformed_string += transformed_string;

  // Double each entry in the input list
  std::vector<double> transformed_list(
      input->double_list.ptr, input->double_list.ptr + input->double_list.len);
  for (auto& val : transformed_list) {
    val *= 2;
  }

  // Set the return value for the list
  ret->doubled_list.len = transformed_list.size();
  const auto double_list_size = ret->doubled_list.len * sizeof(double);
  ret->doubled_list.ptr = (double*)malloc(double_list_size);
  memcpy(ret->doubled_list.ptr, transformed_list.data(), double_list_size);

  // Set the return value for the string
  ret->doubled_string.len = transformed_string.size();
  const auto string_size = ret->doubled_string.len * sizeof(uint8_t);
  ret->doubled_string.ptr = (uint8_t*)malloc(string_size);
  memcpy(ret->doubled_string.ptr, transformed_string.data(), string_size);
}
