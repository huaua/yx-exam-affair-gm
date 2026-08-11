/**
 * @name 登录模块
 */

export namespace Login {
  export interface ReqLoginForm {
    username: string;
    password: string;
    code: string;
    uuid: string;
  }
  export interface ResLogin {
    token: string;
  }
  export interface ResAuthButtons {
    [key: string]: string[];
  }
  export interface ResCaptcha {
    captcha: {
      uuid: string;
      img: string;
    };
  }
  export interface ResUserInfo {
    roles: string[];
    user: {
      userId: string;
      deptId: string;
      userName: string;
      nickName: string;
      avatar: string;
      isUsed: boolean;
      token: string;
    };
  }
}
