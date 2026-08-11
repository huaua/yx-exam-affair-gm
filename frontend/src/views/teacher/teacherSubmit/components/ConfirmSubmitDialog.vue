<!-- 监考老师管理-监考老师上报-上报Dialog-确认上报Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="380"
  >
    <el-form ref="formRef" label-width="60px" label-suffix=" :" :rules="rules" :model="form">
      <el-form-item label="姓名">
        <el-input v-model="dialogProps.row.tchName" disabled />
      </el-form-item>
      <el-form-item label="日期" prop="selectedDateList">
        <el-checkbox-group v-model="form!.selectedDateList">
          <el-checkbox v-for="item in canSelectDateList" :key="item.taskTime" :value="item.taskTime" checked>
            {{ dayjs(item.taskTime).format("MM-DD") }}
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="remark" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="ConfirmSubmitDialog">
import { Teacher } from "@/api/interface";
import { FormInstance } from "element-plus";
import { ref, reactive, computed } from "vue";
import dayjs from "dayjs";

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResCanSubmitTeacherList>;
  dateList: Teacher.ResCurTaskDateList[];
  success?: (params: any, remark: any) => void;
}

const form = reactive<{ selectedDateList: string[] }>({
  selectedDateList: []
});
const remark = ref("");

const rules = reactive({
  selectedDateList: [{ required: true, message: "请选择日期" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {},
  dateList: []
});

const canSelectDateList = computed(() => {
  const usedDateList = dialogProps.value.row.reportTchList || [];
  return dialogProps.value.dateList.filter(item => {
    return !usedDateList.includes(item.taskTime);
  });
});

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      dialogProps.value.success!(form.selectedDateList, remark.value);
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

const handleClosed = () => {
  form.selectedDateList = [];
  remark.value = "";
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
