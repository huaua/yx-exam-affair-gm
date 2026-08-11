import { ReqPage, ResPage, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 公共模块
 */

// 获取字典列表
export const getDictList = (params: ReqPage<Common.ReqDictionary>) => {
  return http.get<ResPage<Common.ResDictionary>>(PROXY_TAG + `/api/auth/sys_dict_data/get_dicts`, params, { cancel: false });
};

// 上传（导入Excel）
export const upLoad = (params: any) => {
  return http.post(PROXY_TAG + `/common/upload_private`, params);
};
