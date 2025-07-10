package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
)

type InstallEventRepository interface {
	Create(ctx context.Context, event *model.InstallEvent) error
	CreateBatch(ctx context.Context, events []*model.InstallEvent) error
}

type installEventRepository struct {
	ch clickhouse.Conn
}

func NewInstallEventRepository(ch clickhouse.Conn) InstallEventRepository {
	return &installEventRepository{ch: ch}
}

// 单条插入（内部调用批量插入）
func (r *installEventRepository) Create(ctx context.Context, event *model.InstallEvent) error {
	return r.CreateBatch(ctx, []*model.InstallEvent{event})
}

// 批量插入
func (r *installEventRepository) CreateBatch(ctx context.Context, events []*model.InstallEvent) error {
	if len(events) == 0 {
		return nil
	}

	batch, err := r.ch.PrepareBatch(ctx, `
		INSERT INTO install_events (
			app_id, app_name, app_version, app_type,
			event_id, event_date, event_time,
			device_id, channel_id, install_ip,
			install_type, install_result,
			os_language, os_timezone, os_name, os_version, os_build, os_family,
			signature_status, signature_version, signature_params
		)
	`)
	if err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to prepare batch", err)
	}

	for _, event := range events {
		err = batch.Append(
			event.AppID, event.AppName, event.AppVersion, uint8(event.AppType),
			event.EventID, event.EventDate, event.EventTime,
			event.DeviceID, event.ChannelID, event.InstallIP,
			uint8(event.InstallType), uint8(event.InstallResult),
			event.OSLanguage, event.OSTimezone, event.OSName, event.OSVersion, event.OSBuild, event.OSFamily,
			event.SignatureStatus, event.SignatureVersion, event.SignatureParams,
		)
		if err != nil {
			return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to append event to batch", err)
		}
	}

	if err = batch.Send(); err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to send batch", err)
	}

	return nil
}