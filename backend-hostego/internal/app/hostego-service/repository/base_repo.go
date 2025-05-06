package repository

import (
	"backend-hostego/internal/app/hostego-service/constants"
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/pkg/db/postgres"
	"backend-hostego/internal/pkg/logger"
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log = logger.GetLogger()

const ORDER_BY = "created_at DESC"

type BaseRepo struct {
	Db *gorm.DB
}

func NewBaseRepo() *BaseRepo {
	return &BaseRepo{}
}

func (baseRepo *BaseRepo) NewRelicDbWrapper() (*gorm.DB, *newrelic.Transaction) {
	return baseRepo.Db, nil
}

// GetFirstRecord Attaching GetFirstRecord to BaseRepo via ptr
func (baseRepo *BaseRepo) GetFirstRecord(model interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	if tx := db.Table(tableName).First(model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to get entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

// GetRecords if a value is empty or null gorm does not include it in where clause, this should be handled in app
func (baseRepo *BaseRepo) GetRecords(model interface{}, filterModel interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).Find(model, filterModel)
	if tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to get records error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) GetRecordsByCondition(model interface{}, whereClause string, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).Where(whereClause).Find(model)
	if tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to get records error:" + tx.Error.Error()})
		return tx.Error
	}

	return nil
}

func (baseRepo *BaseRepo) GetRecordsByConditionAndOrder(result interface{}, whereClause string, orderVal string, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()
	err := db.Table(tableName).
		Where(whereClause).
		Order(orderVal).
		Find(result).Error

	if err != nil {
		log.Errorf("Encountered error while fetching data from DB: %v", err)
		return err
	}
	return nil
}

