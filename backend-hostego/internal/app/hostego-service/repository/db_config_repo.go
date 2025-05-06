package repository

// import (
// 	"backend-hostego/internal/app/hostego-service/models"
// )

// type DbConfigRepo struct {
// 	*BaseRepo
// }

// func NewDbConfigRepo(baseRepo *BaseRepo) *DbConfigRepo {
// 	return &DbConfigRepo{BaseRepo: baseRepo}
// }

// func (dbConfigRepo *DbConfigRepo) GetConfigByKey(key string) (models.DbConfig, error) {
// 	config := models.DbConfig{Key: key}

// 	if tx := dbConfigRepo.db.Where(&models.DbConfig{Key: key}).First(&config); tx.Error != nil {
// 		return config, tx.Error
// 	}

// 	return config, nil
// }
