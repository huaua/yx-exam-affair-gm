<!-- 动态输入框组件 -->
<script lang="ts" setup>
import { Plus, Minus } from "@element-plus/icons-vue";

defineProps<{ currentIndex: number; disabled?: boolean; length: number }>();

const emits = defineEmits(["onPlus", "onMinus"]);

const handlePlus = () => {
  emits("onPlus");
};

const handleMinus = (currentIndex: number) => {
  emits("onMinus", currentIndex);
};
</script>

<template>
  <div class="dynamic-input-wrapper">
    <slot></slot>
    <div v-if="!disabled" class="button-wrapper">
      <el-button style="margin-left: 12px" type="info" :icon="Plus" circle @click="handlePlus" />
      <template v-if="length > 1">
        <el-button type="info" :icon="Minus" circle @click="handleMinus(currentIndex)" />
      </template>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dynamic-input-wrapper {
  display: flex;
  flex: 1;
  margin-bottom: 8px;
  .button-wrapper {
    display: flex;
    width: 140px;
    margin-top: 6px;
    button {
      width: 18px;
      height: 18px;
    }
  }
}
</style>
