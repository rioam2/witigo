#include <cstdlib>
#include <cstring>
#include <string>
#include <vector>

#include "bindings/basic_example.c"
#include "bindings/basic_example.h"

void exports_basic_example_string_func(basic_example_string_t* input,
                                       basic_example_string_t* ret) {
  std::string transformed_string((char*)input->ptr, input->len);
  transformed_string += " - modified by C++";
  ret->len = transformed_string.size();
  ret->ptr = (uint8_t*)realloc(ret->ptr, ret->len);
  memcpy(ret->ptr, transformed_string.data(), transformed_string.size());
}

void exports_basic_example_record_func(basic_example_customer_t* input,
                                       basic_example_customer_t* ret) {
  exports_basic_example_string_func(&input->name, &ret->name);
}

void exports_basic_example_nested_record_func(basic_example_nested_t* input,
                                              basic_example_nested_t* ret) {
  exports_basic_example_record_func(&input->customer, &ret->customer);
}

void exports_basic_example_tuple_func(basic_example_tuple2_string_u32_t* input,
                                      basic_example_tuple2_string_u32_t* ret) {
  exports_basic_example_string_func(&input->f0, &ret->f0);
  ret->f1 = input->f1;
}

void exports_basic_example_list_func(basic_example_list_u64_t* input,
                                     basic_example_list_u64_t* ret) {
  std::vector<uint64_t> transformed_list(input->ptr, input->ptr + input->len);
  transformed_list.push_back(99);
  ret->len = transformed_list.size();
  ret->ptr = (uint64_t*)realloc(ret->ptr, ret->len * sizeof(ret->ptr[0]));
  memcpy(ret->ptr, transformed_list.data(), ret->len * sizeof(ret->ptr[0]));
}

bool exports_basic_example_option_func(uint64_t* maybe_input, uint64_t* ret) {
  if (maybe_input) {
    *ret = *maybe_input;
    return true;
  } else {
    return false;
  }
}

bool exports_basic_example_result_func(basic_example_result_u64_string_t* input,
                                       uint64_t* ret,
                                       basic_example_string_t* err) {
  if (input->is_err) {
    *err = input->val.err;
    return false;
  } else {
    *ret = input->val.ok;
    return true;
  }
}

void exports_basic_example_variant_func(
    basic_example_allowed_destinations_t* input,
    basic_example_allowed_destinations_t* ret) {
  if (input->tag == BASIC_EXAMPLE_ALLOWED_DESTINATIONS_NONE) {
    ret->tag = BASIC_EXAMPLE_ALLOWED_DESTINATIONS_NONE;
  } else if (input->tag == BASIC_EXAMPLE_ALLOWED_DESTINATIONS_ANY) {
    ret->tag = BASIC_EXAMPLE_ALLOWED_DESTINATIONS_ANY;
  } else if (input->tag == BASIC_EXAMPLE_ALLOWED_DESTINATIONS_RESTRICTED) {
    ret->tag = BASIC_EXAMPLE_ALLOWED_DESTINATIONS_RESTRICTED;
    ret->val.restricted.len = input->val.restricted.len;
    ret->val.restricted.ptr = (basic_example_string_t*)realloc(
        ret->val.restricted.ptr,
        input->val.restricted.len * sizeof(basic_example_string_t));
    memcpy(ret->val.restricted.ptr, input->val.restricted.ptr,
           input->val.restricted.len * sizeof(basic_example_string_t));
  }
}

basic_example_color_t exports_basic_example_enum_func(
    basic_example_color_t input) {
  switch (input) {
    case BASIC_EXAMPLE_COLOR_HOT_PINK:
      return BASIC_EXAMPLE_COLOR_HOT_PINK;
    case BASIC_EXAMPLE_COLOR_LIME_GREEN:
      return BASIC_EXAMPLE_COLOR_LIME_GREEN;
    case BASIC_EXAMPLE_COLOR_NAVY_BLUE:
      return BASIC_EXAMPLE_COLOR_NAVY_BLUE;
    default:
      // Handle unexpected values gracefully
      return BASIC_EXAMPLE_COLOR_HOT_PINK;  // Default case
  }
}

int64_t exports_basic_example_int64_func(int64_t input) {
  // Example transformation: simply return the input incremented by 1
  return input + 1;
}