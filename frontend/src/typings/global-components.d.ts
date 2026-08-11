/* eslint-disable */
/* prettier-ignore */
// @ts-nocheck
// Read more: https://github.com/vuejs/core/pull/3399
export {};

declare module "vue" {
  export interface GlobalComponents {
    InputInteger: (typeof import("../components/InputInteger.vue"))["default"];
  }
}
