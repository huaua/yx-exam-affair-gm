<template>
  <div class="table-box college-report">
    <div class="filter-bar">
      <el-select v-model="taskId" filterable clearable placeholder="任务名称-请选择" style="width: 340px" @change="onTaskChange">
        <el-option v-for="item in tasks" :key="item.judgeTaskId" :label="item.taskName" :value="item.judgeTaskId" />
      </el-select>
      <el-select
        v-model="directionId"
        filterable
        clearable
        :disabled="!taskId"
        placeholder="评分方向-请选择"
        style="width: 220px"
      >
        <el-option
          v-for="item in directionOptions"
          :key="item.requirementId"
          :label="item.directionName"
          :value="item.requirementId"
        />
      </el-select>
      <el-select v-model="reportStatus" clearable placeholder="上报状态-请选择" style="width: 160px">
        <el-option label="未上报" value="1" />
        <el-option label="已上报" value="2" />
      </el-select>
      <el-button type="primary" @click="onQuery">查询</el-button>
    </div>
    <div class="toolbar"><el-button type="warning" :disabled="!taskId" @click="exportTask">导出</el-button></div>
    <el-empty v-if="directions.length === 0" description="暂无学院上报任务" />
    <el-table v-else :data="directions" border>
      <el-table-column
        prop="taskName"
        label="任务名称"
        min-width="180"
        align="center"
        header-align="center"
        show-overflow-tooltip
      />
      <el-table-column prop="directionName" label="评分方向" min-width="180" align="center" header-align="center" />
      <el-table-column label="评分日期" width="210" align="center"
        ><template #default="s">{{ formatScoreDate(s.row) }}</template></el-table-column
      >
      <el-table-column prop="expertNum" label="要求上报数" width="110" align="center" />
      <el-table-column prop="reportedNum" label="已上报数" width="100" align="center" />
      <el-table-column prop="approvedNum" label="审核通过数" width="110" align="center" />
      <el-table-column label="上报状态" width="100" align="center">
        <template #default="s">{{ Number(s.row.reportedNum) > 0 ? "已上报" : "未上报" }}</template>
      </el-table-column>
      <el-table-column label="操作" width="210" fixed="right" align="center">
        <template #default="s">
          <el-button link type="primary" @click="openReport(s.row)">{{ reportAction(s.row) }}</el-button>
          <el-button
            v-if="Number(s.row.reportedNum) > 0 || Number(s.row.rejectedNum) > 0"
            link
            type="primary"
            @click="openView(s.row)"
            >查看</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="reportVisible"
      title="评委专家上报"
      width="1280px"
      destroy-on-close
      :close-on-click-modal="false"
      @closed="fetchAll"
    >
      <div class="task-info">
        <div class="task-info-row">
          <span>任务名称：{{ currentDirection?.taskName }}</span>
          <span>评分日期：{{ formatScoreDate(currentDirection) }}</span>
        </div>
        <div class="task-info-row">
          <span>评分方向：{{ currentDirection?.directionName }}</span>
          <span>要求上报数：{{ currentDirection?.expertNum }}</span>
          <span>已选择数：{{ selectedCount }}</span>
          <span>已上报数：{{ currentDirection?.reportedNum }}</span>
          <span>审核通过数：{{ currentDirection?.approvedNum }}</span>
        </div>
      </div>
      <div class="report-cols">
        <div class="report-col report-col-left">
          <el-table :data="reportedRows" border height="430">
            <el-table-column prop="name" label="姓名" width="85" show-overflow-tooltip />
            <el-table-column prop="sex" label="性别" width="56" align="center" header-cell-class-name="hdr-nowrap" />
            <el-table-column prop="phoneNo" label="手机号" width="110" show-overflow-tooltip />
            <el-table-column label="职称" width="75" align="center" show-overflow-tooltip>
              <template #default="s">{{ titleText(s.row.title) }}</template>
            </el-table-column>
            <el-table-column
              prop="goodSubjects"
              label="擅长科目"
              width="120"
              header-cell-class-name="hdr-nowrap"
              show-overflow-tooltip
            />
            <el-table-column label="上报状态" min-width="100" align="center">
              <template #default="s">{{ submitStateText(s.row) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="70" align="center">
              <template #default="s">
                <el-button v-if="s.row.submitState === '1'" link type="primary" @click="clearDraft(s.row)">清除</el-button>
                <span v-else>--</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="report-col">
          <div class="dialog-filter">
            <el-input
              v-model="expertQuery.name"
              clearable
              placeholder="姓名-请输入"
              style="width: 140px"
              @keyup.enter="loadExperts"
            />
            <el-select v-model="expertQuery.title" clearable placeholder="职称-请选择" style="width: 140px" @change="loadExperts">
              <el-option v-for="item in expertTitleEnum" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select
              v-model="expertQuery.goodSubjects"
              clearable
              placeholder="擅长科目-请选择"
              style="width: 140px"
              @change="loadExperts"
            >
              <el-option v-for="item in goodSubjectEnum" :key="item.value" :label="item.label" :value="item.label" />
            </el-select>
            <el-button type="primary" @click="loadExperts">查询</el-button>
          </div>
          <el-table :data="filteredCandidates" border height="380">
            <el-table-column label="选择" width="60" align="center" header-cell-class-name="hdr-nowrap">
              <template #default="s">
                <el-checkbox :model-value="false" @change="toggleDraft(s.row)" />
              </template>
            </el-table-column>
            <el-table-column prop="name" label="姓名" width="80" show-overflow-tooltip />
            <el-table-column prop="sex" label="性别" width="56" align="center" header-cell-class-name="hdr-nowrap" />
            <el-table-column prop="phoneNo" label="手机号" width="110" show-overflow-tooltip />
            <el-table-column label="职称" width="75" align="center" show-overflow-tooltip>
              <template #default="s">{{ titleText(s.row.title) }}</template>
            </el-table-column>
            <el-table-column prop="goodSubjects" label="擅长科目" min-width="95" show-overflow-tooltip />
          </el-table>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer-center">
          <el-button type="primary" @click="submitReport">提交</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog
      v-model="viewVisible"
      title="已上报专家查看"
      width="1040px"
      top="8vh"
      destroy-on-close
      :close-on-click-modal="false"
    >
      <div class="view-filter">
        <el-input v-model="viewQuery.name" clearable placeholder="姓名-请输入" style="width: 160px" />
        <el-select v-model="viewQuery.auditState" clearable placeholder="上报状态-请选择" style="width: 160px">
          <el-option label="已上报，待审核" value="1" />
          <el-option label="审核通过" value="2" />
          <el-option label="审核不通过" value="3" />
        </el-select>
      </div>
      <el-table :data="filteredReportedExperts" border height="560">
        <el-table-column prop="name" label="姓名" width="80" />
        <el-table-column prop="idCard" label="证件号" width="200" />
        <el-table-column prop="sex" label="性别" width="60" />
        <el-table-column prop="phoneNo" label="手机号" width="115" />
        <el-table-column label="职称" width="90" show-overflow-tooltip
          ><template #default="s">{{ titleText(s.row.title) }}</template></el-table-column
        >
        <el-table-column prop="goodSubjects" label="擅长科目" min-width="100" show-overflow-tooltip />
        <el-table-column label="上报状态" width="150" header-cell-class-name="hdr-nowrap"
          ><template #default="s">{{ auditText(s.row.auditState) }}</template></el-table-column
        >
        <el-table-column label="操作" width="90" align="center">
          <template #default="s"
            ><el-button v-if="s.row.auditState === '1'" link type="primary" @click="openReplace(s.row)">替换</el-button></template
          >
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="replaceVisible" title="专家替换" width="760px" destroy-on-close :close-on-click-modal="false">
      <p>请在本学院符合要求的专家中选择一位替换当前待审核专家。</p>
      <el-input
        v-model="replaceName"
        clearable
        placeholder="姓名-请输入"
        style="width: 180px; margin-bottom: 12px"
        @keyup.enter="loadReplaceCandidates"
      />
      <el-table
        :data="replaceCandidates"
        border
        height="340"
        highlight-current-row
        @current-change="row => (selectedReplacement = row)"
      >
        <el-table-column width="55" align="center"
          ><template #default="s"
            ><el-radio :model-value="selectedReplacement?.expertId" :label="s.row.expertId"><span /></el-radio></template
        ></el-table-column>
        <el-table-column prop="name" label="姓名" /><el-table-column prop="sex" label="性别" width="70" />
        <el-table-column prop="phoneNo" label="手机号" width="120" />
        <el-table-column label="职称"
          ><template #default="s">{{ titleText(s.row.title) }}</template></el-table-column
        >
        <el-table-column prop="goodSubjects" label="擅长科目" min-width="140" />
      </el-table>
      <template #footer
        ><el-button @click="replaceVisible = false">取消</el-button
        ><el-button type="primary" @click="confirmReplace">提交</el-button></template
      >
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="judgeResponse">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useDownload } from "@/hooks/useDownload";
import { useExpertTitleEnum, useGoodSubjectEnum } from "@/hooks/useEnum";
import {
  exportCollegeReportedExperts,
  getCollegeReportDirections,
  getCollegeReportedExperts,
  getCollegeReportExperts,
  getCollegeReportTasks,
  removeCollegeReportDraft,
  replaceCollegeReportedExpert,
  saveCollegeReportDraft,
  submitCollegeReport
} from "@/api/modules/judgeExpert";

