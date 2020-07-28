package storage

//
//import (
//	"errors"
//	"io"
//	"strings"
//	. "github.com/xmkuban/utils/components/storage/define"
//	"github.com/xmkuban/utils/geoip"
//)
//
//type UploadResult struct {
//	URL  string
//	Md5  string
//	Size int64
//	Err  error
//}
//
//type GetFormUploadURLResult struct {
//	URL  string
//	Info map[string]string
//	Err  error
//}
//
//type GetUploadSignURLResult struct {
//	URL    string
//	Expire int
//	Err    error
//}
//
//type GetMultiPartUploadInfoResult struct {
//	Info *MultiPartUploadInfo
//	Err  error
//}
//
//type MultiRegion struct {
//	regions map[string]StorageInterface
//	def     StorageInterface
//}
//
//func (mr *MultiRegion) SetRegion(region string, inst StorageInterface, isDefault bool) {
//	region = strings.ToUpper(region)
//	mr.regions[region] = inst
//	if isDefault {
//		mr.def = inst
//	}
//}
//
//func (mr *MultiRegion) GetStorage(region string) StorageInterface {
//	region = strings.ToUpper(region)
//
//	inst, ok := mr.regions[region]
//	if !ok {
//		return mr.def
//	}
//	return inst
//}
//
//func (mr *MultiRegion) GetStorageByIP(ip string) StorageInterface {
//	geo, err := geoip.New()
//	if err != nil {
//		return mr.def
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.GetStorage(c.Short)
//}
//
//func (mr *MultiRegion) Upload2Region(region string, r io.Reader, path string) (*UploadResult, error) {
//	res := new(UploadResult)
//
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return nil, errors.New("storage not init proper")
//	}
//
//	res.URL, res.Md5, res.Size, res.Err = inst.Upload(r, path)
//	return res, nil
//}
//
//func (mr *MultiRegion) Upload2RegionWithContentType(region string, r io.Reader, path string, contentType string) (*UploadResult, error) {
//	res := new(UploadResult)
//
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return nil, errors.New("storage not init proper")
//	}
//
//	res.URL, res.Md5, res.Size, res.Err = inst.UploadWithContentType(r, path, contentType)
//	return res, nil
//}
//
//func (mr *MultiRegion) Upload2RegionByIP(ip string, r io.Reader, path string) (*UploadResult, error) {
//	geo, err := geoip.New()
//	if err != nil {
//		return nil, err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.Upload2Region(c.Short, r, path)
//}
//
//func (mr *MultiRegion) Upload2RegionWithContentTypeByIP(ip string, r io.Reader, path string, contentType string) (*UploadResult, error) {
//	geo, err := geoip.New()
//	if err != nil {
//		return nil, err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.Upload2RegionWithContentType(c.Short, r, path, contentType)
//}
//
//func (mr *MultiRegion) GetFormUploadURL(region string, key string, metaData map[string]string) (*GetFormUploadURLResult, error) {
//	res := new(GetFormUploadURLResult)
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return nil, errors.New("storage not init proper")
//	}
//
//	res.URL, res.Info, res.Err = inst.GetFormUploadURL(key, metaData)
//	return res, nil
//}
//
//func (mr *MultiRegion) GetFormUploadURLByIP(ip string, key string, metaData map[string]string) (*GetFormUploadURLResult, error) {
//	geo, err := geoip.New()
//	if err != nil {
//		return nil, err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.GetFormUploadURL(c.Short, key, metaData)
//}
//
//func (mr *MultiRegion) GetUploadSignURL(region string, key string, metaData map[string]string) (*GetUploadSignURLResult, error) {
//	res := new(GetUploadSignURLResult)
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return nil, errors.New("storage not init proper")
//	}
//
//	res.URL, res.Expire, res.Err = inst.GetUploadSignURL(key, metaData)
//	return res, nil
//}
//
//func (mr *MultiRegion) GetUploadSignURLByIP(ip string, key string, metaData map[string]string) (*GetUploadSignURLResult, error) {
//	geo, err := geoip.New()
//	if err != nil {
//		return nil, err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.GetUploadSignURL(c.Short, key, metaData)
//}
//
//func (mr *MultiRegion) GetMultiPartUploadInfoWithContentType(region, key, contentType string, fileSize int64, partSize int64, metaData map[string]string) (*GetMultiPartUploadInfoResult, error) {
//	res := new(GetMultiPartUploadInfoResult)
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return nil, errors.New("storage not init proper")
//	}
//
//	res.Info, res.Err = inst.GetMultiPartUploadInfoWithContentType(key, contentType, fileSize, partSize, metaData)
//	return res, nil
//}
//
//func (mr *MultiRegion) GetMultiPartUploadInfoWithContentTypeByIP(ip, key, contentType string, fileSize int64, partSize int64, metaData map[string]string) (*GetMultiPartUploadInfoResult, error) {
//	geo, err := geoip.New()
//	if err != nil {
//		return nil, err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.GetMultiPartUploadInfoWithContentType(c.Short, key, contentType, fileSize, partSize, metaData)
//}
//
//func (mr *MultiRegion) CompleteMultiPartUpload(region, key string, etags MultiPartUploadEtags) error {
//	inst := mr.GetStorage(region)
//
//	if inst == nil {
//		return errors.New("storage not init proper")
//	}
//
//	return inst.CompleteMultiPartUpload(key, etags)
//}
//
//func (mr *MultiRegion) CompleteMultiPartUploadByIP(ip, key string, etags MultiPartUploadEtags) error {
//	geo, err := geoip.New()
//	if err != nil {
//		return err
//	}
//
//	c := geo.LookupByStr(ip)
//	return mr.CompleteMultiPartUpload(c.Short, key, etags)
//}
