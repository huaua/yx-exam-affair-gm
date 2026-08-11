<!-- 基础管理-专家评价管理-新增/编辑Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="600"
  >
    <el-form ref="formRef" label-width="110px" label-suffix=" :" :rules="rules" :model="dialogProps.row" inline>
      <el-form-item label="评分任务" prop="judgeTaskId">
        <el-select v-model="dialogProps.row!.judgeTaskId" placeholder="请选择评分任务" :disabled="isEditMode" clearable>
          <el-option
            v-for="item in dialogProps.canEvaluateJudgeTaskEnum"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="相关评价" prop="appraisalContent" class="w-full">
        <el-input
          v-model="dialogProps.row!.appraisalContent"
          type="textarea"
          placeholder="请填写相关评价"
          maxlength="200"
          :autosize="{ minRows: 4, maxRows: 8 }"
          show-word-limit
          clearable
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AddOrEditDialog">
import { JudgeExpert } from "@/api/interface";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive, computed } from "vue";
import { addExpertEvaluation, editExpertEvaluation } from "@/api/modules/judgeExpert";

interface DialogProps {
  type: string;
  expertId: string; // 专家ID
  canEvaluateJudgeTaskEnum: any[]; // 可评价的评分任务列表
  row: Partial<JudgeExpert.ResExpertEvaluationList>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  expertId: "",
  canEvaluateJudgeTaskEnum: [],
  row: {}
});

// 表单校验规则
const rules = reactive({
  judgeTaskId: [{ required: true, message: "请选择评分任务" }],
  appraisalContent: [{ required: true, message: "请输入相关评价" }]
});

const isEditMode = computed(() => dialogProps.value.type === "edit");

const dialogTitle = computed(() => (isEditMode.value ? "编辑评价" : "新增评价"));

// 提交数据（新增/编辑）
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params: any = {
        appraisalId: isEditMode.value ? dialogProps.value.row?.appraisalId : undefined,
        expertId: !isEditMode.value ? dialogProps.value?.expertId : undefined,
        judgeTaskId: !isEditMode.value ? dialogProps.value.row?.judgeTaskId : undefined,
        appraisalContent: dialogProps.value.row?.appraisalContent || ""
      };
      const api = isEditMode.value ? editExpertEvaluation : addExpertEvaluation;
      await api(params);
      ElMessage.success({ message: `${dialogTitle.value}成功！` });
      dialogProps.value.getTableList!();
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
:deep(.el-form-item__content) {
  min-width: 230px;
}
.w-full {
  display: flex;
}
</style>
