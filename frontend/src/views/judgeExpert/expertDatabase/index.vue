<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList" row-key="expertId">
      <template #tableHeader="{ selectedList }">
        <el-button type="primary" :icon="CirclePlus" @click="openDialog('add')">新增</el-button>
        <el-button type="primary" :icon="Upload" plain @click="batchAdd">导入</el-button>
        <el-button type="primary" :icon="Download" plain @click="onDownloadClick">导出</el-button>
        <el-button
          v-if="!isAdminRole"
          type="success"
          :icon="Promotion"
          :disabled="!selectedList || selectedList.filter((r: any) => r.auditStatus === '1').length === 0"
          @click="onBatchSubmit(selectedList)"
          >批量提交</el-button
        >
        <el-button v-if="!isAdminRole" type="success" plain :icon="Select" @click="onSubmitAll">全部提交</el-button>
      </template>
      <template #operation="scope">
        <el-button
          v-if="!(isAdminRole && scope.row.auditStatus === '1')"
          type="primary"
          link
          :icon="Edit"
          @click="openDialog('edit', scope.row)"
          >编辑</el-button
        >
        <el-button v-if="!isAdminRole && scope.row.auditStatus === '1'" type="primary" link @click="onSubmitForAudit(scope.row)"
          >提交</el-button
        >
        <el-button v-if="isAdminRole && scope.row.auditStatus === '2'" type="warning" link @click="openAudit(scope.row)"
          >审核</el-button
        >
        <el-button v-if="isAdminRole && scope.row.auditStatus === '3'" type="success" link @click="openEvaluation(scope.row)"
          >评价</el-button
        >
      </template>
    </ProTable>
    <AddOrEditDialog ref="dialogRef" />
    <AuditDialog ref="auditDialogRef" />
    <EvaluationDialog ref="evaluationDialogRef" />
    <ImportExcel ref="importDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="expertDatabase">
import { reactive, ref } from "vue";
import { CirclePlus, Download, Edit, Upload, Promotion, Select } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { ColumnProps, ProTableInstance } from "@/components/ProTable/interface";
import ProTable from "@/components/ProTable/index.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import AuditDialog from "./components/AuditDialog.vue";
import EvaluationDialog from "./components/EvaluationDialog.vue";
import { useRole } from "@/hooks/useRole";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { useExpertTitleEnum, useExpertTypeEnum } from "@/hooks/useEnum";
import { useDownload } from "@/hooks/useDownload";
import { useHandleData } from "@/hooks/useHandleData";
import { educationType, onOrOffStatus } from "@/utils/dict";
import {
  addExpert,
  changeExpertStatus,
  editExpert,
  exportExpert,
  getExpertDatabaseList,
  importExpert,
  submitExpertForAudit,
  submitExpertForAuditBatch
} from "@/api/modules/judgeExpert";

const { isAdminRole } = useRole();
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const { expertTypeEnum } = useExpertTypeEnum();
const { expertTitleEnum } = useExpertTitleEnum();
const proTable = ref<ProTableInstance>();
let searchedParams: Record<string, any> = {};

const genderEnum = [
  { label: "男", value: "男" },
  { label: "女", value: "女" }
];
const categoryEnum = [
  { label: "校内", value: "1" },
  { label: "校外", value: "2" }
];
const auditStatusEnum = [
  { label: "未提交", value: "1" },
  { label: "待审核", value: "2" },
  { label: "审核通过", value: "3" },
  { label: "审核不通过", value: "4" }
];
const evaluatedEnum = [
  { label: "是", value: "1" },
  { label: "否", value: "2" }
];

const enumLabel = (items: any[], value: any) => items.find(item => String(item.value) === String(value))?.label || "--";
const multiEnumLabel = (items: any[], value?: string) =>
  value
    ? value
        .split(",")
        .map(v => enumLabel(items, v))
        .filter(v => v !== "--")
        .join("、")
    : "--";

