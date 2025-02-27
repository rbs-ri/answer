package tag

import (
	"context"

	"github.com/answerdev/answer/internal/base/data"
	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/entity"
	tagcommon "github.com/answerdev/answer/internal/service/tag_common"
	"github.com/answerdev/answer/internal/service/unique"
	"github.com/answerdev/answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"
)

// tagRelRepo tag rel repository
type tagRelRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewTagRelRepo new repository
func NewTagRelRepo(data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo) tagcommon.TagRelRepo {
	return &tagRelRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddTagRelList add tag list
func (tr *tagRelRepo) AddTagRelList(ctx context.Context, tagList []*entity.TagRel) (err error) {
	for _, item := range tagList {
		item.ObjectID = uid.DeShortID(item.ObjectID)
	}
	_, err = tr.data.DB.Context(ctx).Insert(tagList)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range tagList {
		item.ObjectID = uid.EnShortID(item.ObjectID)
	}
	return
}

// RemoveTagRelListByObjectID delete tag list
func (tr *tagRelRepo) RemoveTagRelListByObjectID(ctx context.Context, objectID string) (err error) {
	objectID = uid.DeShortID(objectID)
	_, err = tr.data.DB.Context(ctx).Where("object_id = ?", objectID).Update(&entity.TagRel{Status: entity.TagRelStatusDeleted})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (tr *tagRelRepo) HideTagRelListByObjectID(ctx context.Context, objectID string) (err error) {
	objectID = uid.DeShortID(objectID)
	_, err = tr.data.DB.Where("object_id = ?", objectID).Cols("status").Update(&entity.TagRel{Status: entity.TagRelStatusHide})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (tr *tagRelRepo) ShowTagRelListByObjectID(ctx context.Context, objectID string) (err error) {
	objectID = uid.DeShortID(objectID)
	_, err = tr.data.DB.Where("object_id = ?", objectID).Cols("status").Update(&entity.TagRel{Status: entity.TagRelStatusAvailable})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// RemoveTagRelListByIDs delete tag list
func (tr *tagRelRepo) RemoveTagRelListByIDs(ctx context.Context, ids []int64) (err error) {
	_, err = tr.data.DB.Context(ctx).In("id", ids).Update(&entity.TagRel{Status: entity.TagRelStatusDeleted})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetObjectTagRelWithoutStatus get object tag relation no matter status
func (tr *tagRelRepo) GetObjectTagRelWithoutStatus(ctx context.Context, objectID, tagID string) (
	tagRel *entity.TagRel, exist bool, err error,
) {
	objectID = uid.DeShortID(objectID)
	tagRel = &entity.TagRel{}
	session := tr.data.DB.Context(ctx).Where("object_id = ?", objectID).And("tag_id = ?", tagID)
	exist, err = session.Get(tagRel)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	tagRel.ObjectID = uid.EnShortID(tagRel.ObjectID)
	return
}

// EnableTagRelByIDs update tag status to available
func (tr *tagRelRepo) EnableTagRelByIDs(ctx context.Context, ids []int64) (err error) {
	_, err = tr.data.DB.Context(ctx).In("id", ids).Update(&entity.TagRel{Status: entity.TagRelStatusAvailable})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetObjectTagRelList get object tag relation list all
func (tr *tagRelRepo) GetObjectTagRelList(ctx context.Context, objectID string) (tagListList []*entity.TagRel, err error) {
	objectID = uid.DeShortID(objectID)
	tagListList = make([]*entity.TagRel, 0)
	session := tr.data.DB.Context(ctx).Where("object_id = ?", objectID)
	session.In("status", []int{entity.TagRelStatusAvailable, entity.TagRelStatusHide})
	err = session.Find(&tagListList)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range tagListList {
		item.ObjectID = uid.EnShortID(item.ObjectID)
	}
	return
}

// BatchGetObjectTagRelList get object tag relation list all
func (tr *tagRelRepo) BatchGetObjectTagRelList(ctx context.Context, objectIds []string) (tagListList []*entity.TagRel, err error) {
	for num, item := range objectIds {
		objectIds[num] = uid.DeShortID(item)
	}
	tagListList = make([]*entity.TagRel, 0)
	session := tr.data.DB.Context(ctx).In("object_id", objectIds)
	session.Where("status = ?", entity.TagRelStatusAvailable)
	err = session.Find(&tagListList)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, item := range tagListList {
		item.ObjectID = uid.EnShortID(item.ObjectID)
	}
	return
}

// CountTagRelByTagID count tag relation
func (tr *tagRelRepo) CountTagRelByTagID(ctx context.Context, tagID string) (count int64, err error) {
	count, err = tr.data.DB.Context(ctx).Count(&entity.TagRel{TagID: tagID, Status: entity.AnswerStatusAvailable})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
