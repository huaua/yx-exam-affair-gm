<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList" row-key="expertId">
      <template #tableHeader="{ selectedList }">
        <el-button type="primary" :icon="CirclePlus" @click="addDialogRef?.acceptParams(proTable?.getTableList)">新增</el-button>
        <el-button type="primary" plain :icon="Upload" @click="openImport">导入</el-button>
        <el-button type="primary" plain :icon="Download" @click="download">导出</el-button>
        <el-button type="danger" plain :icon="Delete" :disabled="!selectedList.length" @click="remove(selectedList)"
          >删除所选</el-button
        >
        <el-button type="danger" :icon="Delete" @click="removeAll">删除全部</el-button>
      </template>
      <template #operation="scope">
        <el-button type="danger" link @click="removeOne(scope.row)">删除</el-button>
      </template>
    </ProTable>
    <AddDialog ref="addDialogRef" />
    <ImportExcel ref="importDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="expertBan">
import { reactive, ref } from "vue";
import { CirclePlus, Delete, Download, Upload } from "@element-plus/icons-vue";
import { ElMessageBox } from "element-plus";
import ProTable from "@/components/ProTable/index.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import { ColumnProps, ProTableInstance } from "@/components/ProTable/interface";
import { JudgeExpert } from "@/api/interface";
import { useDownload } from "@/hooks/useDownload";
import { useHandleData } from "@/hooks/useHandleData";
import { educationType } from "@/utils/dict";
import { useExpertTitleEnum } from "@/hooks/useEnum";
import {
  deleteAllBannedExperts,
  deleteBannedExperts,
  exportBannedExperts,
  getBannedExpertList,
  importBannedExperts
} from "@/api/modules/judgeExpert";
import AddDialog from "./components/AddDialog.vue";

const { expertTitleEnum } = useExpertTitleEnum();
const categoryEnum = [
  { label: "校内", value: "1" },
  { label: "校外", value: "2" }
];

const proTable = ref<ProTableInstance>();
const addDialogRef = ref<InstanceType<typeof AddDialog>>();
const importDialogRef = ref<InstanceType<typeof ImportExcel>>();
let searchedParams: Record<string, any> = {};

const columns = reactive<ColumnProps<JudgeExpert.ResExpertDatabaseList>[]>([
  { type: "selection", fixed: "left", width: 50 },
  { prop: "name", label: "姓名", width: 100, search: { order: 1, el: "input" } },
  { prop: "idCard", label: "证件号", width: 180, search: { order: 2, el: "input" } },
  {
    prop: "sex",
    label: "性别",
    width: 60,
    search: { order: 3, el: "select" },
    enum: [
      { label: "男", value: "男" },
      { label: "女", value: "女" }
    ]
  },
  { prop: "phoneNo", label: "手机号", width: 125 },
  {
    prop: "expertCategory",
    label: "专家类别",
    width: 90,
    enum: categoryEnum,
    search: { order: 4, el: "select", props: { placeholder: "请选择" } },
    render: scope => <span>{scope.row.expertCategory === "2" ? "校外" : "校内"}</span>
  },
  {
    prop: "deptOrUnit",
    label: "所属学院/所在单位",
    minWidth: 150,
    render: scope => <span>{scope.row.expertCategory === "2" ? scope.row.unitName || "--" : scope.row.deptName || "--"}</span>
  },
  {
    prop: "initEdu",
    label: "初始学历",
    width: 90,
    render: scope => <span>{educationType.find(i => i.value === scope.row.initEdu)?.label || "--"}</span>
  },
  { prop: "initMajor", label: "所学专业", minWidth: 110 },
  {
    prop: "goodSubjects",
    label: "擅长科目",
    minWidth: 110,
    search: { order: 6, el: "input", props: { placeholder: "请输入擅长科目" } }
  },
  {
    prop: "title",
    label: "职称",
    width: 90,
    enum: expertTitleEnum,
    search: { order: 5, el: "select", props: { placeholder: "请选择" } }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 80 }
]);

const getTableList = (params: any) => {
  searchedParams = { ...params };
  delete searchedParams.curPage;
  delete searchedParams.pageSize;
  return getBannedExpertList(params);
};
const remove = async (rows: JudgeExpert.ResExpertDatabaseList[]) => {
  await useHandleData(deleteBannedExperts, { expertIds: rows.map(item => Number(item.expertId)) }, "删除禁止抽取专家");
  proTable.value?.clearSelection();
  proTable.value?.getTableList();
};
const removeOne = (row: JudgeExpert.ResExpertDatabaseList) => remove([row]);
const removeAll = () =>
  ElMessageBox.confirm("确认删除全部禁止抽取记录？此操作不会修改专家库中的启用/禁用状态。", "二次确认", {
    type: "warning",
    confirmButtonText: "确认删除全部"
  }).then(async () => {
    await deleteAllBannedExperts();
    proTable.value?.clearSelection();
    proTable.value?.getTableList();
  });
const openImport = () =>
  importDialogRef.value?.acceptParams({
    title: "导入禁止抽取专家",
    templateUrl: "static/禁止抽取专家导入模板.xlsx",
    saveUploadApi: importBannedExperts,
    getTableList: proTable.value?.getTableList
  });
const download = () =>
  ElMessageBox.confirm("确认导出当前禁止抽取专家数据？", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportBannedExperts, "禁止抽取专家导出", searchedParams)
  );
</script>
