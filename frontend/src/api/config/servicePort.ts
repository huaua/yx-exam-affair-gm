// 后端微服务模块前缀
export const PORT1 = "/geeker";
export const PORT2 = "/hooks";
export const LANJING = "/lanjing";

const env = import.meta.env;
const VITE_PROXY_TAG = env.VITE_PROXY_TAG || "";
export const PROXY_TAG = `${VITE_PROXY_TAG}`; // 代理标签
