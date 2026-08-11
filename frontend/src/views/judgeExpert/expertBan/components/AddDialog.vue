<template>
  <el-dialog v-model="visible" title="禁止抽取专家新增" width="760px" :close-on-click-modal="false">
    <el-form label-width="110px">
      <el-form-item label="证件号" required>
        <el-select
          v-model="expertId"
          filterable
          remote
          clearable
          placeholder="请选择或输入证件号查询"
          :remote-method="searchExperts"
          :loading="loading"
          style="width: 100%"
          @change="onExpertChange"
        >
          <el-option
            v-for="item in options"
            :key="item.expertId"
            :label="`${item.idCard}（${item.name}）`"
            :value="item.expertId"
          />
        </el-select>
      </el-form-item>

      <template v-if="detail">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名">
              <el-input :model-value="detail.name" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别">
              <el-input :model-value="detail.sex" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="证件号">
              <el-input :model-value="detail.idCard" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input :model-value="detail.phoneNo" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="专家类别">
              <el-input :model-value="categoryLabel(detail.expertCategory)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="detail.expertCategory === '2' ? '所在单位' : '所属学院'">
              <el-input
                :model-value="detail.expertCategory === '2' ? detail.unitName || '--' : detail.deptName || '--'"
                disabled
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="专家类型">
              <el-input :model-value="multiEnumLabel(expertTypeEnum, detail.type)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="职称">
              <el-input :model-value="enumLabel(expertTitleEnum, detail.title)" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="初始学历">
              <el-input
                :model-value="`${enumLabel(educationType, detail.initEdu)}${detail.initMajor ? ' - ' + detail.initMajor : ''}`"
                disabled
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所学专业">
              <el-input :model-value="detail.initMajor || '--'" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最终学历">
              <el-input
                :model-value="
                  detail.finalEdu && detail.finalEdu !== '0'
                    ? `${enumLabel(educationType, detail.finalEdu)}${detail.finalMajor ? ' - ' + detail.finalMajor : ''}`
                    : '--'
                "
                disabled
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所学专业">
              <el-input :model-value="detail.finalMajor || '--'" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="擅长科目">
              <el-input :model-value="detail.goodSubjects || '--'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="工商银行卡号">
              <el-input :model-value="detail.bankCard || '--'" disabled />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="专业背景">
              <el-input :model-value="detail.intro || '--'" type="textarea" :rows="2" disabled />
            </el-form-item>
          </el-col>
        </el-row>
      </template>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="submit">确认禁止</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { JudgeExpert } from "@/api/interface";
import { addBannedExpert, getBanExpertCandidates } from "@/api/modules/judgeExpert";
import { useExpertTitleEnum, useExpertTypeEnum } from "@/hooks/useEnum";
import { educationType } from "@/utils/dict";

const visible = ref(false);
const loading = ref(false);
const saving = ref(false);
const expertId = ref("");
const options = ref<JudgeExpert.ResExpertDatabaseList[]>([]);
const detail = ref<JudgeExpert.ResExpertDatabaseList | null>(null);
let refresh: (() => void) | undefined;

const { expertTypeEnum } = useExpertTypeEnum();
const { expertTitleEnum } = useExpertTitleEnum();

const categoryEnum = [
  { label: "校内", value: "1" },
  { label: "校外", value: "2" }
];

const enumLabel = (items: any[], value?: string) => {
  if (value === undefined || value === null || value === "") return "--";
  return items.find(item => String(item.value) === String(value))?.label || "--";
};

const multiEnumLabel = (items: any[], value?: string) => {
  if (!value) return "--";
  return (
    value
      .split(",")
      .map(v => enumLabel(items, v))
      .filter(v => v !== "--")
      .join("、") || "--"
  );
};

const categoryLabel = (value?: string) => enumLabel(categoryEnum, value);

const searchExperts = async (keyword = "") => {
  loading.value = true;
  try {
    const all: any[] = [];
    let curPage = 1;
    const PAGE_SIZE = 200;
    // 循环拉取全部候选（避免固定 pageSize 截断），直到某次返回不足一页
    while (true) {
      const params: any = { curPage, pageSize: PAGE_SIZE };
      if (keyword.trim()) params.idCard = keyword.trim();
      const data = await getBanExpertCandidates(params);
      const list: any[] = data.list || [];
      all.push(...list);
      if (list.length < PAGE_SIZE) break;
      curPage++;
    }
    options.value = all;
  } finally {
    loading.value = false;
  }
};

const onExpertChange = (val: string) => {
  detail.value = options.value.find(item => item.expertId === val) || null;
};

const submit = async () => {
  if (!expertId.value) {
    ElMessage.warning("请选择证件号");
    return;
  }
  saving.value = true;
  try {
    await addBannedExpert({ expertId: expertId.value });
    ElMessage.success("禁止抽取成功");
    visible.value = false;
    refresh?.();
  } finally {
    saving.value = false;
  }
};

const acceptParams = (getTableList?: () => void) => {
  expertId.value = "";
  options.value = [];
  detail.value = null;
  refresh = getTableList;
  visible.value = true;
  searchExperts();
};

defineExpose({ acceptParams });
</script>
