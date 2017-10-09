package common

import(
	"crypto/md5"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "io"
)

func GetGuid()string{
	b:=make([]byte,48)
	if _,err:=io.ReadFull(rand.Reader,b);err!=nil{
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

func getMd5String(s string)string{
	h:=md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

