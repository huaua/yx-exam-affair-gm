<!-- 考场编排-考场监考分配-分配Dialog -->
<template>
  <el-dialog
    v-model="drawerVisible"
    title="分配"
    class="assign-dialog"
    align-center
    :lock-scroll="true"
    :close-on-click-modal="false"
    width="88%"
    :destroy-on-close="true"
    @closed="handleClosed"
  >
    <!-- 顶部信息栏 + 搜索 -->
    <!-- prettier-ignore -->
    <div class="dialog-header-bar">
      <div class="header-info">
        <span>考试征集任务：<b>{{ taskName }}</b></span>
        <span>日期：<b>{{ taskDate }}</b></span>
      </div>
      <div class="header-search">
        <el-input
          v-model="searchObj.nameOrJobNo"
          placeholder="教职工姓名/工号-请输入"
          clearable
          style="width: 200px"
        />
        <el-select
          v-model="searchObj.deptId"
          placeholder="所属学院-请选择"
          clearable
          style="width: 200px"
        >
          <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-select
          v-model="searchObj.invigilationExperience"
          placeholder="监考经验-请选择"
          clearable
          style="width: 160px"
        >
          <el-option label="有" value="1" />
          <el-option label="无" value="2" />
        </el-select>
      </div>
    </div>

    <div class="dialog-body">
      <!-- 左侧：教室列表 -->
      <div class="panel-left">
        <ProTable
          ref="proTable1"
          :columns="columns1"
          :data="tableData1"
          :pagination="false"
          :tool-button="false"
          height="calc(100vh - 340px)"
          :row-class-name="tableRowClassName"
          @row-click="onClassroomRowClick"
        >
          <!-- 监考老师列 -->
          <template #arrangeTchName="scope">
            <div v-if="scope.row.arrangeTchReportId?.length" class="tch-tags">
              <el-tag
                v-for="tchId in scope.row.arrangeTchReportId"
                :key="tchId"
                closable
                size="small"
                type="success"
                effect="plain"
                @close="onRemoveTeacher(scope.row, tchId)"
                >{{ getTeacherNameById(tchId) }}
              </el-tag>
            </div>
            <span v-else>--</span>
          </template>
          <!-- 操作列 -->
          <template #operation="scope">
            <!-- eslint-disable vue/html-closing-bracket-newline -->
            <el-button
              v-if="scope.row.arrangeTchReportId?.length"
              type="danger"
              link
              size="small"
              @click.stop="onClearRoom(scope.row)"
              >清除</el-button
            >
            <span v-else>--</span>
          </template>
        </ProTable>
      </div>

      <!-- 右侧：可用监考老师列表 -->
      <div class="panel-right">
        <ProTable
          ref="proTable2"
          :columns="columns2"
          :data="filteredTableData2"
          :pagination="false"
          :tool-button="false"
          height="calc(100vh - 340px)"
          @row-dblclick="onTeacherDblClick"
        />
      </div>
    </div>

    <template #footer>
      <el-button @click="drawerVisible = false">取消</el-button>
      <el-button type="primary" :disabled="isDisabledSubmit" :loading="submitting" @click="handleSubmit">确定</el-button>
    </template>

    <!-- 监考经验记录弹窗 -->
    <MonitorRecordsDialog ref="recordsDialogRef" />
  </el-dialog>
</template>

<script setup lang="tsx" name="AssignDialog">
import { Arrange, Teacher } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ElMessage, ElMessageBox } from "element-plus";
import { ref, reactive, computed } from "vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import { getTeacherAssignList, saveTeacherAssign, cancelSingleAssign } from "@/api/modules/arrange";
import { getTeacherRecruitViewList } from "@/api/modules/teacher";
import MonitorRecordsDialog from "@/views/teacher/teacher/components/MonitorRecordsDialog.vue";

interface TaskObj {
  roomTaskId: string;
  roomTaskName: string;
  taskTime: string;
}

interface DialogProps {
  title: string;
  mode: string;
  taskObj: TaskObj;
  row: Partial<Arrange.ResTeacherAssignList>;
  getTableList?: () => void;
}

enum TchAssignState {
  UNASSIGN = "1",
  ASSIGNED = "2"
}

const drawerVisible = ref(false);
const submitting = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  mode: "batch",
  taskObj: { roomTaskId: "", roomTaskName: "", taskTime: "" },
  row: {}
});

// 当前选中的教室行（用于指定教室分配模式）
const selectedRoomIndex = ref<number>(-1);

