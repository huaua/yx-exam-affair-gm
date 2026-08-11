<!-- 考场编排-编排查看 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :search-col="{ xs: 4, sm: 5, md: 5, lg: 5, xl: 6 }" :request-api="getTableList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button v-show="isShowExportBtn" type="primary" :icon="Download" plain @click="onDownloadClick">导出</el-button>
      </template>
    </ProTable>
    <ArrangeDialog ref="arrangeDialogRef" />
    <ViewArrangeDialog ref="viewArrangeDialogRef" />
    <PreImportDialog ref="preImportDialogRef" />
    <ImportExcel ref="importDialogRef" />
  </div>
</template>

<script setup lang="tsx" name="arrangeView">
import { Arrange } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Download } from "@element-plus/icons-vue";
import { useDownload } from "@/hooks/useDownload";
import { useAllRecruitTaskEnum } from "@/hooks/useCustomEnum";
import ProTable from "@/components/ProTable/index.vue";
import { getArrangeViewList, exportArrangeInfo } from "@/api/modules/arrange";

interface SearchParams {
  [key: string]: any;
}

const arrangeTypeEnum = [
  {
    label: "按组",
    value: "1"
  },
  {
    label: "按位",
    value: "2"
  }
];

let searchedParams: SearchParams | undefined = {};
const { allRecruitTaskEnum } = useAllRecruitTaskEnum();
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<Arrange.ResClassroomArrangeList>[]>([
  {
    prop: "_xxx",
    label: "学院-专业-方向",
    render: scope => {
      const deptName = String(scope.row.deptName || "");
      const profName = String(scope.row.profName || "");
      const directionName = String(scope.row.directionName || "");
      return `${deptName}-${profName}${directionName ? `-${directionName}` : ""}`;
    }
  },
  { prop: "taskDay", label: "日期", width: 130, render: scope => (scope.row.taskDay && scope.row.taskDay.slice(0, 10)) || "--" },
  { prop: "arrangeType", label: "类型", enum: arrangeTypeEnum, width: 80 },
  { prop: "groupNum", label: "组数", width: 80 },
  { prop: "stuNum", label: "考生数", width: 80 },
  { prop: "roomName", label: "教室名称", search: { order: 5, el: "input", props: { placeholder: "请输入" } } },
  { prop: "arrangeNoStr", label: "编号", width: 150 },
  { prop: "tchNum", label: "监考老师数量", width: 90 },
  { prop: "arrangeTchName", label: "监考老师", showOverflowTooltip: false, width: 150 },
  {
    prop: "roomTaskId",
    label: "考场征集任务",
    enum: allRecruitTaskEnum,
    search: { el: "tree-select", order: 1, props: { filterable: true, placeholder: "请选择" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "deptName",
    label: "学院",
    search: { order: 2, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "profName",
    label: "专业",
    search: { order: 3, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "directionName",
    label: "方向",
    search: { order: 4, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  },
  {
    prop: "tchName",
    label: "监考老师",
    search: { order: 6, el: "input", props: { placeholder: "请输入" } },
    isShow: false,
    isSetting: false
  }
]);

const isShowExportBtn = computed(() => proTable.value?.tableData?.length);

const getTableList = (params: any) => {
  let newParams = JSON.parse(JSON.stringify(params));
  if (!newParams.roomTaskId) {
    ElMessage.warning("请选择考场征集任务");
    return Promise.resolve({ list: [] });
  }
  searchedParams = { ...newParams };
  searchedParams?.curPage && delete searchedParams.curPage;
  searchedParams?.pageSize && delete searchedParams.pageSize;
  return getArrangeViewList(newParams);
};

// 导出
const onDownloadClick = () => {
  ElMessageBox.confirm("确认导出编排数据?", "温馨提示", { type: "warning" }).then(() =>
    useDownload(exportArrangeInfo, "专业考场监考老师编排导出", searchedParams)
  );
};
</script>
