<!-- 任务征集-评委征集任务管理-新增/修改Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="760"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="form">
      <el-row v-if="!isEditMode" :gutter="20">
        <el-col :span="12">
          <el-form-item label="任务年份" prop="taskYear">
            <el-date-picker v-model="form.taskYear" type="year" value-format="YYYY" placeholder="请选择任务年份" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="任务类型" prop="judgeTaskType">
            <el-select v-model="form!.judgeTaskType" placeholder="请选择类型" clearable>
              <el-option v-for="item in expertType" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="任务名称" prop="judgeTaskName">
        <el-input v-model="form!.judgeTaskName" placeholder="请输入任务名称" maxlength="50" show-word-limit clearable />
      </el-form-item>
      <el-form-item v-if="!isEditMode" label="任务模式" prop="taskMode">
        <el-radio-group v-model="form.taskMode">
          <el-radio value="1">直接抽取</el-radio><el-radio value="2">学院上报</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-row v-if="!isEditMode && form.taskMode === '1'" :gutter="20">
        <el-col :span="8">
          <el-form-item label="评分轮次" prop="scoreRounds">
            <el-input-number v-model="form.scoreRounds" :min="1" :precision="0" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="科目数量" prop="subjectCount">
            <el-input-number v-model="form.subjectCount" :min="1" :precision="0" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="每轮专家数" prop="expertsPerRound">
            <el-input-number v-model="form.expertsPerRound" :min="1" :precision="0" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="评分日期" prop="dates">
        <el-date-picker
          v-model="form!.dates"
          type="daterange"
          range-separator="To"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD HH:mm:ss"
        />
      </el-form-item>
      <el-form-item v-if="!isEditMode && form.taskMode === '1'" label="" prop="drawSubject">
        <el-checkbox label="按擅长科目抽取" v-model="form!.drawSubject" true-value="1" false-value="2" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AddOrEditDialog">
import { FormInstance, ElMessage } from "element-plus";
import { ref, reactive, computed } from "vue";
import { isValidLength } from "@/utils/eleValidate";
import { expertType } from "@/utils/dict";
import { addJudgeRecruitTask, editJudgeRecruitTask } from "@/api/modules/judgeExpert";

interface DialogProps {
  type: string;
  row: any;
  getTableList?: () => void;
}

const form = reactive<{
  judgeTaskName: string;
  judgeTaskType: string;
  dates: string[];
  drawSubject: string;
  taskYear: string;
  taskMode: string;
  scoreRounds: number;
  subjectCount: number;
  expertsPerRound: number;
}>({
  judgeTaskName: "",
  judgeTaskType: "",
  dates: [],
  drawSubject: "2",
  taskYear: String(new Date().getFullYear()),
  taskMode: "1",
  scoreRounds: 1,
  subjectCount: 1,
  expertsPerRound: 1
});

const rules = reactive({
  judgeTaskName: [
    { required: true, message: "请输入任务名称" },
    {
      validator: (...args: [any, any]) => isValidLength(1, 50, ...args),
      message: "任务名称必须为1到50个字符"
    }
  ],
  judgeTaskType: [{ required: true, message: "请选择类型" }],
  taskYear: [{ required: true, message: "请选择任务年份" }],
  taskMode: [{ required: true, message: "请选择任务模式" }],
  dates: [{ required: true, message: "请选择评分日期" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add", // 默认新增;
  row: {
    judgeTaskId: "",
    judgeTaskName: "",
    taskStartTime: "", // 开始日期
    taskEndTime: "", // 结束日期
    drawSubject: "2"
  }
});

const isEditMode = computed(() => dialogProps.value.type === "edit");

const dialogTitle = computed(() => (isEditMode.value ? "编辑任务" : "新增任务"));

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params = {
        judgeTaskId: isEditMode.value ? dialogProps.value.row.judgeTaskId : undefined,
        judgeTaskName: form.judgeTaskName,
        taskYear: form.taskYear,
        taskMode: form.taskMode,
        scoreRounds: String(form.scoreRounds),
        subjectCount: String(form.subjectCount),
        expertsPerRound: String(form.expertsPerRound),
        judgeTaskType: form.judgeTaskType,
        taskStartTime: form.dates[0],
        taskEndTime: form.dates[1],
        drawSubject: form.drawSubject
      };
      const api = isEditMode.value ? editJudgeRecruitTask : addJudgeRecruitTask;
      await api(params);
      ElMessage.success({ message: `${dialogTitle.value}成功！` });
      dialogProps.value.getTableList!();
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

const handleClosed = () => {
  form.judgeTaskName = "";
  form.judgeTaskType = "";
  form.dates = [];
  form.drawSubject = "2";
  form.taskYear = String(new Date().getFullYear());
  form.taskMode = "1";
  form.scoreRounds = 1;
  form.subjectCount = 1;
  form.expertsPerRound = 1;
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  if (isEditMode.value) {
    form.judgeTaskName = params.row.judgeTaskName;
    form.judgeTaskType = params.row.judgeTaskType;
    form.dates = [params.row.taskStartTime, params.row.taskEndTime];
    form.drawSubject = params.row.drawSubject || "2"; // 默认值为 "2"
  }
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
:deep(.el-input) {
  width: 100%;
}
</style>
