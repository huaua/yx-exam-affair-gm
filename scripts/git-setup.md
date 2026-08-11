# 国美考务系统 - Git 仓库恢复 / 换机配置指南

仓库地址：`git@github.com:huaua/yx-exam-affair-gm.git`（SSH）
主分支：`main`
提交身份：user.name=`YWJLDP`, user.email=`1373820144@qq.com`

---

## 一、全新克隆（换机 / 重装后）

```powershell
# 1. 安装 Git for Windows（默认装到 C:\Program Files\Git）
#    安装时勾选 "Add Git to PATH"

# 2. 克隆仓库（SSH，走 22 端口，无需代理）
git clone git@github.com:huaua/yx-exam-affair-gm.git D:\code\YWJLDP
cd D:\code\YWJLDP

# 3. 安装依赖并启动（见 README.md）
```

---

## 二、配置 SSH key（让 SSH 推送通过）

```powershell
# 生成密钥（无密码短语，便于脚本自动 push）
ssh-keygen -t ed25519 -C "1373820144@qq.com" -f "$env:USERPROFILE\.ssh\id_ed25519" -N ""

# 复制公钥，粘贴到 GitHub -> Settings -> SSH and GPG keys -> New SSH key
Get-Content "$env:USERPROFILE\.ssh\id_ed25519.pub"

# 验证
ssh -T git@github.com
# 看到 "Hi huaua! You've successfully authenticated" 即成功
```

> 注意：`ssh -T` 返回 `exit 1` 是正常的（GitHub 不提供 shell），只要看到上面的欢迎语就说明认证 OK。

---

## 三、网络与代理（GitHub 443 被防火墙拦截时）

- **SSH 推送（默认）**走 22 端口，不需要代理，推荐优先使用。
- 若要用 **HTTPS 方式**推送，需本地开代理（Clash / V2Ray 等，端口通常为 `127.0.0.1:7897`）：

```powershell
git config --global http.proxy  http://127.0.0.1:7897
git config --global https.proxy http://127.0.0.1:7897
```

关闭代理后若 HTTPS 推送失败，切回 SSH 即可（remote 已是 SSH 地址）。

---

## 四、日常提交与推送（一键脚本）

脚本：`scripts/git-commit.ps1`

```powershell
# 默认信息提交（含日期）
.\scripts\git-commit.ps1

# 指定提交信息
.\scripts\git-commit.ps1 "修复评委抽取导出空列问题"

# 提交并推送到远程
.\scripts\git-commit.ps1 "修复xx" -Push
```

脚本行为：进入仓库根目录 -> `git add -A` -> 有改动才提交 ->（可选 `-Push`）推到 `origin`。

---

## 五、将 GitHub 默认分支由 master 改为 main

当前远程默认分支仍是 `master`（首次推送时建的）。建议改为 `main`：

1. 打开 https://github.com/huaua/yx-exam-affair-gm/settings/branches
2. Default branch 改为 `main` -> Update
3. 然后本地执行 `git push origin --delete master` 删掉旧的 master 分支

（改默认分支后才能删除 master，否则 GitHub 拒绝删除当前默认分支）

---

## 六、.gitignore 要点（勿误提交）

已忽略：`node_modules/`、`dist/`、`.pnpm-store/`、`.runtime/`（运行环境工具链）、
`backend/*.exe`、`*.log`、`.env`、`.codebuddy/`、
`国美考务-源码+演示库/`（演示库源码副本，巨量重复）、临时 diff/构建产物。

> **切勿删除** `.runtime/go/**`、`.runtime/mysql/bin/**` 等运行环境 exe，它们是本地运行依赖，且已被 git 忽略。
