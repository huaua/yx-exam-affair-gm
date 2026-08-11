package myutils

import (
	"archive/zip"
	"bytes"
	"regexp"
	"strconv"
	"yx-exam-affair-gm/app/hint"

	"github.com/sealsee/web-base/public/IOFile"
)

func CompressZipFiles(files map[string]string) string {
	// 创建buf用于存储zip中的文件
	var zipBuffer *bytes.Buffer = new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	// 遍历文件map  真实文件名->服务器文件地址
	for k, v := range files {
		// 多文件，打包成一个zip
		zw, _ := zipWriter.Create(k)
		data, err := IOFile.GetConfig().Download(v)
		if err != nil {
			panic(hint.SYS_FILE_READ_FAIL)
		}

		zw.Write(data)
	}
	// 处理压缩文件
	zipWriter.Close()
	// 上传zip文件得到profile地址
	// zipFilePath, err := IOFile.GetConfig().PublicUploadFile(IOFile.NewFileParamsNameBuffer(fileUtils.GetFileNameRandom(0, "zip"), zipBuffer))
	zipFileUrl, err := IOFile.GetConfig().Upload(zipBuffer, "0", "zip", true)
	if err != nil {
		panic(hint.SYS_FILE_READ_FAIL)
	}
	return zipFileUrl
}

func StrToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		// 可能字符串 str 不是合法的整数格式，处理错误
		panic(err)
	}
	return i
}

func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return i
}

func StrToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// 获取数字对应的中文大写数字，支持1-6
func GetNumberZh(num int) string {
	if num == 1 {
		return "一"
	} else if num == 2 {
		return "二"
	} else if num == 3 {
		return "三"
	} else if num == 4 {
		return "四"
	} else if num == 5 {
		return "五"
	} else if num == 6 {
		return "六"
	} else {
		return ""
	}
}

// 判断数字是否为正整数
func IsPositiveInteger(str string) bool {
	if len(str) <= 0 {
		return false
	}
	matched, _ := regexp.MatchString(`^[1-9]\d*$`, str)
	return matched
}

// 判断字符串是否只包含数字和字母
func IsNumberAndLetter(str string) bool {
	if len(str) <= 0 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, str)
	return matched
}
