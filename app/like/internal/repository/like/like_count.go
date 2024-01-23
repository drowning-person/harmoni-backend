package like

import (
	"context"
	"fmt"

	objectv1 "harmoni/api/common/object/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	polike "harmoni/app/like/internal/infrastructure/po/like"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

func byObject(object *objectv1.Object) data.ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("object_id = ?", object.GetId()).
			Where("object_type = ?", object.GetType())
	}
}

func genLikeKey(objectID int64, objectType objectv1.ObjectType) string {
	return fmt.Sprintf("like:%s:%d", objectType.Format(), objectID)
}

func (r *LikeRepo) ObjectLikeCount(ctx context.Context, object *objectv1.Object) (*entitylike.LikeCount, error) {
	lc := &polike.LikeCount{}
	err := r.data.DB(ctx).Model(lc).
		Scopes(byObject(object)).
		First(lc).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return &entitylike.LikeCount{
		Count:  lc.Counts,
		Object: object,
	}, nil
}

func (r *LikeRepo) ListObjectLikeCount(ctx context.Context, objectIDs []int64, objectType objectv1.ObjectType) (entitylike.LikeCountList, error) {
	lcList := make([]*polike.LikeCount, 0, 10)
	err := r.data.DB(ctx).Model(lcList).
		Where("object_id IN ?", objectIDs).
		Where("object_type = ?", objectType).
		Find(lcList).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return lo.Map(lcList, func(lc *polike.LikeCount, _ int) *entitylike.LikeCount {
		return &entitylike.LikeCount{
			Count: lc.Counts,
			Object: &objectv1.Object{
				Id:   lc.ObjectID,
				Type: lc.OjbectType,
			},
		}
	}), nil
}

func (r *LikeRepo) AddLikeCount(ctx context.Context, object *objectv1.Object, count int64) error {
	err := r.data.DB(ctx).Model(&polike.LikeCount{}).
		Scopes(byObject(object)).
		UpdateColumn("counts", gorm.Expr("counts + ?", count)).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
