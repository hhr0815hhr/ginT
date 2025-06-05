package database

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] interface {
	GetOne(fields string, query interface{}, args ...interface{}) (*T, error)
	GetList(dest *[]T, fields string, query interface{}, args ...interface{}) error
	Add(entity *T) error
	AddAll(entities []*T) error
	Update(entity *T, where, updates map[string]interface{}) error
	Delete(entity *T, query interface{}, args ...interface{}) error
	WithTrx(trxHandle *gorm.DB) Repository[T]
	SelfUpdate(entity *T) error
	Exist(query interface{}, args ...interface{}) bool
	Upsert(entity *T) error
	Count(conditions ...QueryCondition) int64
	ForPage(fields string, page, limit int, order string, conditions ...QueryCondition) ([]T, int64, error)
}

type BaseRepository[T any] struct {
	Db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{Db: db}
}

func (r *BaseRepository[T]) WithTrx(trxHandle *gorm.DB) Repository[T] {
	if trxHandle == nil {
		return r
	}
	return &BaseRepository[T]{Db: trxHandle}
}

func (r *BaseRepository[T]) Exist(query interface{}, args ...interface{}) bool {
	var count int64
	if err := r.Db.Model(new(T)).Where(query, args...).Limit(1).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (r *BaseRepository[T]) GetOne(fields string, query interface{}, args ...interface{}) (*T, error) {
	var dest = new(T)
	builder := r.Db.Where(query, args...)
	var err error
	if fields == "" || fields == "*" {
		err = builder.First(dest).Error
	} else {
		err = builder.Select(fields).First(dest).Error
	}
	return dest, err
}

func (r *BaseRepository[T]) GetList(dest *[]T, fields string, query interface{}, args ...interface{}) error {
	builder := r.Db.Where(query, args...)
	if fields == "" || fields == "*" {
		return builder.Find(dest).Error
	}
	return builder.Select(fields).Find(dest).Error
}

func (r *BaseRepository[T]) GetList2(fields string, conditions ...QueryCondition) (*[]T, error) {
	var model T
	tx := r.Db.Model(&model)
	queryTx := buildQuery(tx, conditions...)
	var dest []T
	if fields == "" || fields == "*" {
		err := queryTx.Find(&dest).Error
		return &dest, err
	}
	err := queryTx.Find(&dest).Error
	return &dest, err
}

//func (r *BaseRepository[T]) ForPage(dest *[]T, fields string, page, limit int, query interface{}, args ...interface{}) (int64, error) {
//	var count int64
//	tx := r.Db.Model(new(T)) // 基于模型创建查询
//
//	// 构建查询条件
//	if query != nil {
//		tx = tx.Where(query, args...)
//	}
//
//	// 统计总数
//	if err := tx.Count(&count).Error; err != nil {
//		return 0, err
//	}
//
//	// 应用字段选择
//	if fields != "" {
//		tx = tx.Select(fields)
//	}
//	// 应用分页
//	offset := (page - 1) * limit
//	tx = tx.Offset(offset).Limit(limit)
//	// 执行查询
//	if err := tx.Find(dest).Error; err != nil {
//		return 0, err
//	}
//
//	return count, nil
//}

func (r *BaseRepository[T]) Add(entity *T) error {
	return r.Db.Create(entity).Error
}

func (r *BaseRepository[T]) AddAll(entities []*T) error {
	return r.Db.CreateInBatches(entities, 100).Error
}

func (r *BaseRepository[T]) SelfUpdate(entity *T) error {
	return r.Db.Save(entity).Error
}

func (r *BaseRepository[T]) UpdateWithConditions(updates map[string]interface{}, conditions ...QueryCondition) error {
	var model T
	tx := r.Db.Model(&model)
	queryTx := buildQuery(tx, conditions...)
	return queryTx.Updates(updates).Error
}

func (r *BaseRepository[T]) Update(entity *T, where, updates map[string]interface{}) error {
	return r.Db.Model(entity).Where(where).Updates(updates).Error
}

func (r *BaseRepository[T]) Delete(entity *T, query interface{}, args ...interface{}) error {
	return r.Db.Where(query, args...).Delete(entity).Error
}

func (r *BaseRepository[T]) Upsert(entity *T) error {
	return r.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(entity).Error
}

func (r *BaseRepository[T]) ForPage(fields string, page, limit int, order string, conditions ...QueryCondition) ([]T, int64, error) {
	var (
		results []T
		count   int64
		model   T
	)

	tx := r.Db.Model(&model)

	// 构建动态查询条件
	queryTx := buildQuery(tx, conditions...)

	// 1. 先执行Count查询 (在应用分页和字段选择之前)
	if err := queryTx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 如果总数为0，直接返回，避免无效的Find查询
	if count == 0 {
		return results, 0, nil
	}

	// 2. 再执行Find查询
	// 应用字段选择
	if fields != "" {
		queryTx = queryTx.Select(fields)
	}

	// 应用分页
	offset := (page - 1) * limit
	if err := queryTx.Offset(offset).Limit(limit).Order(order).Find(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, count, nil
}

func (r *BaseRepository[T]) Count(conditions ...QueryCondition) int64 {
	var model T
	tx := r.Db.Model(&model)
	queryTx := buildQuery(tx, conditions...)
	var cnt int64
	queryTx.Count(&cnt)
	return cnt
}

// buildQuery 是一个辅助函数，用于根据条件列表构建GORM查询
func buildQuery(tx *gorm.DB, conditions ...QueryCondition) *gorm.DB {
	for _, c := range conditions {
		if c.Field == "" || c.Operator == "" {
			continue
		}

		// 处理 OR 查询
		if c.Operator == "OR" {
			// OR条件的值应该是一个 []QueryCondition 切片
			if subConditions, ok := c.Value.([]QueryCondition); ok {
				orQuery := tx.Session(&gorm.Session{NewDB: true}) // 创建一个新的Session以构建OR子查询
				orQuery = buildQuery(orQuery, subConditions...)
				tx = tx.Or(orQuery)
			}
			continue
		}

		// 构建查询表达式
		queryStr := fmt.Sprintf("%s %s ?", c.Field, c.Operator)
		tx = tx.Where(queryStr, c.Value)
	}
	return tx
}

// QueryCondition 定义了单个查询条件
type QueryCondition struct {
	Field    string      // 字段名, e.g., "name"
	Operator string      // 操作符, e.g., "=", "LIKE", "IN", "BETWEEN", "OR"
	Value    interface{} // 对应的值
}
