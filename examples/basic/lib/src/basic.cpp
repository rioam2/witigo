#include <cstdlib>
#include <cstring>

#include "bindings/basic_example.c"
#include "bindings/basic_example.h"

void basic_example_string_func(basic_example_string_t* input,
                               basic_example_string_t* ret) {
  ret->ptr = (uint8_t*)realloc(ret->ptr, input->len);
  memcpy(ret->ptr, input->ptr, input->len);
}

void basic_example_record_func(basic_example_customer_t* input,
                               basic_example_customer_t* ret) {
  basic_example_string_func(&input->name, &ret->name);
}

void basic_example_tuple_func(basic_example_tuple2_string_u32_t* input,
                              basic_example_tuple2_string_u32_t* ret) {
  basic_example_string_func(&input->f0, &ret->f0);
  ret->f1 = input->f1;
}

void basic_example_list_func(basic_example_list_u64_t* input,
                             basic_example_list_u64_t* ret) {
  ret->len = input->len;
  ret->ptr = (uint64_t*)realloc(ret->ptr, input->len * sizeof(uint64_t));
  memcpy(ret->ptr, input->ptr, input->len * sizeof(uint64_t));
}

bool basic_example_option_func(uint64_t* maybe_input, uint64_t* ret) {
  if (maybe_input) {
    *ret = *maybe_input;
    return true;
  } else {
    return false;
  }
}

bool basic_example_result_func(basic_example_result_u64_string_t* input,
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
