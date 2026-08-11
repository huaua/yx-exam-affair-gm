import { Common } from "@/api/interface";
import { ref } from "vue";
import { useRole } from "@/hooks/useRole";
import { admissionClassroomType, admissionExpertType } from "@/utils/dict";
import { getDepartmentList } from "@/api/modules/department";
import { getRecruitManageList } from "@/api/modules/classroom";
import { getAllJudgeRecruitTaskList } from "@/api/modules/judgeExpert";

export const useDepartmentEnum = (type: "all" | "onlyDepartment" | "expertDataBase" = "all") => {
  const { isAdminRole } = useRole();
  const departmentEnum = ref<Common.enumDict[]>([]);
  // 获取全部学院列表
  const fetchDepartmentEnum = async () => {
    const { list = [] } = await getDepartmentList({ curPage: 1, pageSize: 1000 });
    const deptList = list.map((item: any) => ({
      label: item.deptName,
      value: item.deptId
    }));
    let additionList = [...admissionClassroomType];
    if (type === "onlyDepartment") {
      additionList = [];
    }
    if (type === "expertDataBase") {
      additionList = [...admissionExpertType];
    }
    departmentEnum.value = [...additionList, ...deptList];
  };

  (isAdminRole.value || type === "expertDataBase") && fetchDepartmentEnum();
  return {
    departmentEnum
  };
};

export interface TaskEnumItem extends Common.enumDict {
  date: string;
}
export const useAllRecruitTaskEnum = () => {
  const allRecruitTaskEnum = ref<TaskEnumItem[]>([]);
  const _getList = async () => {
    const { list = [] } = await getRecruitManageList({ curPage: 1, pageSize: 1000 });
    allRecruitTaskEnum.value = list.map((item: any) => ({
      label: item.taskName,
      value: item.taskId,
      date: item.taskDay?.slice(0, 10)
    }));
  };
  _getList();
  return {
    allRecruitTaskEnum
  };
};

export const useAllJudgeRecruitTaskEnum = () => {
  const allJudgeRecruitTaskEnum = ref<Common.enumDict[]>([]);
  const _getList = async () => {
    const { list = [] } = await getAllJudgeRecruitTaskList({});
    allJudgeRecruitTaskEnum.value = list.map((item: any) => ({
      label: item.judgeTaskName,
      value: item.judgeTaskId
    }));
  };
  _getList();
  return {
    allJudgeRecruitTaskEnum
  };
};
