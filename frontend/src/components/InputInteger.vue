<!-- 正整数输入框 -->
<template>
  <el-input v-model="_value" @input="onInput" maxlength="9" clearable />
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue?: number | string;
    min?: number;
    max?: number;
  }>(),
  {
    min: 1
  }
);
const emit = defineEmits(["update:modelValue"]);

const _value = ref(props.modelValue);

// props.modelValue被重新赋值时
watch(
  () => props.modelValue,
  newVal => {
    _value.value = newVal;
  }
);

const onInput = (val: string | undefined) => {
  const result = verifyValue(val); // 拦截输入,进行校验
  _value.value = result;
  emit("update:modelValue", result);
};

const verifyValue = (value: string | undefined): number | string | undefined => {
  const { max, min } = props;
  let result = value;
  let newVal = Number(value);
  if (isNil(value) || Number.isNaN(newVal) || hasDot(value)) {
    return props.modelValue; // 保持输入前的数值
  }
  if (value === "") {
    return undefined;
  }
  if (max && newVal > max) {
    result = String(max);
  }
  if (min && newVal < min) {
    result = String(min);
  }
  return result;
};

// 判断null或undefined
const isNil = (value: any) => {
  return value == null;
};

// 判断是否包含.字符
const hasDot = (value: any) => {
  return String(value).includes(".");
};
</script>
