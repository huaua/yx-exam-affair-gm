<!-- 基础管理-教室管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getClassroomList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openDialog('add')">新增教室</el-button>
        <el-button type="primary" :icon="Upload" plain @click="batchAdd">导入教室</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button
          v-if="!isAdminRole || scope.row.deptId === DeptId.ADMIN"
          type="primary"
          link
          :icon="Edit"
          @click="openDialog('edit', scope.row)"
        >
          编辑
        </el-button>
      </template>
    </ProTable>
    <AddOrEditDialog ref="dialogRef" />
    <ImportExcel ref="importDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="classroom">
import { Classroom } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Upload, Edit } from "@element-plus/icons-vue";
import { useRole } from "@/hooks/useRole";
import { useLevelEnum } from "@/hooks/useEnum";
import { useHandleData } from "@/hooks/useHandleData";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import { onOrOffStatus } from "@/utils/dict";
import { getClassroomList, changeClassroomStatus, addClassroom, editClassroom, importClassroom } from "@/api/modules/classroom";

enum DeptId {
  ADMIN = "-100"
}

enum ClassrommStatusEnum {
  ON = "1", // 正常
  OFF = "2" // 停用
}

type MapString = {
  [key: string]: string;
};

type MapFunction = {
  [key: string]: (params: any) => Promise<any>;
};

const { levelEnum } = useLevelEnum();
const { departmentEnum } = useDepartmentEnum();
const proTable = ref<ProTableInstance>();
const { isAdminRole } = useRole();

// 表格配置项
const columns = reactive<ColumnProps<Classroom.ResClassroomList>[]>([
  { prop: "roomName", label: "教室名称", search: { order: 1, el: "input", props: { placeholder: "请输入" } } },
  { prop: "campusName", label: "所属校区" },
  {
    prop: "levelCode",
    label: "层级",
    enum: levelEnum,
    search: { order: 2, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: 110
  },
  {
    prop: "xx",
    label: "组数/按组容量",
    render: scope => (
      <div>
        {scope.row.groupNum || "--"}/{scope.row.groupCapacity || "--"}
      </div>
    ),
    width: 140
  },
  { prop: "capacity", label: "按位容量", width: 90 },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    search: { order: 3, el: "tree-select", props: { filterable: true, placeholder: "请选择" }, isShow: isAdminRole },
    isShow: false,
    isSetting: false
  },
  {
    prop: "deptName",
    label: "所属学院"
  },
  {
    prop: "state",
    label: "状态",
    enum: onOrOffStatus,
    search: { order: 4, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    render: scope => {
      return (
        <>
          {!isAdminRole.value || scope.row.deptId === DeptId.ADMIN ? (
            <el-switch
              model-value={scope.row.state}
              active-text={scope.row.state === "1" ? "启用" : "停用"}
              active-value={"1"}
              inactive-value={"2"}
              onClick={() => changeStatus(scope.row)}
            />
          ) : (
            <el-tag type={scope.row.state === "1" ? "success" : "danger"}>{scope.row.state === "1" ? "启用" : "禁用"}</el-tag>
          )}
        </>
      );
    },
    width: 120
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 切换状态
const changeStatus = async (row: Classroom.ResClassroomList) => {
  await useHandleData(
    changeClassroomStatus,
    {
      roomId: row.roomId,
      state: row.state === ClassrommStatusEnum.ON ? ClassrommStatusEnum.OFF : ClassrommStatusEnum.ON
    },
    `${row.state === ClassrommStatusEnum.ON ? "停用" : "启用"}【${row.roomName}】教室`
  );
  proTable.value?.getTableList();
};

// 导入
const importDialogRef = ref<InstanceType<typeof ImportExcel> | null>(null);
const batchAdd = () => {
  const params = {
    title: "导入教室",
    templateUrl: "static/教室导入模板.xlsx",
    saveUploadApi: importClassroom,
    getTableList: proTable.value?.getTableList
  };
  importDialogRef.value?.acceptParams(params);
};

// 打开 Dialog(新增、编辑)
const dialogRef = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openDialog = (type: string, row: Partial<Classroom.ResClassroomList> = {}) => {
  const params = {
    type,
    title: getDialogTitle(type),
    row: { ...row },
    api: getDialogApi(type),
    getTableList: proTable.value?.getTableList
  };
  dialogRef.value?.acceptParams(params);
};

const getDialogTitle = (type: string) => {
  const dialogTitleMap: MapString = {
    add: "新增教室",
    edit: "修改教室"
  };
  return dialogTitleMap[type];
};

const getDialogApi = (type: string) => {
  const dialogApiMap: MapFunction = {
    add: addClassroom,
    edit: editClassroom
  };
  return dialogApiMap[type];
};
</script>
