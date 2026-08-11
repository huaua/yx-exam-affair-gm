<!-- 考生管理-考生管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getStudentList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" @click="onSyncStudentClick">信息同步</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" link :icon="View" @click="openViewDialog(scope.row)">查看</el-button>
      </template>
    </ProTable>
    <PreSyncDialog ref="preSyncDialogRef" />
    <SyncDialog ref="syncDialogRef" />
    <ViewDialog ref="viewDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="student">
import { Student } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { StudentSyncEnum } from "@/enums";
import { ref, reactive, onMounted } from "vue";
import { View } from "@element-plus/icons-vue";
import ProTable from "@/components/ProTable/index.vue";
import PreSyncDialog from "./components/PreSyncDialog.vue";
import SyncDialog from "./components/SyncDialog.vue";
import ViewDialog from "./components/ViewDialog.vue";
import { getStudentList, checkStudentInfoSyncState, getExamList } from "@/api/modules/student";

// ProTable 实例
const proTable = ref<ProTableInstance>();

const examEnum = ref<{ value: string; label: string }[]>([]);

// 表格配置项
const columns = reactive<ColumnProps<Student.ResStudentList>[]>([
  { type: "index", label: "序号", width: 60 },
  { prop: "shengFenMC", label: "省份", search: { order: 2, el: "input", props: { placeholder: "请输入" } } },
  { prop: "xingMing", label: "姓名", search: { order: 3, el: "input", props: { placeholder: "请输入" } } },
  { prop: "shenFenZH", label: "证件号", search: { order: 4, el: "input", props: { placeholder: "请输入" } } },
  { prop: "shouJi", label: "手机号" },
  { prop: "zhunKaoZH", label: "准考证号", search: { order: 5, el: "input", props: { placeholder: "请输入" } } },
  { prop: "zhuanYeMC", label: "专业" },
  { prop: "kaoShiRQSM", label: "考试时间" },
  {
    prop: "kaoShiID",
    label: "专业考试",
    enum: examEnum,
    search: { order: 1, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

onMounted(() => {
  fetchExamEnum();
});

// 点击-信息同步
const onSyncStudentClick = async () => {
  const res = await checkStudentInfoSyncState({});
  const { state, percent, kaoShiID, kaoShiMC, errMsg = "" } = res?.obj || {};
  if (state === StudentSyncEnum.NOT_STARTED) {
    openPreSyncDialog();
  } else {
    openSyncDialog({ kaoShiID, kaoShiMC, state, percent, errMsg });
  }
};

// 打开任务选择弹窗
const preSyncDialogRef = ref<InstanceType<typeof PreSyncDialog> | null>(null);
const openPreSyncDialog = () => {
  preSyncDialogRef.value?.acceptParams({ openSyncDialog });
};

// 打开-同步弹窗
const syncDialogRef = ref<InstanceType<typeof SyncDialog> | null>(null);
const openSyncDialog = ({
  kaoShiID,
  kaoShiMC,
  state,
  percent,
  errMsg
}: {
  kaoShiID: string;
  kaoShiMC: string;
  state: string;
  percent: string;
  errMsg: string;
}) => {
  const params = {
    title: "查看",
    kaoShiID,
    kaoShiMC,
    state,
    percent,
    errMsg,
    openPreSyncDialog,
    getTableList: proTable.value?.getTableList
  };
  syncDialogRef.value?.acceptParams(params);
};

// 打开-查看
const viewDialogRef = ref<InstanceType<typeof ViewDialog> | null>(null);
const openViewDialog = (row: any = {}) => {
  const params = {
    row: { ...row }
  };
  viewDialogRef.value?.acceptParams(params);
};

// 获取考试枚举
const fetchExamEnum = async () => {
  const res = await getExamList({ reSync: false });
  examEnum.value = (res.list || []).map(item => ({
    value: String(item.kaoShiID),
    label: item.kaoShiMC
  }));
};
</script>
