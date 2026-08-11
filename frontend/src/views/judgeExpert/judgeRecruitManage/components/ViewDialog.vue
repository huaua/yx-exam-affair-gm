<!-- 任务征集-评委征集任务管理-查看Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="1100"
  >
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="drawTimes" label="抽取次数" align="center" width="90" />
      <el-table-column prop="createTime" label="抽取时间" align="center" width="160" />
      <el-table-column prop="subjectName" label="擅长科目" align="center" />
      <el-table-column prop="onCampusTotalNum" label="校内抽取数" align="center" width="100" />
      <el-table-column prop="onCampusAttendNum" label="参加评委数" align="center" width="100" />
      <el-table-column prop="onCampusUnattendNum" label="不参加评委数" align="center" width="120" />
      <el-table-column prop="offCampusTotalNum" label="校外抽取数" align="center" width="100" />
      <el-table-column prop="offCampusAttendNum" label="参加评委数" align="center" width="100" />
      <el-table-column prop="offCampusUnattendNum" label="不参加评委数" align="center" width="120" />
    </el-table>
  </el-dialog>
</template>

<script setup lang="ts" name="ViewDialog">
import { ref } from "vue";
import { viewJudgeRecruitTask } from "@/api/modules/judgeExpert";

interface DialogProps {
  title: string;
  judgeTaskId: string;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "查看",
  judgeTaskId: ""
});

const tableData = ref([]);

const getTableData = async () => {
  try {
    const params = {
      judgeTaskId: dialogProps.value.judgeTaskId
    };
    const res = await viewJudgeRecruitTask(params);
    tableData.value = res.list || [];
  } catch (error) {
    console.log(error);
  }
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  getTableData();
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
