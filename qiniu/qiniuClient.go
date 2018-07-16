package qiniu

import (
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"gitlab.com/qiniuyun/constant"
	"fmt"
	"log"
	"context"
)

// UploadQiniu Upload file to Qiniu...
func UploadQiniu(localFile, bucket, key string) (storage.PutRet, error) {
	accessKey := constant.AccessKey
	secretKey := constant.SecretKey
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "banner",
		},
	}
	//err := formUploader.Put(context.Background(), &ret, upToken, key, localFile, 1024*1024, &putExtra)
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	log.Printf("key:%v hash%v", ret.Key, ret.Hash)
	return ret, nil
}

// GetQiniuFileList Get file from Qiniu...
func GetQiniuFileList(bucket string) ([]storage.ListItem, error) {
	qb := &qbox.Mac{
		AccessKey: constant.AccessKey,
		SecretKey: []byte(constant.SecretKey),
	}
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	r := storage.NewBucketManager(qb, &cfg)
	f, _, _, _, err := r.ListFiles(bucket, "", "", "", 1000)
	if err != nil {
		log.Fatalf("获取空间文件列表失败", err)
		return f, err
	}
	return f, nil
}

// DeleteOneFileList Delete file from Qiniu...
func DeleteOneFileList(bucket, key string) error {
	qb := &qbox.Mac{
		AccessKey: constant.AccessKey,
		SecretKey: []byte(constant.SecretKey),
	}
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	r := storage.NewBucketManager(qb, &cfg)
	err := r.Delete(bucket, key)
	if err != nil {
		log.Printf("删除文件失败", err)
		return err
	}
	return nil
}
