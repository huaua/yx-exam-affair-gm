<!-- 监考老师管理-监考老师征集任务-新增/编辑Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    :width="getDialogWidth()"
    @closed="handleClosed"
  >
    <el-tabs v-if="tabList.length" type="border-card" v-model="activeName">
      <el-tab-pane v-for="(tab, i) in tabList" :key="i" :label="tab.dateStr" :name="i">
        <div class="tips">
          <div v-if="isAddMode">当天征集人数为：{{ sumByKey(list![i], "needNum") }}</div>
          <div v-if="isEditMode">要求征集人数：{{ sumByKey(list![i], "needNum") }}</div>
          <div v-if="isEditMode">已征集人数：{{ sumByKey(list![i], "collectNum") }}</div>
          <div v-if="isEditMode" class="tips-force">
            <span>强制征集 : </span>
            <el-switch
              v-model="tab.tabForceFlag"
              :active-value="ForceFlagEnum.ACTIVE"
              :inactive-value="ForceFlagEnum.INACTIVE"
              @change="e => onCurTabForceFlagChange(i, e)"
            />
          </div>
        </div>
        <el-form :ref="collectRef" class="form" label-width="90px" label-suffix=" :" :inline="true" :model="list![i]">
          <DynamicInput
            v-for="(item, index) in list![i]"
            :key="item._id"
            :current-index="index"
            :length="list![i].length"
            @on-plus="handlePlus"
            @on-minus="handleMinus"
          >
            <el-form-item label="学院" :prop="`${index}.deptId`" :rules="rules.deptId">
              <el-select
                v-model="item.deptId"
                @change="onDeptChange(item)"
                placeholder="请选择学院"
                :disabled="shouldDisabledDeptSelect(item.deptId)"
                filterable
                clearable
              >
                <el-option
                  v-for="department in curTabDepartmentEnum"
                  :key="department.value"
                  :label="department.label"
                  :value="department.value"
                  :disabled="department.disabled"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="征集人数" :prop="`${index}.needNum`" :rules="rules.needNum">
              <InputInteger v-model="item.needNum" />
            </el-form-item>
            <el-form-item label="强制征集" :prop="`${index}.forceFlag`">
              <el-switch v-model="item.forceFlag" :active-value="ForceFlagEnum.ACTIVE" :inactive-value="ForceFlagEnum.INACTIVE" />
            </el-form-item>
            <el-form-item v-if="isEditMode" label="已征集人数" :prop="`${index}.collectNum`">
              <InputInteger v-model="item.collectNum" disabled />
            </el-form-item>
          </DynamicInput>
        </el-form>
      </el-tab-pane>
    </el-tabs>
    <el-empty v-else description="暂无数据" />
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AddOrEditDialog">
import { Common } from "@/api/interface";
import { ElMessage } from "element-plus";
import { ref, reactive, computed } from "vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { cloneDeep } from "lodash-es";
import { nanoid } from "nanoid";
import DynamicInput from "@/components/DynamicInput.vue";
import { editTeacherRecruitTaskDetail } from "@/api/modules/teacher";

interface FormItem {
  _id: string;
  taskDetailId?: string;
  taskTime?: string; // 征集日期
  deptId?: string; // 学院id
  deptName?: string;
  needNum?: number | undefined; // 征集人数
  collectNum?: number | undefined; // 已征集人数
  forceFlag?: string; // 是否强制征集: 1-强制 2-不强制
}

interface TabItem {
  tabForceFlag: string; // 是否强制征集所有
  date: string;
  dateStr: string;
}

interface DialogProps {
  type: string;
  title: string;
  row: any;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

interface DepartmentEnumItem extends Common.enumDict {
  disabled?: boolean;
}

enum ForceFlagEnum {
  ACTIVE = "1",
  INACTIVE = "2"
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  title: "",
  row: {
    tchTaskId: "",
    tchTaskName: "",
    dates: ["", ""] // 征集日期 [开始日期，结束日期]
  }
});

