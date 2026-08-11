<!-- 带编辑按钮的正整数输入框 -->
<template>
  <InputInteger v-if="_enabled" v-model="_value" v-bind="$attrs" :min="0" @blur="onBlur" width="10px" />
  <template v-else>
    <div>{{ $attrs.modelValue }}</div>
    <el-button
      type="primary"
      link
      :icon="EditPen"
      @click="
        () => {
          _enabled = true;
        }
      "
    />
  </template>
</template>

<script lang="ts" setup>
import { ref, useAttrs, watch } from "vue";
import { EditPen } from "@element-plus/icons-vue";

defineOptions({
  inheritAttrs: false
});

const props = withDefaults(
  defineProps<{
    notEmpty?: boolean;
  }>(),
  {
    notEmpty: false // 不允许为空，为空时自动转为0
  }
);

const emit = defineEmits(["save", "update:modelValue"]);

const attrs = useAttrs();
const _enabled = ref(false);
const _value: any = ref(attrs.modelValue);

// attrs.modelValue 被重新赋值时
watch(
  () => attrs.modelValue,
  newVal => {
    _value.value = newVal;
  }
);

const onBlur = () => {
  _enabled.value = false;
  if (props.notEmpty && _value.value == null) {
    _value.value = 0;
  }
  emit("update:modelValue", Number(_value.value));
  emit("save", Number(_value.value));
};
</script>
