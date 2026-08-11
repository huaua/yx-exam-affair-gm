<template>
  <el-dialog
    v-model="dialogVisible"
    :title="dialogProps.title"
    :close-on-click-modal="false"
    destroy-on-close
    align-center
    width="900px"
  >
    <el-alert
      v-if="!isAdminRole && dialogProps.type === 'edit' && dialogProps.row.auditStatus === '4' && dialogProps.row.auditRemark"
      type="error"
      :closable="false"
      class="audit-alert"
      :title="`上次审核不通过原因：${dialogProps.row.auditRemark}`"
    />
    <el-form ref="formRef" :model="dialogProps.row" :rules="rules" label-width="130px" label-suffix="：" inline>
      <el-form-item label="姓名" prop="name">
        <el-input v-model="dialogProps.row.name" maxlength="30" placeholder="请输入姓名" clearable />
      </el-form-item>
      <el-form-item label="性别" prop="sex">
        <el-select v-model="dialogProps.row.sex" placeholder="请选择性别" clearable>
          <el-option label="男" value="男" />
          <el-option label="女" value="女" />
        </el-select>
      </el-form-item>
      <el-form-item label="证件号" prop="idCard">
        <el-input
          v-model="dialogProps.row.idCard"
          maxlength="18"
          placeholder="请输入证件号"
          :disabled="dialogProps.type === 'edit'"
          clearable
        />
      </el-form-item>
      <el-form-item label="手机号" prop="phoneNo">
        <el-input v-model="dialogProps.row.phoneNo" maxlength="11" placeholder="请输入手机号" clearable />
      </el-form-item>
      <el-form-item v-if="isAdminRole" label="专家类别" prop="expertCategory">
        <el-select
          v-model="dialogProps.row.expertCategory"
          :disabled="!isAdminRole"
          placeholder="请选择专家类别"
          @change="onCategoryChange"
        >
          <el-option label="校内" value="1" />
          <el-option v-if="isAdminRole" label="校外" value="2" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="dialogProps.row.expertCategory === '1' && isAdminRole" label="所属学院" prop="deptId">
        <el-select v-model="dialogProps.row.deptId" filterable placeholder="请选择所属学院">
          <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="dialogProps.row.expertCategory === '2'" label="所在单位" prop="unitName">
        <el-input v-model="dialogProps.row.unitName" maxlength="30" placeholder="请输入所在单位" clearable />
      </el-form-item>
      <el-form-item label="专家类型" prop="_typeValues">
        <el-select v-model="dialogProps.row._typeValues" multiple placeholder="请选择专家类型">
          <el-option v-for="item in expertTypeEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="职称" prop="title">
        <el-select v-model="dialogProps.row.title" placeholder="请选择职称" clearable>
          <el-option v-for="item in expertTitleEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="初始学历" prop="initEdu">
        <el-select v-model="dialogProps.row.initEdu" placeholder="请选择初始学历" clearable>
          <el-option v-for="item in educationType" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="dialogProps.row.initEdu" label="初始学历所学专业" prop="initMajor" class="nowrap-label">
        <el-input v-model="dialogProps.row.initMajor" maxlength="20" placeholder="请输入所学专业" clearable />
      </el-form-item>
      <el-form-item label="最终学历" prop="finalEdu">
        <el-select
          v-model="dialogProps.row.finalEdu"
          placeholder="请选择最终学历（非必填）"
          clearable
          @clear="dialogProps.row.finalMajor = ''"
        >
          <el-option v-for="item in educationType" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item v-if="dialogProps.row.finalEdu" label="最终学历所学专业" prop="finalMajor" class="nowrap-label">
        <el-input v-model="dialogProps.row.finalMajor" maxlength="20" placeholder="请输入所学专业" clearable />
      </el-form-item>
      <el-form-item label="擅长科目" prop="_goodSubjectLabels">
        <el-select v-model="dialogProps.row._goodSubjectLabels" multiple filterable placeholder="请选择擅长科目">
          <el-option v-for="item in goodSubjectEnum" :key="item.value" :label="item.label" :value="item.label" />
        </el-select>
      </el-form-item>
      <el-form-item label="工商银行卡号" prop="bankCard">
        <el-input v-model="dialogProps.row.bankCard" maxlength="32" placeholder="请输入银行卡号（非必填）" clearable />
      </el-form-item>
      <el-form-item label="专家介绍" prop="intro" class="full-row">
        <el-input
          v-model="dialogProps.row.intro"
          type="textarea"
          maxlength="1000"
          show-word-limit
          :rows="4"
          placeholder="请输入专家介绍"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <template v-if="!isAdminRole">
        <el-button @click="handleSubmit('draft')">暂存</el-button>
        <el-button type="primary" @click="handleSubmit('submit')">提交审核</el-button>
      </template>
      <el-button v-else type="primary" @click="handleSubmit()">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="ExpertAddOrEditDialog">
