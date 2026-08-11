### tm_admit

1. 下载依赖: go mod tidy

2. 编译打包: go build

3. 本地运行: go run tm_admit_main.go

4. 后台运行: setid go run tm_admit_main.go

---

#### 镜像加速
1. 七牛 CDN
go env -w  GOPROXY=https://goproxy.cn,direct

2. 阿里云
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

3. 官方
go env -w  GOPROXY=https://goproxy.io,direct

##### Windows下设置 GOPROXY 的命令为：
go env -w GOPROXY=https://goproxy.cn,direct
##### MacOS或Linux 下设置 GOPROXY 的命令为：
export GOPROXY=https://goproxy.cn,direct

---

