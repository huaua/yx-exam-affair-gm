<!-- 监考老师管理-监考老师上报 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTeacherSumbitList" :pagination="false">
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button v-if="scope.row?.collectNum" type="primary" link @click="openPreAppendDialog(scope.row)"> 去追加 </el-button>
        <el-button v-else type="primary" link @click="openSubmitDialog(scope.row)"> 去上报 </el-button>
      </template>
    </ProTable>
    <SubmitDialog ref="dialogRef1" />
    <PreAppendDialog ref="dialogRef2" />
    <AppendDialog ref="dialogRef3" />
  </div>
</template>

<script setup lang="tsx" name="teacherSubmit">
import { Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import ProTable from "@/components/ProTable/index.vue";
import SubmitDialog from "./components/SubmitDialog.vue";
import PreAppendDialog from "./components/PreAppendDialog.vue";
import AppendDialog from "./components/AppendDialog.vue";
import { getTeacherSumbitList } from "@/api/modules/teacher";

const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Teacher.ResTeacherRecruitManageList>[]>([
  { prop: "tchTaskName", label: "考试任务" },
  {
    prop: "xx",
    label: "征集日期",
    render: scope => (
      <div>
        {formatTime(scope.row.taskStartTime) || "--"}~{formatTime(scope.row.taskEndTime) || "--"}
      </div>
    )
  },
  { prop: "needNum", label: "要求上报总数" },
  {
    prop: "collectNum",
    label: "已上报数量"
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 打开上报
const dialogRef1 = ref<InstanceType<typeof SubmitDialog> | null>(null);
const openSubmitDialog = (row: Partial<Teacher.ResTeacherRecruitManageList> = {}) => {
  const params = {
    title: "上报",
    row: { ...row },
    getTableList: proTable.value?.getTableList
  };
  dialogRef1.value?.acceptParams(params);
};

// 打开预追加
const dialogRef2 = ref<InstanceType<typeof PreAppendDialog> | null>(null);
const openPreAppendDialog = (row: Partial<Teacher.ResTeacherRecruitManageList> = {}) => {
  const params = {
    title: "选择日期",
    row: { ...row },
    success: ({ selectedDateObj }) => {
      openAppendDialog(row, selectedDateObj);
    }
  };
  dialogRef2.value?.acceptParams(params);
};

// 打开追加
const dialogRef3 = ref<InstanceType<typeof AppendDialog> | null>(null);
const openAppendDialog = (row: Partial<Teacher.ResTeacherRecruitManageList> = {}, selectedDateObj) => {
  const params = {
    title: "追加",
    row: { ...row },
    selectedDateObj,
    getTableList: proTable.value?.getTableList
  };
  dialogRef3.value?.acceptParams(params);
};

const formatTime = fullTimeStr => {
  if (!fullTimeStr) {
    return "";
  }
  return fullTimeStr.slice(0, 10);
};
</script>
