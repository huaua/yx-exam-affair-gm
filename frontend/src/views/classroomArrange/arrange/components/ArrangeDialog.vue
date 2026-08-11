<!-- 考场编排-考场排考编排-编排Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="1200"
  >
    <div class="main">
      <div class="main-left">
        <div class="table-header">
          <div class="item-1">考试征集任务：{{ detailObj?.roomTaskName || "--" }}</div>
          <div>日期：{{ detailObj?.taskDay || "--" }}</div>
          <div>学院：{{ detailObj?.deptName || "--" }}</div>
          <div>专业：{{ detailObj?.profName || "--" }}</div>
          <div>方向：{{ detailObj?.directionName || "--" }}</div>
          <div>类型：{{ arrangeTypeStr }}</div>
          <div>考生数：{{ detailObj?.stuNum || "--" }}</div>
          <div>已选考场：{{ selectedRoomNum || "--" }}考场 / {{ selectedStuNum || "--" }}人</div>
        </div>
        <el-table :data="tableData1" height="400px" border>
          <el-table-column prop="roomName" label="教室名称" align="center" width="180" />
          <el-table-column prop="stuNum-roomCapacity" label="考生数/容量" align="center" width="150">
            <template #default="scope">
              <div class="flex-center">
                {{ scope.row.stuNum }}/<EditInputInteger v-model="scope.row.roomCapacity" @save="onInputSave" not-empty />
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="arrangeNoStr" label="考场编号" align="center" width="130">
            <template #default="scope">
              <InputSerialNumber v-model="scope.row.arrangeNoStr" :table-data="tableData1" @change="onSerialNumberChange" />
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" align="center" width="110">
            <template #default="scope">
              <el-button type="danger" link @click.prevent="onDeleteClick(scope.$index, scope.row)">删除</el-button>
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
          @row-dblclick="onRowDblClick"
        >
          <!-- 表格 header -->
          <template #tableHeader>
            <el-form-item label="教室名称: ">
              <el-input v-model="searchObj.roomName" placeholder="请输入" clearable />
            </el-form-item>
            <el-form-item label="所属学院">
              <el-select v-model="searchObj.deptId" placeholder="请选择" clearable>
                <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
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

<script setup lang="tsx" name="ArrangeDialog">
import { Arrange } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ElMessage } from "element-plus";
import { ref, reactive, computed } from "vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import EditInputInteger from "@/components/EditInputInteger.vue";
import InputSerialNumber from "./InputSerialNumber.vue";
import { getArrangeDetail, getArrangeDetailList, getCanArrangeClassroomList, saveArrangeDetailList } from "@/api/modules/arrange";

interface DialogProps {
  title: string;
  mode: string;
  row: Partial<Arrange.ResClassroomArrangeList>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  mode: "add",
  row: {}
});

// ===================== 左侧表格逻辑 ======================
let maxSerialNo = 0;
const tableData1 = ref<Arrange.ResArrangeDetailList[]>([]);
const detailObj = reactive({
  roomTaskName: "",
  taskDay: "",
  deptName: "",
  profName: "",
  directionName: "",
  arrangeType: "",
  stuNum: 0,
  maxArrangeRoomNo: 0,
  maxArrangeRoomNoStr: "000"
});

const arrangeTypeStr = computed(() => {
  return detailObj?.arrangeType === "1" ? "按组" : "按位";
});

const selectedRoomNum = computed(() => {
  return tableData1.value.length;
});

const selectedStuNum = computed(() => {
  return tableData1.value.reduce((acc, item) => acc + Number(item.stuNum), 0);
});

const getTableData1 = async () => {
  const params = {
    infoId: dialogProps.value.row.infoId!
  };
  const { list } = await getArrangeDetailList(params);
  tableData1.value = list || [];
};