func (baseRepo *BaseRepo) Update(rctx dto.ReqCtx, model interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()
	log := rctx.GoBricksLog

	if tx := db.Table(tableName).Updates(model); tx.Error != nil {
		log.Errorf("failed to update entity error:" + tx.Error.Error())
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) UpdateById(model interface{}, tableName string, userID int64) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	if tx := db.Table(tableName).Where("user_id = ?", userID).Updates(model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to update entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) UpdateByWhereCondition(model interface{}, whereClause string, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	if tx := db.Table(tableName).Where(whereClause).Updates(model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to update entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) Create(rctx dto.ReqCtx, model interface{}, tableName string) error {
	log := rctx.GoBricksLog
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	if tx := db.Table(tableName).Create(model); tx.Error != nil {
		log.Errorf("Failed to Create Entry: %v, model: %v", tx.Error, model)
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) GetRecordsWithFilterAndLimitAndOrderBy(model interface{}, filterModel interface{}, limit int, orderBy string, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).
		Where("updated_by <> ? AND created_by <> ?", constants.AUTOMATION, constants.AUTOMATION).
		Order(orderBy).
		Limit(limit).
		Find(model, filterModel)

	if tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to get records error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) Delete(model interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	//Passing nil here as we want to delete based on condition filter not value
	if tx := db.Table(tableName).Unscoped().Delete(nil, model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to delete entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) GetTransaction() *gorm.DB {
	tx := baseRepo.Db.Begin()
	return tx
}

func (baseRepo *BaseRepo) Rollback(rCtx dto.ReqCtx, tx *gorm.DB) {
	if tx != nil {
		tx.Rollback()
	}
}

func (baseRepo *BaseRepo) GetRecordsWithTxnAndWithSkipLock(rCtx dto.ReqCtx, tx interface{}, model interface{}, filterModel interface{}) error {
	db := tx.(*gorm.DB)

	result := db.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Find(model, filterModel)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (baseRepo *BaseRepo) SaveWithTxn(rCtx dto.ReqCtx, tx *gorm.DB, value interface{}) (err error) {

	tx.Save(value)
	commitError := tx.Commit().Error
	if commitError != nil {
		log.Errorf("Commit error: %v", commitError)
		return commitError
	}
	return nil
}

func (baseRepo *BaseRepo) Save(model interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()
	log.Info("In Save Operation")
	if tx := db.Table(tableName).Save(model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to update entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (baseRepo *BaseRepo) BulkSave(rCtx dto.ReqCtx, values interface{}, tableName string) (err error) {
	// Begin a database transaction
	txn := baseRepo.GetTransaction()
	defer baseRepo.Rollback(rCtx, txn)
	return baseRepo.BulkSaveWithTxn(rCtx, txn, values, tableName)
}

func (baseRepo *BaseRepo) BulkSaveWithTxn(rCtx dto.ReqCtx, tx *gorm.DB, values interface{}, tableName string) (err error) {
	for _, value := range values.([]interface{}) {
		tx.Table(tableName).Save(value)
	}

	commitError := tx.Commit().Error
	if commitError != nil {
		log.Errorf("Commit error: %v", commitError)
		return commitError
	}
	return nil
}

func (baseRepo *BaseRepo) FirstOrCreate(model interface{}, tableName string) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	if tx := db.Table(tableName).FirstOrCreate(model); tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to create entity error:" + tx.Error.Error()})
		return tx.Error
	}
	return nil
}
func (baseRepo *BaseRepo) SetDb(db *gorm.DB) {
	baseRepo.Db = db
}

func (baseRepo *BaseRepo) CustomQuery(rCtx dto.ReqCtx, selectColumn string, tableName string, whereClause string, groupBy string, model interface{}) error {
	log := rCtx.GoBricksLog
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).
		Select(selectColumn).
		Where(whereClause).
		Group(groupBy).
		Find(model)
	if tx.Error != nil {
		log.Errorf("failed to execute query: ", tx.Error)
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (baseRepo *BaseRepo) UpdateWithCondition(rCtx dto.ReqCtx, updateColumn string, tableName string, whereClause string, value string) error {
	log := rCtx.GoBricksLog
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).
		Where(whereClause).
		Update(updateColumn, value)
	if tx.Error != nil {
		log.Errorf("failed to execute query: ", tx.Error)
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (baseRepo *BaseRepo) GetCountsByQuery(query string, args []interface{}, result interface{}) error {
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Raw(query, args...).Scan(result)
	if tx.Error != nil {
		log.Error(dto.LoggerWithFields{Message: "failed to execute query: " + tx.Error.Error()})
		return tx.Error
	}
	return nil
}

func (repo *BaseRepo) dbInstance(ctx context.Context, tx interface{}) (db *gorm.DB) {
	session := postgres.GetDB(ctx)
	if tx != nil {
		db = tx.(*gorm.DB)
	} else {
		db = session
	}
	return
}

func (repo *BaseRepo) RunTransaction(rCtx *dto.ReqCtx, callback func(tx *gorm.DB) error) (err error) {
	session := postgres.GetDB(rCtx.Context)
	txn := session.Begin()
	err = callback(txn)
	log := rCtx.GoBricksLog
	if err != nil {
		log.Errorf("Rolling back the transaction: %v", err)
		rollbackError := txn.Rollback().Error
		if rollbackError != nil {
			log.Errorf("Rollback error: %v", rollbackError)
			return rollbackError
		}
		return err
	} else {
		commitError := txn.Commit().Error
		if commitError != nil {
			log.Errorf("Commit error: %v", commitError)
			return commitError
		}
		return nil
	}
}

func (baseRepo *BaseRepo) SaveWithTxnUncommited(rCtx *dto.ReqCtx, tx *gorm.DB, tableName string, whereClause string, value interface{}) (err error) {
	err = tx.Table(tableName).Where(whereClause).Save(value).Error
	return err
}

func (baseRepo *BaseRepo) SaveModelWithTxnUncommited(rCtx *dto.ReqCtx, tx *gorm.DB, tableName string, value interface{}) (err error) {
	err = tx.Table(tableName).Create(value).Error
	return err
}

func (baseRepo *BaseRepo) UpdateModelWithTxnUncommited(rCtx *dto.ReqCtx, tx *gorm.DB, tableName string, value interface{}) (err error) {
	err = tx.Table(tableName).Updates(value).Error
	return err
}

func (baseRepo *BaseRepo) UpdateModelWithCondition(rCtx dto.ReqCtx, tableName string, whereClause string, model interface{}) error {
	log := rCtx.GoBricksLog
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).
		Where(whereClause).
		Updates(model)
	if tx.Error != nil {
		log.Errorf("failed to execute query: ", tx.Error)
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (baseRepo *BaseRepo) GetRecordsByConditionWithError(rctx *dto.ReqCtx, model interface{}, whereClause string, tableName string) error {
	log := rctx.GoBricksLog
	db, newRelicTxn := baseRepo.NewRelicDbWrapper()
	defer newRelicTxn.End()

	tx := db.Table(tableName).Where(whereClause).Find(model)
	if tx.Error != nil {
		log.Errorf("Transaction Failed for table: %v, whereClause : %v", tableName, whereClause)
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
