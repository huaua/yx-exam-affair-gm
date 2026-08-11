<template>
  <div v-if="columns.length" ref="rootRef" class="card table-search">
    <el-form ref="formRef" :model="searchParam" label-width="0">
      <Grid ref="gridRef" compact :collapsed="collapsed" :gap="[12, 0]" :cols="searchCol">
        <GridItem v-for="(item, index) in columns" :key="item.prop" v-bind="getResponsive(item)" :index="index">
          <el-form-item>
            <SearchFormItem :column="item" :search-param="searchParam" />
          </el-form-item>
        </GridItem>
        <GridItem suffix>
          <div class="operation">
            <el-button type="primary" :icon="Search" @click="search"> 搜索 </el-button>
            <el-button :icon="Delete" @click="reset"> 重置 </el-button>
            <el-button v-if="showCollapse" type="primary" link class="search-isOpen" @click="collapsed = !collapsed">
              {{ collapsed ? "展开" : "合并" }}
              <el-icon class="el-icon--right">
                <component :is="collapsed ? ArrowDown : ArrowUp"></component>
              </el-icon>
            </el-button>
          </div>
        </GridItem>
      </Grid>
    </el-form>
  </div>
</template>
<script setup lang="ts" name="SearchForm">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { ColumnProps } from "@/components/ProTable/interface";
import { BreakPoint } from "@/components/Grid/interface";
import { Delete, Search, ArrowDown, ArrowUp } from "@element-plus/icons-vue";
import SearchFormItem from "./components/SearchFormItem.vue";
import Grid from "@/components/Grid/index.vue";
import GridItem from "@/components/Grid/components/GridItem.vue";

interface ProTableProps {
  columns?: ColumnProps[]; // 搜索配置列
  searchParam?: { [key: string]: any }; // 搜索参数
  searchCol: number | Record<BreakPoint, number>;
  search: (params: any) => void; // 搜索方法
  reset: (params: any) => void; // 重置方法
}

// 默认值
const props = withDefaults(defineProps<ProTableProps>(), {
  columns: () => [],
  searchParam: () => ({})
});

// 获取响应式设置
const getResponsive = (item: ColumnProps) => {
  return {
    span: item.search?.span,
    offset: item.search?.offset ?? 0,
    xs: item.search?.xs,
    sm: item.search?.sm,
    md: item.search?.md,
    lg: item.search?.lg,
    xl: item.search?.xl
  };
};

// 是否默认折叠搜索项（默认合并状态）
const collapsed = ref(true);
// 是否显示展开/合并按钮（仅当查询项超过一行时）
const showCollapse = ref(false);
const gridRef = ref();
const rootRef = ref<HTMLElement>();

// 根据实际渲染宽度，计算是否需要折叠按钮；
// 1) 所有项能在一行内放下：不显示折叠按钮（showCollapse=false），全部展示。
// 2) 放不下：默认合并（collapsed=true）到第一行能放下的项数 + 显示"展开"按钮；
//    点开后全部展示 + 按钮变为"合并"。
// 注意：测量时必须先临时 setHiddenIndex(-1) 全部展示，否则被 v-show=false 的项不参与布局，
//     会导致 firstRowCount 偏小、折叠后只显示"前 5 项 + 按钮"而非"前 6 项 + 按钮"。
const recalcCollapse = async () => {
  const gridEl = gridRef.value?.$el as HTMLElement | undefined;
  if (!gridEl) return;
  const gridChildren = Array.from(gridEl.children) as HTMLElement[];
  if (gridChildren.length <= 1) {
    showCollapse.value = false;
    gridRef.value?.setHiddenIndex?.(-1);
    return;
  }
  // 临时全部展示用于测量首行能放下的项数
  gridRef.value?.setHiddenIndex?.(-1);
  await nextTick();
  const firstTop = gridChildren[0].offsetTop;
  let firstRowCount = 0;
  for (const el of gridChildren) {
    if (Math.abs(el.offsetTop - firstTop) < 1) firstRowCount++;
    else break;
  }
  // 减掉 suffix（按钮组），suffix 始终放最后
  const dataItemCount = firstRowCount - 1;
  if (dataItemCount <= 0 || dataItemCount >= props.columns.length) {
    // 所有数据项均能放下，无需折叠
    showCollapse.value = false;
    gridRef.value?.setHiddenIndex?.(-1);
    return;
  }
  // 存在换行：默认合并（折叠）到第一行能放下的项数，并显示"展开"按钮
  showCollapse.value = true;
  if (collapsed.value) {
    gridRef.value?.setHiddenIndex?.(dataItemCount);
  } else {
    gridRef.value?.setHiddenIndex?.(-1);
  }
};

// 监听 collapsed 变化，同步设置 hiddenIndex
watch(collapsed, val => {
  if (!showCollapse.value) return;
  if (val) {
    // 折叠到一行
    recalcCollapse();
  } else {
    gridRef.value?.setHiddenIndex?.(-1);
  }
});

// 监听 columns 变化重新计算
watch(
  () => props.columns.length,
  () => nextTick(recalcCollapse)
);

let resizeObserver: ResizeObserver | null = null;

onMounted(() => {
  nextTick(recalcCollapse);
  if (rootRef.value && typeof ResizeObserver !== "undefined") {
    resizeObserver = new ResizeObserver(() => recalcCollapse());
    resizeObserver.observe(rootRef.value);
  }
});

onBeforeUnmount(() => {
  resizeObserver?.disconnect();
  resizeObserver = null;
});

defineExpose({ recalcCollapse });
</script>
