package s3

//
//import (
//	"crypto/md5"
//	"encoding/hex"
//	"errors"
//	"fmt"
//	"github.com/minio/minio-go"
//	"io"
//	"io/ioutil"
//	"math"
//	"net/url"
//	"time"
//
//	"github.com/xmkuban/utils/logger"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/credentials"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/s3"
//	"github.com/aws/aws-sdk-go/service/s3/s3manager"
//
//	"bytes"
//
//	. "github.com/xmkuban/utils/components/storage/define"
//)
//
//func GenFeedbackLogFile(userId int, deviceUniqueId, clientUniqueId string) string {
//	return fmt.Sprintf("%d_%s_%s_%d.zip", userId, deviceUniqueId, clientUniqueId, time.Now().Unix())
//}
//func GenThumbnailFile(userID int, liveID, suffix string) string {
//	return fmt.Sprintf("%d_%s_%d.%s", userID, liveID, time.Now().Unix(), suffix)
//}
//
//var S3Service S3
//
//type S3 struct {
//	region string
//	bucket string
//
//	keyId     string
//	accessKey string
//
//	PreSignTimeout time.Duration
//	UseAccelerate  bool
//}
//
//type S3Config struct {
//	Region    string `split_words:"true" yaml:"region"`
//	Bucket    string `split_words:"bucket" yaml:"bucket"`
//	KeyID     string `split_words:"key_id" yaml:"key_id"`
//	AccessKey string `split_words:"true" yaml:"access_key"`
//}
//
//func (c *S3Config) SetCredential(args ...interface{}) error {
//	if len(args) < 4 {
//		return errors.New("the count of arguments is not enough")
//	}
//	var ok bool
//	regionRaw := args[0]
//	bucketRaw := args[1]
//	keyIDRaw := args[2]
//	accessRaw := args[3]
//
//	c.Region, ok = regionRaw.(string)
//	if !ok || c.Region == "" {
//		return errors.New("region need to be string")
//	}
//
//	c.Bucket, ok = bucketRaw.(string)
//	if !ok || c.Bucket == "" {
//		return errors.New("bucket need to be string")
//	}
//
//	c.KeyID, ok = keyIDRaw.(string)
//	if !ok || c.KeyID == "" {
//		return errors.New("key_id need to be string")
//	}
//
//	c.AccessKey, ok = accessRaw.(string)
//	if !ok || c.AccessKey == "" {
//		return errors.New("access_key need to be string")
//	}
//
//	return nil
//}
//
//func (c *S3Config) GenerateInst() StorageInterface {
//	return NewWithAutherization(c.Region, c.Bucket, c.KeyID, c.AccessKey)
//}
//
//func New(region, bucket string) *S3 {
//	logger.Debugf("new s3 uploader region:%s,bucket:%s", region, bucket)
//	return &S3{
//		region:         region,
//		bucket:         bucket,
//		PreSignTimeout: 15 * time.Minute,
//	}
//}
//
//func NewWithAutherization(region, bucket, keyId, accessKey string) *S3 {
//	logger.Debugf("new s3 uploader region:%s,bucket:%s with autherization", region, bucket)
//
//	return &S3{
//		region:         region,
//		bucket:         bucket,
//		keyId:          keyId,
//		accessKey:      accessKey,
//		PreSignTimeout: 15 * time.Minute,
//	}
//}
//
//func (s *S3) Upload(r io.Reader, path string) (s3Url string, md5Str string, size int64, err error) {
//	data, err := ioutil.ReadAll(r)
//	if err != nil {
//		return
//	}
//	size = int64(len(data))
//	hash := md5.New()
//	_, err = hash.Write(data)
//	if err != nil {
//		return
//	}
//	md5sum := hash.Sum(nil)
//	md5Str = hex.EncodeToString(md5sum)
//
//	sess, err := s.getSess()
//	if err != nil {
//		return
//	}
//
//	uploader := s3manager.NewUploader(sess)
//	result, err := uploader.Upload(&s3manager.UploadInput{
//		Body:   bytes.NewReader(data),
//		Bucket: aws.String(s.bucket),
//		Key:    aws.String(path),
//	})
//
//	if err != nil {
//		return
//	}
//	s3Url, _ = url.QueryUnescape(result.Location)
//	return
//}
//
//func (s *S3) UploadWithContentType(r io.Reader, path string, contentType string) (s3Url, md5Str string, size int64, err error) {
//	data, err := ioutil.ReadAll(r)
//	if err != nil {
//		return
//	}
//	size = int64(len(data))
//	hash := md5.New()
//	_, err = hash.Write(data)
//	if err != nil {
//		return
//	}
//	md5sum := hash.Sum(nil)
//	md5Str = hex.EncodeToString(md5sum)
//	sess, err := s.getSess()
//	if err != nil {
//		return
//	}
//
//	uploader := s3manager.NewUploader(sess)
//	result, err := uploader.Upload(&s3manager.UploadInput{
//		Body:        bytes.NewReader(data),
//		Bucket:      aws.String(s.bucket),
//		Key:         aws.String(path),
//		ContentType: aws.String(contentType),
//	})
//
//	if err != nil {
//		return
//	}
//	s3Url, _ = url.QueryUnescape(result.Location)
//	return
//}
//
//func (s *S3) GetFormUploadURL(key string, metaData map[string]string) (uri string, info map[string]string, err error) {
//	keyID, accessKey := s.getCredentials()
//	if keyID == "" || accessKey == "" {
//		return uri, info, errors.New("missing s3 credentials info")
//	}
//	policy := minio.NewPostPolicy()
//	policy.SetKey(key)
//	policy.SetBucket(s.bucket)
//	policy.SetExpires(time.Now().UTC().Add(15 * time.Minute))
//	var s3Client *minio.Client
//	if s.region == "" {
//		s3Client, err = minio.New("s3.amazonaws.com", keyID, accessKey, true)
//	} else {
//		s3Client, err = minio.NewWithRegion("s3.amazonaws.com", keyID, accessKey, true, s.region)
//	}
//
//	if s.UseAccelerate {
//		s3Client.SetS3TransferAccelerate("s3-accelerate.amazonaws.com")
//	}
//
//	if err != nil {
//		return uri, info, err
//	}
//	_uri, _info, err := s3Client.PresignedPostPolicy(policy)
//	if err != nil {
//		return uri, info, err
//	}
//
//	uri = _uri.String()
//	info = _info
//	return
//}
//
//func (s *S3) GetFormUploadURLWithContentType(key, contentType string, metaData map[string]string) (uri string, info map[string]string, err error) {
//	keyID, accessKey := s.getCredentials()
//	if keyID == "" || accessKey == "" {
//		return uri, info, errors.New("missing s3 credentials info")
//	}
//	policy := minio.NewPostPolicy()
//	policy.SetKey(key)
//	policy.SetBucket(s.bucket)
//	policy.SetContentType(contentType)
//	policy.SetExpires(time.Now().UTC().Add(15 * time.Minute))
//	var s3Client *minio.Client
//	if s.region == "" {
//		s3Client, err = minio.New("s3.amazonaws.com", keyID, accessKey, true)
//	} else {
//		s3Client, err = minio.NewWithRegion("s3.amazonaws.com", keyID, accessKey, true, s.region)
//	}
//
//	if s.UseAccelerate {
//		s3Client.SetS3TransferAccelerate("s3-accelerate.amazonaws.com")
//	}
//
//	if err != nil {
//		return uri, info, err
//	}
//	_uri, _info, err := s3Client.PresignedPostPolicy(policy)
//	if err != nil {
//		return uri, info, err
//	}
//
//	uri = _uri.String()
//	info = _info
//	return
//}
//
//func (s *S3) GetUploadSignURL(key string, metaData map[string]string) (u string, expire int, err error) {
//	sess, err := session.NewSession(&aws.Config{Region: aws.String(s.region)})
//	if err != nil {
//		return u, 0, err
//	}
//
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
//		//temp remove
//		//ACL:      aws.String(acl),
//		Bucket:   aws.String(s.bucket),
//		Key:      aws.String(key),
//		Metadata: s.genMetaData(metaData),
//	})
//
//	u, err = req.Presign(s.PreSignTimeout)
//
//	expire = int(s.PreSignTimeout / time.Second)
//
//	return u, expire, err
//}
//
//func (s *S3) GetUploadURL(key string, metaData map[string]string) (uri string, info map[string]interface{}, err error) {
//
//	var expire int
//	uri, expire, err = s.GetUploadSignURL(key, metaData)
//	info = make(map[string]interface{})
//	info["expire"] = expire
//	return
//}
//
//func (s *S3) GetDownloadSignURL(key string, expire int64) (u string, err error) {
//	sess, err := s.getSess()
//	if err != nil {
//		return u, err
//	}
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: aws.String(s.bucket),
//		Key:    aws.String(key),
//	})
//
//	u, err = req.Presign(time.Duration(expire) * time.Second)
//	return
//}
//
//func (s *S3) GetDownlaodURL(key string) (uri string, err error) {
//	// 传0为永久的链接，返回原始的url，但是由于bucket的设置是非public实际上是没有访问权限
//	return s.GetDownloadSignURL(key, 0)
//}
//
//func (s *S3) genMetaData(oData map[string]string) (res map[string]*string) {
//	res = make(map[string]*string)
//	if oData == nil {
//		return res
//	}
//	for k, v := range oData {
//		res[k] = aws.String(v)
//	}
//	return
//}
//
//func (s *S3) Set(key string, value interface{}) {
//	if key == "accelerate" {
//		valueBool, _ := value.(bool)
//		s.UseAccelerate = valueBool
//	}
//}
//
//func (s *S3) getSess() (*session.Session, error) {
//	if s.accessKey != "" && s.keyId != "" {
//		cred := credentials.NewStaticCredentials(s.keyId, s.accessKey, "")
//		sess, err := session.NewSession(&aws.Config{Region: aws.String(s.region), Credentials: cred})
//		return sess, err
//	} else {
//		sess, err := session.NewSession(&aws.Config{Region: aws.String(s.region)})
//		return sess, err
//	}
//}
//
//func (s *S3) getCredentials() (keyID string, secretKey string) {
//	if s.accessKey != "" && s.keyId != "" {
//		return s.keyId, s.accessKey
//	}
//
//	creds := credentials.NewEnvCredentials()
//
//	credValue, err := creds.Get()
//	if err != nil {
//		logger.Error("get aws credentials fail")
//		return "", ""
//	}
//
//	return credValue.AccessKeyID, credValue.SecretAccessKey
//}
//
//func (s *S3) GetMultiPartUploadInfo(key string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error) {
//	sess, err := s.getSess()
//	if err != nil {
//		return info, err
//	}
//
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	req, output := svc.CreateMultipartUploadRequest(&s3.CreateMultipartUploadInput{
//		Bucket:   aws.String(s.bucket),
//		Key:      aws.String(key),
//		Metadata: s.genMetaData(metaData),
//	})
//
//	err = req.Send()
//	if err != nil {
//		return
//	}
//
//	info = new(MultiPartUploadInfo)
//	info.UploadID = *output.UploadId
//	info.MultiPartEndpointsMap = make(map[int]string)
//
//	blocks := math.Ceil(float64(fileSize) / float64(partSize))
//	info.PartsNum = int(blocks)
//	for i := 0; i < info.PartsNum; i++ {
//		contentLength := int64(0)
//		if fileSize >= partSize {
//			contentLength = partSize
//		} else {
//			contentLength = fileSize
//		}
//		fileSize = fileSize - partSize
//		partReq, _ := svc.UploadPartRequest(&s3.UploadPartInput{
//			Bucket:        aws.String(s.bucket),
//			PartNumber:    aws.Int64(int64(i + 1)),
//			UploadId:      output.UploadId,
//			Key:           aws.String(key),
//			ContentLength: aws.Int64(contentLength),
//		})
//		partUploadURL, err := partReq.Presign(24 * time.Hour)
//		if err != nil {
//			return nil, err
//		}
//
//		info.MultiPartEndpointsMap[i+1] = partUploadURL
//	}
//	return
//}
//
//func (s *S3) GetMultiPartUploadInfoWithContentType(key string, contentType string, fileSize int64, partSize int64, metaData map[string]string) (info *MultiPartUploadInfo, err error) {
//	sess, err := s.getSess()
//	if err != nil {
//		return info, err
//	}
//
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	req, output := svc.CreateMultipartUploadRequest(&s3.CreateMultipartUploadInput{
//		Bucket:      aws.String(s.bucket),
//		Key:         aws.String(key),
//		Metadata:    s.genMetaData(metaData),
//		ContentType: aws.String(contentType),
//	})
//
//	err = req.Send()
//	if err != nil {
//		return
//	}
//
//	info = new(MultiPartUploadInfo)
//	info.UploadID = *output.UploadId
//	info.MultiPartEndpointsMap = make(map[int]string)
//
//	blocks := math.Ceil(float64(fileSize) / float64(partSize))
//	info.PartsNum = int(blocks)
//	for i := 0; i < info.PartsNum; i++ {
//		contentLength := int64(0)
//		if fileSize >= partSize {
//			contentLength = partSize
//		} else {
//			contentLength = fileSize
//		}
//		fileSize = fileSize - partSize
//		partReq, _ := svc.UploadPartRequest(&s3.UploadPartInput{
//			Bucket:        aws.String(s.bucket),
//			PartNumber:    aws.Int64(int64(i + 1)),
//			UploadId:      output.UploadId,
//			Key:           aws.String(key),
//			ContentLength: aws.Int64(contentLength),
//		})
//		partUploadURL, err := partReq.Presign(24 * time.Hour)
//		if err != nil {
//			return nil, err
//		}
//
//		info.MultiPartEndpointsMap[i+1] = partUploadURL
//	}
//	return
//}
//
//func (s *S3) AbortMultiPartUploadInfo(key string, uploadID string) error {
//	sess, err := s.getSess()
//	if err != nil {
//		return err
//	}
//
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	req, _ := svc.AbortMultipartUploadRequest(&s3.AbortMultipartUploadInput{
//		Bucket:   aws.String(s.bucket),
//		Key:      aws.String(key),
//		UploadId: aws.String(uploadID),
//	})
//
//	return req.Send()
//}
//
//func (s *S3) CompleteMultiPartUpload(key string, etags MultiPartUploadEtags) error {
//	sess, err := s.getSess()
//	if err != nil {
//		return err
//	}
//
//	svc := s3.New(sess, &aws.Config{S3UseAccelerate: aws.Bool(s.UseAccelerate)})
//	multipartInfo := &s3.CompletedMultipartUpload{}
//	for i, etag := range etags.Etags {
//		multipartInfo.Parts = append(multipartInfo.Parts, &s3.CompletedPart{
//			ETag:       aws.String(etag),
//			PartNumber: aws.Int64(int64(i + 1)),
//		})
//	}
//
//	req, _ := svc.CompleteMultipartUploadRequest(&s3.CompleteMultipartUploadInput{
//		Key:             aws.String(key),
//		UploadId:        aws.String(etags.UploadID),
//		Bucket:          aws.String(s.bucket),
//		MultipartUpload: multipartInfo,
//	})
//
//	return req.Send()
//}
//func (q *S3) GetUploadToken(expire uint32) (info *UploadTokenInfo, err error) {
//	return nil, errors.New("not supported")
//}
//
//func (q *S3) GetUploadTokenWithCallback(url, body, bodyType string, expire uint32) (info *UploadTokenInfo, err error) {
//	return nil, errors.New("not supported")
//}
