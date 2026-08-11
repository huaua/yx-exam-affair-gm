<!-- 考场编排-考场排考编排-查看考场Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    @closed="handleClosed"
    width="600"
  >
    <div class="main">
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
      <el-table :data="tableData" height="400px" border>
        <el-table-column prop="roomName" label="教室名称" align="center" />
        <el-table-column prop="stuNum-roomCapacity" label="考生数/容量" align="center" width="150">
          <template #default="scope">
            <div class="flex-center">{{ scope.row.stuNum }}/{{ scope.row.roomCapacity }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="arrangeNoStr" label="考场编号" align="center" width="130" />
      </el-table>
    </div>
  </el-dialog>
</template>

<script setup lang="tsx" name="ViewArrangeDialog">
import { Arrange } from "@/api/interface";
import { ref, reactive, computed } from "vue";
import { getArrangeDetail, getArrangeDetailList } from "@/api/modules/arrange";

interface DialogProps {
  title: string;
  row: Partial<Arrange.ResClassroomArrangeList>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {}
});

const tableData = ref<any>([]);
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
  return tableData.value.length;
});

const selectedStuNum = computed(() => {
  return tableData.value.reduce((acc, item) => acc + Number(item.stuNum), 0);
});

const getTableData = async () => {
  const params = {
    infoId: dialogProps.value.row.infoId!
  };
  const { list } = await getArrangeDetailList(params);
  tableData.value = list || [];
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
  detailObj.stuNum = Number(obj.stuNum);
  detailObj.roomTaskName = obj.roomTaskName;
  detailObj.arrangeType = obj.arrangeType;
  detailObj.maxArrangeRoomNoStr = obj.maxArrangeRoomNoStr;
  detailObj.maxArrangeRoomNo = obj.maxArrangeRoomNo;
  typeProperty.value = obj.arrangeType === "1" ? "groupCapacity" : "capacity";
};

const typeProperty = ref("capacity");

const handleClosed = () => {
  tableData.value = [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  _getArrangeDetail();
  getTableData();
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.main {
  display: flex;
  flex-direction: column;
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
</style>
