<!-- 监考老师管理-监考老师上报-新增（任务+日期选择）Dialog -->
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
        <el-input v-model="dialogProps.row.tchTaskName" disabled />
      </el-form-item>
      <el-form-item label="征集日期" prop="selectedDateObj">
        <el-select v-model="form!.selectedDateObj" value-key="taskTime" placeholder="请选择征集日期" clearable>
          <el-option v-for="item in dateList" :key="item.taskTime" :value="item" :label="item.taskTime.slice(0, 10)" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="PreAppendDialog">
import { Teacher } from "@/api/interface";
import { FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { getCurTaskDateList } from "@/api/modules/teacher";

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResTeacherRecruitManageList>;
  success?: (params: { selectedDateObj: Teacher.ResCurTaskDateList }) => void;
}

interface Form {
  selectedDateObj: Teacher.ResCurTaskDateList | undefined;
}

const dateList = ref<Teacher.ResCurTaskDateList[]>([]);

const form = reactive<Form>({
  selectedDateObj: undefined
});

const rules = reactive({
  selectedDateObj: [{ required: true, message: "请选择征集日期" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {}
});

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      dialogProps.value.success!({ selectedDateObj: form.selectedDateObj! });
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

const handleClosed = () => {
  form.selectedDateObj = undefined;
};

// 获取日期list
const getDateList = async () => {
  const params = {
    tchTaskId: dialogProps.value.row.tchTaskId!,
    onlyQueryNeedReportTimeList: true
  };
  const { list = [] } = await getCurTaskDateList(params);
  return list || [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  dateList.value = await getDateList();
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
