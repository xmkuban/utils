package storage

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	. "github.com/xmkuban/utils/storage/define"
	"github.com/xmkuban/utils/storage/qiniu"
	//"github.com/xmkuban/utils/components/storage/s3"
)

var storages map[string]StorageInterface
var countryStorages map[string]StorageInterface
var storageConfTypeMap map[string]reflect.Type

func init() {
	storages = make(map[string]StorageInterface)
	countryStorages = make(map[string]StorageInterface)

	storageConfTypeMap = make(map[string]reflect.Type)
	//storageConfTypeMap["s3"] = reflect.TypeOf(s3.S3Config{})
	storageConfTypeMap["qiniu"] = reflect.TypeOf(qiniu.QiniuConfig{})
}

//func GetS3(key, region, bucket string) StorageInterface {
//	inst, ok := storages[key]
//	if ok {
//		return inst
//	}
//	s3inst := s3.New(region, bucket)
//	storages[key] = s3inst
//	return s3inst
//}
//
//func GetS3WithCredentials(key, region, bucket, keyID, accessKey string) StorageInterface {
//	inst, ok := storages[key]
//	if ok {
//		return inst
//	}
//
//	s3inst := s3.NewWithAutherization(region, bucket, keyID, accessKey)
//	storages[key] = s3inst
//	return s3inst
//}

func GetQiniu(key, accesskey, secretkey, bucket string, zone string) StorageInterface {
	inst, ok := storages[key]
	if ok {
		return inst
	}
	qninst := qiniu.New(accesskey, secretkey, bucket, zone)
	storages[key] = qninst
	return qninst
}

func GetQiniuWithDomain(key, accesskey, secretkey, bucket, zone, domain string) StorageInterface {
	inst, ok := storages[key]
	if ok {
		return inst
	}
	qninst := qiniu.New(accesskey, secretkey, bucket, zone)
	qninst.Domain = domain
	storages[key] = qninst
	return qninst
}

func GenFeedbackFilename(userID int, deviceUniqueID, clientUniqueID string) string {
	return fmt.Sprintf("%d_%s_%s_%d.zip", userID, deviceUniqueID, clientUniqueID, time.Now().Unix())
}

func MustGetStorage(countryShort string) StorageInterface {
	countryShort = strings.ToUpper(countryShort)
	inst, ok := countryStorages[countryShort]
	if ok && inst != nil {
		return inst
	}

	inst, ok = countryStorages["DEFAULT"]
	if ok && inst != nil {
		return inst
	}

	return nil
}

func NewConfByType(confType string) StorageConfigInterface {
	v, ok := storageConfTypeMap[confType]
	if !ok {
		return nil
	}

	cfg := reflect.New(v).Interface().(StorageConfigInterface)
	return cfg
}

func GetStorage(countryShort string, useFor string) StorageInterface {
	countryShort = strings.ToUpper(countryShort)
	inst, ok := countryStorages[countryShort]
	if ok && inst != nil {
		return inst
	}

	return nil
}

func SetStorage(countryShort string, useFor string, inst StorageInterface) error {
	countryShort = strings.ToUpper(countryShort)
	if inst == nil {
		return errors.New("storage instance is nil")
	}
	if countryShort != "DEFAULT" || len(countryShort) != 2 {
		return errors.New("contry short name must be default or two letters")
	}

	countryStorages[countryShort] = inst

	return nil
}
