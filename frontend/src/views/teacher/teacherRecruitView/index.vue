<!-- 监考老师管理-征集老师审核 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button v-show="isShowExportBtn" type="primary" :icon="Download" plain @click="onDownloadClick">导出</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <!-- 招办：状态为已上报,待审核时显示 通过 / 不通过 -->
        <template v-if="isAdminRole && scope.row.state === RecruitStatus.PENDING_AUDIT">
          <el-button type="primary" link @click="onAudit(scope.row, 'pass')">通过</el-button>
          <el-button type="danger" link @click="onAudit(scope.row, 'reject')">不通过</el-button>
        </template>
        <!-- 招办：状态为待分配 / 不通过时显示 重置 -->
        <template
          v-if="isAdminRole && (scope.row.state === RecruitStatus.PENDING_ASSIGN || scope.row.state === RecruitStatus.REJECTED)"
        >
          <el-button type="warning" link @click="onReset(scope.row)">重置</el-button>
        </template>
      </template>
    </ProTable>
    <MonitorRecordsDialog ref="monitorDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="teacherRecruitView">
import { Teacher, Common } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive, computed, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Download } from "@element-plus/icons-vue";
import { useDownload } from "@/hooks/useDownload";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { useRole } from "@/hooks/useRole";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import {
  getTeacherRecruitViewList,
  getTeacherRecruitManageList,
  exportTeacherRecruitViewList,
  auditTeacherReport,
  resetTeacherReport
} from "@/api/modules/teacher";
import MonitorRecordsDialog from "../teacher/components/MonitorRecordsDialog.vue";

interface SearchParams {
  [key: string]: any;
}

interface TaskEnumItem extends Common.enumDict {
  taskStartTime: string;
  taskEndTime: string;
}

// 征集状态枚举（与后端 cst.TCH_REPORT_* 对应）
enum RecruitStatus {
  PENDING_AUDIT = "1", // 已上报,待审核
  PENDING_ASSIGN = "2", // 待分配
  REJECTED = "3", // 不通过
  ASSIGNED = "4" // 已分配
}

const RecruitStatusEnum = [
  { label: "已上报,待审核", value: RecruitStatus.PENDING_AUDIT },
  { label: "待分配", value: RecruitStatus.PENDING_ASSIGN },
  { label: "不通过", value: RecruitStatus.REJECTED },
  { label: "已分配", value: RecruitStatus.ASSIGNED }
];

const GenderEnum = [
  { label: "男", value: "男" },
  { label: "女", value: "女" }
];

const InvigilationExperienceEnum = [
  { label: "有", value: "1" },
  { label: "无", value: "2" }
];

let searchedParams: SearchParams | undefined = {};
const { isAdminRole } = useRole();
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const proTable = ref<ProTableInstance>();
const taskEnum = ref<TaskEnumItem[]>([]);
const taskTimeEnum = ref<Common.enumDict[]>([]);

const isShowExportBtn = computed(() => proTable.value?.tableData?.length);

const onTaskChange = value => {
  if (proTable.value) {
    proTable.value.searchParam.taskTime = undefined;
  }
  if (!value) {
    taskTimeEnum.value = [];
    return;
  }
  const curItem = taskEnum.value.find(item => item.value === value);
  if (!curItem) {
    taskTimeEnum.value = [];
    return;
  }
  const filledDateList = fillDateRange([curItem.taskStartTime, curItem.taskEndTime]);
  taskTimeEnum.value = filledDateList.map((date: string) => ({
    label: formatTime(date),
    value: date
  }));
};

