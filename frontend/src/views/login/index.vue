<template>
  <div class="login-container flx-center">
    <div class="login-box">
      <SwitchDark class="dark" />
      <div class="login-left">
        <img class="login-left-img" src="@/assets/images/login_left.png" alt="login" />
      </div>
      <div class="login-form">
        <div class="login-logo">
          <img class="login-icon" src="@/assets/images/logo.png" alt="" />
          <h2 class="logo-text">{{ TITLE }}</h2>
        </div>
        <LoginForm />
      </div>
      <div class="version">{{ version }}</div>
    </div>
  </div>
</template>

<script setup lang="ts" name="login">
import LoginForm from "./components/LoginForm.vue";
import SwitchDark from "@/components/SwitchDark/index.vue";
import { useVersion } from "@/hooks/useVersion";

const TITLE = import.meta.env.VITE_GLOB_APP_TITLE;
const { version } = useVersion();
</script>

<style scoped lang="scss">
@import "./index";
.version {
  position: absolute;
  right: 16px;
  bottom: 12px;
  color: var(--el-text-color-secondary);
}
</style>
