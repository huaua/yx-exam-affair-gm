<!-- 基础管理-教室管理-详情Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="500"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="dialogProps.row">
      <el-form-item label="教室名称" prop="roomName">
        <el-input v-model="dialogProps.row!.roomName" placeholder="请输入内容" maxlength="20" clearable />
      </el-form-item>
      <el-form-item label="所属校区" prop="campusName">
        <el-input v-model="dialogProps.row!.campusName" placeholder="请输入内容" maxlength="15" clearable />
      </el-form-item>
      <el-form-item label="层次" prop="levelCode">
        <el-select v-model="dialogProps.row!.levelCode" placeholder="请选择层次" :disabled="isDisabledLevelField" clearable>
          <el-option v-for="item in levelEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="组数" prop="groupNum">
        <InputInteger v-model="dialogProps.row!.groupNum" placeholder="请输入" />
      </el-form-item>
      <el-form-item v-if="isShowGroupCapacityField" label="按组容量" prop="groupCapacity">
        <InputInteger v-model="dialogProps.row!.groupCapacity" placeholder="请输入" />
      </el-form-item>
      <el-form-item label="按位容量" prop="capacity">
        <InputInteger v-model="dialogProps.row!.capacity" placeholder="请输入" />
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
import { ref, reactive, computed, watch } from "vue";
import { useLevelEnum } from "@/hooks/useEnum";

interface DialogProps {
  type: string;
  title: string;
  row: Partial<Classroom.ResClassroomList>;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

const { levelEnum } = useLevelEnum();

const rules = reactive({
  roomName: [{ required: true, message: "请输入教室名称" }],
  levelCode: [{ required: true, message: "请选择层次" }],
  groupCapacity: [{ required: true, message: "请输入按组容量" }]
});

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  title: "",
  row: {}
});

const isShowGroupCapacityField = computed(() => {
  const groupNum = dialogProps.value.row?.groupNum;
  return groupNum && groupNum > 0;
});
const isDisabledLevelField = computed(() => dialogProps.value.type === "edit");

// 清空组数时，清空按组容量
watch(
  () => dialogProps.value.row?.groupNum,
  newVal => {
    if (!newVal) {
      dialogProps.value.row!.groupCapacity = undefined;
    }
  }
);

// 提交数据（新增/编辑）
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  if (!dialogProps.value.row?.groupNum && !dialogProps.value.row?.capacity) {
    ElMessage.error({ message: "未设置容量" });
    return;
  }
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params = {
        roomId: dialogProps.value.row?.roomId || undefined,
        roomName: dialogProps.value.row?.roomName,
        campusName: dialogProps.value.row?.campusName,
        levelCode: dialogProps.value.row?.levelCode,
        groupNum: dialogProps.value.row?.groupNum == null ? undefined : String(dialogProps.value.row?.groupNum),
        groupCapacity: dialogProps.value.row?.groupCapacity == null ? undefined : String(dialogProps.value.row?.groupCapacity),
        capacity: dialogProps.value.row?.capacity == null ? undefined : String(dialogProps.value.row?.capacity)
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
