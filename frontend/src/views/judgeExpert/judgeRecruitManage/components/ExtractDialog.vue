<!-- 任务征集-评委征集任务管理-抽取Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="800"
  >
    <div class="times" v-if="isNextMode">已抽取次数：{{ extractTimes }}次</div>
    <el-form ref="formRef" label-width="auto" label-suffix=" :" :inline="true" :model="dialogProps.row">
      <el-form-item v-if="isNextMode" label="参加本次校内评委" prop="onCampusAttendNum">
        <InputInteger v-model="dialogProps.row!.onCampusAttendNum" placeholder="0" disabled />
      </el-form-item>
      <el-form-item v-if="isNextMode" label="参加本次校外评委" prop="offCampusAttendNum">
        <InputInteger v-model="dialogProps.row!.offCampusAttendNum" placeholder="0" disabled />
      </el-form-item>
      <template v-if="!isGoodSubjectMode">
        <el-form-item label="抽取校内评委" prop="onCampusNeedNum">
          <InputInteger v-model="dialogProps.row!.onCampusNeedNum" placeholder="请输入" style="width: 196px" />
        </el-form-item>
        <el-form-item label="抽取校外评委" prop="offCampusNeedNum">
          <InputInteger v-model="dialogProps.row!.offCampusNeedNum" placeholder="请输入" style="width: 196px" />
        </el-form-item>
      </template>
    </el-form>
    <el-table v-if="isGoodSubjectMode" :data="tableData" style="width: 100%">
      <el-table-column prop="subjectName" label="科目名称" />
      <el-table-column prop="onCampusNeedNum" label="校内评委" width="180">
        <template #default="scope">
          <InputInteger v-model="scope.row.onCampusNeedNum" placeholder="请输入" />
        </template>
      </el-table-column>
      <el-table-column prop="offCampusNeedNum" label="校外评委" width="180">
        <template #default="scope">
          <InputInteger v-model="scope.row.offCampusNeedNum" placeholder="请输入" />
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="ExtractDialog">
import { FormInstance, ElMessage } from "element-plus";
import { ref, computed } from "vue";
import { useGoodSubjectEnum } from "@/hooks/useEnum";
import { extractJudge, extractJudgeBySubject } from "@/api/modules/judgeExpert";

interface DialogProps {
  mode: "first" | "next";
  row: any;
  getTableList?: () => void;
}

interface TableRow {
  judgeTaskId: string;
  subjectName: string;
  onCampusNeedNum: string;
  offCampusNeedNum: string;
}

const { goodSubjectEnum } = useGoodSubjectEnum();

// TODO 这里最好添加下表单校验，时间紧张后续优化，先走服务端校验

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  mode: "first",
  row: {
    judgeTaskId: "",
    judgeTaskName: "",
    drawSubject: "2"
  }
});

const tableData = ref<TableRow[]>([]);

const isNextMode = computed(() => dialogProps.value.mode === "next");

const dialogTitle = computed(() => (isNextMode.value ? "追加抽取" : "抽取"));

const isGoodSubjectMode = computed(() => dialogProps.value.row?.drawSubject === "1");

const extractTimes = computed(() => dialogProps.value.row?.drawTimes || "0");

// 提交数据
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      let params: any = {};
      if (isGoodSubjectMode.value && tableData.value.length > 0) {
        params = tableData.value.map(item => ({
          judgeTaskId: item.judgeTaskId,
          subjectName: item.subjectName,
          onCampusNeedNum: item.onCampusNeedNum || "0",
          offCampusNeedNum: item.offCampusNeedNum || "0"
        }));
      } else {
        params = {
          judgeTaskId: dialogProps.value.row.judgeTaskId,
          onCampusNeedNum: dialogProps.value.row.onCampusNeedNum || "0",
          offCampusNeedNum: dialogProps.value.row.offCampusNeedNum || "0"
        };
      }
      const api = isGoodSubjectMode.value ? extractJudgeBySubject : extractJudge;
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
  if (isGoodSubjectMode.value) {
    tableData.value = goodSubjectEnum.value.map(item => ({
      judgeTaskId: params.row.judgeTaskId,
      subjectName: item.label,
      onCampusNeedNum: "",
      offCampusNeedNum: ""
    }));
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
.times {
  margin-bottom: 10px;
  font-size: 22px;
  font-weight: bold;
}
</style>
