<!-- 考生管理-考生管理-同步Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    title="信息同步"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="400"
  >
    <div class="dialog-content">
      <div>{{ stateMessage }}</div>
      <div v-if="errMsg" class="error-msg">失败原因：{{ errMsg }}</div>
      <el-progress v-if="state !== StudentSyncEnum.NOT_STARTED" :percentage="percent" />
      <template v-if="state !== StudentSyncEnum.IN_PROGRESS">
        <el-button type="primary" :icon="Refresh" @click="reStartSync">重新同步</el-button>
      </template>
    </div>
  </el-dialog>
</template>

<script setup lang="ts" name="SyncDialog">
import { StudentSyncEnum } from "@/enums";
import { Refresh } from "@element-plus/icons-vue";
import { ref, computed } from "vue";
import { syncStudentInfo, checkStudentInfoSyncState } from "@/api/modules/student";

interface DialogProps {
  kaoShiID: string;
  kaoShiMC: string;
  state: string;
  percent: string;
  errMsg: string;
  openPreSyncDialog?: () => void;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  kaoShiID: "",
  kaoShiMC: "",
  state: StudentSyncEnum.NOT_STARTED,
  percent: "0",
  errMsg: ""
});

const state = ref(""); // 同步状态
const percent = ref(0); // 同步进度 1-100
const errMsg = ref(""); // 同步失败消息

// 状态消息
const stateMessage = computed(() => {
  switch (state.value) {
    case StudentSyncEnum.NOT_STARTED:
      return "未开始同步";
    case StudentSyncEnum.IN_PROGRESS:
      return `当前正在进行【${dialogProps.value.kaoShiMC}】考生信息同步`;
    case StudentSyncEnum.SUCCESS:
      return `【${dialogProps.value.kaoShiMC}】考生信息同步成功`;
    case StudentSyncEnum.FAILED:
      return `【${dialogProps.value.kaoShiMC}】考生信息同步失败`;
    default:
      return "未知错误";
  }
});

// 更新状态
const updateState = (res: any) => {
  state.value = res.obj?.state || StudentSyncEnum.NOT_STARTED;
  percent.value = Number(res.obj?.percent) || 0;
  errMsg.value = res.obj?.errMsg || "";
};

// 执行同步
const doSync = async () => {
  try {
    const params = {
      kaoShiID: dialogProps.value.kaoShiID,
      kaoShiMC: dialogProps.value.kaoShiMC,
      reSync: true
    };
    const res = await syncStudentInfo(params);
    updateState(res);
    if (state.value === StudentSyncEnum.IN_PROGRESS) {
      startPolling();
    } else if (state.value === StudentSyncEnum.SUCCESS) {
      dialogProps.value?.getTableList?.();
    }
  } catch (error) {
    console.log(error);
  }
};

// 轮询获取同步状态
const startPolling = () => {
  setTimeout(async () => {
    const res = await checkStudentInfoSyncState({});
    updateState(res);

    if (state.value === StudentSyncEnum.IN_PROGRESS) {
      startPolling();
    } else if (state.value === StudentSyncEnum.SUCCESS) {
      dialogProps.value?.getTableList?.();
    }
  }, 1000);
};

const reStartSync = () => {
  resetState();
  dialogVisible.value = false; // 关闭当前弹窗
  dialogProps.value?.openPreSyncDialog?.(); // 打开预同步弹窗
};

// 重置状态
const resetState = () => {
  state.value = StudentSyncEnum.NOT_STARTED;
  percent.value = 0;
  errMsg.value = "";
};

const handleClosed = () => {
  resetState();
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  state.value = dialogProps.value?.state;
  percent.value = Number(dialogProps.value?.percent) || 0;
  errMsg.value = dialogProps.value?.errMsg || "";
  if (state.value === StudentSyncEnum.NOT_STARTED) {
    doSync();
  } else if (state.value === StudentSyncEnum.IN_PROGRESS) {
    startPolling();
  }
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  align-items: center;
  padding: 10px;
  :deep(.el-progress) {
    width: 100%; // 设置进度条宽度为父容器的宽度
  }
  .error-msg {
    color: red;
  }
}
</style>
