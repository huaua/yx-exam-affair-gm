<!-- 任务征集-考场征集管理-详情Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="430"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="dialogProps.row">
      <el-form-item label="考试任务" prop="taskName">
        <el-input v-model="dialogProps.row!.taskName" placeholder="请输入内容" maxlength="30" clearable />
      </el-form-item>
      <el-form-item label="征集日期" prop="taskDay">
        <el-date-picker v-model="dialogProps.row!.taskDay" type="date" placeholder="请选择日期" />
      </el-form-item>
      <el-form-item label="所属层次" prop="levelCode">
        <el-select v-model="dialogProps.row!.levelCode" placeholder="请选择所属层次" clearable>
          <el-option v-for="item in levelEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="考场数量" prop="needNum">
        <InputInteger v-model="dialogProps.row!.needNum" placeholder="请输入数量" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AddOrEditDialog">
import { Classroom } from "@/api/interface";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import dayjs from "dayjs";
import { isValidLength } from "@/utils/eleValidate";
import { useLevelEnum } from "@/hooks/useEnum";

interface DialogProps {
  type: string;
  title: string;
  row: Partial<Classroom.ResRecruitManageList>;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

const { levelEnum } = useLevelEnum();

const rules = reactive({
  taskName: [
    { required: true, message: "请输入考试任务" },
    {
      validator: (...args: [any, any]) => isValidLength(2, 30, ...args),
      message: "考试任务名称必须为2到30位"
    }
  ],
  taskDay: [{ required: true, message: "请选择征集日期" }],
  levelCode: [{ required: true, message: "请选择所属层次" }],
  needNum: [{ required: true, message: "请输入考场数量" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  title: "",
  row: {}
});

// 提交数据（新增/编辑）
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params = {
        taskId: dialogProps.value.row?.taskId || undefined,
        taskName: dialogProps.value.row?.taskName,
        taskDay: dayjs(dialogProps.value.row?.taskDay).format("YYYY-MM-DD HH:mm:ss"),
        levelCode: dialogProps.value.row?.levelCode,
        needNum: String(dialogProps.value.row?.needNum)
      };
      await dialogProps.value.api!(params);
      ElMessage.success({ message: `${dialogProps.value.title}成功！` });
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
:deep(.el-form-item) {
  width: 358px;
}
:deep(.el-date-editor--date) {
  width: 258px;
}
</style>
