<!-- 监考老师管理-教职工管理-审核Dialog -->
<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${dialogProps.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :align-center="true"
    width="640"
  >
    <!-- 当前信息（首次审核时直接展示，无标题） / 修改后信息（再次提交审核时带标题） -->
    <div class="audit-block" :class="{ 'has-frame': !isFirstAudit }">
      <div class="block-title" v-if="!isFirstAudit">修改后信息：</div>
      <div class="field-grid">
        <div class="field-row">
          <div class="field">
            <span class="label">姓名：</span>
            <el-input :model-value="dialogProps.row?.tchName || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">性别：</span>
            <el-input :model-value="dialogProps.row?.sex || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">工号：</span>
            <el-input :model-value="dialogProps.row?.jobNo || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">证件号：</span>
            <el-input :model-value="dialogProps.row?.idCard || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">所属学院：</span>
            <el-input :model-value="dialogProps.row?.deptName || dialogProps.row?.deptId || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">职务：</span>
            <el-input :model-value="dialogProps.row?.duties || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">工商银行卡号：</span>
            <el-input :model-value="dialogProps.row?.bankcard || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">专业或行政：</span>
            <el-input :model-value="dialogProps.row?.profOffice || '--'" disabled />
          </div>
        </div>
        <div class="field-row single">
          <div class="field">
            <span class="label">专业背景：</span>
            <el-input
              :model-value="dialogProps.row?.education || '--'"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              disabled
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 修改前信息（仅非首次审核时展示，用于对比上次通过的数据） -->
    <div class="audit-block has-frame" v-if="!isFirstAudit && latestSnapshot">
      <div class="block-title">修改前信息：</div>
      <div class="field-grid">
        <div class="field-row">
          <div class="field">
            <span class="label">姓名：</span>
            <el-input :model-value="snapshotFields.tchName || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">性别：</span>
            <el-input :model-value="snapshotFields.sex || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">工号：</span>
            <el-input :model-value="snapshotFields.jobNo || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">证件号：</span>
            <el-input :model-value="snapshotFields.idCard || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">所属学院：</span>
            <el-input :model-value="snapshotFields.deptName || snapshotFields.deptId || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">职务：</span>
            <el-input :model-value="snapshotFields.duties || '--'" disabled />
          </div>
        </div>
        <div class="field-row">
          <div class="field">
            <span class="label">工商银行卡号：</span>
            <el-input :model-value="snapshotFields.bankcard || '--'" disabled />
          </div>
          <div class="field">
            <span class="label">专业或行政：</span>
            <el-input :model-value="snapshotFields.profOffice || '--'" disabled />
          </div>
        </div>
        <div class="field-row single">
          <div class="field">
            <span class="label">专业背景：</span>
            <el-input
              :model-value="snapshotFields.education || '--'"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              disabled
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="danger" :disabled="auditing" @click="onSubmit('reject')">不通过</el-button>
      <el-button type="primary" :loading="auditing" @click="onSubmit('pass')">通过</el-button>
    </template>
  </el-dialog>

  <!-- 审核原因子弹窗（仅在不通过时展示） -->
  <el-dialog
    v-model="rejectDialogVisible"
    title="审核原因"
    :close-on-click-modal="false"
    :align-center="true"
    width="420"
    append-to-body
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" label-suffix=" :">
      <el-form-item label="不通过原因" prop="reason">
        <el-input
          v-model="form.reason"
          type="textarea"
          :rows="4"
          :maxlength="500"
          show-word-limit
          placeholder="请填写审核不通过原因"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="rejectDialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="auditing" @click="confirmReject">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="AuditDialog">
import { Teacher } from "@/api/interface";
import { ElMessage, FormInstance } from "element-plus";
import { ref, reactive } from "vue";
import { auditTeacher, getLatestAuditHistory, getNextPendingTeacher } from "@/api/modules/teacher";

interface DialogProps {
  title: string;
  row: Partial<Teacher.ResTeacherList> & { tchId: string };
  isFirstAudit?: boolean;
  getTableList?: () => void;
  // 列表页当前查询条件，用于连续审核时按同样的筛选范围取下一条
  searchParam?: Record<string, any>;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  title: "",
  row: {} as any
});

const form = reactive({
  action: "pass" as "pass" | "reject",
  reason: ""
});

const formRef = ref<FormInstance>();
const rules = reactive({
  reason: [{ required: true, message: "请填写审核不通过原因", trigger: "blur" }]
});

const isFirstAudit = ref(true);
const auditing = ref(false);

const latestSnapshot = ref<Teacher.ResAuditHistory | null>(null);
const snapshotFields = reactive<Record<string, any>>({});

