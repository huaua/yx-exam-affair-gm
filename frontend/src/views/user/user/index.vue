<!-- 用户管理-用户管理 -->
<template>
  <div class="table-box">
    <ProTable ref="proTable" :columns="columns" :request-api="getUserList">
      <!-- 表格 header 按钮 -->
      <template #tableHeader>
        <el-button type="primary" :icon="CirclePlus" @click="openDialog('add')">新增用户</el-button>
      </template>
      <!-- 表格操作 -->
      <template #operation="scope">
        <el-button type="primary" link :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
        <el-button type="primary" link :icon="Refresh" @click="openDialog('password', scope.row)">修改密码</el-button>
      </template>
    </ProTable>
    <UserDialog ref="dialogRef" />
  </div>
</template>

<script setup lang="tsx" name="user">
import { User } from "@/api/interface";
import { ProTableInstance, ColumnProps } from "@/components/ProTable/interface";
import { ref, reactive } from "vue";
import { CirclePlus, Edit, Refresh } from "@element-plus/icons-vue";
import { useHandleData } from "@/hooks/useHandleData";
import ProTable from "@/components/ProTable/index.vue";
import UserDialog from "./components/UserDialog.vue";
import { userStatus } from "@/utils/dict";
import { getUserList, changeUserStatus, addUser, editUser, setUserPassword } from "@/api/modules/user";

type MapString = {
  [key: string]: string;
};

type MapFunction = {
  [key: string]: (params: any) => Promise<any>;
};

// ProTable 实例
const proTable = ref<ProTableInstance>();

// 表格配置项
const columns = reactive<ColumnProps<User.ResUserList>[]>([
  { prop: "nickName", label: "姓名" },
  { prop: "userName", label: "账号" },
  { prop: "deptName", label: "所属学院" },
  {
    prop: "state",
    label: "状态",
    enum: userStatus,
    render: scope => {
      return (
        <el-switch
          model-value={scope.row.state}
          active-text={scope.row.state === "1" ? "启用" : "禁用"}
          active-value={"1"}
          inactive-value={"0"}
          onClick={() => changeStatus(scope.row)}
        />
      );
    }
  },
  { prop: "operation", label: "操作", fixed: "right", width: 200 }
]);

// 切换用户状态
const changeStatus = async (row: User.ResUserList) => {
  await useHandleData(
    changeUserStatus,
    { userId: row.userId, state: row.state === "0" ? "1" : "0" },
    `${row.state === "0" ? "启用" : "禁用"}【${row.userName}】用户`
  );
  proTable.value?.getTableList();
};

// 打开 Dialog(新增、编辑、修改密码)
const dialogRef = ref<InstanceType<typeof UserDialog> | null>(null);
const openDialog = (type: string, row: Partial<User.ResUserList> = {}) => {
  const params = {
    type,
    title: getDialogTitle(type),
    row: { ...row },
    api: getDialogApi(type),
    getTableList: proTable.value?.getTableList
  };
  dialogRef.value?.acceptParams(params);
};

const getDialogTitle = (type: string) => {
  const dialogTitleMap: MapString = {
    add: "新增用户",
    edit: "修改用户",
    password: "修改密码"
  };
  return dialogTitleMap[type];
};

const getDialogApi = (type: string) => {
  const dialogApiMap: MapFunction = {
    add: addUser,
    edit: editUser,
    password: setUserPassword
  };
  return dialogApiMap[type];
};
</script>
