<!-- 监考老师管理-教职工管理-新增编辑Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="786"
  >
    <el-form ref="formRef" label-width="110px" label-suffix=" :" :rules="rules" :model="dialogProps.row" inline>
      <el-form-item label="姓名" prop="tchName">
        <el-input v-model="dialogProps.row!.tchName" placeholder="请填写姓名" maxlength="18" clearable></el-input>
      </el-form-item>
      <el-form-item label="性别" prop="sex">
        <el-select v-model="dialogProps.row!.sex" placeholder="请选择性别" clearable>
          <el-option v-for="item in genderEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="工号" prop="jobNo">
        <el-input v-model="dialogProps.row!.jobNo" placeholder="请填写工号" maxlength="20" clearable />
      </el-form-item>
      <el-form-item label="证件号" prop="idCard">
        <el-input
          v-model="dialogProps.row!.idCard"
          placeholder="请填写证件号"
          :disabled="dialogProps.type === 'edit'"
          :clearable="dialogProps.type !== 'edit'"
        />
      </el-form-item>
      <el-form-item v-if="isAdminRole" label="所属学院" prop="deptId">
        <el-select v-model="dialogProps.row!.deptId" placeholder="请选择所属学院" clearable>
          <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="职务" prop="duties">
        <el-input v-model="dialogProps.row!.duties" placeholder="请填写职务" maxlength="15" clearable />
      </el-form-item>
      <el-form-item label="工商银行卡号" prop="bankcard">
        <el-input v-model="dialogProps.row!.bankcard" placeholder="请填写工商银行卡号" maxlength="19" clearable />
      </el-form-item>
      <el-form-item label="专业或行政" prop="profOffice">
        <el-input v-model="dialogProps.row!.profOffice" placeholder="请填写专业或行政" maxlength="15" clearable />
      </el-form-item>
      <el-form-item label="专业背景" prop="education" class="w-full">
        <el-input
          v-model="dialogProps.row!.education"
          type="textarea"
          placeholder="请填写专业背景"
          maxlength="1000"
          :autosize="{ minRows: 2, maxRows: 8 }"
          show-word-limit
          clearable
        />
      </el-form-item>
      <el-form-item
        v-if="!isAdminRole && dialogProps.type === 'edit' && dialogProps.row?.auditStatus === '4' && dialogProps.row?.auditRemark"
        label="不通过原因"
        class="w-full"
      >
        <div class="audit-reason-text">{{ dialogProps.row.auditRemark }}</div>
      </el-form-item>
      <el-form-item
        v-if="isAdminRole && dialogProps.type === 'edit' && snapshotFields.tchName && snapshotVisible"
        label="上次审核通过数据"
        class="w-full snapshot-item"
      >
        <div class="snapshot-block">
          <div class="snapshot-title">（参考）以下为上次审核通过时保存的数据</div>
          <div class="row">
            <span class="label">姓名：</span><span>{{ snapshotFields.tchName }}</span>
          </div>
          <div class="row">
            <span class="label">性别：</span><span>{{ snapshotFields.sex }}</span>
          </div>
          <div class="row">
            <span class="label">工号：</span><span>{{ snapshotFields.jobNo }}</span>
          </div>
          <div class="row">
            <span class="label">所属学院：</span><span>{{ snapshotFields.deptName || snapshotFields.deptId || "--" }}</span>
          </div>
          <div class="row">
            <span class="label">职务：</span><span>{{ snapshotFields.duties || "--" }}</span>
          </div>
          <div class="row">
            <span class="label">专业或行政：</span><span>{{ snapshotFields.profOffice || "--" }}</span>
          </div>
          <div class="row">
            <span class="label">专业背景：</span><span class="multiline">{{ snapshotFields.education || "--" }}</span>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button v-if="isAdminRole" type="primary" @click="handleSubmit()">确定</el-button>
      <template v-else>
        <el-button type="primary" plain @click="handleSubmit('draft')">暂存</el-button>
        <el-button type="primary" @click="handleSubmit('submit')">提交审核</el-button>
      </template>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AddOrEditDialog">
import { Teacher } from "@/api/interface";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { useRole } from "@/hooks/useRole";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { isValidLength, isAlphanumeric } from "@/utils/eleValidate";
import { getLatestAuditHistory } from "@/api/modules/teacher";

interface DialogProps {
  type: string;
  title: string;
  row: Partial<Teacher.ResTeacherList>;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}
const genderEnum = [
  {
    label: "男",
    value: "男"
  },
  {
    label: "女",
    value: "女"
  }
];
const { isAdminRole } = useRole();
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  title: "",
  row: {}
});

// 招办编辑时，用于对比的“上次审核通过”快照
const snapshotVisible = ref(false);
const snapshotFields = reactive<Record<string, any>>({});

const deptNameOf = (deptId: any): string => {
  if (deptId === undefined || deptId === null || deptId === "") return "";
  const matched = (departmentEnum.value as any[])?.find(item => String(item.value) === String(deptId));
  return matched?.label ?? "";
};

