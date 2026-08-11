<!-- 监考老师管理-监考老师上报-上报Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="1300"
  >
    <div class="main">
      <div class="main-left">
        <div class="table-header">
          <span>考试任务：{{ dialogProps.row.tchTaskName || "--" }}</span>
          <span>
            征集日期：{{ formatTime(dialogProps.row.taskStartTime) || "--" }}~{{
              formatTime(dialogProps.row.taskEndTime) || "--"
            }}
          </span>
        </div>
        <el-table ref="proTable1" :data="tableData1" height="400px" border>
          <el-table-column prop="tchName" label="姓名" align="center" width="90" />
          <el-table-column prop="jobNo" label="工号" align="center" width="90" />
          <el-table-column v-for="item in dateList" :key="item.taskTime" :prop="item.taskTime" align="center" width="90">
            <template #header>
              <div>{{ formatTime(item.taskTime, "MM-DD") }}</div>
              <div>({{ sumByDate(item.taskTime) || "--" }}/{{ item.needNum || "--" }})</div>
            </template>
            <template #default="scope">
              {{ scope.row[item.taskTime] ? "✔" : "" }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" align="center" width="90" />
          <el-table-column fixed="right" label="操作" align="center" width="80">
            <template #default="scope">
              <el-button type="primary" link @click.prevent="onDeleteClick(scope.$index, scope.row)">清除</el-button>
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
  <ConfirmSubmitDialog ref="dialogRef" />
</template>

<script setup lang="tsx" name="SubmitDialog">
import { Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ElMessageBox, ElMessage } from "element-plus";
import { ref, h, reactive, computed } from "vue";
import ProTable from "@/components/ProTable/index.vue";
import ConfirmSubmitDialog from "./ConfirmSubmitDialog.vue";
import {
  getCanSubmitTeacherList,
  getCurTaskDateList,
  submitTeacher,
  getCacheSubmitTeacher,
  cacheSubmitTeacher,
  nocacheSubmitTeacher
} from "@/api/modules/teacher";

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResTeacherRecruitManageList>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {}
});

// ===================== 左侧表格逻辑 ======================

const proTable1 = ref<ProTableInstance>();
const tableData1 = ref<any>([]);
const dateList = ref<Teacher.ResCurTaskDateList[]>([]);

// 点击-清除
const onDeleteClick = (index: number, row) => {
  _nocacheSubmitTeacher(row); // 移除服务端缓存
  // 新增元素（右侧表格）
  dateList.value.forEach(item => {
    delete row[item.taskTime];
  });
  tableData2.value?.push(row);
  // 移除当前行（左侧表格）
  tableData1.value.splice(index, 1);
};

// 获取表格1日期数据
const getDateList = async () => {
  const params = {
    tchTaskId: dialogProps.value.row.tchTaskId!,
    onlyQueryNeedReportTimeList: true
  };
  try {
    const { list } = await getCurTaskDateList(params);
    dateList.value = list || [];
  } catch (error) {
    dialogVisible.value = false;
  }
};

const getTableData1 = async () => {
  const params = { tchTaskId: dialogProps.value.row.tchTaskId! };
  const { list } = await getCacheSubmitTeacher(params);
  tableData1.value =
    list?.map(item => {
      const row = { ...item };
      const taskTimes = row?.taskTimes || [];
      taskTimes.forEach(dateItem => {
        row[dateItem] = dateItem;
      });
      return row;
    }) || [];
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
  console.log("123result321:", result);
  return result;
});

const getTableData2 = async () => {
  const params = {
    tchTaskId: dialogProps.value.row.tchTaskId!
  };
  const { list } = await getCanSubmitTeacherList(params);
  return list || [];
};

// 表格行被双击
const onRowDblClick = (row: any) => {
  openConfirmSubmitDialog(row);
};

// 打开确认上报
const dialogRef = ref<InstanceType<typeof ConfirmSubmitDialog> | null>(null);
const openConfirmSubmitDialog = (row: Partial<Teacher.ResCanSubmitTeacherList> = {}) => {
  const params = {
    title: "确认上报",
    row: { ...row },
    dateList: dateList.value,
    success: (selectedDateList, remark) => {
      handleAddDate(row, selectedDateList, remark);
    }
  };
  dialogRef.value?.acceptParams(params);
};