// 点击-清除
const onDeleteClick = (index: number, row) => {
  // 新增元素（右侧表格）
  remainStuNum += row.stuNum || 0;
  tableData2.value?.push(row);
  // 移除当前行（左侧表格）
  tableData1.value.splice(index, 1);
  // 重新填充考生数
  refillStuNum();
  // 更新最大考场编号
  updateMaxSerialNo();
};

// 编辑-考场编号
const onSerialNumberChange = () => {
  updateMaxSerialNo();
};

// 更新最大考场编号
const updateMaxSerialNo = () => {
  const curMaxSerialNumber = Math.max(...tableData1.value.map(item => Number(item.arrangeNoStr)));
  maxSerialNo = Math.max(curMaxSerialNumber, detailObj.maxArrangeRoomNo);
};

const onInputSave = () => {
  refillStuNum();
};

// 重新填充考生数
const refillStuNum = () => {
  remainStuNum = detailObj.stuNum;
  const len = tableData1.value.length;
  if (len < 1) {
    return;
  }
  for (let i = 0; i < len; i++) {
    const item = tableData1.value[i];
    if (remainStuNum <= 0) {
      item.stuNum = 0; // 当前项考生数置为 0
      continue; // 剩余人数为0，跳过
    }
    const stuNum = Number(item.stuNum);
    const roomCapacity = Number(item.roomCapacity);
    const canFillStuNum = roomCapacity - stuNum; // 可填充数量
    remainStuNum = remainStuNum - stuNum; // 也可能为负数（当最后一个有剩余时）
    if (canFillStuNum > remainStuNum) {
      item.stuNum = stuNum + remainStuNum;
      remainStuNum = 0;
    } else {
      item.stuNum = stuNum + canFillStuNum; // 填满
      remainStuNum -= canFillStuNum;
    }
  }
};

const _getArrangeDetail = async () => {
  const params = {
    infoId: dialogProps.value.row.infoId!
  };
  const res = await getArrangeDetail(params);
  const obj = res.obj || {};
  detailObj.roomTaskName = obj.roomTaskName;
  detailObj.taskDay = obj.taskDay && obj.taskDay.slice(0, 10);
  detailObj.deptName = obj.deptName;
  detailObj.profName = obj.profName;
  detailObj.directionName = obj.directionName;
  detailObj.stuNum = remainStuNum = Number(obj.stuNum);
  detailObj.roomTaskName = obj.roomTaskName;
  detailObj.arrangeType = obj.arrangeType;
  detailObj.maxArrangeRoomNoStr = obj.maxArrangeRoomNoStr || "000";
  maxSerialNo = detailObj.maxArrangeRoomNo = obj.maxArrangeRoomNo || 0;
  typeProperty.value = obj.arrangeType === "1" ? "groupCapacity" : "capacity";
};

// ===================== 右侧表格逻辑 ======================
let remainStuNum = 0;
const { departmentEnum } = useDepartmentEnum();
const proTable2 = ref<ProTableInstance>();
const tableData2 = ref<Arrange.ResCanArrangeClassroomList[]>([]);
const searchObj = reactive({
  roomName: "",
  deptId: undefined
});
const typeProperty = ref("capacity");

// 表格配置项
const columns2 = reactive<ColumnProps<Arrange.ResClassroomArrangeList>[]>([
  { prop: "roomName", label: "教室名称", width: 210 },
  { prop: "deptId", label: "所属学院", enum: departmentEnum, width: 130 },
  { prop: "groupCapacity", label: "按组容量", width: 90 },
  { prop: "capacity", label: "按位容量", width: 90 }
]);

// 过滤搜索的表格数据
const filteredTableData2 = computed(() => {
  let result = tableData2.value;
  const queryRoomName = searchObj.roomName;
  const queryDeptId = searchObj.deptId;
  if (queryRoomName != null) {
    result = result.filter(classroomObj => {
      return classroomObj.roomName.includes(queryRoomName);
    });
  }
  if (queryDeptId != null) {
    result = result.filter(classroomObj => {
      return classroomObj.deptId.includes(queryDeptId);
    });
  }
  return result;
});

