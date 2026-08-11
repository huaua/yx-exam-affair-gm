<template>
  <el-dialog v-model="visible" title="专家审核" width="980px" align-center :close-on-click-modal="false">
    <div class="compare-grid">
      <section class="info-card">
        <h3>{{ snapshot ? "本次提交信息" : "专家信息" }}</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item v-for="item in fields" :key="item.key" :label="item.label">
            {{ displayValue(row, item.key) }}
          </el-descriptions-item>
        </el-descriptions>
      </section>
      <section v-if="snapshot" class="info-card">
        <h3>上次审核通过信息</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item v-for="item in fields" :key="item.key" :label="item.label">
            {{ displayValue(snapshot, item.key) }}
          </el-descriptions-item>
        </el-descriptions>
      </section>
    </div>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="danger" :disabled="auditing" @click="rejectVisible = true">不通过</el-button>
      <el-button type="primary" :loading="auditing" @click="doAudit('pass')">通过</el-button>
    </template>
  </el-dialog>
  <el-dialog v-model="rejectVisible" title="审核原因" width="460px" append-to-body align-center :close-on-click-modal="false">
    <el-input v-model="reason" type="textarea" maxlength="500" show-word-limit :rows="5" placeholder="请输入审核不通过原因" />
    <template #footer>
      <el-button @click="rejectVisible = false">取消</el-button>
      <el-button type="primary" :loading="auditing" @click="confirmReject">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { auditExpert, getLatestExpertAuditHistory, getNextPendingExpert } from "@/api/modules/judgeExpert";
import { useDepartmentEnum } from "@/hooks/useCustomEnum";
import { useExpertTitleEnum, useExpertTypeEnum } from "@/hooks/useEnum";
import { educationType } from "@/utils/dict";

const fields = [
  { key: "name", label: "姓名" },
  { key: "sex", label: "性别" },
  { key: "idCard", label: "证件号" },
  { key: "phoneNo", label: "手机号" },
  { key: "expertCategory", label: "专家类别" },
  { key: "deptOrUnit", label: "所属学院/所在单位" },
  { key: "type", label: "专家类型" },
  { key: "title", label: "职称" },
  { key: "initEdu", label: "初始学历" },
  { key: "initMajor", label: "初始学历所学专业" },
  { key: "finalEdu", label: "最终学历" },
  { key: "finalMajor", label: "最终学历所学专业" },
  { key: "goodSubjects", label: "擅长科目" },
  { key: "intro", label: "专家介绍" }
];
const visible = ref(false);
const rejectVisible = ref(false);
const reason = ref("");
const auditing = ref(false);
const row = ref<Partial<JudgeExpert.ResExpertDatabaseList>>({});
const snapshot = ref<Record<string, any> | null>(null);
let refresh: (() => void) | undefined;
// 列表页当前查询条件，用于连续审核时按同样的筛选范围取下一条
let searchParam: Record<string, any> = {};
const { departmentEnum } = useDepartmentEnum("onlyDepartment");
const { expertTypeEnum } = useExpertTypeEnum();
const { expertTitleEnum } = useExpertTitleEnum();
const enumLabel = (items: any[], value: any) => items.find(item => String(item.value) === String(value))?.label || "--";

const displayValue = (data: Record<string, any> | null, key: string) => {
  if (!data) return "--";
  if (key === "expertCategory") return String(data.expertCategory || data.ExpertCategory) === "2" ? "校外" : "校内";
  if (key === "deptOrUnit") {
    const deptId = data.deptId || data.DeptId;
    return data.unitName || data.UnitName || data.deptName || data.DeptName || enumLabel(departmentEnum.value, deptId);
  }
  if (key === "type") {
    const value = data.type || data.Type || "";
    return (
      String(value)
        .split(",")
        .map(item => enumLabel(expertTypeEnum.value, item))
        .join("、") || "--"
    );
  }
  if (key === "title") return enumLabel(expertTitleEnum.value, data.title || data.Title);
  if (key === "initEdu" || key === "finalEdu")
    return enumLabel(educationType, data[key] || data[key.charAt(0).toUpperCase() + key.slice(1)]);
  return data[key] ?? data[key.charAt(0).toUpperCase() + key.slice(1)] ?? "--";
};

const loadExpert = async (expert: Partial<JudgeExpert.ResExpertDatabaseList>) => {
  row.value = { ...expert };
  reason.value = "";
  snapshot.value = null;
  const { obj } = await getLatestExpertAuditHistory({ expertId: row.value.expertId! });
  if (obj?.snapshotData) snapshot.value = JSON.parse(obj.snapshotData);
};

const doAudit = async (action: "pass" | "reject", auditReason = "") => {
  if (auditing.value) return;
  auditing.value = true;
  try {
    await auditExpert({ expertId: row.value.expertId!, action, reason: auditReason });
    rejectVisible.value = false;
    refresh?.();
    // 带上列表查询条件 + 当前记录 id，保证只在筛选结果内按顺序流转
    const { obj } = await getNextPendingExpert({ ...searchParam, expertId: row.value.expertId });
    if (obj) {
      await loadExpert(obj);
      ElMessage.success(action === "pass" ? "审核通过成功，已进入下一位" : "审核不通过成功，已进入下一位");
      return;
    }
    visible.value = false;
    ElMessage.success(action === "pass" ? "审核通过成功，当前已无待审核专家" : "审核不通过成功，当前已无待审核专家");
  } finally {
    auditing.value = false;
  }
};
const confirmReject = () => {
  if (!reason.value.trim()) {
    ElMessage.warning("请输入审核不通过原因");
    return;
  }
  doAudit("reject", reason.value.trim());
};
const acceptParams = async (params: {
  row: Partial<JudgeExpert.ResExpertDatabaseList>;
  getTableList?: () => void;
  searchParam?: Record<string, any>;
}) => {
  refresh = params.getTableList;
  searchParam = { ...(params.searchParam || {}) };
  visible.value = true;
  await loadExpert(params.row);
};
defineExpose({ acceptParams });
</script>

<style scoped lang="scss">
.compare-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  max-height: none;
  overflow: visible;
}
.info-card {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  padding: 8px 10px;
}
.info-card h3 {
  margin: 0 0 6px;
  font-size: 14px;
}
:deep(.el-descriptions) {
  --el-descriptions-item-bordered-label-background: #fafafa;
  :deep(.el-descriptions__body) {
    .el-descriptions__table {
      .el-descriptions__cell {
        &.is-bordered-label,
        &.is-bordered-content {
          padding: 4px 10px;
          font-size: 13px;
        }
      }
      tr:not(:last-child) td {
        border-bottom: none;
      }
    }
  }
}
</style>
