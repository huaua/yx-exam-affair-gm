<!-- 监考老师管理-监考老师上报-追加Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="1100"
  >
    <div class="main">
      <div class="main-left">
        <div class="table-header">
          <div>考试任务：{{ dialogProps.row.tchTaskName || "--" }}</div>
          <div>征集日期：{{ formatTime(dialogProps.selectedDateObj?.taskTime) || "--" }}</div>
          <div>要求上报人数：{{ dialogProps.selectedDateObj!.needNum || "--" }}</div>
          <div>当前上报人数：{{ dialogProps.selectedDateObj!.collectNum || "--" }}</div>
        </div>
        <el-table ref="proTable1" :data="tableData1" height="400px" border>
          <el-table-column prop="tchName" label="姓名" align="center" width="110" />
          <el-table-column prop="jobNo" label="工号" align="center" width="160" />
          <el-table-column prop="remark" label="备注" align="center" width="200">
            <template #default="scope">
              <el-input v-model="scope.row.remark" @change="handleScoreChange(scope.row)" size="small"></el-input>
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" align="center" width="140">
            <template #default="scope">
              <el-button type="primary" link @click.prevent="onDeleteClick(scope.$index, scope.row)">取消追加</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div class="main-right">
        <ProTable
          ref="proTable2"
          :columns="columns2"
          :data="filteredTableData2"
          :pagination="false"
          :tool-button="false"
          height="400px"
          width="500"
          @row-dblclick="onRowDblClick"
        >
          <!-- 表格 header -->
          <template #tableHeader>
            <el-form-item label="姓名或工号: ">
              <el-input v-model="searchNameOrJobNo" placeholder="请输入" clearable />
            </el-form-item>
          </template>
        </ProTable>
      </div>
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :disabled="isDisabledSubmit" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="tsx" name="AppendDialog">
import { Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ElMessage } from "element-plus";
import { ref, reactive, computed } from "vue";
import ProTable from "@/components/ProTable/index.vue";
import { getCanSubmitTeacherList, appendTeacher } from "@/api/modules/teacher";

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResTeacherRecruitManageList>;
  selectedDateObj: Teacher.ResCurTaskDateList | undefined;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {},
  selectedDateObj: undefined
});

// ===================== 左侧表格逻辑 ======================

const proTable1 = ref<ProTableInstance>();
const tableData1 = ref<any>([]);

// 点击-清除
const onDeleteClick = (index: number, row) => {
  // 新增元素（右侧表格）
  tableData2.value?.push(row);
  // 移除当前行（左侧表格）
  tableData1.value.splice(index, 1);
};

// ===================== 右侧表格逻辑 ======================
const proTable2 = ref<ProTableInstance>();
const tableData2 = ref<any>([]);
const searchNameOrJobNo = ref<string>();
// 表格配置项
const columns2 = reactive<ColumnProps<Teacher.ResCanSubmitTeacherList>[]>([
  { prop: "tchName", label: "教职工姓名", width: 110 },
  { prop: "sex", label: "性别", width: 90 },
  { prop: "jobNo", label: "工号", width: 120 }
]);

// 过滤搜索的表格数据
const filteredTableData2 = computed(() => {
  let result = tableData2.value;
  const query = searchNameOrJobNo.value;
  if (query != null) {
    result = result.filter(person => {
      return person.tchName.includes(query) || person.jobNo.includes(query);
    });
  }
  return result;
});

const getTableData2 = async () => {
  const params = {
    tchTaskId: dialogProps.value.row.tchTaskId!,
    taskTime: dialogProps.value.selectedDateObj?.taskTime
  };
  const { list } = await getCanSubmitTeacherList(params);
  return list || [];
};

// 表格行被双击
const onRowDblClick = (row: any) => {
  // 新增元素（左侧表格）
  tableData1.value?.push(row);
  // 移除当前行（右侧表格）
  const indexToRemove = tableData2.value.findIndex(item => item.tchId === row.tchId);
  if (indexToRemove !== -1) {
    tableData2.value.splice(indexToRemove, 1);
  }
};

const handleScoreChange = (row: any) => {
  const index = tableData1.value.findIndex((item: any) => item.tchId === row.tchId);
  if (index !== -1) {
    tableData1.value[index].remark = row.remark;
  }
};

// ===================== 共用逻辑 ======================

// 是否禁用提交按钮
const isDisabledSubmit = computed(() => !tableData1.value.length);

// 提交数据
const handleSubmit = async () => {
  const teacherList = tableData1.value.map(item => ({
    tchId: item.tchId,
    tchName: item.tchName,
    jobNo: item.jobNo,
    sex: item.sex,
    remark: item.remark
  }));
  const params = {
    taskDetailId: (dialogProps.value.selectedDateObj && dialogProps.value.selectedDateObj.taskDetailId!) || "",
    tchList: teacherList
  };
  await appendTeacher(params);
  ElMessage.success({ message: `${dialogProps.value.title}成功！` });
  dialogProps.value.getTableList!();
  dialogVisible.value = false;
};

const formatTime = (fullTimeStr, type = "YYYY-MM-DD") => {
  if (!fullTimeStr) {
    return "";
  }
  if (type === "MM-DD") {
    return fullTimeStr.slice(5, 10);
  }
  return fullTimeStr.slice(0, 10);
};

const handleClosed = () => {
  tableData1.value = [];
  tableData2.value = [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  tableData2.value = await getTableData2();
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.main {
  display: flex;
  justify-content: space-between;
  .main-left {
    .table-header {
      display: grid;
      grid-template-columns: 1fr 1fr;
      grid-gap: 6px;
      padding-bottom: 16px;
    }
  }
  .main-right {
    :deep(.el-table__row) {
      cursor: pointer;
    }
    :deep(.el-form-item__content) {
      width: 180px;
    }
  }
}
</style>
