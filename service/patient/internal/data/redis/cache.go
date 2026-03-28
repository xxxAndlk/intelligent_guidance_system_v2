// service/patient/internal/data/redis/cache.go
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	"intelligent-guidance-system/service/patient/internal/biz/dto"
)

// PatientCache 患者缓存
type PatientCache struct {
	client *redis.Client
	log    *log.Helper
}

const (
	PatientCacheKeyPrefix = "patient:"
	PatientCacheTTL       = 30 * time.Minute
)

// NewPatientCache 创建患者缓存
func NewPatientCache(client *redis.Client, logger log.Logger) *PatientCache {
	return &PatientCache{
		client: client,
		log:    log.NewHelper(logger),
	}
}

// Set 设置患者缓存
func (c *PatientCache) Set(ctx context.Context, patientID int64, patient *dto.PatientDetailDTO) error {
	key := c.getKey(patientID)
	data, err := json.Marshal(patient)
	if err != nil {
		c.log.Errorf("序列化患者数据失败: %v", err)
		return err
	}

	result := c.client.Set(ctx, key, data, PatientCacheTTL)
	if result.Err() != nil {
		c.log.Errorf("设置患者缓存失败: %v", result.Err())
		return result.Err()
	}

	c.log.Infof("设置患者缓存成功: key=%s", key)
	return nil
}

// Get 获取患者缓存
func (c *PatientCache) Get(ctx context.Context, patientID int64) (*dto.PatientDetailDTO, error) {
	key := c.getKey(patientID)
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		c.log.Errorf("获取患者缓存失败: %v", err)
		return nil, err
	}

	var patient dto.PatientDetailDTO
	if err := json.Unmarshal([]byte(result), &patient); err != nil {
		c.log.Errorf("反序列化患者数据失败: %v", err)
		return nil, err
	}

	return &patient, nil
}

// Delete 删除患者缓存
func (c *PatientCache) Delete(ctx context.Context, patientID int64) error {
	key := c.getKey(patientID)
	result := c.client.Del(ctx, key)
	if result.Err() != nil {
		c.log.Errorf("删除患者缓存失败: %v", result.Err())
		return result.Err()
	}

	return nil
}

// SetByPhone 设置手机号索引缓存
func (c *PatientCache) SetByPhone(ctx context.Context, phone string, patientID int64) error {
	key := c.getPhoneKey(phone)
	result := c.client.Set(ctx, key, patientID, PatientCacheTTL)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

// GetByPhone 获取手机号索引缓存
func (c *PatientCache) GetByPhone(ctx context.Context, phone string) (int64, error) {
	key := c.getPhoneKey(phone)
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}

	var patientID int64
	if err := json.Unmarshal([]byte(result), &patientID); err != nil {
		return 0, err
	}

	return patientID, nil
}

// getKey 获取缓存key
func (c *PatientCache) getKey(patientID int64) string {
	return PatientCacheKeyPrefix + "id:" + string(patientID)
}

// getPhoneKey 获取手机号缓存key
func (c *PatientCache) getPhoneKey(phone string) string {
	return PatientCacheKeyPrefix + "phone:" + phone
}