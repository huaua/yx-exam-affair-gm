// ? Element 常用表单校验规则

/**
 *  @rule 手机号
 */
export function checkPhoneNumber(rule: any, value: any, callback: any) {
  const regexp = /^(((13[0-9]{1})|(15[0-9]{1})|(16[0-9]{1})|(17[3-8]{1})|(18[0-9]{1})|(19[0-9]{1})|(14[5-7]{1}))+\d{8})$/;
  if (value === "") callback("请输入手机号码");
  if (!regexp.test(value)) {
    callback(new Error("请输入正确的手机号码"));
  } else {
    return callback();
  }
}

/**
 * 密码
 * @description 6-18位必须包含大小写字母、数字组合；特殊字符可选
 */
export function isPassword(rule: any, value: any) {
  const reg = /^(?=.*[0-9])(?=.*[A-Z])(?=.*[a-z])[,._!@#$%^&*0-9a-zA-Z]{6,18}$/;
  return reg.test(value);
}

/**
 * 校验位数
 * @description 校验输入位数(使用时需要再包一层函数，以传入minLength和maxLength)
 */
export function isValidLength(minLength: number, maxLength: number, rule: any, value: any) {
  const reg = new RegExp(`^.{${minLength},${maxLength}}$`);
  return reg.test(value);
}

/**
 * 数字和字母
 * @description 仅允许大小写字母、数字；
 */
export function isAlphanumeric(rule: any, value: any) {
  const reg = /^[a-zA-Z0-9]+$/;
  return reg.test(value);
}
