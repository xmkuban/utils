package qiniu

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"fmt"

	. "github.com/xmkuban/utils/storage/define"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

type Qiniu struct {
	accessKey string
	secretKey string
	bucket    string
	Zone      string
	Domain    string
}

type QiniuConfig struct {
	AccessKey string `split_words:"true" yaml:"access_key"`
	SecretKey string `split_words:"true" yaml:"secret_key"`
	Bucket    string `split_words:"true" yaml:"bucket"`
	Zone      string `split_words:"true" yaml:"zone"`
}

func (c *QiniuConfig) SetCredential(args ...interface{}) error {
	if len(args) < 3 {
		return errors.New("the count of arguments is not enough")
	}
	var ok bool
	accessRaw := args[0]
	secretRaw := args[1]
	bucketRaw := args[2]

	c.AccessKey, ok = accessRaw.(string)
	if !ok || c.AccessKey == "" {
		return errors.New("access_key need to be string")
	}

	c.SecretKey, ok = secretRaw.(string)
	if !ok || c.SecretKey == "" {
		return errors.New("secret_key need to be string")
	}

	c.Bucket, ok = bucketRaw.(string)
	if !ok || c.Bucket == "" {
		return errors.New("bucket need to be string")
	}

	return nil
}

func (c *QiniuConfig) GenerateInst() StorageInterface {
	return New(c.AccessKey, c.SecretKey, c.Bucket, c.Zone)
}

func New(accessKey, secretKey, bucket string, zone string) *Qiniu {
	return &Qiniu{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		Zone:      zone,
	}
}

func (q *Qiniu) Upload(r io.Reader, path string) (uri, md5Str string, size int64, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return uri, md5Str, size, err
	}

	size = int64(len(data))

	hash := md5.New()
	_, err = hash.Write(data)
	if err != nil {
		return
	}
	md5sum := hash.Sum(nil)
	md5Str = hex.EncodeToString(md5sum)

	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	uptoken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	uploader := storage.NewFormUploader(&cfg)
	err = uploader.Put(context.Background(), nil, uptoken, path, bytes.NewReader(data), size, &storage.PutExtra{
		UpHost: q.GetServerUploadHost(),
	})

	if err != nil {
		return
	}

	uri = storage.MakePublicURL(q.Domain, path)
	return
}

func (q *Qiniu) GetServerUploadHost() string {
	if q.Zone == "z0" {
		return "https://up.qiniup.com"
	} else if q.Zone == "z1" {
		return "https://up-z1.qiniup.com"
	} else if q.Zone == "z2" {
		return "https://up-z2.qiniup.com"
	} else if q.Zone == "na0" {
		return "https://up-na0.qiniup.com"
	} else if q.Zone == "as0" {
		return "https://up-as0.qiniup.com"
	}
	return "https://up.qiniup.com"
}

func (q *Qiniu) GetClientUploadHost() string {
	if q.Zone == "z0" {
		return "https://upload.qiniup.com"
	} else if q.Zone == "z1" {
		return "https://upload-z1.qiniup.com"
	} else if q.Zone == "z2" {
		return "https://upload-z2.qiniup.com"
	} else if q.Zone == "na0" {
		return "https://upload-na0.qiniup.com"
	} else if q.Zone == "as0" {
		return "https://upload-as0.qiniup.com"
	}
	return "https://upload.qiniup.com"
}

func (q *Qiniu) UploadWithContentType(r io.Reader, path string, contentType string) (uri, md5Str string, size int64, err error) {
	err = errors.New("not suppored")
	return
}

func (q *Qiniu) GetFormUploadURL(key string, metaData map[string]string) (uri string, info map[string]string, err error) {
	policy := storage.PutPolicy{
		Scope:   fmt.Sprintf("%s:%s", q.bucket, key),
		SaveKey: key,
		Expires: 900, //15 min
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	uptoken := policy.UploadToken(mac)

	info = make(map[string]string)
	info["key"] = key

	for k, v := range metaData {
		info[fmt.Sprintf("x:%s", k)] = v
	}
	uri = fmt.Sprintf("http://%s?token=%s", storage.ZoneHuadong.SrcUpHosts[0], uptoken)
	err = nil
	return
}
func (q *Qiniu) GetFormUploadURLWithContentType(key, contentType string, metaData map[string]string) (uri string, info map[string]string, err error) {
	policy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", q.bucket, key),
		SaveKey:    key,
		Expires:    900, //15 min
		DetectMime: 1,
	}

	mac := qbox.NewMac(q.accessKey, q.secretKey)
	uptoken := policy.UploadToken(mac)

	info = make(map[string]string)
	info["key"] = key

	for k, v := range metaData {
		info[fmt.Sprintf("x:%s", k)] = v
	}
	uri = fmt.Sprintf("http://%s?token=%s", storage.ZoneHuadong.SrcUpHosts[0], uptoken)
	err = nil
	return
}
func (q *Qiniu) GetUploadURL(key string, metaData map[string]string) (uri string, info map[string]interface{}, err error) {
	err = errors.New("qiniu only support form upload")
	return

}

func (q *Qiniu) GetUploadSignURL(key string, metaData map[string]string) (u string, expire int, err error) {
	err = errors.New("qiniu only support form upload")
	return
}

func (q *Qiniu) GetDownloadSignURL(key string, expire int64) (u string, err error) {
	if q.Domain == "" {
		return "", errors.New("domain not set")
	}

	mac := qbox.NewMac(q.accessKey, q.secretKey)
	deadline := time.Now().Unix() + expire
	return storage.MakePrivateURL(mac, q.Domain, key, deadline), nil
}

func (q *Qiniu) GetDownlaodURL(key string) (uri string, err error) {
	if q.Domain == "" {
		return "", errors.New("domain not set")
	}
	uri = storage.MakePublicURL(q.Domain, key)
	return
}

func (q *Qiniu) Set(key string, value interface{}) {
	if key == "domain" {
		valStr, _ := value.(string)
		q.Domain = valStr
	}
}

func (q *Qiniu) GetMultiPartUploadInfo(key string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error) {
	return nil, errors.New("not supported")
}
func (q *Qiniu) GetMultiPartUploadInfoWithContentType(key string, contentType string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error) {
	return nil, errors.New("not supported")
}
func (q *Qiniu) AbortMultiPartUploadInfo(key string, uploadID string) error {
	return errors.New("not supported")
}
func (q *Qiniu) CompleteMultiPartUpload(key string, etags MultiPartUploadEtags) error {
	return errors.New("not supported")
}

func (q *Qiniu) GetUploadToken(expire uint32) (info *UploadTokenInfo, err error) {
	putPolicy := storage.PutPolicy{
		Scope:   q.bucket,
		Expires: uint64(expire),
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)
	info = new(UploadTokenInfo)
	info.Token = upToken
	return
}
func (q *Qiniu) GetUploadTokenWithCallback(url, body, bodyType string, expire uint32) (info *UploadTokenInfo, err error) {
	info = &UploadTokenInfo{}
	putPolicy := storage.PutPolicy{
		Scope:            q.bucket,
		CallbackURL:      url,
		CallbackBody:     body,
		CallbackBodyType: bodyType,
		Expires:          uint64(expire),
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)
	info = new(UploadTokenInfo)
	info.Token = upToken
	return
}
