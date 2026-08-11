<!-- 任务征集-评委征集任务管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getJudgeRecruitManageList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openAddOrEditDialog('add')">新增任务</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <!-- 编辑：学院上报-未发布 / 直接抽取-未开始或抽取中 -->
        <el-button
          v-if="
            (scope.row.taskMode === '2' && scope.row.publishState === '1') ||
            (scope.row.taskMode === '1' &&
              (scope.row.state === JudgeTaskState.NOT_STARTED || scope.row.state === JudgeTaskState.EXTRACTING))
          "
          type="primary"
          link
          @click="openAddOrEditDialog('edit', scope.row)"
        >
          编辑
        </el-button>
        <!-- 发布：学院上报+未发布 -->
        <el-button
          v-if="scope.row.taskMode === '2' && scope.row.publishState === '1'"
          type="primary"
          link
          @click="onPublish(scope.row)"
        >
          发布
        </el-button>
        <!-- 追加轮次/添加轮次：直接抽取模式+未开始/已抽取 -->
        <el-button
          v-if="
            scope.row.taskMode === '1' &&
            (scope.row.state === JudgeTaskState.NOT_STARTED || scope.row.state === JudgeTaskState.EXTRACTING)
          "
          type="primary"
          link
          @click="onAppendRounds(scope.row)"
        >
          {{ scope.row.state === JudgeTaskState.EXTRACTING ? "添加轮次" : "追加轮次" }}
        </el-button>
        <!-- 取消发布：学院上报+已发布 / 直接抽取+抽取中 -->
        <el-button
          v-if="
            (scope.row.taskMode === '2' && scope.row.publishState === '2') ||
            (scope.row.taskMode === '1' && scope.row.state === JudgeTaskState.EXTRACTING)
          "
          type="primary"
          link
          @click="onUnpublish(scope.row)"
        >
          取消发布
        </el-button>
        <!-- 上报导入：学院上报模式（两种状态均可） -->
        <el-button v-if="scope.row.taskMode === '2'" type="primary" link @click="openRequirementImport(scope.row)">
          上报导入
        </el-button>
        <!-- 取消导入：已发布+学院上报+非招办角色 -->
        <el-button
          v-if="scope.row.taskMode === '2' && scope.row.publishState === '2' && !isAdminRole"
          type="primary"
          link
          @click="onCancelImport(scope.row)"
        >
          取消导入
        </el-button>
        <!-- 查看：学院上报（两种状态均可）/ 直接抽取-抽取中或已完成 -->
        <el-button
          v-if="
            scope.row.taskMode === '2' ||
            scope.row.state === JudgeTaskState.EXTRACTING ||
            scope.row.state === JudgeTaskState.COMPLETED
          "
          type="primary"
          link
          @click="openViewDialog(scope.row)"
        >
          查看
        </el-button>
        <!-- 完成征集：已抽取状态 -->
        <el-button v-if="scope.row.state === JudgeTaskState.EXTRACTING" type="danger" link @click="onCancleClick(scope.row)">
          完成征集
        </el-button>
        <!-- 删除：见 canDelete 逻辑 -->
        <el-button v-if="canDelete(scope.row)" type="danger" link @click="onDeleteClick(scope.row)"> 删除 </el-button>
      </template>
    </ProTable>
    <AddOrEditDialog ref="addOrEditDialogRef" />
    <ViewDialog ref="viewDialogRef" />
    <RequirementDialog ref="requirementDialogRef" />
    <ImportExcel ref="importExcelRef" />
  </div>
</template>