const handleAddDate = (row: Partial<Teacher.ResCanSubmitTeacherList> = {}, selectedDateList, remark) => {
  row.remark = remark;
  // 新增元素（左侧表格）
  selectedDateList.forEach(date => {
    row[date] = date;
  });
  tableData1.value?.push(row);
  _cacheSubmitTeacher(row); // 服务端缓存
  // 移除当前行（右侧表格）
  const indexToRemove = tableData2.value.findIndex(item => item.tchId === row.tchId);
  if (indexToRemove !== -1) {
    tableData2.value.splice(indexToRemove, 1);
  }
};

const _cacheSubmitTeacher = row => {
  const params = [generateApiParamForRow(row)];
  cacheSubmitTeacher(params);
};

const _nocacheSubmitTeacher = row => {
  const params = [generateApiParamForRow(row)];
  nocacheSubmitTeacher(params);
};

// ===================== 共用逻辑 ======================

// 是否禁用提交按钮
const isDisabledSubmit = computed(() => !tableData1.value.length);

// 提交数据
const handleSubmit = async () => {
  const isPassed = await validateSubmitNumber();
  if (!isPassed) {
    return;
  }
  const params = tableData1.value.map(item => {
    return generateApiParamForRow(item);
  });
  await submitTeacher(params);
  ElMessage.success({ message: `${dialogProps.value.title}成功！` });
  dialogProps.value.getTableList!();
  dialogVisible.value = false;
};

// 生成接口参数（单行）
const generateApiParamForRow = row => {
  const reportTaskDetailList: any = [];
  dateList.value.forEach(dateObj => {
    if (dateObj.taskTime in row) {
      reportTaskDetailList.push({
        taskDetailId: dateObj.taskDetailId,
        tchTaskId: dialogProps.value.row.tchTaskId,
        taskTime: dateObj.taskTime
      });
    }
  });
  return {
    tchId: row.tchId,
    tchName: row.tchName,
    jobNo: row.jobNo,
    remark: row.remark,
    sex: row.sex,
    reportTaskDetailList
  };
};

const validateSubmitNumber = () => {
  const failForceMessageList: any = []; // 校验失败的强制征集列表
  const failUnforceMessageList: any = []; // 校验失败的非强制征集列表
  dateList.value.forEach(dateObj => {
    const sum = sumByDate(dateObj.taskTime);
    if (sum < dateObj.needNum) {
      if (dateObj.forceFlag === "1") {
        failForceMessageList.push(`${formatTime(dateObj.taskTime, "MM-DD")} 要求上报${dateObj.needNum}人，已上报${sum}人`);
      } else {
        failUnforceMessageList.push(`${formatTime(dateObj.taskTime, "MM-DD")} 要求上报${dateObj.needNum}人，已上报${sum}人`);
      }
    }
  });
  if (failForceMessageList.length > 0) {
    const message = failForceMessageList.map(item => h("div", null, item));
    message.push(h("div", null, "未达到要求上报人数，请继续添加"));
    ElMessageBox.alert(h("div", null, message), "温馨提示", {
      confirmButtonText: "确定",
      type: "warning",
      draggable: true
    });
    return false;
  } else if (failUnforceMessageList.length > 0) {
    const message = failUnforceMessageList.map(item => h("div", null, item));
    message.push(h("div", null, "未达到要求上报人数，确认上报吗"));
    return ElMessageBox.confirm(h("div", null, message), "温馨提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
      draggable: true
    })
      .then(() => {
        return true;
      })
      .catch(() => {
        return false;
      });
  }
  return true;
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

// 计算某天预上报人数
const sumByDate = date => {
  return tableData1.value.reduce((acc, obj) => acc + ((obj[date] && 1) || 0), 0);
};

const handleClosed = () => {
  tableData1.value = [];
  tableData2.value = [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  getTableData1();
  await getDateList();
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
      display: flex;
      padding: 8px 0 22px;
      span + span {
        margin-left: 50px;
      }
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