// 记录每个教室从服务端加载时的原始已分配老师集合（用于区分"已有持久化分配"和"本次会话新增的本地分配"）
const originalAssignedMap = ref<Map<string, Set<string>>>(new Map());

// 任务下所有教室已分配老师的 reportId 全集（用于右侧列表过滤，单条模式也需要排除其他教室已分配的老师）
const allAssignedReportIds = ref<Set<string>>(new Set());

// header 信息快捷计算
const taskName = computed(() => dialogProps.value.taskObj?.roomTaskName || "--");
const taskDate = computed(() => dialogProps.value.taskObj?.taskTime?.slice(0, 10) || "--");

// ===================== 左侧教室表格 ======================
const proTable1 = ref<ProTableInstance>();
const tableData1 = ref<Arrange.ResTeacherAssignList[]>([]);

const columns1 = reactive<ColumnProps<Arrange.ResTeacherAssignList>[]>([
  { prop: "roomName", label: "教室名称", minWidth: 120 },
  { prop: "profDirection", label: "学院-专业-方向", minWidth: 140 },
  { prop: "tchNum", label: "需要人数", width: 80 },
  {
    prop: "arrangeTchName",
    label: "监考老师",
    minWidth: 160,
    showOverflowTooltip: false
  },
  {
    prop: "operation",
    label: "操作",
    width: 70,
    fixed: "right"
  }
]);

const getTableData1 = async () => {
  const params = {
    roomTaskId: dialogProps.value.taskObj.roomTaskId!,
    curPage: 1,
    pageSize: 10000
  };
  const { list } = await getTeacherAssignList(params);
  const fullList = list || [];
  // 从全量数据中收集任务下所有已分配老师的 ID（用于右侧列表过滤）
  const assignedIds = new Set<string>();
  fullList.forEach(room => {
    (room.arrangeTchReportId || []).forEach(id => assignedIds.add(id));
  });
  allAssignedReportIds.value = assignedIds;
  // 单条分配模式：只保留目标教室
  let result = fullList;
  if (dialogProps.value.mode === "single" && dialogProps.value.row?.roomReportId) {
    result = result.filter(item => item.roomReportId === dialogProps.value.row!.roomReportId);
  }
  tableData1.value = result;
  // 记录每个教室的原始已分配老师集合（基于过滤后的列表，仅当前显示教室）
  const map = new Map<string, Set<string>>();
  for (const room of tableData1.value) {
    map.set(String(room.roomReportId), new Set(room.arrangeTchReportId || []));
  }
  originalAssignedMap.value = map;
};

// 点击教室行选中（用于指定教室分配模式）
const onClassroomRowClick = (row: Arrange.ResTeacherAssignList) => {
  const idx = tableData1.value.findIndex(item => item.roomReportId === row.roomReportId);
  selectedRoomIndex.value = idx;
};

const tableRowClassName = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex === selectedRoomIndex.value ? "selected-row" : "";
};

// 获取某教室剩余可分配名额
const getRemainingCapacity = (row: Arrange.ResTeacherAssignList) => {
  return Number(row.tchNum || 0) - (row.arrangeTchReportId?.length || 0);
};

// 清除某教室所有已分配老师
const onClearRoom = async (roomRow: Arrange.ResTeacherAssignList) => {
  try {
    await ElMessageBox.confirm(`确认清除【${roomRow.roomName}】的所有监考分配？`, "提示", { type: "warning" });
    const roomKey = String(roomRow.roomReportId);
    const originalSet = originalAssignedMap.value.get(roomKey);
    // 逐个调用取消接口（仅对已持久化的分配）
    const ids = [...(roomRow.arrangeTchReportId || [])];
    for (const tchReportId of ids) {
      if (originalSet?.has(tchReportId)) {
        await cancelSingleAssign({ roomReportId: roomKey, tchReportId });
      }
    }
    // 从本地数据中移除
    roomRow.arrangeTchReportId = [];
    roomRow.assignmentState = "1";
    // 将老师恢复到可用列表
    ids.forEach(tchReportId => {
      const tchItem = tableData2.value.find(t => t.reportId === tchReportId);
      if (tchItem) {
        tchItem.state = TchAssignState.UNASSIGN;
      }
    });
    ElMessage.success("已清除分配");
  } catch {
    // 用户取消
  }
};

