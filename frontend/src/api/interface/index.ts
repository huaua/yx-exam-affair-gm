export * from "./common";
export * from "./login";
export * from "./user";
export * from "./department";
export * from "./teacher";
export * from "./classroom";
export * from "./arrange";
export * from "./student";
export * from "./judgeExpert";

// 请求响应参数（不包含data）
export interface Result {
  code: string;
  msg: string;
}

// 请求响应参数（包含data）
export interface ResultData<T = any> extends Result {
  data: T;
}

interface Pagination {
  curPage: number;
  pageSize: number;
}

// 分页响应参数
export interface ResPage<T> extends Pagination {
  list: T[];
  totalSize: number;
}

// 分页请求参数
export type ReqPage<T = {}> = T & Pagination;

// 文件上传模块
export namespace Upload {
  export interface ResFileUrl {
    fileUrl: string;
  }
}
