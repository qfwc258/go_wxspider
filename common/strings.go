package common
import(
	"fmt"
	"strings"
	"strconv"
)
//将字符串中的占位符格式化，功能类似c#中的string.format()
func Format(str string,formatStrs ...interface{})string{
	s:=str
	for i:=0;i<len(formatStrs);i++{
		num:=strconv.Itoa(i)
		f:="{"+num+"}"
		s=strings.Replace(s,f,fmt.Sprintf("%v",formatStrs[i]),-1)
	}
	return s	
}