// 移除某个教室的某个老师
const onRemoveTeacher = async (roomRow: Arrange.ResTeacherAssignList, tchReportId: string) => {
  try {
    await ElMessageBox.confirm("确认取消该老师的监考分配？", "提示", { type: "warning" });
    const roomKey = String(roomRow.roomReportId);
    // 只有原本就存在于服务端的分配才需要调 API 取消
    // 本次会话新增的本地分配（未持久化）直接本地移除即可
    const isPersisted = originalAssignedMap.value.get(roomKey)?.has(tchReportId);
    if (isPersisted) {
      await cancelSingleAssign({ roomReportId: roomKey, tchReportId });
    }
    // 从本地数据中移除
    const idx = roomRow.arrangeTchReportId?.indexOf(tchReportId);
    if (idx !== undefined && idx > -1) {
      roomRow.arrangeTchReportId!.splice(idx, 1);
    }
    // 如果该教室没有已分配老师了，更新状态
    if (!roomRow.arrangeTchReportId?.length) {
      roomRow.assignmentState = "1";
    }
    // 将该老师恢复到可用列表
    const tchItem = tableData2.value.find(t => t.reportId === tchReportId);
    if (tchItem) {
      tchItem.state = TchAssignState.UNASSIGN;
    }
    ElMessage.success("已取消分配");
  } catch {
    // 用户取消
  }
};

// 根据reportId获取老师姓名
const getTeacherNameById = (reportId: string) => {
  const item = tableData2.value.find(t => t.reportId === reportId);
  return item?.tchName || "--";
};

// ===================== 右侧老师表格 =====================
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const proTable2 = ref<ProTableInstance>();
const tableData2 = ref<Teacher.ResTeacherRecruitViewList[]>([]);
const searchObj = reactive({
  nameOrJobNo: "",
  deptId: undefined as number | undefined,
  invigilationExperience: "" as string
});

const columns2 = reactive<ColumnProps<Teacher.ResTeacherRecruitViewList>[]>([
  { prop: "tchName", label: "教职工姓名", width: 100 },
  { prop: "sex", label: "性别", width: 60 },
  { prop: "jobNo", label: "工号", width: 100 },
  {
    prop: "invigilationExperience",
    label: "监考经验",
    width: 90,
    render: scope => {
      if (scope.row.invigilationExperience === "1") {
        return (
          <el-button type="primary" link size="small" onClick={() => onViewRecords(scope.row)}>
            查看
          </el-button>
        );
      }
      return <span>--</span>;
    }
  },
  { prop: "remark", label: "备注", minWidth: 100, render: scope => scope.row.remark || "--" },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    minWidth: 120
  }
]);

// 过滤：只显示未分配的老师（实时过滤）
const filteredTableData2 = computed(() => {
  let result = tableData2.value.filter(person => person.state === TchAssignState.UNASSIGN);
  const q = searchObj.nameOrJobNo?.trim();
  if (q) {
    result = result.filter(p => p.tchName.includes(q) || p.jobNo.includes(q));
  }
  if (searchObj.deptId != null) {
    result = result.filter(p => p.deptId === String(searchObj.deptId));
  }
  if (searchObj.invigilationExperience) {
    result = result.filter(p => p.invigilationExperience === searchObj.invigilationExperience);
  }
  return result;
});

const getTableData2 = async () => {
  const params = {
    taskTime: dialogProps.value.taskObj.taskTime!,
    curPage: 1,
    pageSize: 10000
  };
  const { list } = await getTeacherRecruitViewList(params);
  // 初始化所有老师为未分配状态
  tableData2.value = (list || []).map(item => ({ ...item, state: TchAssignState.UNASSIGN }));
};

// 根据左侧教室的已分配老师同步右侧老师状态（使用任务下全量已分配集合）
const syncAssignedState = () => {
  tableData2.value.forEach(tch => {
    if (allAssignedReportIds.value.has(tch.reportId)) {
      tch.state = TchAssignState.ASSIGNED;
    }
  });
};

