// 国美考务系统 - 角色操作说明截图脚本
// 用法: node shot.mjs
// 依赖: puppeteer-core (已安装), 系统已装 Chrome, 前端运行在 http://127.0.0.1:8848
import puppeteer from "puppeteer-core";
import fs from "fs";

const BASE = "http://127.0.0.1:8848";
const API = "http://127.0.0.1:30000";
const CHROME = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe";
const OUT = "d:/2026/codex/YWJLDP/docs/shots";
fs.mkdirSync(OUT, { recursive: true });

const accounts = [
  { name: "admin", user: "admin", pwd: "Admin123", role: "admissions" },
  { name: "zgh", user: "zgh", pwd: "Ysk002", role: "department" }
];

const pages = {
  admin: [
    ["01-首页", "/home/index"],
    ["02-考场征集任务管理", "/classroom/classroomRecruitManage"],
    ["03-考场征集确认", "/classroom/classroomRecruit"],
    ["04-征集考场查看", "/classroom/usedClassroom"],
    ["05-教职工管理", "/teacher/teacher"],
    ["06-监考老师征集任务", "/teacher/teacherRecruitManage"],
    ["07-征集老师审核", "/teacher/teacherRecruitView"],
    ["08-考场排考编排", "/classroomArrange/arrange"],
    ["09-考场监考分配", "/classroomArrange/assign"],
    ["10-考场编排查看", "/classroomArrange/arrangeView"],
    ["11-专家库管理", "/judgeExpert/expertDatabase"],
    ["12-禁止抽取管理", "/judgeExpert/expertBan"],
    ["13-评委征集任务管理", "/judgeExpert/judgeRecruitManage"],
    ["14-评委抽取管理", "/judgeExpert/judgeDrawManage"],
    ["15-评委上报审核", "/judgeExpert/judgeReportAudit"]
  ],
  zgh: [
    ["01-首页", "/home/index"],
    ["02-教室管理", "/classroom/classroom"],
    ["03-教职工管理", "/teacher/teacher"],
    ["04-专家库管理", "/judgeExpert/expertDatabase"],
    ["05-考场征集", "/classroom/classroomRecruit"],
    ["06-监考老师上报", "/teacher/teacherSubmit"],
    ["07-监考老师上报查看", "/teacher/teacherSubmitView"],
    ["08-评委专家上报", "/judgeExpert/judgeResponse"]
  ]
};

async function login(user, pwd) {
  const cap = await fetch(API + "/captcha").then(r => r.json());
  const uuid = cap.data.captcha.uuid;
  const loginRes = await fetch(API + "/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: user, password: pwd, code: "8888", uuid })
  }).then(r => r.json());
  const token = loginRes.data.token;
  const info = await fetch(API + "/get_info", { headers: { Authorization: token } }).then(r => r.json());
  let roles = info.data.roles;
  // 后端超级管理员(userId=1)硬编码返回 roles=["admin"]，但前端菜单枚举是 "admissions"，
  // 这里把 admin 改写为 admissions 以让菜单正确渲染（仅影响侧边菜单路由，页面内按钮权限靠后端 token 校验，不受影响）。
  if (user === "admin") roles = ["admissions"];
  const u = info.data.user;
  const userInfo = { userId: u.userId, userName: u.userName, nickName: u.nickName, avatar: u.avatar };
  // pinia persist: { token, roles, userInfo }
  return { token, store: { token, roles, userInfo } };
}

const browser = await puppeteer.launch({
  executablePath: CHROME,
  headless: "new",
  args: ["--no-sandbox", "--disable-setuid-sandbox", "--window-size=1680,1000"]
});

for (const acc of accounts) {
  const dir = `${OUT}/${acc.name}`;
  fs.mkdirSync(dir, { recursive: true });
  const { store } = await login(acc.user, acc.pwd);
  const page = await browser.newPage();
  await page.setViewport({ width: 1680, height: 1000 });
  // 先访问站点，建立 origin 以便写 localStorage
  await page.goto(BASE, { waitUntil: "networkidle2", timeout: 30000 });
  await page.evaluate((val) => {
    localStorage.setItem("geeker-user", JSON.stringify(val));
  }, store);
  // 刷新页面，让 pinia 从 localStorage 恢复 token/roles，让路由守卫放行
  await page.reload({ waitUntil: "networkidle2", timeout: 30000 });
  await new Promise(r => setTimeout(r, 1200));
  const list = pages[acc.name];
  for (const [title, path] of list) {
    const url = BASE + "/#" + path;
    try {
      await page.goto(url, { waitUntil: "networkidle2", timeout: 30000 });
      await new Promise(r => setTimeout(r, 1800)); // 等表格/菜单渲染
      await page.screenshot({ path: `${dir}/${title}.png`, fullPage: false });
      console.log(`OK  ${acc.name}/${title}`);
    } catch (e) {
      console.error(`FAIL ${acc.name}/${title}: ${e.message}`);
    }
  }
  await page.close();
}

await browser.close();
console.log("ALL DONE");
