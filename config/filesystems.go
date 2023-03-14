package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config
	config.Add("filesystems", map[string]any{
		// Default Filesystem Disk
		//
		// Here you may specify the default filesystem disk that should be used
		// by the framework. The "local" disk, as well as a variety of cloud
		// based disks are available to your application. Just store away!
		"default": config.Env("FILESYSTEM_DISK", "local"),

		// Filesystem Disks
		//
		// Here you may configure as many filesystem "disks" as you wish, and you
		// may even configure multiple disks of the same driver. Defaults have
		// been set up for each driver as an example of the required values.
		//
		// Supported Drivers: "local", "s3", "oss", "cos", "minio", "custom"
		"disks": map[string]any{
			"local": map[string]any{
				"driver": "local",
				"root":   "storage/app",
				"url":    config.Env("APP_URL", "").(string) + "/storage",
			},
			"s3": map[string]any{
				"driver": "s3",
				"key":    config.Env("AWS_ACCESS_KEY_ID"),
				"secret": config.Env("AWS_ACCESS_KEY_SECRET"),
				"region": config.Env("AWS_REGION"),
				"bucket": config.Env("AWS_BUCKET"),
				"url":    config.Env("AWS_URL"),
			},
			"oss": map[string]any{
				"driver":   "oss",
				"key":      config.Env("ALIYUN_ACCESS_KEY_ID"),
				"secret":   config.Env("ALIYUN_ACCESS_KEY_SECRET"),
				"bucket":   config.Env("ALIYUN_BUCKET"),
				"url":      config.Env("ALIYUN_URL"),
				"endpoint": config.Env("ALIYUN_ENDPOINT"),
			},
			"cos": map[string]any{
				"driver": "cos",
				"key":    config.Env("TENCENT_ACCESS_KEY_ID"),
				"secret": config.Env("TENCENT_ACCESS_KEY_SECRET"),
				"bucket": config.Env("TENCENT_BUCKET"),
				"url":    config.Env("TENCENT_URL"),
			},
			"minio": map[string]any{
				"driver":   "minio",
				"key":      config.Env("MINIO_ACCESS_KEY_ID"),
				"secret":   config.Env("MINIO_ACCESS_KEY_SECRET"),
				"region":   config.Env("MINIO_REGION"),
				"bucket":   config.Env("MINIO_BUCKET"),
				"url":      config.Env("MINIO_URL"),
				"endpoint": config.Env("MINIO_ENDPOINT"),
				"ssl":      config.Env("MINIO_SSL", false),
			},
		},
	})
}
