package requests

import (
	"github.com/axgle/mahonia"
	"strings"
	"unsafe"
)

// ConvertString 是将 str 以 code 解码后重新编码为 UTF-8
// 如原字符为 GBK, 需要转成 UTF-8, 则使用 ConvertString(str, "GBK")
func ConvertString(str string, code string) string {
	enc := mahonia.NewEncoder(code)
	return enc.ConvertString(str)
}

// ConvertBytes 是将 bytes 以 code 解码后重新编码为 UTF-8
// 如原字符为 GBK, 需要转成 UTF-8, 则使用 ConvertBytes(bytes, "GBK")
func ConvertBytes(bytes *[]byte, code string) string {
	str := *(*string)(unsafe.Pointer(bytes))
	if strings.ToUpper(code) == "UTF-8" {
		return str
	}

	return ConvertString(str, code)
}
