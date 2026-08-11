<!-- 考生管理-考生管理-同步任务选择Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    title="信息同步"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="500"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="form" label-position="top">
      <el-form-item label="选择需要同步的考试" prop="kaoShiObj">
        <div class="item-content">
          <el-select v-model="form!.kaoShiObj" placeholder="请选择考试" clearable value-key="kaoShiID">
            <el-option v-for="item in examEnum" :key="item.kaoShiID" :label="item.kaoShiMC" :value="item" />
          </el-select>
          <el-button class="item-content-btn" type="primary" @click="handleReSync">考试信息同步</el-button>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="PreSyncDialog">
import { FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { StudentSyncEnum } from "@/enums";
import { getExamList } from "@/api/modules/student";

interface DialogProps {
  openSyncDialog?: (params: { kaoShiID: string; kaoShiMC: string; state: string; percent: string; errMsg: string }) => void;
}

const examEnum = ref<{ kaoShiID: string; kaoShiMC: string }[]>([]);

const form = reactive({
  kaoShiObj: undefined as { kaoShiID: string; kaoShiMC: string } | undefined
});

const rules = reactive({
  kaoShiObj: [{ required: true, message: "请选择考试" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({});

// 表单引用
const formRef = ref<FormInstance>();
// 提交表单
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    // 打开同步弹窗
    dialogProps.value?.openSyncDialog?.({
      kaoShiID: form.kaoShiObj!.kaoShiID,
      kaoShiMC: form.kaoShiObj!.kaoShiMC,
      state: StudentSyncEnum.NOT_STARTED,
      percent: "0",
      errMsg: ""
    });
    dialogVisible.value = false;
  });
};

// 获取考试枚举
const fetchExamEnum = async (reSync = false) => {
  const res = await getExamList({ reSync });
  examEnum.value = (res.list || []).map(item => ({
    kaoShiID: String(item.kaoShiID),
    kaoShiMC: item.kaoShiMC
  }));
};

// 考试信息重新同步
const handleReSync = () => {
  fetchExamEnum(true);
  handleClosed();
};

// 弹窗关闭时重置表单
const handleClosed = () => {
  form.kaoShiObj = undefined;
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  fetchExamEnum(false);
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.item-content {
  display: flex;
  width: 100%;
  &-btn {
    margin-left: 16px;
  }
}
</style>
