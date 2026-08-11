import { defineStore } from "pinia";
import { UserState } from "@/stores/interface";
import piniaPersistConfig from "@/stores/helper/persist";

export const useUserStore = defineStore({
  id: "geeker-user",
  state: (): UserState => ({
    token: "",
    roles: [],
    userInfo: { userId: "", userName: "", nickName: "", avatar: "" }
  }),
  getters: {},
  actions: {
    // Set Token
    setToken(token: string) {
      this.token = token;
    },
    // Set Token
    setRoles(roles: string[]) {
      this.roles = roles;
    },
    // Set setUserInfo
    setUserInfo(userInfo: UserState["userInfo"]) {
      this.userInfo = userInfo;
    }
  },
  persist: piniaPersistConfig("geeker-user")
});