const tasks = ref<any[]>([]),
  taskId = ref(""),
  directions = ref<any[]>([]),
  experts = ref<any[]>([]),
  reportedExperts = ref<any[]>([]),
  allDirections = ref<any[]>([]),
  directionOptions = ref<any[]>([]),
  directionId = ref(""),
  reportStatus = ref("");
const reportVisible = ref(false),
  viewVisible = ref(false),
  replaceVisible = ref(false);
const currentDirection = ref<any>(),
  replacingRow = ref<any>(),
  replaceCandidates = ref<any[]>([]),
  selectedReplacement = ref<any>();
const replaceName = ref("");
const expertQuery = reactive({ name: "", title: "", goodSubjects: "" });
const viewQuery = reactive({ name: "", auditState: "" });
const { expertTitleEnum } = useExpertTitleEnum();
const { goodSubjectEnum } = useGoodSubjectEnum();
const titleText = (val: string) => expertTitleEnum.value.find((e: any) => String(e.value) === String(val))?.label || val;
// 左侧 = 已暂存 + 全部已上报（含审核不通过），列表和"已选择数"均显示总数
const reportedRows = computed(() => [...experts.value.filter(s => s.submitState === "1"), ...reportedExperts.value]);
const candidateRows = computed(() => experts.value.filter(s => s.submitState !== "1" && s.submitState !== "2"));
const filteredReportedExperts = computed(() => {
  let list = reportedExperts.value;
  const name = viewQuery.name.trim();
  if (name) list = list.filter((s: any) => String(s.name).includes(name));
  if (viewQuery.auditState) list = list.filter((s: any) => String(s.auditState) === String(viewQuery.auditState));
  return list;
});
// 当前已选择数（全部，含审核不通过）
const selectedCount = computed(() => reportedRows.value.length);
// 有效已选择数（仅待审核+审核通过），用于选人时判断是否还能继续追加
const validSelectedCount = computed(() => reportedRows.value.filter((s: any) => s.auditState !== "3").length);
// 右侧查询条件只在本地过滤候选列表，不影响左侧
const filteredCandidates = computed(() => {
  let list = candidateRows.value;
  const name = expertQuery.name.trim();
  if (name) list = list.filter((s: any) => String(s.name).includes(name));
  if (expertQuery.title) list = list.filter((s: any) => String(s.title) === String(expertQuery.title));
  if (expertQuery.goodSubjects) list = list.filter((s: any) => String(s.goodSubjects).includes(expertQuery.goodSubjects));
  return list;
});
const submitStateText = (row: any) => {
  if (row.submitState === "1") return "未提交";
  return auditText(row.auditState);
};
const currentTask = computed(() => tasks.value.find(item => item.judgeTaskId === taskId.value));
const formatScoreDate = (row: any) =>
  row.taskStartTime && row.taskEndTime
    ? `${String(row.taskStartTime).slice(0, 10)}至${String(row.taskEndTime).slice(0, 10)}`
    : "--";