<script setup lang="tsx" name="judgeRecruitManage">
import { JudgeExpert } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { JudgeTaskState } from "@/enums";
import { ref, reactive } from "vue";
import { useRole } from "@/hooks/useRole";
import { ElMessageBox } from "element-plus";
import { CirclePlus } from "@element-plus/icons-vue";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import AddOrEditDialog from "./components/AddOrEditDialog.vue";
import ViewDialog from "./components/ViewDialog.vue";
import RequirementDialog from "./components/RequirementDialog.vue";
import ImportExcel from "@/components/ImportExcel/index.vue";
import { judgeDrawStatus, expertType } from "@/utils/dict";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";
import {
  getJudgeRecruitManageList,
  cancelJudgeRecruitTask,
  deleteJudgeRecruitTask,
  appendJudgeTaskRounds,
  publishJudgeRecruitTask,
  importJudgeTaskRequirements
} from "@/api/modules/judgeExpert";

const { isAdminRole } = useRole();
const proTable = ref<ProTableInstance>();

// 表格配置项（按图二调整：任务年份|名称|模式|类型|评分时间|轮次/科目数|抽取数|上报数|状态|操作）
const columns = reactive<ColumnProps<JudgeExpert.ResJudgeRecruitManageList>[]>([
  {
    prop: "taskYear",
    label: "任务年份",
    width: 90,
    search: { order: 0, el: "date-picker", props: { type: "year", valueFormat: "YYYY", placeholder: "请选择年份" } }
  },
  {
    prop: "judgeTaskName",
    label: "任务名称",
    minWidth: 140,
    search: { order: 1, el: "input", props: { placeholder: "请输入" } }
  },
  {
    prop: "taskMode",
    label: "任务模式",
    enum: [
      { label: "直接抽取", value: "1" },
      { label: "学院上报", value: "2" }
    ],
    search: { order: 5, el: "select", props: { placeholder: "请选择" } },
    width: 90
  },
  {
    prop: "judgeTaskType",
    label: "任务类型",
    enum: expertType,
    search: { order: 2, el: "select", props: { placeholder: "请选择" } },
    width: 90
  },
  {
    prop: "_scoreTime",
    label: "评分时间",
    width: 180,
    render: scope => (
      <div>
        {formatTime(scope.row.taskStartTime) || "--"}至{formatTime(scope.row.taskEndTime) || "--"}
      </div>
    )
  },
  {
    prop: "_roundsSubjects",
    label: "评分轮次/科目数量",
    width: 180,
    render: scope => {
      if (scope.row.taskMode === "1") {
        return (
          <div>
            {scope.row.scoreRounds || "-"}/{scope.row.subjectCount || "-"}
          </div>
        );
      }
      return <div>-</div>;
    }
  },
  {
    prop: "requirementCount",
    label: "需求数",
    width: 80,
    render: scope => <div>{scope.row.requirementCount ?? "--"}</div>
  },
  {
    prop: "reportCount",
    label: "上报数",
    width: 80,
    render: scope => {
      // 学院上报：上报数；直接抽取：已抽取评委数
      if (scope.row.taskMode === "2") {
        return <div>{scope.row.reportCount || "--"}</div>;
      }
      return <div>{scope.row.drawCount || "--"}</div>;
    }
  },
  {
    prop: "state",
    label: "任务状态",
    enum: judgeDrawStatus,
    width: 90,
    render: scope => {
      // 学院上报模式：按发布状态展示 未发布/已发布
      if (scope.row.taskMode === "2") {
        return (
          <el-tag type={scope.row.publishState === "2" ? "success" : "warning"}>
            {scope.row.publishState === "2" ? "已发布" : "未发布"}
          </el-tag>
        );
      }
      // 直接抽取模式：存在已抽取评委数据则显示 已抽取，否则 未开始
      const drawn = Number(scope.row.drawCount || 0);
      return <el-tag type={drawn > 0 ? "warning" : "info"}>{drawn > 0 ? "已抽取" : "未开始"}</el-tag>;
    }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 追加轮次 / 添加轮次：点取消/关闭时静默处理，避免 reject 冒泡触发"未知错误"
const onAppendRounds = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  const result = await ElMessageBox.prompt("请输入需要追加的评分轮次数量", "追加轮次", {
    inputPattern: /^[1-9]\d*$/,
    inputErrorMessage: "请输入大于0的正整数"
  }).catch(() => null);
  if (!result) return;
  await appendJudgeTaskRounds({ judgeTaskId: row.judgeTaskId, scoreRounds: result.value });
  proTable.value?.getTableList();
};

const onPublish = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  await useHandleData(publishJudgeRecruitTask, { judgeTaskId: row.judgeTaskId }, `发布【${row.judgeTaskName}】`);
  proTable.value?.getTableList();
};

// 取消发布
const onUnpublish = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  await useHandleData(
    () => http.post(`${PROXY_TAG}/api/auth/ea_judge_task_info/unpublish`, { judgeTaskId: row.judgeTaskId }),
    {},
    `取消发布【${row.judgeTaskName}】`
  );
  proTable.value?.getTableList();
};

