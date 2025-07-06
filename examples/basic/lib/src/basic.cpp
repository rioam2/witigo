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
