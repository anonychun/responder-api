package lark_message

import (
	"context"
	"uuid"

	"github.com/anonychun/responder-api/internal/bootstrap"
	dbSql "github.com/anonychun/responder-api/internal/db/sql"
	"github.com/anonychun/responder-api/internal/entity"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type IRepository interface {
	ExistsByMessageId(ctx context.Context, messageId string) (bool, error)
	Create(ctx context.Context, larkMessage *entity.LarkMessage) error
	FindByMessageId(ctx context.Context, messageId string) (*entity.LarkMessage, error)
	MarkReplied(ctx context.Context, id uuid.UUID) error
}

type Repository struct {
	sqlDB dbSql.IDB
}

var _ IRepository = (*Repository)(nil)

func NewRepository(i do.Injector) (*Repository, error) {
	return &Repository{
		sqlDB: do.MustInvoke[*dbSql.PostgresDB](i),
	}, nil
}

func (r *Repository) ExistsByMessageId(ctx context.Context, messageId string) (bool, error) {
	return r.sqlDB.DB(ctx).NewSelect().Model(&entity.LarkMessage{}).Where("message_id = ?", messageId).Exists(ctx)
}

func (r *Repository) Create(ctx context.Context, larkMessage *entity.LarkMessage) error {
	_, err := r.sqlDB.DB(ctx).NewInsert().Model(larkMessage).Exec(ctx)
	return err
}

func (r *Repository) FindByMessageId(ctx context.Context, messageId string) (*entity.LarkMessage, error) {
	larkMessage := &entity.LarkMessage{}
	err := r.sqlDB.DB(ctx).NewSelect().Model(larkMessage).Where("message_id = ?", messageId).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return larkMessage, nil
}

func (r *Repository) MarkReplied(ctx context.Context, id uuid.UUID) error {
	_, err := r.sqlDB.DB(ctx).NewUpdate().
		Model(&entity.LarkMessage{}).
		Set("replied_at = now()").
		Set("updated_at = now()").
		Where("id = ? AND replied_at IS NULL", id).
		Exec(ctx)

	return err
}
