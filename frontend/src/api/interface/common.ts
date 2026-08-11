/**
 * @name 公共模块
 */

export namespace Common {
  export type emptyObject = Record<string, never>;
  export interface ReqDictionary {
    dictType: string;
  }
  export interface ResDictionary {
    id: string;
    code: string; // 代码
    name: string; // 名称
    remark: string; // 描述
    state: string; // 状态
    action: string;
  }

  export interface enumDict {
    label: string;
    value: string;
  }
}
