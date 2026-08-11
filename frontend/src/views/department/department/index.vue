<!-- 院系管理-院系管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getDepartmentList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openDialog('新增')">新增</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" link :icon="Edit" @click="openDialog('编辑', scope.row)">编辑</el-button>
      </template>
    </ProTable>
    <DepartmentDialog ref="dialogRef" />
  </div>
</template>

<script setup lang="tsx" name="Department">
import { Department } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Edit } from "@element-plus/icons-vue";
import ProTable from "@/components/ProTable/index.vue";
import DepartmentDialog from "./components/DepartmentDialog.vue";
import { getDepartmentList, addDepartment, editDepartment } from "@/api/modules/department";

// ProTable 实例
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Department.ResDepartmentList>[]>([
  { prop: "deptCode", label: "学院代码" },
  { prop: "deptName", label: "学院名称" },
  { prop: "campusName", label: "所在校区" },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 打开 Dialog(新增、编辑)
const dialogRef = ref<InstanceType<typeof DepartmentDialog> | null>(null);
const openDialog = (title: string, row: Partial<Department.ResDepartmentList> = {}) => {
  const params = {
    title,
    row: [{ ...row }],
    api: title === "新增" ? addDepartment : editDepartment,
    getTableList: proTable.value?.getTableList
  };
  dialogRef.value?.acceptParams(params);
};
</script>
