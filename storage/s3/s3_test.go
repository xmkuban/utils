package s3

//
//import (
//	"testing"
//)
//
//func TestGenMultiUploadInfo(t *testing.T) {
//	s3 := NewWithAutherization("ap-northeast-2", "chriswutest", "AKIAIFMY2Q3BEQJDMQUA", "1Kd8bZdPPpaV74UHk/j7mXSxyNeKadWmCG2L8zJn")
//	//s3 := New("ap-northeast-2", "chriswutest")
//	info, err := s3.GetMultiPartUploadInfo("test.zip", 2411381700000, 5242880, map[string]string{})
//	if err != nil {
//		t.Errorf("get multipart upload info fail:%s", err)
//		return
//	}
//	t.Logf("split %d block(s)", info.PartsNum)
//}
