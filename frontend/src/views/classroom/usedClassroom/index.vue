<!-- 考场管理-征集考场查看 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :pagination="false" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" plain @click="onExport">导出</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button v-if="scope.row.state === RecruitState.UNCONFIRM" type="primary" link @click="onConfirmClick(scope.row)">
          确认征集
        </el-button>
        <el-button v-else-if="scope.row.state === RecruitState.CONFIRM" type="danger" link @click="onResetClick(scope.row)">
          重置
        </el-button>
        <span v-else>--</span>
      </template>
    </ProTable>

    <!-- 确认征集弹框 -->
    <el-dialog v-model="confirmDialogVisible" title="征集确认" width="420px" :close-on-click-modal="false" destroy-on-close>
      <div style="display: flex; align-items: center; gap: 8px">
        <span>监考老师数量：</span>
        <el-input-number v-model="confirmTchNum" :min="1" :max="99" controls-position="right" />
      </div>
      <template #footer>
        <el-button @click="confirmDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="confirmLoading" @click="handleConfirmSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="tsx" name="usedClassroom">
import { Classroom } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { useLevelEnum } from "@/hooks/useEnum";
import { useDepartmentEnum, useAllRecruitTaskEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import { getClassroomListByTaskId, confirmRecruit, resetRecruit, exportTaskRoomInfo } from "@/api/modules/classroom";
import { useDownload } from "@/hooks/useDownload";

// 征集状态枚举
enum RecruitState {
  UNCONFIRM = "1", // 已上报（待确认）
  CONFIRM = "2" // 已征集（已确认）
}

const recruitStateEnum = [
  { label: "已上报", value: RecruitState.UNCONFIRM },
  { label: "已征集", value: RecruitState.CONFIRM }
];

const { levelEnum } = useLevelEnum();
const { departmentEnum } = useDepartmentEnum();
const { allRecruitTaskEnum: taskEnum } = useAllRecruitTaskEnum();
const proTable = ref<ProTableInstance>();

// 选择考试任务
const handleChangeTaskId = () => {
  proTable.value?.search();
};

// 表格配置项
const columns = reactive<ColumnProps<Classroom.ResClassroomListByTaskId>[]>([
  {
    prop: "roomName",
    label: "教室名称",
    search: { order: 2, el: "input", props: { placeholder: "请输入" } }
  },
  { prop: "campusName", label: "所属校区" },
  {
    prop: "levelCode",
    label: "层级",
    enum: levelEnum,
    width: 100
  },
  {
    prop: "xx",
    label: "组数/按组容量",
    render: scope => (
      <div>
        {scope.row.groupNum || "--"}/{scope.row.groupCapacity || "--"}
      </div>
    ),
    width: 130
  },
  { prop: "capacity", label: "按位容量", width: 90 },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    search: { order: 3, el: "tree-select", props: { filterable: true, placeholder: "请选择" } }
  },
  {
    prop: "tchNum",
    label: "监考老师数量",
    width: 120,
    render: scope => (scope.row.tchNum && Number(scope.row.tchNum) > 0 ? scope.row.tchNum : "--")
  },
  {
    prop: "state",
    label: "状态",
    enum: recruitStateEnum,
    width: 90,
    search: { order: 4, el: "select", props: { placeholder: "请选择", clearable: true } }
  },
  {
    prop: "taskId",
    label: "考试任务",
    enum: taskEnum,
    search: { el: "tree-select", order: 1, props: { filterable: true, placeholder: "请选择", onChange: handleChangeTaskId } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 120 }
]);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  // 不再强制过滤 state，由用户通过搜索条件控制；不传 state 时后端返回全部
  if (!newParams.taskId) {
    ElMessage.warning("请选择考试任务");
    return Promise.resolve({ list: [] });
  }
  return getClassroomListByTaskId(newParams);
};

// ===================== 确认征集 =====================
const confirmDialogVisible = ref(false);
const confirmLoading = ref(false);
const confirmTchNum = ref(2);
const confirmRow = ref<Classroom.ResClassroomListByTaskId | null>(null);

const onConfirmClick = (row: Classroom.ResClassroomListByTaskId) => {
  confirmRow.value = row;
  confirmTchNum.value = row.tchNum ? Number(row.tchNum) : 2;
  confirmDialogVisible.value = true;
};

const handleConfirmSubmit = async () => {
  if (!confirmRow.value) return;
  confirmLoading.value = true;
  try {
    await confirmRecruit({
      id: confirmRow.value.id,
      tchNum: String(confirmTchNum.value)
    });
    ElMessage.success("征集确认成功");
    confirmDialogVisible.value = false;
    proTable.value?.getTableList();
  } finally {
    confirmLoading.value = false;
  }
};

// ===================== 重置 =====================
const onResetClick = async (row: Classroom.ResClassroomListByTaskId) => {
  try {
    await ElMessageBox.confirm(`确认重置【${row.roomName}】的征集状态？`, "温馨提示", { type: "warning" });
    await resetRecruit({ id: row.id });
    ElMessage.success("重置成功");
    proTable.value?.getTableList();
  } catch {
    // 用户取消
  }
};

// ===================== 导出 =====================
const onExport = async () => {
  const taskId = proTable.value?.searchParam?.taskId;
  if (!taskId) {
    ElMessage.warning("请先选择考试任务");
    return;
  }
  try {
    await ElMessageBox.confirm("确认导出当前任务的全部考场数据？", "导出提示", { type: "info" });
    await useDownload(exportTaskRoomInfo, `征集考场数据_${taskId}`, { taskId });
  } catch {
    // 用户取消或错误
  }
};
</script>
