<!-- 考场编排-导入-考场征集任务选择Dialog -->
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
    <el-form ref="formRef" label-width="118px" label-suffix=" :" :rules="rules" :model="form">
      <el-form-item label="考场征集任务" prop="selectedRoomTaskObj">
        <el-select v-model="form!.selectedRoomTaskObj" value-key="value" placeholder="请选择任务" clearable>
          <el-option v-for="item in allRecruitTaskEnum" :key="item.value" :label="item.label" :value="item" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="PreImportDialog">
import { FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { TaskEnumItem, useAllRecruitTaskEnum } from "@/hooks/useCustomEnum";

interface DialogProps {
  title: string;
  success?: (params: { selectedRoomTaskObj: TaskEnumItem }) => void;
}

const form = reactive<{ selectedRoomTaskObj: TaskEnumItem | undefined }>({
  selectedRoomTaskObj: undefined
});

const rules = reactive({
  selectedRoomTaskObj: [{ required: true, message: "请选择考场征集任务" }]
});

const { allRecruitTaskEnum } = useAllRecruitTaskEnum();
const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: ""
});

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      dialogProps.value.success!({ selectedRoomTaskObj: form.selectedRoomTaskObj! });
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

const handleClosed = () => {
  form.selectedRoomTaskObj = undefined;
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
