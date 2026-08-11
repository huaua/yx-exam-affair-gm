<!-- 考场编排-考场监考分配 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :search-col="{ xs: 4, sm: 5, md: 5, lg: 5, xl: 6 }" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button v-show="isShowBatchAssignBtn" type="primary" plain @click="openAssignDialog('batch')">监考分配</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" link @click="openAssignDialog('single', scope.row)"> 分配 </el-button>
        <el-button v-if="scope.row.assignmentState === AssignState.ASSIGNED" type="danger" link @click="onCancelClick(scope.row)">
          取消分配
        </el-button>
      </template>
    </ProTable>
    <AssignDialog ref="assignDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="assign">
import { Arrange } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive, computed } from "vue";
import { ElMessage } from "element-plus";
import { useHandleData } from "@/hooks/useHandleData";
import { useAllRecruitTaskEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import AssignDialog from "./components/AssignDialog.vue";
import { getTeacherAssignList, cacelAssign } from "@/api/modules/arrange";

enum AssignState {
  UNASSIGN = "1",
  ASSIGNED = "2"
}

const assignStateEnum = [
  {
    label: "未分配",
    value: AssignState.UNASSIGN
  },
  {
    label: "已分配",
    value: AssignState.ASSIGNED
  }
];

const { allRecruitTaskEnum } = useAllRecruitTaskEnum();
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Arrange.ResTeacherAssignList>[]>([
  { prop: "roomName", label: "教室名称", search: { order: 3, el: "input", props: { placeholder: "请输入" } } },
  {
    prop: "profDirection",
    label: "学院-专业-方向",
    render: scope => (
      <>{(scope.row.profDirection && scope.row.profDirection.split("\n").map(item => <div>{item}</div>)) || "--"}</>
    )
  },
  {
    prop: "taskDay",
    label: "日期",
    width: 130,
    render: scope => (scope.row.taskDay && scope.row.taskDay.slice(0, 10)) || "--"
  },
  { prop: "tchNum", label: "需要人数", width: 100 },
  {
    prop: "arrangeTchName",
    label: "监考老师",
    showOverflowTooltip: false,
    width: 150,
    render: scope => scope?.row?.arrangeTchName?.join("、") || "--"
  },
  { prop: "arrangeNoStr", label: "考场编号", width: 150 },
  {
    prop: "assignmentState",
    label: "分配状态",
    enum: assignStateEnum,
    search: { order: 2, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: "100"
  },
  {
    prop: "roomTaskId",
    label: "考场征集任务",
    enum: allRecruitTaskEnum,
    search: { el: "tree-select", order: 1, props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "deptName",
    label: "学院",
    search: { order: 4, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "profName",
    label: "专业",
    search: { order: 5, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "directionName",
    label: "方向",
    search: { order: 6, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 150 }
]);

const hasSearched = ref(false);

const isShowBatchAssignBtn = computed(() => hasSearched.value && !!proTable.value?.searchParam?.roomTaskId);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.roomTaskId) {
    ElMessage.warning("请选择考场征集任务");
    return Promise.resolve({ list: [] });
  }
  hasSearched.value = true;
  return getTeacherAssignList(newParams);
};

// 分配
const assignDialogRef = ref<InstanceType<typeof AssignDialog> | null>(null);
const openAssignDialog = (mode = "batch", row: Partial<Arrange.ResTeacherAssignList> = {}) => {
  const roomTaskId = proTable.value?.searchParam?.roomTaskId || "";
  const taskObj = {
    roomTaskId,
    roomTaskName: getTaskNameById(roomTaskId),
    taskTime: getTaskTimeById(roomTaskId)
  };
  const params = {
    title: "分配",
    mode,
    taskObj,
    row: { ...row },
    getTableList: proTable.value?.getTableList
  };
  assignDialogRef.value?.acceptParams(params);
};

// 取消分配
const onCancelClick = async (row: Arrange.ResTeacherAssignList) => {
  await useHandleData(cacelAssign, { roomReportId: row.roomReportId }, `取消分配【${row.profDirection}】`);
  proTable.value?.getTableList();
};

const getTaskNameById = (id: string) => {
  const findedItem = allRecruitTaskEnum.value.find(item => item.value === id);
  return findedItem?.label || "";
};

const getTaskTimeById = (id: string) => {
  const findedItem = allRecruitTaskEnum.value.find(item => item.value === id);
  return (findedItem?.date && `${findedItem?.date} 00:00:00`) || "";
};
</script>
