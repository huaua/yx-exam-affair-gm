<template>
  <el-dialog v-model="visible" title="任务查看" width="1250px" :close-on-click-modal="false" class="requirement-dialog">
    <div v-loading="loading" class="requirement-content">
      <el-empty v-if="!loading && rows.length === 0" description="暂无学院上报要求数据" />
      <div v-for="row in rows" :key="row.requirementId" class="requirement-row">
        <div class="field field-code">
          <span class="label">方向代码：</span>
          <el-input :model-value="row.directionCode" disabled />
        </div>
        <div class="field field-name">
          <span class="label">方向名称：</span>
          <el-input :model-value="row.directionName" disabled />
        </div>
        <div class="field field-dept">
          <span class="label">上报学院：</span>
          <el-select :model-value="row.deptName" disabled>
            <el-option :label="row.deptName" :value="row.deptName" />
          </el-select>
        </div>
        <div class="field field-number">
          <span class="label">上报专家数：</span>
          <el-input-number v-model="row.expertNum" :min="1" :precision="0" :controls="false" />
        </div>
        <el-button type="primary" link class="delete-button" @click="removeRow(row)">删除</el-button>
      </div>
    </div>
    <template #footer>
      <div class="dialog-footer"><el-button type="primary" :loading="saving" @click="handleConfirm">确定</el-button></div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { getJudgeTaskRequirements, batchEditJudgeTaskRequirements, deleteJudgeTaskRequirement } from "@/api/modules/judgeExpert";

const visible = ref(false);
const loading = ref(false);
const saving = ref(false);
const taskId = ref("");
const rows = ref<JudgeExpert.JudgeTaskReportRequirement[]>([]);
// 父列表刷新回调（保存/删除后更新需求数等）
const getTableList = ref<(() => void) | null>(null);

const refresh = async () => {
  loading.value = true;
  try {
    const res = await getJudgeTaskRequirements({ judgeTaskId: taskId.value });
    rows.value = (res?.list || []).map(item => ({ ...item, expertNum: Number(item.expertNum) }));
  } finally {
    loading.value = false;
  }
};

const acceptParams = async (params: any) => {
  const row = params?.row ?? params;
  taskId.value = row.judgeTaskId;
  getTableList.value = params?.getTableList ?? null;
  visible.value = true;
  await refresh();
};

const handleConfirm = async () => {
  if (rows.value.length === 0) {
    visible.value = false;
    return;
  }
  if (rows.value.some(item => !/^[1-9]\d*$/.test(String(item.expertNum)))) {
    ElMessage.warning("上报专家数必须为大于0的正整数");
    return;
  }
  saving.value = true;
  try {
    await batchEditJudgeTaskRequirements(
      {
        items: rows.value.map(item => ({
          requirementId: item.requirementId,
          expertNum: String(item.expertNum)
        }))
      },
      { showErrorMsg: false }
    );
    ElMessage.success("保存成功");
    visible.value = false;
    getTableList.value?.();
  } catch {
    // 错误信息已在全局拦截器处理，这里不再重复弹窗
  } finally {
    saving.value = false;
  }
};

const removeRow = async (row: JudgeExpert.JudgeTaskReportRequirement) => {
  await ElMessageBox.confirm(`确认删除【${row.directionName}-${row.deptName}】的上报要求吗？`, "提示", { type: "warning" });
  await deleteJudgeTaskRequirement({ requirementId: row.requirementId });
  ElMessage.success("删除成功");
  await refresh();
  getTableList.value?.();
};

defineExpose({ acceptParams });
</script>

<style scoped lang="scss">
.requirement-content {
  min-height: 120px;
  padding: 4px 0;
  background: #f2f4f7;
}
.requirement-row {
  display: flex;
  align-items: center;
  gap: 26px;
  margin: 10px 12px;
  padding: 22px 28px;
  background: #fff;
  border-radius: 4px;
}
.field {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
}
.label {
  flex: 0 0 auto;
  margin-right: 10px;
  color: #303133;
  font-size: 16px;
  white-space: nowrap;
}
.field-code :deep(.el-input) {
  width: 95px;
}
.field-name :deep(.el-input) {
  width: 195px;
}
.field-dept :deep(.el-select) {
  width: 195px;
}
.field-number :deep(.el-input-number) {
  width: 95px;
}
.delete-button {
  margin-left: auto;
  font-size: 16px;
}
.dialog-footer {
  display: flex;
  justify-content: center;
}
.dialog-footer .el-button {
  min-width: 140px;
}
:deep(.el-input.is-disabled .el-input__wrapper) {
  background-color: #f5f7fa;
}
</style>
