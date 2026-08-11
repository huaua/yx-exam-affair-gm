<!-- 用户管理-用户管理-详情Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="500"
  >
    <el-form ref="formRef" label-width="100px" label-suffix=" :" :rules="rules" :model="dialogProps.row">
      <el-form-item v-if="authMode(['add', 'edit', 'password'])" label="账号" prop="userName">
        <el-input
          v-model="dialogProps.row!.userName"
          placeholder="请填写账号"
          :disabled="isDisabledUserName"
          maxlength="18"
          clearable
        ></el-input>
      </el-form-item>
      <el-form-item v-if="authMode(['add', 'password'])" :label="passwordLabel" prop="password">
        <el-input
          v-model="dialogProps.row!.password"
          :placeholder="`请填写${passwordLabel}`"
          autocomplete="new-password"
          show-password
        ></el-input>
      </el-form-item>
      <el-form-item v-if="authMode(['add', 'password'])" label="确认密码" prop="confirmPassword">
        <el-input
          v-model="dialogProps.row!.confirmPassword"
          placeholder="请再次输入密码"
          autocomplete="new-password"
          show-password
        ></el-input>
      </el-form-item>
      <el-form-item v-if="authMode(['add', 'edit'])" label="姓名" prop="nickName">
        <el-input v-model="dialogProps.row!.nickName" placeholder="请填写姓名" maxlength="15" clearable></el-input>
      </el-form-item>
      <el-form-item v-if="authMode(['add'])" label="账号权限" prop="roleKey">
        <el-radio-group v-model="dialogProps.row!.roleKey">
          <el-radio value="department">学院</el-radio>
          <el-radio value="admissions">招办</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="isShowDepartment" label="所属学院" prop="deptId">
        <el-select v-model="dialogProps.row!.deptId" placeholder="请选择所属学院" clearable>
          <el-option v-for="item in departmentEnum" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="UserDialog">
import { User } from "@/api/interface";
import { RoleKeyEnum } from "@/enums";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive, computed } from "vue";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { isValidLength, isPassword } from "@/utils/eleValidate";

const SPLIT_DEPARTMENT_CODE = "0"; // 院系的deptId > 0

interface AddUserParams extends User.ResUserList {
  confirmPassword: string;
}

interface DialogProps {
  type: string;
  title: string;
  row: Partial<AddUserParams>;
  api?: (params: any) => Promise<any>;
  getTableList?: () => void;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  type: "add",
  title: "",
  row: {}
});
const { departmentEnum } = useDepartmentEnum("onlyDepartment");

// 自定义确认密码校验规则
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== dialogProps.value?.row?.password) {
    callback(new Error("两次输入的密码不一致"));
  } else {
    callback();
  }
};

// 表单校验规则
const rules = reactive({
  userName: [
    { required: true, message: "请填写账号" },
    {
      validator: (...args: [any, any]) => isValidLength(3, 18, ...args),
      message: "账号必须为3到18位",
      trigger: "blur"
    }
  ],
  password: [
    { required: true, message: "请填写密码" },
    { validator: isPassword, message: "密码长度为6到18位必须包含数字、小写和大写字母" }
  ],
  confirmPassword: [
    { required: true, message: "请再次输入密码", trigger: "blur" },
    { validator: validateConfirmPassword, trigger: "blur" }
  ],
  nickName: [
    { required: true, message: "请填写姓名" },
    {
      validator: (...args: [any, any]) => isValidLength(2, 15, ...args),
      message: "姓名必须为2到15位",
      trigger: "blur"
    }
  ],
  roleKey: [{ required: true, message: "请选择账号权限" }],
  deptId: [{ required: true, message: "请选择所属学院" }]
});

const isDisabledUserName = computed(() => dialogProps.value.type !== "add");
const isShowDepartment = computed(
  () =>
    (dialogProps.value.type === "add" && dialogProps.value?.row?.roleKey === RoleKeyEnum.DEPARTMENT) ||
    (dialogProps.value.type === "edit" && dialogProps.value.row.deptId! > SPLIT_DEPARTMENT_CODE)
);
const passwordLabel = computed(() => (dialogProps.value.type === "password" ? "新密码" : "密码"));

const authMode = (modeList: string[] = []) => {
  return modeList.includes(dialogProps.value.type);
};

// 提交数据（新增/编辑/重置密码）
const formRef = ref<FormInstance>();
const handleSubmit = () => {
  formRef.value!.validate(async valid => {
    if (!valid) return;
    try {
      const params: Partial<User.ResUserList> = {};
      if (dialogProps.value.type === "password") {
        params.userId = dialogProps.value.row?.userId;
        params.password = dialogProps.value.row?.password;
      } else if (dialogProps.value.type === "edit") {
        params.userId = dialogProps.value.row?.userId;
        params.nickName = dialogProps.value.row?.nickName;
        isShowDepartment.value && (params.deptId = dialogProps.value.row?.deptId);
      } else if (dialogProps.value.type === "add") {
        params.userName = dialogProps.value.row?.userName; // 用户名
        params.nickName = dialogProps.value.row?.nickName; // 用户昵称
        params.password = dialogProps.value.row?.password; // 密码
        params.roleKey = dialogProps.value.row?.roleKey; // 角色代码： department-学院 admissions-招办
        isShowDepartment.value && (params.deptId = dialogProps.value.row?.deptId); // 所属学院id （可选，roleKey=department时必填）
      }
      await dialogProps.value.api!(params);
      ElMessage.success({ message: `${dialogProps.value.title}成功！` });
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
  if (dialogProps.value.type === "password") {
    dialogProps.value.row.password = "";
  }
  dialogVisible.value = true;
};

defineExpose({
  acceptParams
});
</script>
