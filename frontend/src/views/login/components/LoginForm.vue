<template>
  <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="large">
    <el-form-item prop="username">
      <el-input v-model="loginForm.username" placeholder="用户名">
        <template #prefix>
          <el-icon class="el-input__icon">
            <user />
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input v-model="loginForm.password" type="password" placeholder="密码" show-password autocomplete="new-password">
        <template #prefix>
          <el-icon class="el-input__icon">
            <lock />
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <div class="img-code">
      <el-form-item prop="code" class="code">
        <el-input v-model="loginForm.code" autocomplete="off" placeholder="请输入图片验证码">
          <template #prefix>
            <el-icon><Select /></el-icon>
          </template>
        </el-input>
      </el-form-item>
      <img class="img" :src="captchaUrl" @click="refreshCaptcha" alt="" />
    </div>
  </el-form>
  <div class="login-btn">
    <el-button :icon="CircleClose" round size="large" @click="resetForm(loginFormRef)"> 重置 </el-button>
    <el-button :icon="UserFilled" round size="large" type="primary" :loading="loading" @click="login(loginFormRef)">
      登录
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import { HOME_URL } from "@/config";
import { getTimeState } from "@/utils";
import { Login } from "@/api/interface";
import { ElNotification } from "element-plus";
import { loginApi, getCaptcha, getInfo } from "@/api/modules/login";
import { useUserStore } from "@/stores/modules/user";
import { useTabsStore } from "@/stores/modules/tabs";
import { useKeepAliveStore } from "@/stores/modules/keepAlive";
import { initDynamicRouter } from "@/routers/modules/dynamicRouter";
import { CircleClose, UserFilled } from "@element-plus/icons-vue";
import type { ElForm } from "element-plus";
// import md5 from "md5";

const TITLE = import.meta.env.VITE_GLOB_APP_TITLE;
const router = useRouter();
const userStore = useUserStore();
const tabsStore = useTabsStore();
const keepAliveStore = useKeepAliveStore();

type FormInstance = InstanceType<typeof ElForm>;
const loginFormRef = ref<FormInstance>();
const loginRules = reactive({
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
  code: [{ required: true, message: "请输入图片验证码", trigger: "blur" }]
});

const loading = ref(false);
const loginForm = reactive<Login.ReqLoginForm>({
  username: "",
  password: "",
  code: "",
  uuid: ""
});

// 验证码图片地址
const captchaUrl = ref<string>("");

// 获取图片验证码
const getCaptchaInfo = () => {
  getCaptcha()
    .then(res => {
      if (res) {
        captchaUrl.value = res?.captcha?.img;
        loginForm.uuid = res?.captcha?.uuid;
      }
    })
    .catch(err => {
      console.log("获取图片验证码失败err :>> ", err);
    });
};

// 刷新验证码
const refreshCaptcha = () => {
  captchaUrl.value = "";
  getCaptchaInfo();
};

// login
const login = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate(async valid => {
    if (!valid) return;
    loading.value = true;
    try {
      // 1.执行登录接口
      const { token } = await loginApi({ ...loginForm });
      userStore.setToken(token);
      const { roles, user } = await getInfo();
      userStore.setRoles(roles);
      userStore.setUserInfo(user);

      // 2.添加动态路由
      await initDynamicRouter();

      // 3.清空 tabs、keepAlive 数据
      tabsStore.setTabs([]);
      keepAliveStore.setKeepAliveName([]);

      // 4.跳转到首页
      router.push(HOME_URL);
      ElNotification({
        title: getTimeState(),
        message: `欢迎登录 ${TITLE}`,
        type: "success",
        duration: 3000
      });
    } finally {
      loading.value = false;
    }
  });
};

// resetForm
const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};

onMounted(() => {
  // 获取图片验证码地址
  getCaptchaInfo();
  // 监听 enter 事件（调用登录）
  document.onkeydown = (e: KeyboardEvent) => {
    if (e.code === "Enter" || e.code === "enter" || e.code === "NumpadEnter") {
      if (loading.value) return;
      login(loginFormRef.value);
    }
  };
});

onBeforeUnmount(() => {
  document.onkeydown = null;
});
</script>

<style scoped lang="scss">
@import "../index";
.img-code {
  display: flex;
  :deep(.el-form-item) {
    margin-bottom: 0 !important;
  }
  .code {
    flex-grow: 1;
  }
  img {
    display: block;
    flex-shrink: 0;
    margin-left: 20px;
    border: 1px solid #cccccc;
  }
}
</style>