const columns = reactive<ColumnProps<JudgeExpert.ResExpertDatabaseList>[]>([
  { prop: "name", label: "姓名", width: 100, search: { order: 1, el: "input", props: { placeholder: "请输入姓名" } } },
  { prop: "idCard", label: "证件号", width: 175, search: { order: 3, el: "input", props: { placeholder: "请输入证件号" } } },
  {
    prop: "sex",
    label: "性别",
    width: 55,
    headerRender: () => <span style={{ whiteSpace: "nowrap" }}>性别</span>,
    enum: genderEnum,
    search: { order: 2, el: "select", props: { placeholder: "请选择" } }
  },
  { prop: "phoneNo", label: "手机号", width: 125 },
  {
    prop: "expertCategory",
    label: "专家类别",
    width: 85,
    isShow: isAdminRole,
    isSetting: isAdminRole,
    enum: categoryEnum,
    headerRender: () => <span style={{ whiteSpace: "nowrap" }}>专家类别</span>,
    render: scope => <span>{enumLabel(categoryEnum, scope.row.expertCategory)}</span>,
    search: { order: 4, el: "select", props: { placeholder: "请选择" }, isShow: isAdminRole }
  },
  {
    prop: "deptOrUnit",
    label: "所属学院/所在单位",
    minWidth: 125,
    isShow: isAdminRole,
    isSetting: isAdminRole,
    headerRender: () => (
      <el-tooltip content="所属学院/所在单位" placement="top">
        <span class="sle" style={{ display: "inline-block", maxWidth: "100%" }}>
          所属学院/所在单位
        </span>
      </el-tooltip>
    ),
    render: scope => <span>{scope.row.expertCategory === "2" ? scope.row.unitName || "--" : scope.row.deptName || "--"}</span>
  },
  {
    prop: "type",
    label: "专家类型",
    width: isAdminRole.value ? 110 : 120,
    enum: expertTypeEnum,
    render: scope => <span>{multiEnumLabel(expertTypeEnum.value, scope.row.type)}</span>,
    search: { order: 5, el: "select", props: { placeholder: "请选择" } }
  },
  {
    prop: "title",
    label: "职称",
    width: 80,
    enum: expertTitleEnum,
    search: { order: 8, el: "select", props: { placeholder: "请选择" } }
  },
  {
    prop: "initEduMajor",
    label: "初始学历-所学专业",
    width: isAdminRole.value ? 125 : 150,
    headerRender: () => (
      <el-tooltip content="初始学历-所学专业" placement="top">
        <span class="sle" style={{ display: "inline-block", maxWidth: "100%" }}>
          初始学历-所学专业
        </span>
      </el-tooltip>
    ),
    render: scope => <span>{`${enumLabel(educationType, scope.row.initEdu)}-${scope.row.initMajor || "--"}`}</span>
  },
  {
    prop: "finalEduMajor",
    label: "最终学历-所学专业",
    width: 150,
    isShow: !isAdminRole.value,
    isSetting: !isAdminRole.value,
    headerRender: () => <span style={{ whiteSpace: "nowrap" }}>最终学历-所学专业</span>,
    render: scope => (
      <span>
        {scope.row.finalEdu && scope.row.finalEdu !== "0"
          ? `${enumLabel(educationType, scope.row.finalEdu)}-${scope.row.finalMajor || "--"}`
          : "--"}
      </span>
    )
  },
  {
    prop: "goodSubjects",
    label: "擅长科目",
    minWidth: isAdminRole.value ? 105 : 70,
    search: { order: 12, el: "input", props: { placeholder: "请输入擅长科目" } }
  },
  {
    prop: "auditStatus",
    label: "审核状态",
    width: 85,
    enum: auditStatusEnum,
    search: { order: 13, el: "select", props: { placeholder: "请选择" } }
  },
  {
    prop: "allowDraw",
    label: "状态",
    width: 125,
    isShow: isAdminRole,
    isSetting: isAdminRole,
    enum: onOrOffStatus,
    search: { order: 9, el: "select", props: { placeholder: "请选择" }, isShow: isAdminRole },
    render: scope =>
      isAdminRole.value ? (
        <el-switch
          model-value={scope.row.allowDraw}
          active-value="1"
          inactive-value="2"
          active-text={scope.row.allowDraw === "1" ? "启用" : "禁用"}
          onClick={() => onChangeStatus(scope.row)}
        />
      ) : (
        <span>{enumLabel(onOrOffStatus, scope.row.allowDraw)}</span>
      )
  },
  {
    prop: "unitName",
    label: "所在单位",
    isShow: false,
    isSetting: false,
    search: { order: 6, el: "input", props: { placeholder: "请输入所在单位" }, isShow: isAdminRole }
  },
  {
    prop: "deptId",
    label: "所属学院",
    isShow: false,
    isSetting: false,
    enum: departmentEnum,
    search: { order: 7, el: "tree-select", props: { filterable: true, placeholder: "请选择所属学院" }, isShow: isAdminRole }
  },
  {
    prop: "initEdu",
    label: "初始学历",
    isShow: false,
    isSetting: false,
    enum: educationType,
    search: { order: 10, el: "select", props: { placeholder: "请选择" } }
  },
  {
    prop: "finalEdu",
    label: "最终学历",
    isShow: false,
    isSetting: false,
    enum: educationType,
    search: { order: 11, el: "select", props: { placeholder: "请选择" } }
  },
  {
    prop: "isEvaluated",
    label: "是否评价",
    isShow: false,
    isSetting: false,
    enum: evaluatedEnum,
    search: { order: 14, el: "select", props: { placeholder: "请选择" }, isShow: isAdminRole }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 130 }
]);

