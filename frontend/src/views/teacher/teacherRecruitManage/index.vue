<!-- 监考老师管理-监考老师征集任务 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTeacherRecruitManageList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openPreAddDialog">新增</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button
          v-if="scope.row.state === TaskState.UNPUBLISH"
          type="primary"
          link
          :icon="Edit"
          @click="openDialog('edit', scope.row)"
        >
          编辑
        </el-button>
      </template>
    </ProTable>
    <PreAddDialog ref="dialogRef1" />
    <AddOrEditDialog ref="dialogRef2" />
  </div>
</template>

<script setup lang="tsx" name="teacherRecruitManage">
import { Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Edit } from "@element-plus/icons-vue";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import PreAddDialog from "./components/PreAddDialog.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import { recruitStatus } from "@/utils/dict";
import {
  getTeacherRecruitManageList,
  changeTeacherRecruitManageStatus,
  addTeacherRecruitTask,
  editTeacherRecruitTask
} from "@/api/modules/teacher";

enum TaskState {
  UNPUBLISH = "1", // 未发布
  PUBLISHED = "2" // 已发布
}

type MapString = {
  [key: string]: string;
};

type MapFunction = {
  [key: string]: (params: any) => Promise<any>;
};

const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Teacher.ResTeacherRecruitManageList>[]>([
  { prop: "tchTaskName", label: "考试任务", search: { order: 1, el: "input", props: { placeholder: "请输入" } } },
  {
    prop: "xx",
    label: "征集日期",
    render: scope => (
      <div>
        {formatTime(scope.row.taskStartTime) || "--"}~{formatTime(scope.row.taskEndTime) || "--"}
      </div>
    )
  },
  {
    prop: "state",
    label: "征集状态",
    enum: recruitStatus,
    search: { order: 2, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    render: scope => {
      return (
        <el-switch
          model-value={scope.row.state}
          active-value={TaskState.PUBLISHED}
          inactive-value={TaskState.UNPUBLISH}
          onClick={() => changeStatus(scope.row)}
        />
      );
    },
    width: 120
  },
  { prop: "needNum", label: "总征集数", width: 120 },
  {
    prop: "collectNum",
    label: "已征集数",
    width: 120
  },
  { prop: "auditPassNum", label: "审核通过数", width: 120 },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 切换状态
const changeStatus = async (row: Teacher.ResTeacherRecruitManageList) => {
  await useHandleData(
    changeTeacherRecruitManageStatus,
    {
      tchTaskId: row.tchTaskId,
      state: row.state === TaskState.UNPUBLISH ? TaskState.PUBLISHED : TaskState.UNPUBLISH
    },
    `${row.state === TaskState.UNPUBLISH ? "发布" : "取消发布"}【${row.tchTaskName}】`
  );
  proTable.value?.getTableList();
};

// 打开新增
const dialogRef1 = ref<InstanceType<typeof PreAddDialog> | null>(null);
const openPreAddDialog = () => {
  const params = {
    title: "新增任务",
    success: ({ tchTaskName, dates }) => {
      openDialog("add", { tchTaskName, dates });
    }
  };
  dialogRef1.value?.acceptParams(params);
};

const dialogRef2 = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openDialog = (type: string, row: any = {}) => {
  const params = {
    type,
    title: getDialogTitle(type),
    row: {
      tchTaskId: row.tchTaskId,
      tchTaskName: row.tchTaskName,
      dates: row.dates
    },
    api: getDialogApi(type),
    getTableList: proTable.value?.getTableList
  };
  dialogRef2.value?.acceptParams(params);
};

const getDialogTitle = (type: string) => {
  const dialogTitleMap: MapString = {
    add: "新增任务",
    edit: "编辑任务"
  };
  return dialogTitleMap[type];
};

const getDialogApi = (type: string) => {
  const dialogApiMap: MapFunction = {
    add: addTeacherRecruitTask,
    edit: editTeacherRecruitTask
  };
  return dialogApiMap[type];
};

const formatTime = fullTimeStr => {
  if (!fullTimeStr) {
    return "";
  }
  return fullTimeStr.slice(0, 10);
};
</script>
