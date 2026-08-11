<template>
  <el-dialog v-model="visible" title="专家评价" width="560px" align-center :close-on-click-modal="false">
    <div class="expert-name">专家姓名：{{ row.name }}</div>
    <el-input v-model="evaluation" type="textarea" maxlength="2000" show-word-limit :rows="8" placeholder="请输入评价内容" />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="submit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { evaluateExpert, getExpertDetail } from "@/api/modules/judgeExpert";
const visible = ref(false);
const evaluation = ref("");
const row = ref<Partial<JudgeExpert.ResExpertDatabaseList>>({});
let refresh: (() => void) | undefined;
const submit = async () => {
  await evaluateExpert({ expertId: row.value.expertId!, evaluation: evaluation.value });
  ElMessage.success("评价保存成功");
  visible.value = false;
  refresh?.();
};
const acceptParams = async (params: { row: Partial<JudgeExpert.ResExpertDatabaseList>; getTableList?: () => void }) => {
  row.value = { ...params.row };
  refresh = params.getTableList;
  visible.value = true;
  // 打开时拉取最新数据，确保显示已有评价内容
  const res = await getExpertDetail({ expertId: String(params.row.expertId) });
  if (res.obj) {
    evaluation.value = res.obj.evaluation || "";
  } else {
    evaluation.value = params.row.evaluation || "";
  }
};
defineExpose({ acceptParams });
</script>
<style scoped>
.expert-name {
  margin-bottom: 12px;
  color: #606266;
}
</style>
