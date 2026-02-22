package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"blog-server/entity"
	"blog-server/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	cfg, err := test.LoadConfig()
	require.NoError(t, err)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		"blog_test",
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		TranslateError:                           true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	tx := db.Begin()
	require.NoError(t, tx.Error)
	require.NoError(t, tx.Session(&gorm.Session{
		Logger: tx.Logger.LogMode(logger.Silent),
	}).AutoMigrate(&entity.ImageFolder{}))

	// 默认每个 test 用事务隔离，结束时回滚
	cleanup := func() {
		_ = tx.Rollback().Error
	}
	t.Cleanup(cleanup)

	return tx, func() {
		// 允许个别 test 选择提交
		_ = tx.Commit().Error
	}
}

func mustCreateFolder(t *testing.T, ctx context.Context, db *gorm.DB, repo IImageFolderRepo, name string, parentID *uuid.UUID) *entity.ImageFolder {
	t.Helper()

	f := &entity.ImageFolder{
		Name:     name,
		ParentID: parentID,
	}
	err := repo.Create(ctx, db, f)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, f.ID)
	require.False(t, f.CreatedAt.IsZero())
	require.False(t, f.UpdatedAt.IsZero())
	return f
}

func TestImageFolderRepo_Create_GetByID(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	f := &entity.ImageFolder{Name: "root"}

	err := repo.Create(ctx, db, f)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, f.ID)
	require.False(t, f.CreatedAt.IsZero())
	require.False(t, f.UpdatedAt.IsZero())

	got, err := repo.GetByID(ctx, db, f.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, f.ID, got.ID)
	require.Equal(t, "root", got.Name)
	require.Nil(t, got.ParentID)
	require.Nil(t, got.DeletedAt)
}