const getTableData2 = async () => {
  const params = {
    infoId: dialogProps.value.row.infoId!
  };
  const { list } = await getCanArrangeClassroomList(params);
  tableData2.value = list || [];
};

// 表格行被双击
const onRowDblClick = (row: any) => {
  if (remainStuNum < 1) {
    ElMessage.warning("已无剩余容量");
    return;
  }
  if (!row[typeProperty.value]) {
    ElMessage.warning(`该教室无剩余${arrangeTypeStr.value}容量`);
    return;
  }
  // 新增元素（左侧表格）
  row.roomCapacity = Number(row[typeProperty.value]);
  row.stuNum = remainStuNum > row.roomCapacity ? row.roomCapacity : remainStuNum;
  remainStuNum = remainStuNum - row.roomCapacity > 0 ? remainStuNum - row.roomCapacity : 0;
  row.arrangeNoStr = getSerialStr(); // 获取排号
  tableData1.value?.push(row);
  // 移除当前行（右侧表格）
  const indexToRemove = tableData2.value.findIndex(item => item.id === row.id);
  if (indexToRemove !== -1) {
    tableData2.value.splice(indexToRemove, 1);
  }
  // 重新计算考生数
  refillStuNum();
};

const getSerialStr = () => {
  maxSerialNo = maxSerialNo + 1;
  const len = detailObj.maxArrangeRoomNoStr?.length;
  return String(maxSerialNo).padStart(len, "0");
};

// 计算初始剩余考生数（编辑时）
const calcInitRemainStuNum = () => {
  const arrangedStuNum = tableData1.value.reduce((acc, item) => acc + Number(item.stuNum), 0) || 0;
  remainStuNum = detailObj.stuNum - arrangedStuNum;
};

// ===================== 共用逻辑 ======================

// 是否禁用提交按钮
const isDisabledSubmit = computed(() => !tableData1.value.length);

// 提交数据
const handleSubmit = async () => {
  const hasEmptyArrangeNoStr = tableData1.value.some(item => item.arrangeNoStr == null || item.arrangeNoStr === "");
  if (hasEmptyArrangeNoStr) {
    ElMessage.error("请先为所有教室分配编号");
    return;
  }
  try {
    const arrangeRoomList = tableData1.value
      .filter(item => Number(item.stuNum) > 0)
      .map(item => ({
        roomReportId: item.id || item.roomReportId,
        roomTaskId: item.taskId || item.roomTaskId,
        roomId: item.roomId,
        roomName: item.roomName,
        stuNum: String(item.stuNum),
        roomCapacity: String(item.roomCapacity),
        arrangeNoStr: String(item.arrangeNoStr),
        needTchNum: String(item.tchNum || item.needTchNum) // 左右表格字段名不一致
      }));
    const params = {
      infoId: dialogProps.value.row.infoId!,
      arrangeRoomList
    };
    await saveArrangeDetailList(params);
    ElMessage.success({ message: `${dialogProps.value.title}成功！` });
    dialogProps.value.getTableList!();
    dialogVisible.value = false;
  } catch (error) {
    console.log(error);
  }
};

const handleClosed = () => {
  tableData1.value = [];
  tableData2.value = [];
  searchObj.roomName = "";
  searchObj.deptId = undefined;
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  try {
    const promises = [_getArrangeDetail(), getTableData2()];
    if (dialogProps.value.mode === "edit") {
      promises.push(
        getTableData1().then(() => {
          calcInitRemainStuNum();
        })
      );
    }
    // 使用 Promise.all 并发执行所有异步操作
    await Promise.all(promises);
  } catch (error) {
    console.error("Error occurred:", error);
  }
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
      grid-template-columns: 2fr 2fr 2.5fr;
      grid-gap: 6px;
      padding-bottom: 16px;
      .item-1 {
        grid-column: 1 / 3;
      }
    }
    .flex-center {
      display: flex;
      align-items: center;
      justify-content: center;
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
