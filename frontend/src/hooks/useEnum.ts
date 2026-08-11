import { Ref } from "vue";
import { ref } from "vue";
import { getDictList } from "@/api/modules/common";

const _setFormatedEnum = (refObj: Ref, dictType: string) => {
  getDictList({
    curPage: 1, //页数，默认1
    pageSize: 1100, //每页数据量，默认10
    dictType
  })
    .then((res: any) => {
      const list = res?.list || [];
      refObj.value = list.map((item: any) => ({
        label: item.dictLabel,
        value: item.dictValue
      }));
    })
    .catch(error => {
      console.error("error :>>> ", error);
    });
};

/**
 * @description 所属层次
 */
export const useLevelEnum = () => {
  const levelEnum = ref<any>([]);
  _setFormatedEnum(levelEnum, "biz_level");
  return {
    levelEnum
  };
};

/**
 * @description 擅长专业
 */
export const useGoodSubjectEnum = () => {
  const goodSubjectEnum = ref<any>([]);
  _setFormatedEnum(goodSubjectEnum, "biz_dict_subjects");
  return {
    goodSubjectEnum
  };
};

export const useExpertTypeEnum = () => {
  const expertTypeEnum = ref<any>([]);
  _setFormatedEnum(expertTypeEnum, "biz_expert_type");
  return { expertTypeEnum };
};

export const useExpertTitleEnum = () => {
  const expertTitleEnum = ref<any>([]);
  _setFormatedEnum(expertTitleEnum, "biz_expert_title");
  return { expertTitleEnum };
};
