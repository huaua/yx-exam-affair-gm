package main

import (
	_ "yx-exam-affair-gm/app/biz/bizApi"
	_ "yx-exam-affair-gm/app/sync/api"
	_ "yx-exam-affair-gm/app/user/userApi"

	"github.com/sealsee/web-base/public/app"
)

func main() {
	app.RunServerGraceful()
}
