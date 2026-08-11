<!-- 任务征集-考场征集管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openDialog('add')">新增任务</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button
          :type="scope.row.state === PublishStatusEnum.UNPUBLISH ? 'primary' : 'danger'"
          link
          @click="changeStatus(scope.row)"
        >
          {{ scope.row.state === PublishStatusEnum.UNPUBLISH ? "确认" : "取消" }}发布
        </el-button>
        <el-button
          v-if="scope.row.state === PublishStatusEnum.PUBLISHED"
          type="primary"
          link
          @click="onConfirmRecruitClick(scope.row)"
        >
          确认征集
        </el-button>
        <el-button
          v-else-if="scope.row.collectNum == null || scope.row.collectNum === '0'"
          type="primary"
          link
          :icon="Edit"
          @click="openDialog('edit', scope.row)"
        >
          编辑
        </el-button>
      </template>
    </ProTable>
    <AddOrEditDialog ref="addOrEditDialogRef" />
    <ConfirmRecruitDialog ref="confirmRecruitDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="classroomRecruitManage">
import { Classroom } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Edit } from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { useLevelEnum } from "@/hooks/useEnum";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import ConfirmRecruitDialog from "./components/ConfirmRecruitDialog.vue";
import { recruitStatus } from "@/utils/dict";
import {
  getRecruitManageList,
  changeRecruitTaskStatus,
  addRecruitTask,
  editRecruitTask,
  confirmRecruitBatch
} from "@/api/modules/classroom";

enum PublishStatusEnum {
  UNPUBLISH = "1",
  PUBLISHED = "2"
}

type MapString = {
  [key: string]: string;
};

type MapFunction = {
  [key: string]: (params: any) => Promise<any>;
};

const { levelEnum } = useLevelEnum();

// ProTable 实例
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Classroom.ResRecruitManageList>[]>([
  { prop: "taskName", label: "考试任务", search: { order: 1, el: "input", props: { placeholder: "请输入" } } },
  {
    prop: "taskDay",
    label: "征集日期",
    search: { order: 2, el: "date-picker", props: { placeholder: "请选择" } },
    render: scope => {
      return dayjs(scope.row.taskDay).format("YYYY-MM-DD");
    }
  },
  {
    prop: "levelCode",
    label: "所属层级",
    enum: levelEnum,
    search: { order: 3, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: 120
  },
  {
    prop: "state",
    label: "征集状态",
    enum: recruitStatus,
    search: { order: 4, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: 120
  },
  { prop: "needNum", label: "需征集数", width: 120 },
  { prop: "collectNum", label: "已征集数", width: 120 },
  { prop: "comfirmNum", label: "已确认数", width: 120 },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  newParams.taskDay && (newParams.taskDay = dayjs(newParams.taskDay).format("YYYY-MM-DD"));
  return getRecruitManageList(newParams);
};

// 切换发布状态
const changeStatus = async (row: Classroom.ResRecruitManageList) => {
  await useHandleData(
    changeRecruitTaskStatus,
    {
      taskId: row.taskId,
      state: row.state === PublishStatusEnum.UNPUBLISH ? PublishStatusEnum.PUBLISHED : PublishStatusEnum.UNPUBLISH
    },
    `${row.state === PublishStatusEnum.UNPUBLISH ? "发布" : "取消发布"}【${row.taskName}】任务`
  );
  proTable.value?.getTableList();
};

// 确认征集
const confirmRecruitDialogRef = ref<InstanceType<typeof ConfirmRecruitDialog> | null>(null);
const onConfirmRecruitClick = (row: Partial<Classroom.ResRecruitManageList> = {}) => {
  const params = {
    title: "确认征集",
    row: { ...row },
    api: confirmRecruitBatch,
    getTableList: proTable.value?.getTableList
  };
  confirmRecruitDialogRef.value?.acceptParams(params);
};

// 打开 Dialog(新增、编辑)
const addOrEditDialogRef = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openDialog = (type: string, row: Partial<Classroom.ResRecruitManageList> = {}) => {
  const params = {
    type,
    title: getDialogTitle(type),
    row: { ...row },
    api: getDialogApi(type),
    getTableList: proTable.value?.getTableList
  };
  addOrEditDialogRef.value?.acceptParams(params);
};

const getDialogTitle = (type: string) => {
  const dialogTitleMap: MapString = {
    add: "新增任务",
    edit: "修改任务"
  };
  return dialogTitleMap[type];
};

const getDialogApi = (type: string) => {
  const dialogApiMap: MapFunction = {
    add: addRecruitTask,
    edit: editRecruitTask
  };
  return dialogApiMap[type];
};
</script>
