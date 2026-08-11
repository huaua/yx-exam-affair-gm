<!-- 任务征集-评委查询 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader v-if="isShowHeader">
        该评分任务：确定参加：{{ countObject.attendNum ?? "--" }}人；不参加：{{ countObject.unattendNum ?? "--" }}人；待定：{{
          countObject.notSureNum ?? "--"
        }}人
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button
          v-if="scope.row.confirmState !== JudgeResponseState.NOT_RESPONDED"
          type="primary"
          :icon="Refresh"
          link
          @click="onRestClick(scope.row)"
        >
          重置
        </el-button>
      </template>
    </ProTable>
  </div>
</template>

<script setup lang="tsx" name="judgeQuery">
import { JudgeExpert } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { ElMessage } from "element-plus";
import { Refresh } from "@element-plus/icons-vue";
import { JudgeResponseState } from "@/enums";
import { useHandleData } from "@/hooks/useHandleData";
import { useGoodSubjectEnum } from "@/hooks/useEnum";
import { useAllJudgeRecruitTaskEnum } from "@/hooks/useCustomEnum";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import { judgeResponseStatus } from "@/utils/dict";
import { getJudgeQueryList, restJudgeResponse } from "@/api/modules/judgeExpert";

const { allJudgeRecruitTaskEnum } = useAllJudgeRecruitTaskEnum();
const { goodSubjectEnum } = useGoodSubjectEnum();
const { departmentEnum } = useDepartmentEnum("expertDataBase");
const proTable = ref<ProTableInstance>();
const isShowHeader = ref(false);
const countObject = ref({
  attendNum: "--",
  unattendNum: "--",
  notSureNum: "--"
});

// 表格配置项
const columns = reactive<ColumnProps<JudgeExpert.ResJudgeResponseList>[]>([
  { prop: "expertName", label: "姓名", width: 90 },
  { prop: "sex", label: "性别", width: 60 },
  { prop: "phoneNo", label: "手机号", width: 120 },
  { prop: "typeStr", label: "类型", width: 80 },
  {
    prop: "goodSubjects",
    label: "擅长科目",
    enum: goodSubjectEnum,
    isFilterEnum: false,
    search: { order: 3, el: "tree-select", props: { placeholder: "请输入" } },
    fieldNames: { label: "label", value: "label" }
  },
  {
    prop: "subjectName",
    label: "抽取科目"
  },
  {
    prop: "taskTimeRemark",
    label: "评分日期",
    width: 180
  },
  {
    prop: "deptId",
    label: "所属学院",
    enum: departmentEnum,
    search: { order: 2, el: "tree-select", props: { filterable: true, placeholder: "请选择" } }
  },
  {
    prop: "judgeTaskId",
    label: "任务",
    enum: allJudgeRecruitTaskEnum,
    search: { order: 1, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  { prop: "drawTimes", label: "抽取次数", width: 90 },
  {
    prop: "confirmState",
    label: "应答状态",
    tag: true,
    enum: judgeResponseStatus,
    search: { order: 4, el: "tree-select", props: { filterable: true, placeholder: "请选择" } },
    width: 90
  },
  { prop: "operation", label: "操作", fixed: "right", width: 150 }
]);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  isShowHeader.value = !!newParams.judgeTaskId;
  if (!newParams.judgeTaskId) {
    ElMessage.warning("请选择任务");
    return Promise.resolve({ list: [] });
  }
  return getJudgeQueryList(newParams).then(res => {
    countObject.value = {
      attendNum: res?.list?.attendNum,
      unattendNum: res?.list?.unattendNum,
      notSureNum: res?.list?.notSureNum
    };
    return { list: res?.list?.list || [], page: res?.page || {} };
  });
};

// 点击-重置
const onRestClick = async (row: JudgeExpert.ResJudgeResponseList) => {
  await useHandleData(
    restJudgeResponse,
    {
      id: row.id
    },
    `重置【${row.expertName}】的应答状态`
  );
  proTable.value?.getTableList();
};
</script>
