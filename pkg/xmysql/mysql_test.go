package xmysql_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmysql"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type UserModel struct {
	Id        uint64 `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func TestNew(t *testing.T) {
	c := xconfig.New(
		xconfig.WithSource(apollo.NewSource("http://apollo.dev.jiangyang.me", "go-kit", "default", "application", os.Getenv("APOLLO_ACCESS_KEY_SECRET_GO_KIT"))),
	)
	confStr := c.Get("mysql_conf")
	logger := xlog.New(xlog.WithTrace("name"))
	conn := xmysql.New(confStr, logger)
	db, err := conn.DB()
	if err != nil {
		t.Error(err)
	}
	if err := db.Ping(); err != nil {
		t.Error(err)
	}
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.MD{"name": {"xxxx-xxxx-xxxx-1234"}})

	if err := conn.WithContext(ctx).AutoMigrate(UserModel{}); err != nil {
		t.Error(err)
	}

	user := UserModel{}
	if err := conn.WithContext(ctx).FirstOrCreate(&user).Error; err != nil {
		t.Error(err)
	}
	logger.Info(ctx, user)
}