// 取消导入
const onCancelImport = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  await useHandleData(
    () => http.post(`${PROXY_TAG}/api/auth/ea_judge_task_info/cancel_import`, { judgeTaskId: row.judgeTaskId }),
    {},
    `取消导入【${row.judgeTaskName}】`
  );
  proTable.value?.getTableList();
};

const importExcelRef = ref<InstanceType<typeof ImportExcel> | null>(null);
const openRequirementImport = (row: JudgeExpert.ResJudgeRecruitManageList) => {
  importExcelRef.value?.acceptParams({
    title: "学院评分方向专家数量导入",
    templateUrl: "static/学院评分方向专家数量导入模板.xlsx",
    saveUploadApi: importJudgeTaskRequirements,
    extraData: { judgeTaskId: row.judgeTaskId },
    getTableList: proTable.value?.getTableList
  });
};

// 打开-新增/编辑
const addOrEditDialogRef = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openAddOrEditDialog = (type: string, row: any = {}) => {
  const params = {
    type,
    row,
    getTableList: proTable.value?.getTableList
  };
  addOrEditDialogRef.value?.acceptParams(params);
};

// 打开-查看Dialog
const viewDialogRef = ref<InstanceType<typeof ViewDialog> | null>(null);
const requirementDialogRef = ref<InstanceType<typeof RequirementDialog> | null>(null);
const openViewDialog = (row: any = {}) => {
  if (row.taskMode === "2") {
    requirementDialogRef.value?.acceptParams({ row, getTableList: proTable.value?.getTableList });
    return;
  }
  const params = {
    title: "查看",
    judgeTaskId: row.judgeTaskId || ""
  };
  viewDialogRef.value?.acceptParams(params);
};

// 取消评委征集任务
const onCancleClick = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  await useHandleData(
    cancelJudgeRecruitTask,
    {
      judgeTaskId: row.judgeTaskId
    },
    `取消【${row.judgeTaskName}】`
  );
  proTable.value?.getTableList();
};

// 删除评委征集任务
const onDeleteClick = async (row: JudgeExpert.ResJudgeRecruitManageList) => {
  await useHandleData(
    deleteJudgeRecruitTask,
    {
      judgeTaskId: row.judgeTaskId
    },
    `删除【${row.judgeTaskName}】`
  );
  proTable.value?.getTableList();
};

const formatTime = fullTimeStr => {
  if (!fullTimeStr) {
    return "";
  }
  return fullTimeStr.slice(0, 10);
};

// 删除按钮显示逻辑：
// 学院上报：未发布 且 无已上报专家（reportCount=0）才允许删除
// 直接抽取：无已抽取评委（onCampusAttendNum+offCampusAttendNum=0）才允许删除
const canDelete = (row: JudgeExpert.ResJudgeRecruitManageList) => {
  if (row.taskMode === "2") {
    return row.publishState === "1" && Number(row.reportCount || 0) === 0;
  }
  return Number(row.drawCount || 0) === 0;
};
</script>
