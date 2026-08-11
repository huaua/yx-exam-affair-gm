<!-- 监考老师管理-教职工管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList" row-key="tchId">
      <!-- 表格 header 按钮 -->
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
      <!-- 表格操作 -->
      <template #operation="scope">
        <!-- 招办不可编辑学院未提交(UNSUBMITTED)的数据 -->
        <el-button
          v-if="!(isAdminRole && scope.row.auditStatus === '1')"
          type="primary"
          link
          :icon="Edit"
          @click="openDialog('edit', scope.row)"
        >
          编辑
        </el-button>
        <!-- 学院用户对未提交(UNSUBMITTED)的数据：显示"提交"按钮 -->
        <el-button v-if="!isAdminRole && scope.row.auditStatus === '1'" type="primary" link @click="onSubmitForAudit(scope.row)">
          提交
        </el-button>
        <template v-if="isAdminRole">
          <el-button v-if="scope.row.auditStatus === '2'" type="warning" link @click="openAuditDialog(scope.row)">
            审核
          </el-button>
        </template>
      </template>
    </ProTable>
    <AddOrEditDialog ref="dialogRef" />
    <ImportExcel ref="importDialogRef" />
    <AuditDialog ref="auditDialogRef" />
    <MonitorRecordsDialog ref="monitorDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="teacher">
import { Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Edit, Upload, Download, Promotion, Select } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { useRole } from "@/hooks/useRole";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { useDownload } from "@/hooks/useDownload";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import AuditDialog from "./components/AuditDialog.vue";
import MonitorRecordsDialog from "./components/MonitorRecordsDialog.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import { onOrOffStatus } from "@/utils/dict";
import {
  getTeacherList,
  addTeacher,
  editTeacher,
  importTeacher,
  exportTeacher,
  changeTeacherEnabled,
  submitTchForAudit,
  submitTchForAuditBatch
} from "@/api/modules/teacher";

interface SearchParams {
  [key: string]: any;
}

type MapString = {
  [key: string]: string;
};

type MapFunction = {
  [key: string]: (params: any) => Promise<any>;
};
const { isAdminRole } = useRole();
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const proTable = ref<ProTableInstance>();
let searchedParams: SearchParams = {};

const genderEnum = [
  { label: "男", value: "男" },
  { label: "女", value: "女" }
];

const auditStatusEnum = [
  { label: "未提交", value: "1" },
  { label: "待审核", value: "2" },
  { label: "审核通过", value: "3" },
  { label: "审核不通过", value: "4" }
];

const invigilationExperienceEnum = [
  { label: "有", value: "1" },
  { label: "无", value: "2" }
];

const enabledStatusEnum = onOrOffStatus;

// 表格配置项
const columns = reactive<ColumnProps<Teacher.ResTeacherList>[]>([
  { prop: "tchName", label: "姓名", width: "110", search: { order: 1, el: "input", props: { placeholder: "请输入" } } },
  {
    prop: "sex",
    label: "性别",
    enum: genderEnum,
    search: { order: 4, el: "tree-select", props: { placeholder: "请选择" } },
    width: "70"
  },
  { prop: "jobNo", label: "工号", width: "120", search: { order: 2, el: "input", props: { placeholder: "请输入" } } },
  {
    prop: "idCard",
    label: "证件号",
    search: { order: 3, el: "input", props: { placeholder: "请输入完整证件号" } },
    minWidth: "170"
  },
  {
    prop: "deptName",
    label: "所属学院",
    isShow: isAdminRole,
    minWidth: "140"
  },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    search: { order: 7, el: "tree-select", props: { filterable: true, placeholder: "请选择" }, isShow: isAdminRole },
    isShow: false,
    isSetting: false
  },
  { prop: "duties", label: "职务", minWidth: "150" },
  { prop: "profOffice", label: "专业或行政", minWidth: "150" },
  {
    prop: "invigilationExperience",
    label: "监考经验",
    enum: invigilationExperienceEnum,
    search: { order: 6, el: "tree-select", props: { placeholder: "请选择" } },
    render: scope => {
      if (scope.row.invigilationExperience === "1") {
        return (
          <el-button type="primary" link onClick={() => openMonitorDialog(scope.row)}>
            查看
          </el-button>
        );
      }
      return <span>无</span>;
    },
    width: "95"
  },
  {
    prop: "auditStatus",
    label: "审核状态",
    enum: auditStatusEnum,
    search: { order: 5, el: "tree-select", props: { placeholder: "请选择" } },
    width: "90"
  },
  {
    prop: "enabledStatus",
    label: "状态",
    enum: enabledStatusEnum,
    isShow: isAdminRole,
    search: { order: 8, el: "select", props: { placeholder: "请选择" }, isShow: isAdminRole },
    render: scope => (
      <el-switch
        model-value={scope.row.enabledStatus}
        active-text={scope.row.enabledStatus === "1" ? "启用" : "禁用"}
        active-value={"1"}
        inactive-value={"2"}
        onClick={() => onChangeEnabled(scope.row)}
      />
    ),
    width: "130"
  },
  { prop: "operation", label: "操作", fixed: "right", width: 140 }
]);

