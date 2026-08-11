<!-- 任务征集-考场征集管理-确认征集Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="1300"
  >
    <div class="main">
      <div class="main-left">
        <ProTable ref="proTable1" :columns="columns1" :data="tableData1" :pagination="false" :tool-button="false" height="400px">
          <!-- 表格 header -->
          <template #tableHeader>
            <span>共需考场数量：{{ dialogProps.row.needNum || 0 }} </span>
            <span>已征集考场数量：{{ dialogProps.row.collectNum || 0 }} </span>
            <span>已确认考场数量：{{ confirmedClassroomNum }} </span>
          </template>
          <!-- 表格操作 -->
          <template #operation="scope">
            <el-button type="danger" link @click.prevent="onDeleteClick(scope.$index, scope.row)">清除</el-button>
          </template>
        </ProTable>
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
            <el-form-item label="所属学院: ">
              <el-select v-model="curSelectedDepartmentId" placeholder="请选择" clearable>
                <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </template>
        </ProTable>
      </div>
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="tsx" name="ConfirmRecruitDialog">
import { Classroom } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ElMessage } from "element-plus";
import { ref, reactive, computed } from "vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import InputInteger from "@/components/InputInteger.vue";
import { getClassroomListByTaskId } from "@/api/modules/classroom";

interface DialogProps {
  title: string;
  row: Partial<Classroom.ResRecruitManageList>;
  api?: (params: any) => Promise<any>;
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
// 表格配置项
const columns1 = reactive<ColumnProps<any>[]>([
  { prop: "roomName", label: "教室名称" },
  {
    prop: "tchNum",
    label: "监考老师数量",
    render: scope => {
      return (
        <InputInteger model-value={scope.row.tchNum} onUpdate:modelValue={val => (scope.row.tchNum = val)} style="width: 110px" />
      );
    }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 100 }
]);

const confirmedClassroomNum = computed(() => tableData1.value.length);

// 点击-清除
const onDeleteClick = (index: number, row) => {
  // 新增元素（右侧表格）
  row.tchNum = 0;
  tableData2.value?.push(row);
  // 移除当前行（左侧表格）
  tableData1.value.splice(index, 1);
};

// ===================== 右侧表格逻辑 ======================
const proTable2 = ref<ProTableInstance>();
const tableData2 = ref<any>([]);
const { departmentEnum } = useDepartmentEnum();
const curSelectedDepartmentId = ref<string>();
// 表格配置项
const columns2 = reactive<ColumnProps<Classroom.ResClassroomList>[]>([
  { prop: "roomName", label: "教室名称" },
  { prop: "deptId", label: "所属学院", enum: departmentEnum },
  { prop: "campusName", label: "所属校区" },
  {
    prop: "xx",
    label: "组数/按组容量",
    render: scope => (
      <div>
        {scope.row.groupNum || "--"}/{scope.row.groupCapacity || "--"}
      </div>
    ),
    width: 130
  },
  { prop: "capacity", label: "按位容量", width: 90 }
]);

// 过滤搜索的表格数据
const filteredTableData2 = computed(() => {
  let result = tableData2.value;
  if (curSelectedDepartmentId.value != null) {
    result = result.filter(item => item.deptId === curSelectedDepartmentId.value);
  }
  return result;
});

// 表格行被双击
const onRowDblClick = (row: any) => {
  // 新增元素（左侧表格）
  row.tchNum = 2;
  tableData1.value?.push(row);
  // 移除当前行（右侧表格）
  const indexToRemove = tableData2.value.findIndex(item => item.id === row.id);
  if (indexToRemove !== -1) {
    tableData2.value.splice(indexToRemove, 1);
  }
};

// ===================== 共用逻辑 ======================

// 提交数据
const handleSubmit = async () => {
  try {
    const roomList = tableData1.value
      .filter(item => item.tchNum != null)
      .map(item => ({ id: item.id, tchNum: String(item.tchNum) }));
    const params = {
      taskId: dialogProps.value.row.taskId!,
      roomList
    };
    await dialogProps.value.api!(params);
    ElMessage.success({ message: `${dialogProps.value.title}成功！` });
    dialogProps.value.getTableList!();
    dialogVisible.value = false;
  } catch (error: any) {
    ElMessage.error(error.msg || "操作失败");
  }
};

// 获取表格数据
const getTableData = async (type = "unconfirm") => {
  const params = {
    taskId: dialogProps.value.row.taskId!,
    roomState: "1", // 教室状态 1-启用 2-禁用
    state: type === "unconfirm" ? "1" : "2"
  };
  const { list } = await getClassroomListByTaskId(params);
  return list || [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  tableData1.value = await getTableData("confirmed");
  tableData2.value = await getTableData("unconfirm");
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.main {
  display: flex;
  justify-content: space-between;
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