const loadSnapshot = async (tchId: string) => {
  snapshotVisible.value = false;
  Object.keys(snapshotFields).forEach(k => delete snapshotFields[k]);
  if (!tchId) return;
  try {
    const { obj } = await getLatestAuditHistory({ tchId });
    if (obj && obj.snapshotData) {
      const data = JSON.parse(obj.snapshotData);
      const deptId = data.deptId ?? data.DeptId;
      Object.assign(snapshotFields, {
        tchName: data.tchName || data.TchName || "--",
        sex: data.sex || data.Sex || "--",
        jobNo: data.jobNo || data.JobNo || "--",
        deptId,
        deptName: data.deptName || data.DeptName || deptNameOf(deptId),
        duties: data.duties || data.Duties || "",
        profOffice: data.profOffice || data.ProfOffice || "",
        education: data.education || data.Education || ""
      });
      snapshotVisible.value = true;
    }
  } catch (e) {
    snapshotVisible.value = false;
  }
};

const validateBankCard = (rule: any, value: any) => {
  if (!value) {
    return true;
  }
  return isValidLength(16, 19, rule, value);
};

// 表单校验规则
const rules = reactive({
  tchName: [
    { required: true, message: "请填写姓名" },
    {
      validator: (...args: [any, any]) => isValidLength(2, 18, ...args),
      message: "姓名必须为2到18位",
      trigger: "blur"
    }
  ],
  sex: [{ required: true, message: "请选择性别" }],
  jobNo: [
    { required: true, message: "请填写工号" },
    {
      validator: (...args: [any, any]) => isValidLength(2, 20, ...args),
      message: "工号必须为2到20位",
      trigger: "blur"
    },
    {
      validator: isAlphanumeric,
      message: "工号仅允许输入数字或英文字母",
      trigger: "blur"
    }
  ],
  idCard: [
    { required: true, message: "请填写证件号" },
    {
      validator: isAlphanumeric,
      message: "证件号仅允许输入数字或英文字母",
      trigger: "blur"
    }
  ],
  deptId: [{ required: true, message: "请选择所属学院" }],
  bankcard: [
    {
      validator: validateBankCard,
      message: "卡号必须为16到19位",
      trigger: "blur"
    },
    {
      validator: isAlphanumeric,
      message: "卡号仅允许输入数字或英文字母",
      trigger: "blur"
    }
  ]
});

// 提交数据（新增/编辑/重置密码）
const formRef = ref<FormInstance>();
const handleSubmit = (submitAction?: "draft" | "submit") => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params: Partial<Teacher.ResTeacherList> & { submitAction?: "draft" | "submit" } = {};
      params.tchId = dialogProps.value.row?.tchId || undefined;
      params.tchName = dialogProps.value.row?.tchName;
      params.sex = dialogProps.value.row?.sex;
      params.jobNo = dialogProps.value.row?.jobNo;
      params.idCard = dialogProps.value.row?.idCard;
      isAdminRole.value && (params.deptId = dialogProps.value.row?.deptId); // 所属学院id （可选，管理员角色时必填）
      params.duties = dialogProps.value.row?.duties;
      params.bankcard = dialogProps.value.row?.bankcard;
      params.profOffice = dialogProps.value.row?.profOffice;
      params.education = dialogProps.value.row?.education;
      !isAdminRole.value && (params.submitAction = submitAction);
      await dialogProps.value.api!(params);
      const successMessage = isAdminRole.value
        ? `${dialogProps.value.title}成功！`
        : submitAction === "submit"
          ? "提交成功！"
          : "暂存成功！";
      ElMessage.success({ message: successMessage });
      dialogProps.value.getTableList!();
      dialogVisible.value = false;
    } catch (error) {
      console.log(error);
    }
  });
};

// 接收父组件传过来的参数
const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  // 招办编辑时加载上次审核通过快照用于对比
  if (isAdminRole.value && params.type === "edit" && params.row?.tchId) {
    loadSnapshot(params.row.tchId);
  } else {
    snapshotVisible.value = false;
  }
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
:deep(.el-form-item__content) {
  min-width: 230px;
}
.w-full {
  display: flex;
}

.audit-reason-text {
  width: 100%;
  color: var(--el-color-danger);
  line-height: 1.6;
  word-break: break-all;
}

.snapshot-item :deep(.el-form-item__content) {
  display: block;
}

.snapshot-block {
  width: 100%;
  padding: 10px 12px;
  background: #fff7e6;
  border: 1px dashed #f5b041;
  border-radius: 4px;
  .snapshot-title {
    font-weight: 600;
    color: #d35400;
    margin-bottom: 4px;
  }
  .row {
    display: flex;
    font-size: 13px;
    line-height: 26px;
    .label {
      color: #909399;
      width: 90px;
    }
    .multiline {
      flex: 1;
      white-space: pre-wrap;
      word-break: break-all;
    }
  }
}
</style>
