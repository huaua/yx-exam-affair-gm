<template>
  <div class="table-box report-audit">
    <div class="filter-bar" v-show="showFilter">
      <el-select
        v-model="query.judgeTaskId"
        filterable
        clearable
        placeholder="任务名称-请选择"
        style="width: 300px"
        @change="onTaskChange"
      >
        <el-option v-for="item in tasks" :key="item.judgeTaskId" :label="item.taskName" :value="item.judgeTaskId" />
      </el-select>
      <el-select
        v-model="query.directionCode"
        filterable
        clearable
        :disabled="!query.judgeTaskId"
        placeholder="评分方向-请选择"
        style="width: 180px"
      >
        <el-option
          v-for="item in directions"
          :key="item.requirementId"
          :label="`${item.directionCode}${item.directionName}`"
          :value="item.directionCode"
        />
      </el-select>
      <el-select v-model="query.deptId" filterable clearable placeholder="所属学院-请选择" style="width: 170px">
        <el-option v-for="item in departments" :key="item.deptId" :label="item.deptName" :value="item.deptId" />
      </el-select>
      <el-input v-model="query.name" clearable placeholder="姓名-请输入" style="width: 150px" @keyup.enter="search" />
      <el-select v-model="query.sex" clearable placeholder="性别-请选择" style="width: 135px">
        <el-option label="男" value="男" /><el-option label="女" value="女" />
      </el-select>
      <el-select v-model="query.title" filterable clearable placeholder="职称-请选择" style="width: 150px">
        <el-option v-for="item in expertTitleEnum" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select v-model="query.auditState" clearable placeholder="状态-请选择" style="width: 170px">
        <el-option label="已上报，待审核" value="1" /><el-option label="审核通过" value="2" /><el-option
          label="审核不通过"
          value="3"
        />
      </el-select>
      <el-button type="primary" @click="search()">查询</el-button>
    </div>

    <el-empty v-if="!query.judgeTaskId" description="请选择任务后查看上报专家" />
    <div v-else class="card table-main">
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="primary" plain :icon="Download" :disabled="!query.judgeTaskId" @click="download">导出</el-button>
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle title="刷新" @click="search()" />
          <el-button :icon="Search" circle title="显隐搜索" @click="showFilter = !showFilter" />
        </div>
      </div>
      <el-table :data="pageRows" border height="100%" v-loading="loading" style="flex: 1" :fit="true">
        <el-table-column prop="name" label="姓名" width="90" />
        <el-table-column prop="sex" label="性别" width="60" align="center" />
        <el-table-column prop="phoneNo" label="手机号" width="115" />
        <el-table-column prop="deptName" label="所属学院" min-width="100" show-overflow-tooltip />
        <el-table-column label="职称" width="70" show-overflow-tooltip
          ><template #default="s">{{ titleText(s.row.title) }}</template></el-table-column
        >
        <el-table-column label="初始学历-所学专业" min-width="130" show-overflow-tooltip>
          <template #default="s">{{ eduText(s.row.initEdu) }}-{{ s.row.initMajor || "--" }}</template>
        </el-table-column>
        <el-table-column prop="goodSubjects" label="擅长科目" min-width="90" show-overflow-tooltip />
        <el-table-column label="评价" width="70" align="center">
          <template #default="s"
            ><el-button v-if="s.row.evaluation" link type="primary" @click="viewEvaluation(s.row)">查看</el-button
            ><span v-else>--</span></template
          >
        </el-table-column>
        <el-table-column label="评分方向" width="140" show-overflow-tooltip>
          <template #default="s">{{ `${s.row.directionCode || ""}${s.row.directionName || ""}` }}</template>
        </el-table-column>
        <el-table-column label="评分日期" min-width="150" align="center"
          ><template #default="s">{{ scoreDate(s.row) }}</template></el-table-column
        >
        <el-table-column label="状态" width="120" align="center"
          ><template #default="s">{{ auditText(s.row.auditState) }}</template></el-table-column
        >
        <el-table-column label="操作" width="120" align="center">
          <template #default="s">
            <template v-if="s.row.auditState === '1'">
              <el-button link type="primary" @click="confirmAudit(s.row, '2')">通过</el-button>
              <el-button link type="primary" @click="confirmAudit(s.row, '3')">不通过</el-button>
            </template>
            <el-button v-else link type="primary" @click="confirmReset(s.row)">重置</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page.current"
        v-model:page-size="page.size"
        :page-sizes="[10, 20, 50, 100]"
        :total="rows.length"
        layout="total, sizes, prev, pager, next, jumper"
      />
    </div>

    <el-dialog v-model="evaluationVisible" title="评价信息" width="560px" destroy-on-close>
      <div class="evaluation-name">专家姓名：{{ currentRow?.name }}</div>
      <div class="evaluation-content">{{ currentRow?.evaluation || "暂无评价信息" }}</div>
      <template #footer><el-button type="primary" @click="evaluationVisible = false">关闭</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="judgeReportAudit">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Download, Refresh, Search } from "@element-plus/icons-vue";
