<!-- 监考老师管理-监考老师征集任务-新增（任务+日期选择）Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="400"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="form">
      <el-form-item label="考试任务" prop="tchTaskName">
        <el-input v-model="form!.tchTaskName" placeholder="请输入内容" maxlength="30" clearable />
      </el-form-item>
      <el-form-item label="征集日期" prop="dates">
        <el-date-picker
          v-model="form!.dates"
          type="daterange"
          range-separator="To"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD HH:mm:ss"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="PreAddDialog">
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { isValidLength } from "@/utils/eleValidate";
import { checkTeacherRecruitTaskName } from "@/api/modules/teacher";

interface DialogProps {
  title: string;
  success?: (params: { tchTaskName: string; dates: string[] }) => void;
}

const form = reactive({
  tchTaskName: "",
  dates: []
});

const rules = reactive({
  tchTaskName: [
    { required: true, message: "请输入考试任务" },
    {
      validator: (...args: [any, any]) => isValidLength(2, 30, ...args),
      message: "考试任务名称必须为2到30位"
    }
  ],
  dates: [{ required: true, message: "请选择征集日期" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: ""
});

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    const isWithin7Days = compareDatesWithin7Days(form.dates[0], form.dates[1]);
    if (!isWithin7Days) {
      ElMessage.error({ message: "征集日期不能超过7天" });
      return;
    }
    try {
      const params = {
        tchTaskName: form.tchTaskName,
        taskStartTime: form.dates[0],
        taskEndTime: form.dates[1]
      };
      await checkTeacherRecruitTaskName(params);
      dialogProps.value.success!({ tchTaskName: form.tchTaskName, dates: form.dates });
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

// 是否小于等于 7 天
const compareDatesWithin7Days = (dateStr1, dateStr2) => {
  if (!dateStr1 || !dateStr2) return false;
  // 获取两个日期的时间戳
  const timestamp1 = new Date(dateStr1).getTime();
  const timestamp2 = new Date(dateStr2).getTime();
  const timeDiff = Math.abs(timestamp1 - timestamp2); // 计算时间差的绝对值，单位为毫秒
  const dayDiff = timeDiff / (1000 * 60 * 60 * 24);
  // 判断天数差是否小于等于 7-1 天
  return dayDiff <= 6;
};

const handleClosed = () => {
  form.tchTaskName = "";
  form.dates = [];
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
:deep(.el-input) {
  width: 100%;
}
</style>
