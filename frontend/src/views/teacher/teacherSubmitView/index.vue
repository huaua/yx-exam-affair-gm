<!-- 监考老师管理-监考老师上报查看 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格操作 -->
      <template #operation="scope">
        <!-- 已上报,待审核 / 待分配 显示取消上报；不通过、已分配不显示操作 -->
        <el-button
          v-if="scope.row.state === RecruitStatus.PENDING_AUDIT || scope.row.state === RecruitStatus.PENDING_ASSIGN"
          type="primary"
          link
          @click="onCancelSubmitClick(scope.row)"
        >
          取消上报
        </el-button>
      </template>
    </ProTable>
    <EditDialog ref="editDialogRef" />
    <MonitorRecordsDialog ref="monitorDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="teacherSubmitView">
import { Teacher, Common } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import EditDialog from "./components/EditDialog.vue";
import MonitorRecordsDialog from "../teacher/components/MonitorRecordsDialog.vue";
import { getTeacherSubmitViewList, checkCancelSubmit, cancelSubmit, getTeacherSumbitList } from "@/api/modules/teacher";

interface TaskEnumItem extends Common.enumDict {
  taskStartTime: string;
  taskEndTime: string;
}

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

const proTable = ref<ProTableInstance>();
const taskEnum = ref<TaskEnumItem[]>([]);
const taskTimeEnum = ref<Common.enumDict[]>([]);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.tchTaskId) {
    ElMessage.warning("请选择考试任务");
    return Promise.resolve({ list: [] });
  }
  return getTeacherSubmitViewList(newParams);
};

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
const columns = reactive<ColumnProps<Teacher.ResTeacherSubmitViewList>[]>([
  { prop: "tchName", label: "姓名", search: { order: 4, el: "input", props: { placeholder: "请输入" } }, width: "120" },
  {
    prop: "sex",
    label: "性别",
    enum: GenderEnum,
    search: { order: 6, el: "select", props: { placeholder: "请选择" } },
    width: "70"
  },
  { prop: "jobNo", label: "工号", search: { order: 5, el: "input", props: { placeholder: "请输入" } } },
  { prop: "duties", label: "职务" },
  { prop: "profOffice", label: "专业" },
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
    width: "100"
  },
  {
    prop: "taskTime",
    label: "日期",
    render: scope => <div>{formatTime(scope.row.taskTime) || "--"}</div>
  },
  {
    prop: "state",
    label: "征集状态",
    enum: RecruitStatusEnum,
    search: { order: 8, el: "select", props: { placeholder: "请选择" } },
    width: 130
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
  { prop: "remark", label: "备注" },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

onMounted(() => {
  _getTaskList();
});

const onCancelSubmitClick = async (row: Partial<Teacher.ResTeacherSubmitViewList> = {}) => {
  try {
    const res = await checkCancelSubmit({ reportId: row.reportId! });
    if (res?.obj) {
      // 可以取消
      await useHandleData(cancelSubmit, { reportId: row.reportId! }, `取消上报【${row.tchName}】`);
      proTable.value?.getTableList();
    } else {
      // 不可以取消(替换)
      openDialog("cancel", row);
    }
  } catch (error) {
    console.error("error:", error);
  }
};

// 打开 Dialog(编辑)
const editDialogRef = ref<InstanceType<typeof EditDialog> | null>(null);
const openDialog = (type: string, row: Partial<Teacher.ResTeacherSubmitViewList> = {}) => {
  const params = {
    type,
    title: type === "cancel" ? "取消上报" : "编排替换",
    row: { ...row },
    getTableList: proTable.value?.getTableList
  };
  editDialogRef.value?.acceptParams(params);
};

// 监考记录 Dialog
const monitorDialogRef = ref<InstanceType<typeof MonitorRecordsDialog> | null>(null);
const openMonitorDialog = (row: Partial<Teacher.ResTeacherSubmitViewList>) => {
  monitorDialogRef.value?.acceptParams({
    title: "监考记录",
    row: { ...row } as Partial<Teacher.ResTeacherSubmitViewList> & { tchId: string }
  });
};

const formatTime = fullTimeStr => {
  if (!fullTimeStr) {
    return "";
  }
  return fullTimeStr.slice(0, 10);
};

// 获取所有需要上报的任务列表
const _getTaskList = async () => {
  const { list = [] } = await getTeacherSumbitList({ curPage: 1, pageSize: 1000 });
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
  // 使用一个循环从开始日期到结束日期，依次将每一天加入结果数组
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
</script>