// 表格配置项
const columns = reactive<ColumnProps<Teacher.ResTeacherRecruitViewList>[]>([
  { prop: "tchName", label: "姓名", search: { order: 4, el: "input", props: { placeholder: "请输入" } }, width: "90" },
  {
    prop: "sex",
    label: "性别",
    enum: GenderEnum,
    search: { order: 6, el: "select", props: { placeholder: "请选择" } },
    width: "65"
  },
  { prop: "jobNo", label: "工号", search: { order: 5, el: "input", props: { placeholder: "请输入" } }, width: "100" },
  { prop: "idCard", label: "证件号", minWidth: "50" },
  { prop: "duties", label: "职务", width: "110" },
  { prop: "profOffice", label: "专业或行政", width: "130" },
  {
    prop: "invigilationExperience",
    label: "监考经验",
    enum: InvigilationExperienceEnum,
    search: { order: 7, el: "select", props: { placeholder: "请选择" } },
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
    prop: "taskTime",
    label: "征集日期",
    render: scope => <div>{formatTime(scope.row.taskTime) || "--"}</div>,
    width: "105"
  },
  { prop: "remark", label: "备注", width: "110" },
  {
    prop: "deptId",
    label: "所属学院",
    width: "110",
    enum: departmentEnum,
    search: { order: 3, el: "tree-select", props: { filterable: true, placeholder: "请选择" } }
  },
  {
    prop: "state",
    label: "征集状态",
    enum: RecruitStatusEnum,
    search: { order: 8, el: "select", props: { placeholder: "请选择" } },
    width: "120"
  },
  {
    prop: "tchTaskId",
    label: "任务",
    enum: taskEnum,
    search: { order: 1, el: "tree-select", props: { filterable: true, placeholder: "请选择", onChange: onTaskChange } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "taskTime",
    label: "日期",
    enum: taskTimeEnum,
    search: { order: 2, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 130 }
]);

onMounted(() => {
  _getTaskList();
});

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.tchTaskId) {
    ElMessage.warning("请选择考试任务");
    return Promise.resolve({ list: [] });
  }
  searchedParams = { ...newParams };
  searchedParams?.curPage && delete searchedParams.curPage;
  searchedParams?.pageSize && delete searchedParams.pageSize;
  return getTeacherRecruitViewList(newParams);
};

// 招办审核：通过 / 不通过
const onAudit = async (row: Partial<Teacher.ResTeacherRecruitViewList>, action: "pass" | "reject") => {
  const tip = action === "pass" ? `审核通过【${row.tchName}】` : `审核不通过【${row.tchName}】`;
  await useHandleData(auditTeacherReport, { reportId: row.reportId as string, action }, tip);
  proTable.value?.getTableList();
};

// 招办重置：已审 / 已驳回 重置为 待审核
const onReset = async (row: Partial<Teacher.ResTeacherRecruitViewList>) => {
  await useHandleData(resetTeacherReport, { reportId: row.reportId as string }, `重置【${row.tchName}】征集状态`);
  proTable.value?.getTableList();
};

// 监考记录 Dialog
const monitorDialogRef = ref<InstanceType<typeof MonitorRecordsDialog> | null>(null);
const openMonitorDialog = (row: Partial<Teacher.ResTeacherRecruitViewList>) => {
  monitorDialogRef.value?.acceptParams({
    title: "监考记录",
    row: { ...row } as Partial<Teacher.ResTeacherRecruitViewList> & { tchId: string }
  });
};

const formatTime = fullTimeStr => {
  if (!fullTimeStr) {
    return "";
  }
  return fullTimeStr.slice(0, 10);
};

// 导出
const onDownloadClick = () => {
  const taskName = getTaskNameById(searchedParams?.tchTaskId);
  ElMessageBox.confirm("确认导出征集老师审核数据?", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportTeacherRecruitViewList, taskName, searchedParams)
  );
};

// 获取全部任务列表
const _getTaskList = async () => {
  const { list = [] } = await getTeacherRecruitManageList({ curPage: 1, pageSize: 1000 });
  taskEnum.value = list.map((item: any) => ({
    label: item.tchTaskName,
    value: item.tchTaskId,
    taskStartTime: item.taskStartTime,
    taskEndTime: item.taskEndTime
  }));
};

// 填充日期 ["2024-12-30 00:00:00", "2025-01-01 00:00:00"] --> ["2024-12-30 00:00:00", "2024-12-31 00:00:00", "2025-01-01 00:00:00"]
const fillDateRange = dates => {
  const startDate = new Date(dates[0]);
  const endDate = new Date(dates[1]);
  const result: any = [];
  // 使用一个循环从开始日期到结束日期，依次将每一天加入结果加入结果数组
  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
    // 格式化当前日期为 "YYYY-MM-DD HH:MM:SS"
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0"); // 月份从0开始，补齐2位
    const day = String(d.getDate()).padStart(2, "0");
    const hours = String(d.getHours()).padStart(2, "0");
    const minutes = String(d.getMinutes()).padStart(2, "0");
    const seconds = String(d.getSeconds()).padStart(2, "0");
    result.push(`${year}-${month}-${day} ${hours}:${minutes}:${seconds}`);
  }
  return result;
};

const getTaskNameById = (id: string) => {
  const findedItem = taskEnum.value.find(item => item.value === id);
  return findedItem?.label || "";
};
</script>