const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const activeName = ref(0);
const dateList = ref<any>([]);
const tabList = ref<TabItem[]>([]);
const list = ref<FormItem[][]>([[]]);
let cacheList: FormItem[][] = [[]]; // 编辑前列表
const formRefList = ref<any>([]);
const rules = reactive({
  deptId: [{ required: true, message: "请选择学院" }],
  needNum: [{ required: true, message: "请填写征集人数" }],
  forceFlag: [{ required: true, message: "请选择是否强制征集" }]
});

const isAddMode = computed(() => dialogProps.value.type === "add");
const isEditMode = computed(() => dialogProps.value.type === "edit");

// 学院下拉列表（disabled当前tab已选择的学院）
const curTabDepartmentEnum = computed(() =>
  departmentEnum.value.map(departObj => {
    let result: DepartmentEnumItem = departObj;
    if (list.value[activeName.value]?.some(item => item.deptId === departObj.value)) {
      result.disabled = true;
    } else {
      result.disabled = false;
    }
    return result;
  })
);

const collectRef = el => {
  if (!el || formRefList.value.includes(el)) {
    return;
  }
  formRefList.value.push(el);
};

// 提交数据（新增/编辑）
const handleSubmit = async () => {
  const isValidated = await validateAllForm();
  if (!isValidated) {
    return;
  }
  try {
    const params = {
      tchTaskId: dialogProps.value.row.tchTaskId || undefined // 编辑时必传
    };
    if (isAddMode.value) {
      params["tchTaskName"] = dialogProps.value.row?.tchTaskName;
      params["taskStartTime"] = dialogProps.value.row?.dates[0];
      params["taskEndTime"] = dialogProps.value.row?.dates[1];
      params["tchTaskDetail"] = list.value!.flat().map(item => ({
        deptId: item.deptId,
        deptName: item.deptName,
        forceFlag: item.forceFlag,
        needNum: String(item.needNum),
        taskTime: item.taskTime
      }));
    } else if (isEditMode.value) {
      params["tchTaskDetail"] = list.value!.flat().map(item => ({
        tchTaskId: dialogProps.value.row.tchTaskId || undefined,
        taskDetailId: item.taskDetailId,
        deptId: item.deptId,
        deptName: item.deptName,
        forceFlag: item.forceFlag,
        needNum: String(item.needNum),
        taskTime: item.taskTime
      }));
    }
    await dialogProps.value.api!(params);
    ElMessage.success({ message: `${dialogProps.value.title}成功！` });
    dialogProps.value.getTableList!();
    dialogVisible.value = false;
  } catch (error) {
    console.log(error);
  }
};

const handlePlus = () => {
  list.value![activeName.value].push({
    _id: nanoid(),
    deptId: "",
    needNum: undefined,
    forceFlag: ForceFlagEnum.ACTIVE,
    taskTime: dateList.value[activeName.value]
  });
};

const handleMinus = (index: number) => {
  list.value![activeName.value].splice(index, 1);
};

const onDeptChange = (itemObj: FormItem) => {
  // 添加学院名称
  itemObj["deptName"] = departmentEnum.value.find(item => item.value === itemObj.deptId)?.label;
  // 清空taskDetailId
  itemObj["taskDetailId"] = undefined;
  if (isEditMode.value) {
    // 编辑时，更新对应行数据
    const sameItem = cacheList![activeName.value].find(item => item.deptId === itemObj.deptId);
    itemObj.taskDetailId = sameItem?.taskDetailId;
    itemObj.needNum = sameItem?.needNum;
    itemObj.collectNum = sameItem?.collectNum;
  }
};

const shouldDisabledDeptSelect = (deptId: string | undefined) => {
  if (isAddMode.value) {
    return false;
  }
  return cacheList && cacheList[activeName.value]?.some(item => item.deptId === deptId);
};

// 批量修改当前Tab的强制征集状态
const onCurTabForceFlagChange = (index, val) => {
  list.value[index].forEach(item => {
    item.forceFlag = val;
  });
};

