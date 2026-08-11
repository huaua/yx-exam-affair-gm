import { AxiosResponse } from "axios";
import { ref, computed } from "vue";

const apiVersion = ref(""); // 服务端版本

export const useVersion = () => {
  const VERSION = import.meta.env.VITE_VERSION; // 前端版本

  const version = computed(() => `${VERSION}_${apiVersion.value}`);

  // 设置api版本
  const setApiVersion = (response?: AxiosResponse) => {
    if (response && !apiVersion.value) {
      apiVersion.value = response.headers["Back-End-Version"] || response.headers["back-end-version"] || "";
    }
  };

  return {
    setApiVersion,
    version
  };
};
