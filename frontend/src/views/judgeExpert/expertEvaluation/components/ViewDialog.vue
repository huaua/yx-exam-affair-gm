<!-- 基础管理-专家评价管理-查看Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="860"
  >
    <div class="info">
      <div class="info-item">
        <div class="label">姓名：</div>
        <div class="content">{{ dialogProps.row?.name }}</div>
      </div>
      <div class="info-item">
        <div class="label">身份证：</div>
        <div class="content">{{ dialogProps.row?.idCard }}</div>
      </div>
      <div class="info-item">
        <div class="label">手机号：</div>
        <div class="content">{{ dialogProps.row?.phoneNo }}</div>
      </div>
      <div class="info-item">
        <div class="label">类型：</div>
        <div class="content">{{ mapValueToLabel(dialogProps.row?.type, expertType) }}</div>
      </div>
      <div class="info-item">
        <div class="label">初始学历：</div>
        <div class="content">{{ mapValueToLabel(dialogProps.row?.initEdu, educationType) }}</div>
      </div>
      <div class="info-item">
        <div class="label">最终学历：</div>
        <div class="content">{{ mapValueToLabel(dialogProps.row?.finalEdu, educationType) }}</div>
      </div>
      <div class="info-item">
        <div class="label">擅长科目：</div>
        <div class="content">{{ dialogProps.row?.goodSubjects }}</div>
      </div>
      <div class="info-item">
        <div class="label">参加评分次数：</div>
        <div class="content">{{ judgeTimes }}</div>
      </div>
    </div>
    <div class="operation">
      <el-button type="primary" :icon="CirclePlus" size="small" @click="openAddOrEditDialog('add')"> 新增 </el-button>
    </div>
    <el-table class="table" v-if="tableData.length > 0" :data="tableData" height="300" style="width: 100%">
      <el-table-column prop="judgeTaskName" label="任务名称" align="center" width="150" />
      <el-table-column prop="taskTimeRemark" label="评分日期" align="center" width="120" />
      <el-table-column prop="appraisalContent" label="评价内容" align="center" />
      <el-table-column label="操作" align="center" width="150">
        <template #default="scope">
          <el-button type="primary" link :icon="Edit" @click="openAddOrEditDialog('edit', scope.row)"> 编辑 </el-button>
          <el-button type="danger" link :icon="Delete" @click="onDeleteClick(scope.row)"> 删除 </el-button>
        </template>
      </el-table-column>
    </el-table>
    <AddOrEditDialog ref="addOrEditDialogRef" />
  </el-dialog>
</template>

<script setup lang="ts" name="ViewDialog">
import { JudgeExpert } from "@/api/interface";
import { ref, computed } from "vue";
import { CirclePlus, Edit, Delete } from "@element-plus/icons-vue";
import { useHandleData } from "@/hooks/useHandleData";
import AddOrEditDialog from "./AddOrEditDialog.vue";
import { educationType, expertType } from "@/utils/dict";
import { getExpertEvaluationList, deleteExpertEvaluation, getCanEvaluateJudgeTaskList } from "@/api/modules/judgeExpert";

interface DialogProps {
  title: string;
  row: any;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "查看",
  row: {
    expertId: "",
    name: "",
    idCard: "",
    phoneNo: "",
    type: "",
    initEdu: "",
    finalEdu: "",
    goodSubjects: "" // todo 参加评分次数
  }
});

const canEvaluateJudgeTaskEnum = ref<any>([]);

const tableData = ref<JudgeExpert.ResExpertEvaluationList[]>([]);

const judgeTimes = computed(() => canEvaluateJudgeTaskEnum.value.length);

// 打开-新增/编辑
const addOrEditDialogRef = ref<InstanceType<typeof AddOrEditDialog> | null>(null);
const openAddOrEditDialog = (type: string, row: any = {}) => {
  const params = {
    type,
    expertId: dialogProps.value.row?.expertId || "",
    canEvaluateJudgeTaskEnum: canEvaluateJudgeTaskEnum.value,
    row: { ...row },
    getTableList: getTableData
  };
  addOrEditDialogRef.value?.acceptParams(params);
};

// 删除评价
const onDeleteClick = async (row: JudgeExpert.ResExpertEvaluationList) => {
  try {
    await useHandleData(
      deleteExpertEvaluation,
      {
        appraisalId: row.appraisalId
      },
      `删除评价`
    );
    getTableData();
  } catch (error) {
    console.log(error);
  }
};

const getTableData = async () => {
  try {
    const params = {
      expertId: dialogProps.value.row?.expertId
    };
    const res = await getExpertEvaluationList(params);
    tableData.value = res.list || [];
  } catch (error) {
    console.log(error);
  }
};

// 获取可评价任务列表
const _getCanEvaluateJudgeTaskList = async () => {
  try {
    const params = {
      expertId: dialogProps.value.row?.expertId
    };
    const res = await getCanEvaluateJudgeTaskList(params);
    const list = res?.list || [];
    canEvaluateJudgeTaskEnum.value = list.map(item => ({
      value: item.judgeTaskId,
      label: item.judgeTaskName
    }));
  } catch (error) {
    console.error("获取评分任务列表失败", error);
  }
};

const mapValueToLabel = (value, dictionary) => {
  const item = dictionary.find(entry => entry.value === value);
  return item ? item.label : "未知";
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  getTableData();
  _getCanEvaluateJudgeTaskList();
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
:deep(.el-input) {
  width: 100%;
}
.info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  grid-gap: 16px;
  padding: 20px;
  margin-bottom: 16px;
  &-item {
    display: flex;
    align-items: center;
    .label {
      width: 100px;
      font-weight: bold;
    }
    .content {
      flex: 1;
    }
  }
}
.operation {
  padding: 10px 0;
}
.table {
  margin-bottom: 20px;
}
</style>