const loadSnapshot = async () => {
  const tchId = dialogProps.value.row?.tchId;
  latestSnapshot.value = null;
  Object.keys(snapshotFields).forEach(k => delete snapshotFields[k]);
  if (!tchId) {
    isFirstAudit.value = true;
    return;
  }
  try {
    const { obj } = await getLatestAuditHistory({ tchId });
    if (obj && obj.snapshotData) {
      latestSnapshot.value = obj;
      try {
        const data = JSON.parse(obj.snapshotData);
        Object.assign(snapshotFields, {
          tchName: data.TchName || data.tchName || "--",
          sex: data.Sex || data.sex || "--",
          jobNo: data.JobNo || data.jobNo || "--",
          idCard: data.IdCard || data.idCard || "--",
          deptId: data.DeptId || data.deptId,
          deptName: data.DeptName || data.deptName || "",
          duties: data.Duties || data.duties || "",
          bankcard: data.Bankcard || data.bankcard || "",
          profOffice: data.ProfOffice || data.profOffice || "",
          education: data.Education || data.education || ""
        });
        isFirstAudit.value = false; // 存在上次审核通过的快照 → 非首次审核，展示对比
      } catch (e) {
        isFirstAudit.value = true;
      }
    } else {
      isFirstAudit.value = true;
    }
  } catch (e) {
    isFirstAudit.value = true;
  }
};

// 不通过时弹出子弹窗填写原因
const rejectDialogVisible = ref(false);

const loadTeacher = async (row: Partial<Teacher.ResTeacherList> & { tchId: string }) => {
  dialogProps.value.row = { ...row };
  form.reason = "";
  await loadSnapshot();
};

const doAudit = async (action: "pass" | "reject", reason = "") => {
  if (auditing.value) return;
  auditing.value = true;
  try {
    await auditTeacher({
      tchId: dialogProps.value.row.tchId,
      action,
      reason
    });
    dialogProps.value.getTableList?.();
    // 带上列表查询条件 + 当前记录 id，保证只在筛选结果内按顺序流转
    const { obj } = await getNextPendingTeacher({
      ...(dialogProps.value.searchParam || {}),
      tchId: dialogProps.value.row.tchId
    });
    if (obj) {
      await loadTeacher(obj as Teacher.ResTeacherList & { tchId: string });
      ElMessage.success({ message: action === "pass" ? "审核通过成功，已进入下一位" : "审核不通过成功，已进入下一位" });
      return;
    }
    dialogVisible.value = false;
    ElMessage.success({
      message: action === "pass" ? "审核通过成功，当前已无待审核教职工" : "审核不通过成功，当前已无待审核教职工"
    });
  } catch (e) {
    // 错误信息由 axios 拦截统一处理
  } finally {
    auditing.value = false;
  }
};

const onSubmit = (action: "pass" | "reject") => {
  if (action === "pass") {
    doAudit("pass", "");
  } else {
    form.reason = "";
    rejectDialogVisible.value = true;
  }
};

const confirmReject = () => {
  formRef.value?.validate(async valid => {
    if (!valid) return;
    const reason = form.reason;
    rejectDialogVisible.value = false;
    await doAudit("reject", reason);
  });
};

const acceptParams = (params: DialogProps) => {
  dialogProps.value = params;
  form.reason = "";
  // 显式指定首审判定时优先采用；否则交给 loadSnapshot 依据快照是否存在自动判定
  if (typeof params.isFirstAudit === "boolean") {
    isFirstAudit.value = params.isFirstAudit;
  }
  dialogVisible.value = true;
  loadSnapshot();
};

defineExpose({
  acceptParams
});
</script>

<style lang="scss" scoped>
.audit-block {
  // 首次审核：直接展示，无外框无内边距
  padding: 0;
  margin-bottom: 0;
  border: none;
  background: transparent;

  // 再次提交审核：带外框 + 标题
  &.has-frame {
    background: #fff;
    border: 1px solid #ebeef5;
    border-radius: 4px;
    padding: 14px 16px;
    margin-bottom: 14px;
  }

  .block-title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 12px;
  }

  .field-grid {
    .field-row {
      display: flex;
      gap: 12px;
      margin-bottom: 10px;

      &.single {
        .field {
          flex: 1;
          width: 100%;
        }
      }

      &:last-child {
        margin-bottom: 0;
      }

      .field {
        flex: 1;
        // grid 布局：label 列固定 90px（容纳最长 label「工商银行卡号：」），
        // 右侧 input 列弹性填满，保证所有字段输入框起点和宽度对齐
        display: grid;
        grid-template-columns: 90px 1fr;
        align-items: start;
        column-gap: 8px;
        min-width: 0;

        .label {
          color: #909399;
          font-size: 13px;
          line-height: 32px;
          text-align: right;
          white-space: nowrap;
        }

        // 字段值的 el-input 撑满 grid 第二列
        :deep(.el-input),
        :deep(.el-textarea) {
          width: 100%;
          min-width: 0;
        }

        :deep(.el-input__wrapper) {
          background-color: #f4f4f5;
          box-shadow: 0 0 0 1px #e9e9eb inset;
          padding: 1px 11px;
        }

        :deep(.el-textarea__inner) {
          background-color: #f4f4f5;
          resize: none;
        }
      }
    }
  }
}
</style>
