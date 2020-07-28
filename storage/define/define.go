package define

import "io"

type StorageInterface interface {
	Upload(r io.Reader, path string) (uri, md5 string, size int64, err error)
	UploadWithContentType(r io.Reader, path string, contentType string) (uri, md5 string, size int64, err error)
	GetUploadURL(key string, metaData map[string]string) (uri string, headInfos map[string]interface{}, err error)
	GetUploadSignURL(key string, metaData map[string]string) (u string, expire int, err error)
	GetFormUploadURL(key string, metaData map[string]string) (uri string, formInfos map[string]string, err error)
	GetFormUploadURLWithContentType(key, contentType string, metaData map[string]string) (uri string, formInfos map[string]string, err error)
	GetDownlaodURL(key string) (uri string, err error)
	GetDownloadSignURL(key string, expire int64) (u string, err error)
	GetMultiPartUploadInfo(key string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error)
	GetMultiPartUploadInfoWithContentType(key string, contentType string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error)
	AbortMultiPartUploadInfo(key string, uploadID string) error
	CompleteMultiPartUpload(key string, etags MultiPartUploadEtags) error
	Set(key string, value interface{})
	GetUploadToken(expire uint32) (info *UploadTokenInfo, err error)
	GetUploadTokenWithCallback(url, body, bodyType string, expire uint32) (info *UploadTokenInfo, err error)
}

type StorageConfigInterface interface {
	SetCredential(args ...interface{}) error
	GenerateInst() StorageInterface
}

type MultiPartUploadInfo struct {
	MultiPartEndpointsMap map[int]string `json:"multipart_upload_endpoints"`
	PartsNum              int            `json:"part_num"`
	UploadID              string         `json:"upload_id"`
}

type MultiPartUploadEtags struct {
	Etags    []string `json:"multipart_upload_etags"`
	UploadID string   `json:"multipart_upload_id"`
}

type UploadTokenInfo struct {
	Token string `json:"token"`
}
