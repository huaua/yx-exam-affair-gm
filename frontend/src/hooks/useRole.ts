import { computed } from "vue";
import { useUserStore } from "@/stores/modules/user";

export const useRole = () => {
  const userStore = useUserStore();
  const isDepartmentRole = computed(() => Number(userStore.userInfo.deptId) > 0);
  const isAdminRole = computed(() => !isDepartmentRole.value);
  return {
    isAdminRole,
    isDepartmentRole
  };
};
