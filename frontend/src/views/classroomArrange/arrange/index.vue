<!-- 考场编排-考场排考编排 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="Upload" plain @click="openPreImportDialog">导入</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" link @click="openArrangeDialog(scope.row)"> 编排 </el-button>
        <el-button v-if="scope.row.arrangeState === '2'" type="danger" link @click="onCancelClick(scope.row)">
          取消编排
        </el-button>
        <el-button v-if="scope.row.arrangeState === '2'" type="primary" link @click="onViewClick(scope.row)">
          查看考场
        </el-button>
      </template>
    </ProTable>
    <ArrangeDialog ref="arrangeDialogRef" />
    <ViewArrangeDialog ref="viewArrangeDialogRef" />
    <PreImportDialog ref="preImportDialogRef" />
    <ImportExcel ref="importDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="arrange">
import { Arrange } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { Upload } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { useHandleData } from "@/hooks/useHandleData";
import { useAllRecruitTaskEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import ArrangeDialog from "./components/ArrangeDialog.vue";
import ViewArrangeDialog from "./components/ViewArrangeDialog.vue";
import PreImportDialog from "./components/PreImportDialog.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import { getClassroomArrangeList, importArrangeDeptProfDirection, cacelArrange } from "@/api/modules/arrange";

enum ArrangeState {
  NOT_ARRANGED = "1",
  ARRANGED = "2"
}

const arrangeStateEnum = [
  {
    label: "待编排",
    value: ArrangeState.NOT_ARRANGED
  },
  {
    label: "已编排",
    value: ArrangeState.ARRANGED
  }
];

const arrangeTypeEnum = [
  {
    label: "按组",
    value: "1"
  },
  {
    label: "按位",
    value: "2"
  }
];

const { allRecruitTaskEnum } = useAllRecruitTaskEnum();
const proTable = ref<ProTableInstance>();

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.roomTaskId) {
    ElMessage.warning("请选择考场征集任务");
    return Promise.resolve({ list: [] });
  }
  return getClassroomArrangeList(newParams);
};

// 表格配置项
const columns = reactive<ColumnProps<Arrange.ResClassroomArrangeList>[]>([
  {
    prop: "xxx",
    label: "学院-专业-方向",
    render: scope => {
      const deptName = String(scope.row.deptName || "");
      const profName = String(scope.row.profName || "");
      const directionName = String(scope.row.directionName || "");
      return `${deptName}-${profName}${directionName ? `-${directionName}` : ""}`;
    }
  },
  { prop: "arrangeType", label: "类型", enum: arrangeTypeEnum, width: 150 },
  { prop: "taskDay", label: "日期", width: 180, render: scope => (scope.row.taskDay && scope.row.taskDay.slice(0, 10)) || "--" },
  { prop: "stuNum", label: "考生数", width: 150 },
  { prop: "arrangeStuNum", label: "已编排考生数", width: 180 },
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
    search: { order: 2, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "profName",
    label: "专业",
    search: { order: 3, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "directionName",
    label: "方向",
    search: { order: 4, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "arrangeState",
    label: "编排状态",
    enum: arrangeStateEnum,
    search: { order: 5, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 260 }
]);

// 编排
const arrangeDialogRef = ref<InstanceType<typeof ArrangeDialog> | null>(null);
const openArrangeDialog = (row: Partial<Arrange.ResClassroomArrangeList> = {}) => {
  const params = {
    title: "编排",
    mode: row.arrangeState === "2" ? "edit" : "add",
    row: { ...row },
    getTableList: proTable.value?.getTableList
  };
  arrangeDialogRef.value?.acceptParams(params);
};

// 取消编排
const onCancelClick = async (row: Arrange.ResClassroomArrangeList) => {
  const deptName = String(row.deptName || "");
  const profName = String(row.profName || "");
  const directionName = String(row.directionName || "");
  await useHandleData(
    cacelArrange,
    { infoId: row.infoId },
    `取消编排【${deptName}-${profName}${directionName ? `-${directionName}` : ""}】`
  );
  proTable.value?.getTableList();
};

// 查看考场
const viewArrangeDialogRef = ref<InstanceType<typeof ViewArrangeDialog> | null>(null);
const onViewClick = (row: Partial<Arrange.ResClassroomArrangeList> = {}) => {
  const params = {
    title: "查看考场",
    row: { ...row },
    getTableList: proTable.value?.getTableList
  };
  viewArrangeDialogRef.value?.acceptParams(params);
};

// 点击导入
const preImportDialogRef = ref<InstanceType<typeof PreImportDialog> | null>(null);
const openPreImportDialog = () => {
  const params = {
    title: "导入",
    success: ({ selectedRoomTaskObj }) => {
      batchAdd({ selectedRoomTaskObj });
    }
  };
  preImportDialogRef.value?.acceptParams(params);
};

// 导入
const importDialogRef = ref<InstanceType<typeof ImportExcel> | null>(null);
const batchAdd = ({ selectedRoomTaskObj }) => {
  console.log("🚀 ~ batchAdd ~ roomTaskId:", selectedRoomTaskObj);
  const params = {
    title: "专业方向信息导入",
    templateUrl: "static/学院专业方向导入模板.xlsx",
    selectedRoomTaskObj,
    saveUploadApi: importArrangeDeptProfDirection,
    getTableList: proTable.value?.getTableList
  };
  importDialogRef.value?.acceptParams(params);
};
</script>
