<template>
  <div class="table-box draw-manage">
    <div class="filter-bar">
      <el-select v-model="taskId" filterable clearable placeholder="任务名称-请选择" style="width: 320px" @change="loadRows">
        <el-option v-for="item in tasks" :key="item.judgeTaskId" :label="item.taskName" :value="item.judgeTaskId" />
      </el-select>
      <el-button type="primary" @click="loadRows">查询</el-button>
    </div>
    <div class="toolbar">
      <el-button type="warning" :disabled="!taskId" @click="exportTask">导出</el-button>
    </div>
    <el-empty v-if="!taskId" description="请选择任务后查看抽取数据" />
    <el-table v-else :data="rows" border>
      <el-table-column prop="drawTimes" label="评分轮次" align="center" width="100" />
      <el-table-column label="科目名称" align="center"
        ><template #default="s">{{ s.row.subjectName || "--" }}</template></el-table-column
      >
      <el-table-column label="评分日期" align="center" width="200"
        ><template #default>{{ scoreDate }}</template></el-table-column
      >
      <el-table-column prop="requiredNum" label="需抽取数" align="center" width="100" />
      <el-table-column :label="isAdminRole ? '评分专家数' : '已抽取数'" align="center" width="100"
        ><template #default="s">{{ s.row.submitState === "2" ? s.row.scoringCount : "--" }}</template></el-table-column
      >
      <el-table-column label="备用专家数" align="center" width="110"
        ><template #default="s">{{ s.row.submitState === "2" ? s.row.backupCount : "--" }}</template></el-table-column
      >
      <el-table-column label="抽取状态" align="center" width="100"
        ><template #default="s">{{ s.row.submitState === "2" ? "已抽取" : "未抽取" }}</template></el-table-column
      >
      <el-table-column label="操作" align="center" width="230"
        ><template #default="s">
          <el-button link type="primary" @click="openEdit(s.row)">编辑</el-button>
          <el-button link type="primary" @click="openDraw(s.row)">{{
            s.row.submitState === "2" ? "追加抽取" : "抽取"
          }}</el-button>
          <el-button link type="primary" @click="openView(s.row)">查看</el-button>
        </template></el-table-column
      >
    </el-table>
    <el-dialog v-model="editVisible" title="编辑" width="430px" destroy-on-close :close-on-click-modal="false">
      <el-form label-width="100px">
        <el-form-item label="任务名称"><el-input :model-value="currentTask?.taskName" disabled /></el-form-item>
        <el-form-item label="评分轮次"><el-input :model-value="editRow.drawTimes" disabled /></el-form-item>
        <el-form-item label="科目名称" required
          ><el-input v-model="editRow.subjectName" maxlength="50" placeholder="请输入科目名称"
        /></el-form-item>
        <el-form-item label="抽取专家数" required
          ><el-input-number v-model="editRow.requiredNum" :min="1" :max="999" controls-position="right"
        /></el-form-item>
      </el-form>
      <template #footer
        ><el-button @click="editVisible = false">取消</el-button
        ><el-button type="primary" @click="saveRule">确定</el-button></template
      >
    </el-dialog>
    <DrawDialog ref="drawRef" @changed="loadRows" /><ExpertViewDialog ref="viewRef" />
  </div>
</template>
<script setup lang="ts" name="judgeDrawManage">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { getDrawManageTasks, getDrawManageRows, updateDrawManageRule, exportDrawManageExperts } from "@/api/modules/judgeExpert";
import { useDownload } from "@/hooks/useDownload";
import DrawDialog from "./components/DrawDialog.vue";
import ExpertViewDialog from "./components/ExpertViewDialog.vue";
import { useRole } from "@/hooks/useRole";
const { isAdminRole } = useRole();
const tasks = ref<any[]>([]),
  taskId = ref(""),
  rows = ref<any[]>([]),
  editVisible = ref(false),
  editRow = ref<any>({}),
  drawRef = ref(),
  viewRef = ref();
const currentTask = computed(() => tasks.value.find(x => x.judgeTaskId === taskId.value));
const scoreDate = computed(() =>
  currentTask.value
    ? `${String(currentTask.value.taskStartTime).slice(0, 10)}至${String(currentTask.value.taskEndTime).slice(0, 10)}`
    : "--"
);
onMounted(async () => {
  const r: any = await getDrawManageTasks();
  tasks.value = r.list || [];
});
const loadRows = async () => {
  if (!taskId.value) {
    rows.value = [];
    return;
  }
  const r: any = await getDrawManageRows({ judgeTaskId: taskId.value });
  rows.value = r.list || [];
};
const openEdit = (r: any) => {
  editRow.value = { ...r };
  editVisible.value = true;
};
const saveRule = async () => {
  if (!editRow.value.subjectName?.trim()) {
    ElMessage.warning("请输入科目名称");
    return;
  }
  await updateDrawManageRule(editRow.value);
  ElMessage.success("修改成功");
  editVisible.value = false;
  loadRows();
};
const openDraw = (r: any) => drawRef.value.acceptParams({ row: { ...r }, task: currentTask.value });
const openView = (r: any) => viewRef.value.acceptParams({ row: { ...r } });
const exportTask = () =>
  useDownload(exportDrawManageExperts, currentTask.value?.taskName || "评委抽取", { judgeTaskId: taskId.value });
</script>
<style scoped lang="scss">
.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
.toolbar {
  margin-bottom: 12px;
}
.draw-manage :deep(.el-empty) {
  min-height: 420px;
}
.draw-manage :deep(.el-input-number) {
  width: 100%;
}
</style>
