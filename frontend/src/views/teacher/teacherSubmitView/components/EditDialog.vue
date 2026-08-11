<!-- 监考老师管理-监考老师上报查看-编辑Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="600"
  >
    <div class="tips">
      <div v-if="dialogProps.type === 'cancel'">取消该位教职工后无法达到要求上报数量，需从未上报的教职工中进行替换</div>
      <template v-else-if="dialogProps.type === 'replace'">
        <div>要取消已被编排的教职工，需选择未上报的教职工来替换</div>
        <div>并自动分配至相应考场进行监考任务</div>
      </template>
    </div>
    <!-- 表格 搜索框 -->
    <div>
      <el-form-item label="姓名或工号: ">
        <el-input v-model="searchNameOrJobNo" placeholder="请输入" clearable />
      </el-form-item>
    </div>
    <el-table
      v-if="tableData.length"
      ref="tableRef"
      :data="filteredTableData"
      highlight-current-row
      height="400px"
      @current-change="handleCurrentChange"
      border
    >
      <el-table-column label="" align="center" width="80">
        <template #default="scope">
          <input type="radio" value="true" :checked="scope.row?.tchId === currentRow?.tchId" />
        </template>
      </el-table-column>
      <el-table-column prop="tchName" label="姓名" align="center" />
      <el-table-column prop="sex" label="性别" align="center" />
      <el-table-column prop="jobNo" label="工号" align="center" />
    </el-table>
    <el-empty v-else description="暂无数据" />
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :disabled="isDisabledSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="EditDialog">
import { Teacher } from "@/api/interface";
import { ElTable } from "element-plus";
import { ElMessage } from "element-plus";
import { ref, computed } from "vue";
import { useHandleData } from "@/hooks/useHandleData";
import { getCanSubmitTeacherList, replaceSubmitTeacher } from "@/api/modules/teacher";

interface DialogProps {
  type: string;
  title: string;
  row: Partial<Teacher.ResTeacherSubmitViewList>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "cancel",
  title: "",
  row: {}
});
const tableData = ref<Teacher.ResCanSubmitTeacherList[]>([]);
const currentRow = ref<Teacher.ResCanSubmitTeacherList>();
const tableRef = ref<InstanceType<typeof ElTable>>();
const searchNameOrJobNo = ref<string>();

const isDisabledSubmit = computed(() => !currentRow.value?.tchId);

// 过滤搜索的表格数据
const filteredTableData = computed(() => {
  let result = tableData.value;
  const query = searchNameOrJobNo.value;
  if (query != null) {
    result = result.filter(person => {
      return person.tchName.includes(query) || person.jobNo.includes(query);
    });
  }
  return result;
});

const handleCurrentChange = (val: Teacher.ResCanSubmitTeacherList | undefined) => {
  currentRow.value = val;
};

// 提交数据
const handleSubmit = async () => {
  if (!currentRow.value?.tchId) {
    ElMessage.error({ message: "请选择要替换的教职工！" });
    return;
  }
  try {
    const params = {
      reportId: dialogProps.value.row.reportId!,
      tchId: currentRow.value!.tchId!,
      taskDetailId: dialogProps.value.row.taskDetailId!
    };
    await useHandleData(
      replaceSubmitTeacher,
      params,
      `使用 ${currentRow.value.tchName} 对 ${dialogProps.value.row.tchName} 进行替换`
    );
    dialogProps.value.getTableList!();
    dialogVisible.value = false;
  } catch (error) {
    console.log(error);
  }
};

const getTableData = async () => {
  const params = {
    tchTaskId: dialogProps.value.row.tchTaskId!,
    taskTime: dialogProps.value.row.taskTime
    // taskDetailId: dialogProps.value.row.taskDetailId!
  };
  const { list } = await getCanSubmitTeacherList(params);
  tableData.value = list || [];
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
.tips {
  display: flex;
  flex-direction: column;
  padding-bottom: 16px;
  font-weight: 500;
  div + div {
    margin-top: 6px;
  }
}
</style>
