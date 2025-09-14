#include <cstdlib>
#include <cstring>
#include <string>
#include <vector>

#include "bindings/all_types_example.c"
#include "bindings/all_types_example.h"

void exports_all_types_example_string_func(all_types_example_string_t* input,
                                           all_types_example_string_t* ret) {
  std::string transformed_string((char*)input->ptr, input->len);
  transformed_string += " - modified by C++";
  ret->len = transformed_string.size();
  ret->ptr = (uint8_t*)realloc(ret->ptr, ret->len);
  memcpy(ret->ptr, transformed_string.data(), transformed_string.size());
}

void exports_all_types_example_record_func(all_types_example_customer_t* input,
                                           all_types_example_customer_t* ret) {
  exports_all_types_example_string_func(&input->name, &ret->name);
}

void exports_all_types_example_nested_record_func(
    all_types_example_nested_t* input,
    all_types_example_nested_t* ret) {
  exports_all_types_example_record_func(&input->customer, &ret->customer);
}

void exports_all_types_example_simple_record_func(
    all_types_example_simple_record_t* input,
    all_types_example_simple_record_t* ret) {
  ret->id = input->id + 1;  // Example transformation
}

void exports_all_types_example_big_record_func(
    all_types_example_big_record_t* input,
    all_types_example_big_record_t* ret) {}

void exports_all_types_example_tuple_func(
    all_types_example_tuple2_string_u32_t* input,
    all_types_example_tuple2_string_u32_t* ret) {
  exports_all_types_example_string_func(&input->f0, &ret->f0);
  ret->f1 = input->f1;
}

void exports_all_types_example_list_func(all_types_example_list_u64_t* input,
                                         all_types_example_list_u64_t* ret) {
  std::vector<uint64_t> transformed_list(input->ptr, input->ptr + input->len);
  transformed_list.push_back(99);
  ret->len = transformed_list.size();
  ret->ptr = (uint64_t*)realloc(ret->ptr, ret->len * sizeof(ret->ptr[0]));
  memcpy(ret->ptr, transformed_list.data(), ret->len * sizeof(ret->ptr[0]));
}

bool exports_all_types_example_option_func(uint64_t* maybe_input,
                                           uint64_t* ret) {
  if (maybe_input) {
    *ret = *maybe_input;
    return true;
  } else {
    return false;
  }
}

bool exports_all_types_example_result_func(
    all_types_example_result_u64_string_t* input,
    uint64_t* ret,
    all_types_example_string_t* err) {
  if (input->is_err) {
    *err = input->val.err;
    return false;
  } else {
    *ret = input->val.ok;
    return true;
  }
}

void exports_all_types_example_variant_func(
    all_types_example_allowed_destinations_t* input,
    all_types_example_allowed_destinations_t* ret) {
  if (input->tag == ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_NONE) {
    ret->tag = ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_NONE;
  } else if (input->tag == ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_ANY) {
    ret->tag = ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_ANY;
  } else if (input->tag == ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_RESTRICTED) {
    ret->tag = ALL_TYPES_EXAMPLE_ALLOWED_DESTINATIONS_RESTRICTED;
    ret->val.restricted.len = input->val.restricted.len;
    ret->val.restricted.ptr = (all_types_example_string_t*)malloc(
        ret->val.restricted.len * sizeof(ret->val.restricted.ptr[0]));
    for (size_t i = 0; i < input->val.restricted.len; i++) {
      exports_all_types_example_string_func(&input->val.restricted.ptr[i],
                                            &ret->val.restricted.ptr[i]);
    }
  }
}

all_types_example_color_t exports_all_types_example_enum_func(
    all_types_example_color_t input) {
  switch (input) {
    case ALL_TYPES_EXAMPLE_COLOR_HOT_PINK:
      return ALL_TYPES_EXAMPLE_COLOR_HOT_PINK;
    case ALL_TYPES_EXAMPLE_COLOR_LIME_GREEN:
      return ALL_TYPES_EXAMPLE_COLOR_LIME_GREEN;
    case ALL_TYPES_EXAMPLE_COLOR_NAVY_BLUE:
      return ALL_TYPES_EXAMPLE_COLOR_NAVY_BLUE;
    default:
      // Handle unexpected values gracefully
      return ALL_TYPES_EXAMPLE_COLOR_HOT_PINK;  // Default case
  }
}

int64_t exports_all_types_example_int64_func(int64_t input) {
  // Example transformation: simply return the input incremented by 1
  return input + 1;
}

void exports_all_types_example_no_return_func(bool) {
  // This function intentionally does nothing and has no return value.
  // It can be used to demonstrate a function that performs an action
  // without returning any data.
}

// complex-variant-func logic:
//  - empty      -> stays empty
//  - number(n)  -> number(n+1)
//  - floating(f)   -> floating(f*2)
//  - big(u)     -> big(u-1)
//  - text(s)    -> text(s + "!")
//  - bytes(b)   -> bytes(b concatenated with 0xff)
//  - pair(r)    -> pair({ x: r.x-1, y: r.y+1 })
void exports_all_types_example_complex_variant_func(
    all_types_example_complex_union_t* input,
    all_types_example_complex_union_t* ret) {
  switch (input->tag) {
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_EMPTY: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_EMPTY;
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_NUMBER: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_NUMBER;
      ret->val.number = input->val.number + 1;
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_FLOATING: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_FLOATING;
      ret->val.floating = input->val.floating * 2.0f;
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_BIG: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_BIG;
      ret->val.big = input->val.big - 1;
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_TEXT: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_TEXT;
      // Reallocate with space for '!'
      size_t new_len = input->val.text.len + 1;
      ret->val.text.len = new_len;
      ret->val.text.ptr = (uint8_t*)malloc(new_len);
      memcpy(ret->val.text.ptr, input->val.text.ptr, input->val.text.len);
      ret->val.text.ptr[new_len - 1] = '!';
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_BYTES: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_BYTES;
      size_t new_len = input->val.bytes.len + 1;
      ret->val.bytes.len = new_len;
      ret->val.bytes.ptr = (uint8_t*)malloc(new_len);
      memcpy(ret->val.bytes.ptr, input->val.bytes.ptr, input->val.bytes.len);
      ret->val.bytes.ptr[new_len - 1] = 0xff;
      break;
    }
    case ALL_TYPES_EXAMPLE_COMPLEX_UNION_PAIR: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_PAIR;
      ret->val.pair.x = input->val.pair.x - 1;
      ret->val.pair.y = input->val.pair.y + 1;
      break;
    }
    default: {
      ret->tag = ALL_TYPES_EXAMPLE_COMPLEX_UNION_EMPTY;
      break;
    }
  }
}