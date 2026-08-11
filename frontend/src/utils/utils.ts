import { Ref } from "vue";

interface DictItem {
  label: string;
  value: string;
}

/**
 * 根据字典值获取字典的label
 * @param dictList 字典列表Ref，例如 ref([{label: "xxx", value: "1"}])
 * @param value    要查找的字典值
 * @returns 对应的label，未找到时返回空字符串
 */
export function getDictLabel(dictList: Ref<DictItem[] | undefined | null>, value: string | undefined | null): string {
  if (!dictList?.value || !Array.isArray(dictList.value) || !value) {
    return "";
  }
  const found = dictList.value.find(item => item.value === value);
  return found?.label ?? "";
}

/**
 * 根据字典 label 获取对应的 value
 *
 * @param {Ref<DictItem[] | undefined | null>} dictList - 字典列表 Ref，例如 ref([{label: "xxx", value: "1"}])
 * @param {string | undefined | null} label - 要查找的字典 label
 * @returns {string} 对应的 value，未找到时返回空字符串
 */
export function getDictValue(dictList: Ref<DictItem[] | undefined | null>, label: string | undefined | null): string {
  if (!dictList?.value || !Array.isArray(dictList.value) || !label) {
    return "";
  }
  const found = dictList.value.find(item => item.label === label);
  return found?.value ?? "";
}