onMounted(async () => {
  const res: any = await getCollegeReportTasks();
  tasks.value = res.list || [];
  await fetchAll();
  await loadDirectionOptions();
});
const fetchAll = async () => {
  const params: any = {};
  if (taskId.value) params.judgeTaskId = taskId.value;
  const res: any = await getCollegeReportDirections(params);
  allDirections.value = res.list || [];
  applyFilter();
};
// 方向下拉严格按所选任务过滤；未选任务时下拉禁用 + 空选项
const loadDirectionOptions = async () => {
  if (!taskId.value) {
    directionOptions.value = [];
    directionId.value = "";
    return;
  }
  const res: any = await getCollegeReportDirections({ judgeTaskId: taskId.value });
  directionOptions.value = res.list || [];
  if (directionId.value && !directionOptions.value.find((d: any) => String(d.requirementId) === String(directionId.value))) {
    directionId.value = "";
  }
};
const applyFilter = () => {
  let list = allDirections.value;
  if (taskId.value) list = list.filter((d: any) => String(d.judgeTaskId) === String(taskId.value));
  if (directionId.value) list = list.filter((d: any) => String(d.requirementId) === String(directionId.value));
  if (reportStatus.value === "1") list = list.filter((d: any) => Number(d.reportedNum) === 0);
  else if (reportStatus.value === "2") list = list.filter((d: any) => Number(d.reportedNum) > 0);
  directions.value = list;
};
const onTaskChange = async () => {
  directionId.value = "";
  reportStatus.value = "";
  await loadDirectionOptions();
  applyFilter();
};
const onQuery = () => applyFilter();
const reportAction = (row: any) => (Number(row.reportedNum) > 0 ? "去追加" : "去上报");
const openReport = (row: any) => {
  currentDirection.value = row;
  // 注意：不要改写顶部 taskId 查询框，否则会把列表"任务名称"查询框强制选中当前任务，
  // 且弹窗关闭后 fetchAll 会错误收窄列表。弹窗内任务信息直接用 currentDirection 即可。
  expertQuery.name = "";
  expertQuery.title = "";
  expertQuery.goodSubjects = "";
  reportVisible.value = true;
  loadAll();
};
const loadAll = async () => {
  await Promise.all([loadExperts(), loadReported()]);
};
const loadExperts = async () => {
  // 查询条件只在本地过滤右侧候选列表，避免后端过滤影响左侧已暂存/已上报数据
  const res: any = await getCollegeReportExperts({ requirementId: currentDirection.value.requirementId });
  experts.value = res.list || [];
};
const loadReported = async () => {
  const res: any = await getCollegeReportedExperts({ requirementId: currentDirection.value.requirementId });
  reportedExperts.value = res.list || [];
};
const toggleDraft = async (row: any) => {
  // 右侧候选行一定是未操作过的（submitState=0/不存在），统一走"暂存"
  const limit = Number(currentDirection.value?.expertNum) || 0;
  if (validSelectedCount.value >= limit) {
    ElMessage.warning(`有效选择数量已达要求上报数（${limit}），请先清除不需要的专家或替换待审核的专家`);
    return;
  }
  try {
    await saveCollegeReportDraft({ requirementId: currentDirection.value.requirementId, expertId: row.expertId });
    ElMessage.success("已暂存");
    await loadAll();
  } catch (error) {
    // 全局拦截器已统一弹出后端错误，这里不再重复提示
  }
};
const clearDraft = async (row: any) => {
  if (!row.reportId || row.reportId === "0") {
    ElMessage.warning("该专家暂无可清除记录，请刷新页面");
    return;
  }
  try {
    await removeCollegeReportDraft({ reportId: row.reportId });
    ElMessage.success("已清除");
    await loadAll();
  } catch (error) {
    // 全局拦截器已统一弹出后端错误，这里不再重复提示
  }
};
const submitReport = async () => {
  try {
    await submitCollegeReport({ requirementId: currentDirection.value.requirementId });
    ElMessage.success("提交成功");
    reportVisible.value = false;
  } catch (error) {
    // 全局拦截器已统一弹出后端错误，这里不再重复提示
  }
};
const openView = async (row: any) => {
  currentDirection.value = row;
  viewQuery.name = "";
  viewQuery.auditState = "";
  const res: any = await getCollegeReportedExperts({ requirementId: row.requirementId });
  reportedExperts.value = res.list || [];
  viewVisible.value = true;
};
const auditText = (state: string) => ({ "1": "已上报，待审核", "2": "审核通过", "3": "审核不通过" })[state] || "--";
const openReplace = (row: any) => {
  replacingRow.value = row;
  selectedReplacement.value = undefined;
  replaceName.value = "";
  replaceVisible.value = true;
  loadReplaceCandidates();
};
const loadReplaceCandidates = async () => {
  const res: any = await getCollegeReportExperts({
    requirementId: currentDirection.value.requirementId,
    name: replaceName.value
  });
  // 替换只能选从未上报过的专家；已选择未提交、已上报（含审核不通过）的均不允许再次选择
  replaceCandidates.value = (res.list || []).filter((s: any) => !s.submitState || s.submitState === "0");
};
let replaceSearchTimer: ReturnType<typeof setTimeout> | null = null;
watch(replaceName, () => {
  if (replaceSearchTimer) clearTimeout(replaceSearchTimer);
  replaceSearchTimer = setTimeout(() => loadReplaceCandidates(), 300);
});
const confirmReplace = async () => {
  if (!selectedReplacement.value) {
    ElMessage.warning("请选择替换专家");
    return;
  }
  try {
    await replaceCollegeReportedExpert({
      reportId: replacingRow.value.reportId,
      newExpertId: selectedReplacement.value.expertId
    });
    ElMessage.success("替换成功");
    replaceVisible.value = false;
    openView(currentDirection.value);
    fetchAll();
  } catch (error) {
    // 全局拦截器已统一弹出后端错误，这里不再重复提示
  }
};
const exportTask = () =>
  useDownload(exportCollegeReportedExperts, currentTask.value?.taskName || "评委专家上报", { judgeTaskId: taskId.value });
</script>

<style scoped lang="scss">
.filter-bar,
.dialog-filter,
.view-filter {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}
.filter-bar,
.toolbar,
.dialog-filter,
.view-filter,
.task-info {
  margin-bottom: 14px;
}
.task-info {
  .task-info-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px 24px;
    line-height: 24px;
  }
}
.report-cols {
  display: flex;
  gap: 16px;
  .report-col {
    flex: 1 1 0;
    min-width: 0;
  }
  .report-col-left {
    flex: 1.3 1 0;
  }
}
.report-cols :deep(.el-table .cell) {
  white-space: nowrap;
}
.report-cols :deep(.el-table th.hdr-nowrap .cell) {
  white-space: nowrap;
  padding-left: 4px;
  padding-right: 4px;
}
.dialog-footer-center {
  display: flex;
  justify-content: center;
  width: 100%;
}
.college-report :deep(.el-empty) {
  min-height: 420px;
}
</style>