// 双击老师行 -> 分配到目标教室
const onTeacherDblClick = (row: Teacher.ResTeacherRecruitViewList) => {
  if (row.state === TchAssignState.ASSIGNED) {
    ElMessage.warning("该老师已被分配");
    return;
  }

  // 确定目标教室：优先选中的教室，否则按需求数量自动找第一个有剩余名额的
  let targetRoom: Arrange.ResTeacherAssignList | undefined;

  if (selectedRoomIndex.value >= 0) {
    // 指定教室模式：检查选中教室是否有剩余名额
    const room = tableData1.value[selectedRoomIndex.value];
    if (getRemainingCapacity(room) > 0) {
      targetRoom = room;
    } else {
      ElMessage.warning(`【${room.roomName}】已满员，请选择其他教室或取消选中`);
      return;
    }
  }

  // 自动分配模式：找第一个有剩余名额的教室
  if (!targetRoom) {
    targetRoom = tableData1.value.find(r => getRemainingCapacity(r) > 0);
  }

  if (!targetRoom) {
    ElMessage.warning("无剩余可分配监考老师的教室");
    return;
  }

  // 执行分配（本地操作）
  if (!targetRoom.arrangeTchReportId) {
    targetRoom.arrangeTchReportId = [];
  }
  targetRoom.arrangeTchReportId.push(row.reportId);

  // 如果该教室已满，更新状态
  if (getRemainingCapacity(targetRoom) <= 0) {
    targetRoom.assignmentState = "2"; // ASSIGNED
  }

  // 标记老师为已分配
  const tchItem = tableData2.value.find(t => t.reportId === row.reportId);
  if (tchItem) {
    tchItem.state = TchAssignState.ASSIGNED;
  }
};

// ===================== 监考经验记录查看 =====================
const recordsDialogRef = ref<InstanceType<typeof MonitorRecordsDialog>>();

const onViewRecords = (row: Teacher.ResTeacherRecruitViewList) => {
  const resolvedDeptName =
    row.deptName || departmentEnum.value.find(item => String(item.value) === String(row.deptId))?.label || row.deptId;
  recordsDialogRef.value?.acceptParams({
    title: "监考经验记录",
    row: { tchId: row.tchId, tchName: row.tchName, jobNo: row.jobNo, deptId: row.deptId, deptName: resolvedDeptName }
  });
};

// ===================== 提交 =====================
const isDisabledSubmit = computed(() => !tableData1.value.length);

const handleSubmit = async () => {
  // 校验：至少有一个教室有分配
  const hasAssignment = tableData1.value.some(r => r.arrangeTchReportId?.length);
  if (!hasAssignment) {
    ElMessage.warning("请至少为一个教室分配监考老师");
    return;
  }

  submitting.value = true;
  try {
    const params = tableData1.value
      .filter(item => item.arrangeTchReportId?.length)
      .map(item => ({
        roomTaskId: dialogProps.value.taskObj.roomTaskId!,
        roomReportId: item.roomReportId,
        tchReportIdList: item.arrangeTchReportId || []
      }));
    await saveTeacherAssign(params);
    ElMessage.success(`${dialogProps.value.title}成功！`);
    dialogProps.value.getTableList?.();
    drawerVisible.value = false;
  } finally {
    submitting.value = false;
  }
};

const handleClosed = () => {
  tableData1.value = [];
  tableData2.value = [];
  searchObj.nameOrJobNo = "";
  searchObj.deptId = undefined;
  searchObj.invigilationExperience = "";
  selectedRoomIndex.value = -1;
  submitting.value = false;
};

const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  drawerVisible.value = true;
  await Promise.all([getTableData1(), getTableData2()]);
  syncAssignedState();
};

defineExpose({ acceptParams });
</script>

<style lang="scss" scoped>
:deep(.assign-dialog) {
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 48px);
  margin: 0;

  .el-dialog__header,
  .el-dialog__footer {
    flex-shrink: 0;
  }

  .el-dialog__body {
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }
}

.dialog-header-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px 24px;
  margin-bottom: 12px;
  font-size: 14px;
  .header-info {
    display: flex;
    gap: 20px;
    b {
      color: var(--el-color-primary);
    }
  }
  .header-search {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-left: auto;
  }
}

.dialog-body {
  display: flex;
  gap: 16px;

  .panel-left {
    flex: 1.2;
    min-width: 0;
  }

  .panel-right {
    flex: 1;
    min-width: 0;
  }

  .tch-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    :deep(.el-tag) {
      cursor: pointer;
    }
  }

  :deep(.el-table__row.selected-row) {
    background-color: var(--el-color-primary) !important;
    color: #fff !important;

    & > td {
      background-color: var(--el-color-primary) !important;
      color: #fff !important;
    }

    & > td:first-child {
      box-shadow: inset 6px 0 0 0 #fff;
    }

    &:hover > td {
      background-color: var(--el-color-primary) !important;
    }

    .el-tag {
      color: var(--el-color-primary) !important;
      background-color: #fff !important;
      border-color: #fff !important;
    }

    .el-button {
      color: #fff !important;
    }
  }

  :deep(.el-table__row) {
    cursor: pointer;
  }
}
</style>
