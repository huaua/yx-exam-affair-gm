<!-- 学院管理-学院管理-详情Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}学院`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="880"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :inline="true" :model="dialogProps.row">
      <DynamicInput
        v-for="(item, index) in dialogProps.row"
        :key="index"
        :current-index="index"
        :disabled="isEditMode"
        :length="dialogProps.row.length"
        @on-plus="handlePlus"
        @on-minus="handleMinus"
      >
        <el-form-item label="学院代码" :prop="`${index}.deptCode`" :rules="rules.deptCode">
          <el-input
            v-model="item.deptCode"
            placeholder="请输入内容"
            :disabled="isEditMode"
            minlength="1"
            maxlength="15"
            arable
          ></el-input>
        </el-form-item>
        <el-form-item label="学院名称" :prop="`${index}.deptName`" :rules="rules.deptName">
          <el-input v-model="item.deptName" placeholder="请输入内容" minlength="2" maxlength="15" clearable></el-input>
        </el-form-item>
        <el-form-item label="所在校区" :prop="`${index}.campusName`">
          <el-input v-model="item.campusName" placeholder="请输入内容" minlength="2" maxlength="15" clearable></el-input>
        </el-form-item>
      </DynamicInput>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="DepartmentDialog">
import { Department } from "@/api/interface";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive, computed } from "vue";
import DynamicInput from "@/components/DynamicInput.vue";

type FormItem = Partial<Department.ResDepartmentList>;

interface DialogProps {
  title: string;
  row: FormItem[];
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

const rules = reactive({
  deptCode: [{ required: true, message: "请填写学院代码" }],
  deptName: [{ required: true, message: "请填写学院名称" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: [{}]
});

// 是否为编辑模式
const isEditMode = computed(() => dialogProps.value.title === "编辑");

// 提交数据（新增/编辑）
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      // 新增
      let params: FormItem | FormItem[] = dialogProps.value.row;
      // 编辑
      if (isEditMode.value) {
        params = {
          deptId: dialogProps.value.row[0]?.deptId,
          deptName: dialogProps.value.row[0]?.deptName,
          campusName: dialogProps.value.row[0]?.campusName
        };
      }
      await dialogProps.value.api!(params);
      ElMessage.success({ message: `${dialogProps.value.title}成功！` });
      dialogProps.value.getTableList!();
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

const handlePlus = () => {
  dialogProps.value.row.push({ deptCode: "", deptName: "" });
};

const handleMinus = (index: number) => {
  dialogProps.value.row.splice(index, 1);
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
  width: 400px;
}
</style>