// 学院角色展示勾选列，用于批量提交（仅允许选择未提交状态的数据）
if (!isAdminRole.value) {
  columns.unshift({
    type: "selection",
    width: 55,
    fixed: "left",
    selectable: (row: any) => row.auditStatus === "1"
  } as ColumnProps<JudgeExpert.ResExpertDatabaseList>);
}

const getTableList = (params: any) => {
  const newParams = JSON.parse(JSON.stringify(params));
  // 学院角色仅查看本学院且启用状态的专家
  if (!isAdminRole.value) {
    newParams.allowDraw = "1";
  }
  searchedParams = { ...newParams };
  delete searchedParams.curPage;
  delete searchedParams.pageSize;
  return getExpertDatabaseList(newParams);
};

const dialogRef = ref<InstanceType<typeof AddOrEditDialog>>();
const openDialog = (type: "add" | "edit", row: Partial<JudgeExpert.ResExpertDatabaseList> = {}) => {
  if (type === "edit" && isAdminRole.value && row.auditStatus === "1") {
    ElMessage.warning("学院未提交的专家数据，招办不可编辑");
    return;
  }
  dialogRef.value?.acceptParams({
    type,
    title: type === "add" ? "新增专家" : "编辑专家",
    row: { ...row },
    api: type === "add" ? addExpert : editExpert,
    getTableList: proTable.value?.getTableList
  });
};

const auditDialogRef = ref<InstanceType<typeof AuditDialog>>();
const openAudit = (row: JudgeExpert.ResExpertDatabaseList) =>
  auditDialogRef.value?.acceptParams({
    row,
    getTableList: proTable.value?.getTableList,
    // 连续审核时沿用当前查询条件
    searchParam: { ...(proTable.value?.searchParam || {}) }
  });
const evaluationDialogRef = ref<InstanceType<typeof EvaluationDialog>>();
const openEvaluation = (row: JudgeExpert.ResExpertDatabaseList) =>
  evaluationDialogRef.value?.acceptParams({ row, getTableList: proTable.value?.getTableList });

const onSubmitForAudit = async (row: JudgeExpert.ResExpertDatabaseList) => {
  await useHandleData(submitExpertForAudit, { expertId: row.expertId }, `提交【${row.name}】至待审核`);
  proTable.value?.getTableList();
};

// 学院用户：批量提交勾选的待提交专家（二次确认）
const onBatchSubmit = async (selectedList: any[]) => {
  const pending = (selectedList || []).filter(r => r.auditStatus === "1");
  if (pending.length === 0) {
    ElMessage.warning("请勾选待提交（未提交）状态的专家");
    return;
  }
  const ok = await ElMessageBox.confirm(`确认将选中的 ${pending.length} 条待提交专家提交至待审核？`, "批量提交确认", {
    type: "warning"
  })
    .then(() => true)
    .catch(() => false);
  if (!ok) return;
  const res = await submitExpertForAuditBatch({ expertIds: pending.map(r => Number(r.expertId)) });
  ElMessage.success(`已提交 ${res.affected} 条专家至待审核`);
  proTable.value?.clearSelection();
  proTable.value?.getTableList();
};

// 学院用户：提交本学院全部未提交的专家（二次确认）
const onSubmitAll = async () => {
  const ok = await ElMessageBox.confirm("确认提交本学院全部未提交的专家至待审核？", "全部提交确认", { type: "warning" })
    .then(() => true)
    .catch(() => false);
  if (!ok) return;
  const res = await submitExpertForAuditBatch({ expertIds: [] });
  ElMessage.success(`已提交 ${res.affected} 条专家至待审核`);
  proTable.value?.clearSelection();
  proTable.value?.getTableList();
};
const onChangeStatus = async (row: JudgeExpert.ResExpertDatabaseList) => {
  const allowDraw = row.allowDraw === "1" ? "2" : "1";
  await useHandleData(
    changeExpertStatus,
    { expertId: row.expertId, allowDraw },
    `${allowDraw === "1" ? "启用" : "禁用"}【${row.name}】`
  );
  proTable.value?.getTableList();
};

const importDialogRef = ref<InstanceType<typeof ImportExcel>>();
const batchAdd = () =>
  importDialogRef.value?.acceptParams({
    title: "导入专家",
    templateUrl: isAdminRole.value ? "static/专家导入模板.xlsx" : "static/学院专家导入模板.xlsx",
    saveUploadApi: importExpert,
    getTableList: proTable.value?.getTableList
  });
const onDownloadClick = () =>
  ElMessageBox.confirm("确认导出当前专家数据?", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportExpert, "专家库导出", searchedParams)
  );
</script>

<style scoped>
:deep(.el-table .cell) {
  padding: 0 10px;
}
</style>
