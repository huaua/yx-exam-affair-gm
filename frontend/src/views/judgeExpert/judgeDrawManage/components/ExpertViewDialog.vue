<template>
  <el-dialog v-model="visible" title="抽取专家查看" width="1200px" align-center destroy-on-close :close-on-click-modal="false">
    <div class="filters">
      <el-select v-model="filters.expertCategory" clearable placeholder="专家类别-请选择">
        <el-option label="校内" value="1" />
        <el-option label="校外" value="2" />
      </el-select>
      <el-select v-model="filters.expertRole" clearable placeholder="专家模式-请选择">
        <el-option label="评分专家" value="1" />
        <el-option label="备用专家" value="2" />
      </el-select>
      <el-select v-model="filters.submitState" clearable placeholder="提交状态-请选择">
        <el-option label="已提交" value="2" />
        <el-option label="未提交" value="1" />
      </el-select>
      <el-button type="primary" @click="load">查询</el-button>
    </div>
    <el-table :data="filteredExperts" border height="560">
      <el-table-column prop="name" label="姓名" width="80" />
      <el-table-column prop="idCard" label="证件号" width="170" />
      <el-table-column prop="sex" label="性别" width="55" />
      <el-table-column prop="phoneNo" label="手机号" width="115" />
      <el-table-column label="专家类别" width="95">
        <template #default="s">{{ s.row.expertCategory === "1" ? "校内" : "校外" }}</template>
      </el-table-column>
      <el-table-column label="所属学院/所在单位" width="165" header-class-name="dept-header" show-overflow-tooltip>
        <template #default="s">{{ s.row.deptName || s.row.unitName || "--" }}</template>
      </el-table-column>
      <el-table-column label="职称" width="90" show-overflow-tooltip>
        <template #default="s">{{ getTitleName(s.row.title) }}</template>
      </el-table-column>
      <el-table-column prop="goodSubjects" label="擅长科目" width="115" show-overflow-tooltip />
      <el-table-column label="评价" width="85">
        <template #default="s">
          <span v-if="!s.row.evaluation">--</span>
          <el-button v-else link type="primary" @click="showEvaluation(s.row)">查看</el-button>
        </template>
      </el-table-column>
      <el-table-column label="专家模式" width="90">
        <template #default="s">{{ s.row.expertRole === "1" ? "评分专家" : "备用专家" }}</template>
      </el-table-column>
      <el-table-column label="提交状态" width="90">
        <template #default="s">{{ s.row.submitState === "2" ? "已提交" : "未提交" }}</template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="evalVisible" title="评价内容" width="520px" append-to-body destroy-on-close :close-on-click-modal="false">
      <div class="eval-content">{{ evalContent }}</div>
    </el-dialog>
  </el-dialog>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import { getDrawManageExperts } from "@/api/modules/judgeExpert";
import { useExpertTitleEnum } from "@/hooks/useEnum";
const visible = ref(false),
  evalVisible = ref(false),
  evalContent = ref(""),
  row = ref<any>({}),
  experts = ref<any[]>([]),
  filters = ref<any>({ expertCategory: "", expertRole: "", submitState: "" });
const { expertTitleEnum } = useExpertTitleEnum();
const titleNameMap = computed(() => {
  return Object.fromEntries(expertTitleEnum.value.map((x: any) => [String(x.value), x.label]));
});
const getTitleName = (title: any) => titleNameMap.value[String(title)] || title || "--";
const showEvaluation = (x: any) => {
  evalContent.value = x.evaluation || "";
  evalVisible.value = true;
};
const filteredExperts = computed(() => {
  return experts.value.filter((x: any) => {
    if (x.expertRole === "0") return false;
    if (filters.value.expertCategory && String(x.expertCategory) !== String(filters.value.expertCategory)) return false;
    if (filters.value.expertRole && String(x.expertRole) !== String(filters.value.expertRole)) return false;
    if (filters.value.submitState && String(x.submitState) !== String(filters.value.submitState)) return false;
    return true;
  });
});
const load = async () => {
  const r: any = await getDrawManageExperts({ rulesId: row.value.rulesId });
  experts.value = r.list || [];
};
const acceptParams = async (p: any) => {
  row.value = p.row;
  filters.value = { expertCategory: "", expertRole: "", submitState: "" };
  visible.value = true;
  load();
};
defineExpose({ acceptParams });
</script>
<style scoped lang="scss">
.filters {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
  .el-select {
    width: 180px;
  }
}
.eval-content {
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.6;
  max-height: 60vh;
  overflow-y: auto;
}
:deep(.dept-header .cell) {
  white-space: nowrap;
}
</style>