// 学院角色展示勾选列，用于批量提交（仅允许选择未提交状态的数据）
if (!isAdminRole.value) {
  columns.unshift({
    type: "selection",
    width: 55,
    fixed: "left",
    selectable: (row: any) => row.auditStatus === "1"
  } as ColumnProps<Teacher.ResTeacherList>);
}

const getTableList = (params: any) => {
  const newParams = JSON.parse(JSON.stringify(params));
  searchedParams = { ...newParams };
  delete searchedParams.curPage;
  delete searchedParams.pageSize;
  return getTeacherList(newParams);
};

// 打开 Dialog(新增、编辑)
const dialogRef = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openDialog = (type: string, row: Partial<Teacher.ResTeacherList> = {}) => {
  // 招办不可编辑学院未提交的教职工数据
  if (type === "edit" && isAdminRole.value && row.auditStatus === "1") {
    ElMessage.warning("学院未提交的教职工数据，招办不可编辑");
    return;
  }
  const params = {
    type,
    title: getDialogTitle(type),
    row: { ...row },
    api: getDialogApi(type),
    getTableList: proTable.value?.getTableList
  };
  dialogRef.value?.acceptParams(params);
};

// 审核 Dialog（招办）
const auditDialogRef = ref<InstanceType<typeof AuditDialog> | null>(null);
const openAuditDialog = (row: Partial<Teacher.ResTeacherList>) => {
  auditDialogRef.value?.acceptParams({
    title: "审核教职工",
    row: { ...row } as Partial<Teacher.ResTeacherList> & { tchId: string },
    getTableList: proTable.value?.getTableList,
    // 连续审核时沿用当前查询条件
    searchParam: { ...(proTable.value?.searchParam || {}) }
  });
};

// 监考记录 Dialog
const monitorDialogRef = ref<InstanceType<typeof MonitorRecordsDialog> | null>(null);
const openMonitorDialog = (row: Partial<Teacher.ResTeacherList>) => {
  monitorDialogRef.value?.acceptParams({
    title: "监考记录",
    row: { ...row } as Partial<Teacher.ResTeacherList> & { tchId: string }
  });
};

// 启用/禁用（招办）
const onChangeEnabled = async (row: Partial<Teacher.ResTeacherList>) => {
  const newStatus = row.enabledStatus === "1" ? "2" : "1";
  const action = newStatus === "1" ? "启用" : "禁用";
  await useHandleData(
    changeTeacherEnabled,
    { tchId: row.tchId as string, enabledStatus: newStatus },
    `${action}【${row.tchName}】教职工`
  );
  proTable.value?.getTableList();
};

// 学院用户：提交未提交的教职工至待审核
const onSubmitForAudit = async (row: Partial<Teacher.ResTeacherList>) => {
  await useHandleData(submitTchForAudit, { tchId: row.tchId as string }, `提交【${row.tchName}】至待审核`);
  proTable.value?.getTableList();
};

// 学院用户：批量提交勾选的待提交教职工（二次确认）
const onBatchSubmit = async (selectedList: any[]) => {
  const pending = (selectedList || []).filter(r => r.auditStatus === "1");
  if (pending.length === 0) {
    ElMessage.warning("请勾选待提交（未提交）状态的数据");
    return;
  }
  const ok = await ElMessageBox.confirm(`确认将选中的 ${pending.length} 条待提交数据提交至待审核？`, "批量提交确认", {
    type: "warning"
  })
    .then(() => true)
    .catch(() => false);
  if (!ok) return;
  const res = await submitTchForAuditBatch({ tchIds: pending.map(r => Number(r.tchId)) });
  ElMessage.success(`已提交 ${res.affected} 条数据至待审核`);
  proTable.value?.clearSelection();
  proTable.value?.getTableList();
};

// 学院用户：提交本学院全部未提交的教职工（二次确认）
const onSubmitAll = async () => {
  const ok = await ElMessageBox.confirm("确认提交本学院全部未提交的数据至待审核？", "全部提交确认", { type: "warning" })
    .then(() => true)
    .catch(() => false);
  if (!ok) return;
  const res = await submitTchForAuditBatch({ tchIds: [] });
  ElMessage.success(`已提交 ${res.affected} 条数据至待审核`);
  proTable.value?.clearSelection();
  proTable.value?.getTableList();
};

const getDialogTitle = (type: string) => {
  const dialogTitleMap: MapString = {
    add: "新增教职工",
    edit: "编辑教职工"
  };
  return dialogTitleMap[type];
};

const getDialogApi = (type: string) => {
  const dialogApiMap: MapFunction = {
    add: addTeacher,
    edit: editTeacher
  };
  return dialogApiMap[type];
};

// 导入
const importDialogRef = ref<InstanceType<typeof ImportExcel> | null>(null);
const batchAdd = () => {
  const params = {
    title: "导入教职工",
    templateUrl: isAdminRole.value ? "static/教职工导入模板.xlsx" : "static/学院教职工导入模板.xlsx",
    saveUploadApi: importTeacher,
    getTableList: proTable.value?.getTableList
  };
  importDialogRef.value?.acceptParams(params);
};

const onDownloadClick = () => {
  ElMessageBox.confirm("确认导出当前教职工数据?", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportTeacher, "教职工导出", searchedParams)
  );
};
</script>