const validateAllForm = async () => {
  // 并行验证所有表单，返回结果数组
  const validList = await Promise.all(formRefList.value.map(formRef => formRef?.validate(valid => valid)));
  // 查找第一个失败的表单索引
  const firstFailIndex = validList.indexOf(false);
  if (firstFailIndex !== -1) {
    // 将所有失败的标签信息提取
    const failTabs: TabItem[] = tabList.value.filter((_, index) => !validList[index]);
    // 生成错误的日期字符串
    const errorDateStr = failTabs.map(item => item.dateStr).join("、");
    activeName.value = firstFailIndex;
    ElMessage.error({ message: `请检查${errorDateStr}内容` });
    return false;
  }
  return true;
};

const handleClosed = () => {
  activeName.value = 0;
  formRefList.value = [];
};

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  if (isAddMode.value) {
    dateList.value = fillDateRange(dialogProps.value.row.dates);
    tabList.value = generateTabList(dateList.value);
    list.value = generateEmptyList(dateList.value);
  } else if (isEditMode.value) {
    const res = await editTeacherRecruitTaskDetail({ tchTaskId: dialogProps.value.row.tchTaskId });
    const _list = res?.list || [];
    dateList.value = _list.map(item => item.day);
    tabList.value = generateTabList(dateList.value);
    list.value = _list.map(subList => subList.list.map(item => ({ ...item, _id: nanoid() }))); // 添加_id, 用作循环key
    cacheList = cloneDeep(list.value);
  }
};

defineExpose({
  acceptParams
});

// 填充日期 ["2024-12-30 00:00:00", "2025-01-01 00:00:00"] --> ["2024-12-30 00:00:00", "2024-12-31 00:00:00", "2025-01-01 00:00:00"]
const fillDateRange = dates => {
  const startDate = new Date(dates[0]);
  const endDate = new Date(dates[1]);
  const result: any = [];
  // 使用一个循环从开始日期到结束日期，依次将每一天加入结果数组
  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
    // 格式化当前日期为 "YYYY-MM-DD HH:MM:SS"
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0"); // 月份从0开始，补齐2位
    const day = String(d.getDate()).padStart(2, "0");
    const hours = String(d.getHours()).padStart(2, "0");
    const minutes = String(d.getMinutes()).padStart(2, "0");
    const seconds = String(d.getSeconds()).padStart(2, "0");
    result.push(`${year}-${month}-${day} ${hours}:${minutes}:${seconds}`);
  }
  return result;
};

const generateTabList = dateArray => {
  return dateArray.map(dateString => {
    const date = new Date(dateString);
    // 获取月份和日期（月份从 0 开始，因此要加 1）
    const month = date.getMonth() + 1;
    const day = date.getDate();
    // 返回包含 date 和 dateStr 的对象
    return {
      tabForceFlag: ForceFlagEnum.INACTIVE,
      date: dateString,
      dateStr: `${month}月${day}日`
    };
  });
};

const generateEmptyList = dateArray => {
  return dateArray.map(dateString => {
    return [{ deptId: "", needNum: undefined, forceFlag: ForceFlagEnum.ACTIVE, taskTime: dateString }];
  });
};

const sumByKey = (arr, key) => {
  if (!Array.isArray(arr)) {
    throw new Error("第一个参数必须是 array");
  }
  return arr.reduce((acc, obj) => acc + (Number(obj[key]) || 0), 0);
};

const getDialogWidth = () => {
  let result = "1100";
  if (isAddMode.value) {
    result = "850";
  }
  return result;
};
</script>

<style lang="scss" scoped>
:deep(.el-form-item) {
  width: 210px;
}
:deep(.el-select__wrapper) {
  width: 150px;
}
.tips {
  display: flex;
  align-items: center;
  padding-bottom: 10px;
  padding-left: 30px;
  div {
    min-width: 168px;
    padding-right: 4px;
  }
  &-force {
    display: flex;
    align-items: center;
    margin-left: 103px;
    span {
      margin-right: 12px;
    }
  }
}
.form {
  height: 370px;
  overflow-y: auto;
}
</style>
