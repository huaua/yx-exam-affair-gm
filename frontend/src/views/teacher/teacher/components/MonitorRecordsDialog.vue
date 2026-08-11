<!-- 监考老师管理-教职工管理-监考记录Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogProps.title"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="920"
  >
    <div class="monitor-info" v-if="dialogProps.row">
      <span>监考老师：{{ dialogProps.row.tchName || "--" }}</span>
      <span>工号：{{ dialogProps.row.jobNo || "--" }}</span>
      <!-- 招办角色显示所属学院（departmentEnum 仅管理员/招办有数据） -->
      <span v-if="isAdminRole">所属学院：{{ displayDeptName }}</span>
    </div>
    <el-table :data="records" v-loading="loading" border stripe empty-text="暂无监考记录">
      <el-table-column type="index" label="#" width="50" />
      <el-table-column prop="tchTaskName" label="考试任务" min-width="160" show-overflow-tooltip />
      <el-table-column prop="taskTime" label="监考日期" width="120" :formatter="(row: any) => row.taskTime?.slice(0, 10)" />
      <el-table-column prop="roomName" label="考场" min-width="160" show-overflow-tooltip />
      <el-table-column prop="arrangeNoStr" label="考场编号" width="110" />
      <el-table-column
        prop="taskStartTime"
        label="开始时间"
        min-width="160"
        :formatter="(row: any) => row.taskStartTime?.slice(0, 10)"
      />
    </el-table>
  </el-dialog>
</template>

<script setup lang="ts" name="MonitorRecordsDialog">
import { Teacher } from "@/api/interface";
import { ref, computed } from "vue";
import { getInvigilationRecords } from "@/api/modules/teacher";
import { useRole } from "@/hooks/useRole";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";

const { isAdminRole } = useRole();
// departmentEnum 仅在 isAdminRole 时有数据（useCustomEnum.ts 第29行）
const { departmentEnum } = useDepartmentEnum("onlyDepartment");

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResTeacherList> & { tchId: string };
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {} as any
});
const loading = ref(false);
const records = ref<Teacher.ResInvigilationRecord[]>([]);

// 用 departmentEnum 做 deptId→名称映射（仅招办角色可见）
const displayDeptName = computed(() => {
  const row = dialogProps.value.row;
  if (!row) return "--";
  // 后端返回了 deptName 直接用
  if (row.deptName) return row.deptName;
  const deptId = row.deptId;
  if (deptId == null || deptId === "") return "--";
  const matched = (departmentEnum.value as any[])?.find(item => String(item.value) === String(deptId));
  return matched?.label ?? String(deptId);
});

const loadRecords = async () => {
  const tchId = dialogProps.value.row?.tchId;
  if (!tchId) {
    records.value = [];
    return;
  }
  loading.value = true;
  try {
    const { list } = await getInvigilationRecords({ tchId });
    records.value = list || [];
  } catch (e) {
    records.value = [];
  } finally {
    loading.value = false;
  }
};

const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  records.value = [];
  await loadRecords();
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.monitor-info {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
  font-size: 13px;
  color: #606266;
  background: #f7f8fa;
  border-radius: 4px;
  padding: 10px 14px;
  margin-bottom: 12px;
}
</style>