import { reactive, ref } from "vue";
import { ElMessage, FormInstance } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { useRole } from "@/hooks/useRole";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { useExpertTitleEnum, useExpertTypeEnum, useGoodSubjectEnum } from "@/hooks/useEnum";
import { educationType } from "@/utils/dict";

type Row = Partial<JudgeExpert.ResExpertDatabaseList> & {
  _typeValues?: string[];
  _goodSubjectLabels?: string[];
  submitAction?: "draft" | "submit";
};

interface DialogProps {
  type: "add" | "edit";
  title: string;
  row: Row;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

const { isAdminRole } = useRole();
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const { goodSubjectEnum } = useGoodSubjectEnum();
const { expertTypeEnum } = useExpertTypeEnum();
const { expertTitleEnum } = useExpertTitleEnum();
const dialogVisible = ref(false);
const formRef = ref<FormInstance>();
const dialogProps = ref<DialogProps>({ type: "add", title: "", row: {} });

const requiredSelect = (message: string) => ({ required: true, message, trigger: "change" });
const rules = reactive({
  name: [
    { required: true, message: "请输入姓名", trigger: "blur" },
    { min: 2, max: 30, message: "姓名限制2-30字符", trigger: "blur" }
  ],
  sex: [requiredSelect("请选择性别")],
  idCard: [
    { required: true, message: "请输入证件号", trigger: "blur" },
    { pattern: /^[A-Za-z0-9()（）]{3,18}$/, message: "证件号须为3-18位数字、大小写字母或括号", trigger: "blur" }
  ],
  phoneNo: [
    { required: true, message: "请输入手机号", trigger: "blur" },
    { pattern: /^1\d{10}$/, message: "手机号须为1开头的11位数字", trigger: "blur" }
  ],
  expertCategory: [requiredSelect("请选择专家类别")],
  deptId: [requiredSelect("请选择所属学院")],
  unitName: [
    { required: true, message: "请输入所在单位", trigger: "blur" },
    { min: 1, max: 30, message: "所在单位限制1-30字符", trigger: "blur" }
  ],
  _typeValues: [requiredSelect("请选择专家类型")],
  initEdu: [requiredSelect("请选择初始学历")],
  initMajor: [
    { required: true, message: "请输入初始学历所学专业", trigger: "blur" },
    { min: 1, max: 20, message: "所学专业限制1-20字符", trigger: "blur" }
  ],
  finalMajor: [
    { required: true, message: "请输入最终学历所学专业", trigger: "blur" },
    { min: 1, max: 20, message: "所学专业限制1-20字符", trigger: "blur" }
  ],
  _goodSubjectLabels: [requiredSelect("请选择擅长科目")]
});

const onCategoryChange = () => {
  dialogProps.value.row.deptId = undefined;
  dialogProps.value.row.unitName = "";
};

const handleSubmit = (submitAction?: "draft" | "submit") => {
  formRef.value?.validate(async valid => {
    if (!valid) return;
    const row = dialogProps.value.row;
    const params = {
      expertId: row.expertId,
      name: row.name,
      sex: row.sex,
      idCard: row.idCard,
      phoneNo: row.phoneNo,
      expertCategory: row.expertCategory,
      deptId: row.deptId,
      unitName: row.unitName,
      type: (row._typeValues || []).join(","),
      title: row.title,
      initEdu: row.initEdu,
      initMajor: row.initMajor,
      finalEdu: row.finalEdu || "0",
      finalMajor: row.finalMajor || "",
      goodSubjects: (row._goodSubjectLabels || []).join(","),
      bankCard: row.bankCard,
      intro: row.intro,
      submitAction
    };
    await dialogProps.value.api?.(params);
    ElMessage.success(`${dialogProps.value.title}成功`);
    dialogProps.value.getTableList?.();
    dialogVisible.value = false;
  });
};

const acceptParams = (params: DialogProps) => {
  const row: Row = { ...params.row };
  row.expertCategory = row.expertCategory || "1";
  row._typeValues = row.type ? String(row.type).split(",").filter(Boolean) : [];
  row._goodSubjectLabels = row.goodSubjects ? row.goodSubjects.split(",").filter(Boolean) : [];
  dialogProps.value = { ...params, row };
  dialogVisible.value = true;
};

defineExpose({ acceptParams });
</script>

<style scoped lang="scss">
.audit-alert {
  margin-bottom: 18px;
}
:deep(.el-form-item__content) {
  width: 260px;
}
.full-row {
  display: flex;
  width: calc(100% - 20px);
}
.full-row :deep(.el-form-item__content) {
  width: 690px;
}
.nowrap-label :deep(.el-form-item__label) {
  white-space: nowrap;
}
</style>