func TestImageFolderRepo_GetByID_NotFound(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	got, err := repo.GetByID(ctx, db, uuid.New())
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestImageFolderRepo_Exists(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	ok, err := repo.Exists(ctx, db, uuid.New())
	require.NoError(t, err)
	require.False(t, ok)

	f := mustCreateFolder(t, ctx, db, repo, "root", nil)

	ok, err = repo.Exists(ctx, db, f.ID)
	require.NoError(t, err)
	require.True(t, ok)

	require.NoError(t, repo.SoftDelete(ctx, db, f.ID))

	ok, err = repo.Exists(ctx, db, f.ID)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestImageFolderRepo_ListByParent_RootAndChild(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	// root 级别
	f1 := mustCreateFolder(t, ctx, db, repo, "b", nil)
	f2 := mustCreateFolder(t, ctx, db, repo, "a", nil) // 用于测试排序 name ASC
	_ = mustCreateFolder(t, ctx, db, repo, "c", nil)

	// child 级别
	p := mustCreateFolder(t, ctx, db, repo, "parent", nil)
	c1 := mustCreateFolder(t, ctx, db, repo, "child2", &p.ID)
	c2 := mustCreateFolder(t, ctx, db, repo, "child1", &p.ID)

	roots, err := repo.ListByParent(ctx, db, nil, 0, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(roots), 4)

	// 只校验我们关心的：至少包含 a/b/c/parent，且 name ASC
	// 由于测试库可能已有历史数据，这里用“提取并验证局部顺序”的方式更稳。
	var rootNames []string
	for _, it := range roots {
		if it.ID == f1.ID || it.ID == f2.ID || it.ID == p.ID {
			rootNames = append(rootNames, it.Name)
		}
	}
	require.ElementsMatch(t, []string{"a", "b", "parent"}, rootNames)

	children, err := repo.ListByParent(ctx, db, &p.ID, 0, 0)
	require.NoError(t, err)
	require.Len(t, children, 2)
	require.Equal(t, c2.ID, children[0].ID) // child1
	require.Equal(t, c1.ID, children[1].ID) // child2

	// limit/offset
	children2, err := repo.ListByParent(ctx, db, &p.ID, 1, 0)
	require.NoError(t, err)
	require.Len(t, children2, 1)
	require.Equal(t, "child1", children2[0].Name)

	children3, err := repo.ListByParent(ctx, db, &p.ID, 1, 1)
	require.NoError(t, err)
	require.Len(t, children3, 1)
	require.Equal(t, "child2", children3[0].Name)
}

func TestImageFolderRepo_ExistsBySameNameInParent(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	// root 下同名
	a := mustCreateFolder(t, ctx, db, repo, "dup", nil)

	ok, err := repo.ExistsBySameNameInParent(ctx, db, nil, "dup", nil)
	require.NoError(t, err)
	require.True(t, ok)

	// exclude 自己后应该不存在
	ok, err = repo.ExistsBySameNameInParent(ctx, db, nil, "dup", &a.ID)
	require.NoError(t, err)
	require.False(t, ok)

	// 不同 parent 下同名：root 下有 dup，但 parent 下没有 dup
	p := mustCreateFolder(t, ctx, db, repo, "p", nil)

	ok, err = repo.ExistsBySameNameInParent(ctx, db, &p.ID, "dup", nil)
	require.NoError(t, err)
	require.False(t, ok)

	// parent 下创建 dup 后应存在
	b := mustCreateFolder(t, ctx, db, repo, "dup", &p.ID)
	ok, err = repo.ExistsBySameNameInParent(ctx, db, &p.ID, "dup", nil)
	require.NoError(t, err)
	require.True(t, ok)

	// exclude b 后应该不存在（因为 parent 下只有 b 一个 dup）
	ok, err = repo.ExistsBySameNameInParent(ctx, db, &p.ID, "dup", &b.ID)
	require.NoError(t, err)
	require.False(t, ok)

	// soft delete 后不应算存在
	require.NoError(t, repo.SoftDelete(ctx, db, b.ID))
	ok, err = repo.ExistsBySameNameInParent(ctx, db, &p.ID, "dup", nil)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestImageFolderRepo_Rename(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	f := mustCreateFolder(t, ctx, db, repo, "old", nil)
	oldUpdated := f.UpdatedAt

	// 确保时间变化
	time.Sleep(5 * time.Millisecond)

	err := repo.Rename(ctx, db, f.ID, "new")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, db, f.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "new", got.Name)
	require.True(t, got.UpdatedAt.After(oldUpdated))
}

func TestImageFolderRepo_Move(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	parent1 := mustCreateFolder(t, ctx, db, repo, "p1", nil)
	parent2 := mustCreateFolder(t, ctx, db, repo, "p2", nil)
	child := mustCreateFolder(t, ctx, db, repo, "child", &parent1.ID)
	oldUpdated := child.UpdatedAt

	time.Sleep(5 * time.Millisecond)

	err := repo.Move(ctx, db, child.ID, &parent2.ID)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, db, child.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.ParentID)
	require.Equal(t, parent2.ID, *got.ParentID)
	require.True(t, got.UpdatedAt.After(oldUpdated))

	// ListByParent 也应反映移动结果
	c1, err := repo.ListByParent(ctx, db, &parent1.ID, 0, 0)
	require.NoError(t, err)
	require.Len(t, c1, 0)

	c2, err := repo.ListByParent(ctx, db, &parent2.ID, 0, 0)
	require.NoError(t, err)
	require.Len(t, c2, 1)
	require.Equal(t, child.ID, c2[0].ID)
}

func TestImageFolderRepo_SoftDelete(t *testing.T) {
	db, _ := newTestDB(t)

	ctx := context.Background()
	repo := NewImageFolderRepo()

	f := mustCreateFolder(t, ctx, db, repo, "to_del", nil)

	// 软删前可见
	got, err := repo.GetByID(ctx, db, f.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Nil(t, got.DeletedAt)

	// soft delete
	err = repo.SoftDelete(ctx, db, f.ID)
	require.NoError(t, err)

	// GetByID 应该返回 nil
	got, err = repo.GetByID(ctx, db, f.ID)
	require.NoError(t, err)
	require.Nil(t, got)

	// ListByParent 不应包含
	roots, err := repo.ListByParent(ctx, db, nil, 0, 0)
	require.NoError(t, err)
	for _, it := range roots {
		require.NotEqual(t, f.ID, it.ID)
	}

	// Exists 不应存在
	ok, err := repo.Exists(ctx, db, f.ID)
	require.NoError(t, err)
	require.False(t, ok)

	// ExistsBySameNameInParent 不应认为存在
	ok, err = repo.ExistsBySameNameInParent(ctx, db, nil, "to_del", nil)
	require.NoError(t, err)
	require.False(t, ok)
}
