<template>
  <el-dialog v-model="visible" title="评委抽取" width="1180px" top="5vh" destroy-on-close :close-on-click-modal="false">
    <div class="task-info">
      <span>任务名称：{{ task.taskName }}</span
      ><span>评分日期：{{ dateText }}</span
      ><span>抽取轮次：{{ row.drawTimes }}</span>
    </div>
    <div class="subject">
      <span>科目名称：</span><el-input v-model="row.subjectName" placeholder="请输入科目名称" style="width: 180px" />
    </div>
    <div class="grid">
      <div class="selected">
        <div class="title">评分专家（{{ scoring.length }}/{{ row.requiredNum }}）</div>
        <div class="chips">
          <el-tag v-for="x in scoring" :key="x.expertId" closable @close="remove(x)">{{ x.name }}</el-tag>
        </div>
        <div class="title">备用专家（{{ backup.length }}）</div>
        <div class="chips">
          <el-tag v-for="x in backup" :key="x.expertId" closable type="warning" @close="remove(x)">{{ x.name }}</el-tag>
        </div>
      </div>
      <div>
        <div class="filters">
          <el-input v-model="filters.name" clearable placeholder="请输入姓名" />
          <el-select v-model="filters.expertCategory" clearable placeholder="专家类别">
            <el-option label="校内" value="1" />
            <el-option label="校外" value="2" />
          </el-select>
          <el-select v-model="filters.deptId" clearable placeholder="所属学院">
            <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="filters.sex" clearable placeholder="性别">
            <el-option label="男" value="男" />
            <el-option label="女" value="女" />
          </el-select>
        </div>
        <el-table :data="filteredExperts" height="390" border
          ><el-table-column prop="name" label="姓名" width="85" /><el-table-column
            prop="sex"
            label="性别"
            width="60"
          /><el-table-column prop="phoneNo" label="手机号" width="120" /><el-table-column
            prop="evaluation"
            label="评价"
            min-width="90"
            ><template #default="s">
              <span v-if="!s.row.evaluation">--</span>
              <el-button v-else link type="primary" @click="showEvaluation(s.row)">查看</el-button>
            </template></el-table-column
          ><el-table-column label="专家类别" width="90"
            ><template #default="s">{{ s.row.expertCategory === "1" ? "校内" : "校外" }}</template></el-table-column
          ><el-table-column label="所属学院/所在单位" min-width="160" header-class-name="dept-header"
            ><template #default="s">{{ s.row.deptName || s.row.unitName || "--" }}</template></el-table-column
          ><el-table-column label="选择" width="120"
            ><template #default="s"
              ><el-button link type="primary" :disabled="s.row.expertRole === '1'" @click="choose(s.row, 1)">评分</el-button
              ><el-button link type="warning" :disabled="s.row.expertRole === '2'" @click="choose(s.row, 2)"
                >备用</el-button
              ></template
            ></el-table-column
          ></el-table
        >
      </div>
    </div>
    <el-dialog v-model="evalVisible" title="评价内容" width="520px" append-to-body destroy-on-close :close-on-click-modal="false">
      <div class="eval-content">{{ evalContent }}</div>
    </el-dialog>
    <template #footer
      ><el-button @click="visible = false">取消</el-button><el-button type="primary" @click="submit">确定</el-button></template
    >
  </el-dialog>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import { ElMessage } from "element-plus";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import {
  getDrawManageExperts,
  saveDrawManageSelection,
  removeDrawManageSelection,
  submitDrawManageSelection
} from "@/api/modules/judgeExpert";
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const emit = defineEmits(["changed"]),
  visible = ref(false),
  evalVisible = ref(false),
  evalContent = ref(""),
  row = ref<any>({}),
  task = ref<any>({}),
  experts = ref<any[]>([]),
  filters = ref<any>({ name: "", expertCategory: "", deptId: "", sex: "" });
const showEvaluation = (x: any) => {
  evalContent.value = x.evaluation || "";
  evalVisible.value = true;
};
const scoring = computed(() => experts.value.filter(x => x.expertRole === "1")),
  backup = computed(() => experts.value.filter(x => x.expertRole === "2")),
  dateText = computed(() => `${String(task.value.taskStartTime).slice(0, 10)}至${String(task.value.taskEndTime).slice(0, 10)}`),
  filteredExperts = computed(() => {
    const q = filters.value.name?.trim();
    return experts.value.filter(x => {
      if (x.expertRole === "1" || x.expertRole === "2") return false;
      if (q && !x.name?.includes(q)) return false;
      if (filters.value.expertCategory && String(x.expertCategory) !== String(filters.value.expertCategory)) return false;
      if (filters.value.deptId && String(x.deptId) !== String(filters.value.deptId)) return false;
      if (filters.value.sex && x.sex !== filters.value.sex) return false;
      return true;
    });
  });
const load = async () => {
  const r: any = await getDrawManageExperts({ rulesId: row.value.rulesId });
  experts.value = r.list || [];
};
const choose = async (x: any, role: number) => {
  if (role === 1 && scoring.value.length >= Number(row.value.requiredNum)) {
    ElMessage.warning(`评分专家最多添加${row.value.requiredNum}人`);
    return;
  }
  try {
    await saveDrawManageSelection({ rulesId: row.value.rulesId, expertId: x.expertId, expertRole: String(role) });
    load();
  } catch {
    // 全局拦截器已弹出业务错误提示，这里静默吞掉避免冒泡成“未知错误”
  }
};
const remove = async (x: any) => {
  try {
    await removeDrawManageSelection({ rulesId: row.value.rulesId, expertId: x.expertId });
    load();
  } catch {
    // 全局拦截器已弹出业务错误提示，这里静默吞掉避免冒泡成“未知错误”
  }
};
const submit = async () => {
  if (!row.value.subjectName?.trim()) {
    ElMessage.warning("科目名称必须填写");
    return;
  }
  if (scoring.value.length === 0) {
    ElMessage.warning("未设置评分专家不允许提交");
    return;
  }
  try {
    await submitDrawManageSelection({ rulesId: row.value.rulesId, subjectName: row.value.subjectName });
    ElMessage.success("抽取提交成功");
    visible.value = false;
    emit("changed");
  } catch {
    // 全局拦截器已弹出业务错误提示，这里静默吞掉避免冒泡成“未知错误”
  }
};
const acceptParams = async (p: any) => {
  row.value = p.row;
  task.value = p.task;
  visible.value = true;
  load();
};
defineExpose({ acceptParams });
</script>
<style scoped lang="scss">
.task-info {
  display: flex;
  gap: 32px;
  margin-bottom: 16px;
}
.subject {
  display: flex;
  align-items: center;
  margin-bottom: 14px;
}
.grid {
  display: grid;
  grid-template-columns: 1fr 1.45fr;
  gap: 20px;
}
.selected {
  border: 1px solid #dcdfe6;
  min-height: 430px;
}
.title {
  padding: 10px 14px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  font-weight: 600;
}
.chips {
  min-height: 125px;
  padding: 12px;
  display: flex;
  align-content: flex-start;
  gap: 8px;
  flex-wrap: wrap;
}
.filters {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 10px;
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