import { getDepartmentList } from "@/api/modules/department";
import {
  exportReportAudit,
  getReportAuditDirections,
  getReportAuditList,
  getReportAuditTasks,
  operateReportAudit,
  resetReportAudit
} from "@/api/modules/judgeExpert";
import { useDownload } from "@/hooks/useDownload";
import { useExpertTitleEnum } from "@/hooks/useEnum";
import type { JudgeExpert } from "@/api/interface/judgeExpert";

const tasks = ref<any[]>([]),
  departments = ref<any[]>([]),
  directions = ref<any[]>([]),
  rows = ref<any[]>([]);
const loading = ref(false),
  evaluationVisible = ref(false),
  currentRow = ref<any>(),
  showFilter = ref(true);
const page = reactive({ current: 1, size: 10 });
const query = reactive({ judgeTaskId: "", deptId: "", directionCode: "", name: "", sex: "", title: "", auditState: "" });
const { expertTitleEnum } = useExpertTitleEnum();
const pageRows = computed(() => rows.value.slice((page.current - 1) * page.size, page.current * page.size));
const titleText = (value: string) =>
  expertTitleEnum.value.find((x: any) => String(x.value) === String(value))?.label || value || "--";
const eduText = (value: string) => ({ "1": "专科", "2": "本科", "3": "硕士", "4": "博士" })[value] || "--";
const auditText = (value: string) => ({ "1": "已上报，待审核", "2": "审核通过", "3": "审核不通过" })[value] || "--";
const scoreDate = (row: any) => `${String(row.taskStartTime).slice(0, 10)} 至 ${String(row.taskEndTime).slice(0, 10)}`;
const getQueryParams = (): JudgeExpert.ReqReportAuditList => ({
  judgeTaskId: query.judgeTaskId,
  ...(query.deptId && { deptId: query.deptId }),
  ...(query.directionCode && { directionCode: query.directionCode }),
  ...(query.name && { name: query.name }),
  ...(query.sex && { sex: query.sex }),
  ...(query.title && { title: query.title }),
  ...(query.auditState && { auditState: query.auditState })
});

onMounted(async () => {
  const [taskRes, deptRes]: any[] = await Promise.all([getReportAuditTasks(), getDepartmentList({ curPage: 1, pageSize: 1000 })]);
  tasks.value = taskRes.list || [];
  departments.value = deptRes.list || [];
});
const onTaskChange = async () => {
  query.directionCode = "";
  rows.value = [];
  page.current = 1;
  if (!query.judgeTaskId) {
    directions.value = [];
    return;
  }
  try {
    const res: any = await getReportAuditDirections({ judgeTaskId: query.judgeTaskId });
    directions.value = res.list || [];
  } catch {
    directions.value = [];
  }
  await search().catch(() => null);
};
const search = async (resetPage = true) => {
  if (!query.judgeTaskId) {
    ElMessage.warning("请先选择任务");
    return;
  }
  loading.value = true;
  try {
    // 过滤空字符串，避免后端 int64 + json:",string" 标签用 strconv.ParseInt("") 解析失败导致"参数错误"
    const res: any = await getReportAuditList(getQueryParams());
    rows.value = res.list || [];
    // 仅筛选/切换任务时回到第一页；审核等行内操作保留当前页刷新（参考专家库管理审核行为）
    if (resetPage) page.current = 1;
  } finally {
    loading.value = false;
  }
};
const viewEvaluation = (row: any) => {
  currentRow.value = row;
  evaluationVisible.value = true;
};
const confirmAudit = async (row: any, state: "2" | "3") => {
  const action = state === "2" ? "参加" : "不参加";
  const result = await ElMessageBox.confirm(`【${row.name}】确认${action}【${row.directionName}】评分？`, "温馨提示", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning"
  }).catch(() => null);
  if (!result) return;
  await operateReportAudit({ reportId: row.reportId, auditState: state });
  ElMessage.success("审核成功");
  await search(false);
};
const confirmReset = async (row: any) => {
  const result = await ElMessageBox.confirm(`确认重置【${row.name}】【${row.directionName}】评分状态？`, "温馨提示", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning"
  }).catch(() => null);
  if (!result) return;
  await resetReportAudit({ reportId: row.reportId });
  ElMessage.success("重置成功");
  await search(false);
};
const download = () => {
  if (!query.judgeTaskId) {
    ElMessage.warning("请先选择任务");
    return;
  }
  ElMessageBox.confirm("确认导出当前评委上报审核数据？", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportReportAudit, "评委上报审核导出", getQueryParams())
  );
};
</script>

<style scoped lang="scss">
.report-audit {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.filter-bar {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 16px;
}
.table-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.report-audit :deep(.el-empty) {
  min-height: 460px;
}
.report-audit :deep(.el-table .cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-align: center;
  padding-left: 10px;
  padding-right: 10px;
}
.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.evaluation-name {
  margin-bottom: 12px;
  color: #606266;
}
.evaluation-content {
  min-height: 120px;
  padding: 16px;
  line-height: 1.8;
  white-space: pre-wrap;
  background: #f5f7fa;
  border-radius: 4px;
}
</style>
