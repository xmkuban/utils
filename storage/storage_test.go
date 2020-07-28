package storage

import (
	"testing"
)

func TestS3(t *testing.T) {
	//s3 := GetS3("test", "ap-northeast-2", "chriswutest")
	//url, info, err := s3.GetFormUploadURL("test.zip", map[string]string{"acl": "public-read"})
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	t.Log(url)
	//	t.Logf("%+v", info)
	//}
}

func TestQiniu(t *testing.T) {
	//api := GetQiniu("test", "", "", "course")
	//info, err := api.GetUploadToken(60)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//j, _ := json.Marshal(info)
	//fmt.Println(string(j))
}
