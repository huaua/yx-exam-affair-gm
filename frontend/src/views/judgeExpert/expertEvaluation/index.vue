<!-- 基础管理-专家评价管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getExpertDatabaseList">
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" :icon="View" link @click="onViewClick(scope.row)"> 查看 </el-button>
      </template>
    </ProTable>
    <ViewDialog ref="viewDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="expertEvaluation">
import { JudgeExpert } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { View } from "@element-plus/icons-vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import ViewDialog from "./components/ViewDialog.vue";
import { expertType } from "@/utils/dict";
import { getExpertDatabaseList } from "@/api/modules/judgeExpert";

const { departmentEnum } = useDepartmentEnum("expertDataBase");
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<JudgeExpert.ResExpertDatabaseList>[]>([
  { prop: "name", label: "姓名", width: "130", search: { order: 1, el: "input", props: { placeholder: "请输入" } } },
  { prop: "sex", label: "性别", width: "70" },
  { prop: "phoneNo", label: "手机号", search: { order: 2, el: "input", props: { placeholder: "请输入" } }, width: "160" },
  {
    prop: "type",
    label: "类型",
    enum: expertType,
    search: { order: 4, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: "90"
  },
  { prop: "goodSubjects", label: "擅长科目" },
  // 后端和产品确认，该字段不要
  // {
  //   prop: "taskTimeRemark",
  //   label: "征集日期",
  //   render: scope => <div>{formatTime(scope.row.taskTimeRemark) || "--"}</div>
  // },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    search: { order: 3, el: "tree-select", props: { filterable: true, placeholder: "请选择" } }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 打开-查看Dialog
const viewDialogRef = ref<InstanceType<typeof ViewDialog> | null>(null);
const onViewClick = (row: any = {}) => {
  const params = {
    title: "查看",
    row: { ...row }
  };
  viewDialogRef.value?.acceptParams(params);
};
</script>
