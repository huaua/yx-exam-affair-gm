<!-- 任务征集-考场征集 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格 header -->
      <template #tableHeader v-if="isShowTableHeader"> 该考试任务共上报了 {{ publishedClassroomNum }} 考场 </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <div v-if="scope.row.publishState === PublishStatusEnum.PUBLISHED_CONFIRMED">--</div>
        <el-button
          v-else
          :type="scope.row.publishState === PublishStatusEnum.UNPUBLISH ? 'primary' : 'danger'"
          plain
          @click="changeStatus(scope.row)"
        >
          {{ scope.row.publishState === PublishStatusEnum.UNPUBLISH ? "确认" : "取消" }}上报
        </el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="tsx" name="classroomRecruit">
import { Classroom } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useLevelEnum } from "@/hooks/useEnum";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import {
  getClassroomRecruitList,
  publishClassroom,
  unpublishClassroom,
  getPublishedRecruitTaskAllList
} from "@/api/modules/classroom";

// 考场上报状态
enum PublishStatusEnum {
  UNPUBLISH = "1", // 未上报
  PUBLISHED_UNCONFIRM = "2", // 已上报，未确认
  PUBLISHED_CONFIRMED = "3" // 已上报，已确认
}

/**
 * @description：考场上报状态
 */
const publishStatusEnum = [
  { label: "未上报", value: PublishStatusEnum.UNPUBLISH },
  { label: "已上报", value: PublishStatusEnum.PUBLISHED_UNCONFIRM },
  { label: "已征集", value: PublishStatusEnum.PUBLISHED_CONFIRMED }
];

const { levelEnum } = useLevelEnum();
const taskEnum = ref<any>();
const publishedClassroomNum = ref(0);
const proTable = ref<ProTableInstance>();

const isShowTableHeader = computed(() => {
  return proTable.value?.tableData?.length;
});

// 表格配置项
const columns = reactive<ColumnProps<Classroom.ResClassroomRecruitList>[]>([
  { prop: "roomName", label: "教室名称" },
  { prop: "campusName", label: "所属校区" },
  {
    prop: "levelCode",
    label: "层级",
    enum: levelEnum,
    width: 120
  },
  {
    prop: "xx",
    label: "组数/按组容量",
    render: scope => (
      <div>
        {scope.row.groupNum || "--"}/{scope.row.groupCapacity || "--"}
      </div>
    ),
    width: 140
  },
  { prop: "capacity", label: "按位容量", width: 90 },
  {
    prop: "taskId",
    label: "考试任务",
    enum: taskEnum,
    search: { el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "publishState",
    label: "教室状态",
    enum: publishStatusEnum,
    search: { el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

onMounted(() => {
  _getPublishedRecruitTaskAllList();
});

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.taskId) {
    ElMessage.warning("请选择考试任务");
    publishedClassroomNum.value = 0;
    return Promise.resolve({ list: [] });
  }
  return getClassroomRecruitList(newParams).then(res => {
    publishedClassroomNum.value = res.obj;
    return res;
  });
};

// 切换发布状态
const changeStatus = async (row: Classroom.ResClassroomRecruitList) => {
  const taskId = proTable.value?.searchParam.taskId;
  let api = unpublishClassroom;
  let message = "取消上报";
  if (row.publishState === PublishStatusEnum.UNPUBLISH) {
    api = publishClassroom;
    message = "上报";
  }
  await useHandleData(
    api,
    {
      taskId: taskId,
      roomId: row.roomId
    },
    `${message}【${row.roomName}】考场`
  );
  proTable.value?.getTableList();
};

// 获取所有已发布的任务列表
const _getPublishedRecruitTaskAllList = async () => {
  const { list = [] } = await getPublishedRecruitTaskAllList();
  taskEnum.value = list.map((item: any) => ({
    label: item.taskName,
    value: item.taskId
  }));
};
</script>
