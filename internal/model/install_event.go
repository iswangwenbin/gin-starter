package model

import (
	"time"

	"gorm.io/gorm"
)

type AppType uint8

const (
	Windows AppType = 1
	MacOS   AppType = 2
	IOS     AppType = 3
	Android AppType = 4
)

func (at AppType) String() string {
	switch at {
	case Windows:
		return "Windows"
	case MacOS:
		return "MacOS"
	case IOS:
		return "iOS"
	case Android:
		return "Android"
	default:
		return "Unknown"
	}
}

type InstallType uint8

const (
	FirstInstall  InstallType = 1
	RepeatInstall InstallType = 2
)

func (it InstallType) String() string {
	switch it {
	case FirstInstall:
		return "FirstInstall"
	case RepeatInstall:
		return "RepeatInstall"
	default:
		return "Unknown"
	}
}

type InstallResult uint8

const (
	InstallSuccess InstallResult = 1
	InstallFail    InstallResult = 0
)

func (ir InstallResult) String() string {
	switch ir {
	case InstallSuccess:
		return "Success"
	case InstallFail:
		return "Fail"
	default:
		return "Unknown"
	}
}

type InstallEvent struct {
	BaseModel
	AppID            string            `json:"app_id" gorm:"column:app_id;type:varchar(36);not null;index" validate:"required"`
	AppName          string            `json:"app_name" gorm:"column:app_name;type:varchar(100);not null" validate:"required,max=100"`
	AppVersion       string            `json:"app_version" gorm:"column:app_version;type:varchar(50);not null" validate:"required"`
	AppType          AppType           `json:"app_type" gorm:"column:app_type;type:tinyint;not null;index" validate:"required,min=1,max=4"`
	EventID          string            `json:"event_id" gorm:"column:event_id;type:varchar(36);not null;uniqueIndex" validate:"required"`
	EventDate        time.Time         `json:"event_date" gorm:"column:event_date;type:date;not null;index" validate:"required"`
	EventTime        time.Time         `json:"event_time" gorm:"column:event_time;type:datetime;not null;index" validate:"required"`
	DeviceID         string            `json:"device_id" gorm:"column:device_id;type:varchar(100);not null;index" validate:"required"`
	ChannelID        string            `json:"channel_id" gorm:"column:channel_id;type:varchar(50);not null;index" validate:"required"`
	InstallIP        string            `json:"install_ip" gorm:"column:install_ip;type:varchar(45);not null" validate:"required,ip"`
	InstallType      InstallType       `json:"install_type" gorm:"column:install_type;type:tinyint;not null" validate:"required,min=1,max=2"`
	InstallResult    InstallResult     `json:"install_result" gorm:"column:install_result;type:tinyint;not null;index" validate:"required,min=0,max=1"`
	OSLanguage       string            `json:"os_language" gorm:"column:os_language;type:varchar(10);not null" validate:"required"`
	OSTimezone       string            `json:"os_timezone" gorm:"column:os_timezone;type:varchar(50);not null" validate:"required"`
	OSName           string            `json:"os_name" gorm:"column:os_name;type:varchar(50);not null" validate:"required"`
	OSVersion        string            `json:"os_version" gorm:"column:os_version;type:varchar(50);not null" validate:"required"`
	OSBuild          string            `json:"os_build" gorm:"column:os_build;type:varchar(50);not null" validate:"required"`
	OSFamily         string            `json:"os_family" gorm:"column:os_family;type:varchar(50);not null" validate:"required"`
	SignatureStatus  uint8             `json:"signature_status" gorm:"column:signature_status;type:tinyint;not null"`
	SignatureVersion string            `json:"signature_version" gorm:"column:signature_version;type:varchar(50);not null"`
	SignatureParams  map[string]string `json:"signature_params" gorm:"column:signature_params;type:json"`
}

func (InstallEvent) TableName() string {
	return "install_events"
}

func (e *InstallEvent) BeforeCreate(tx *gorm.DB) error {
	if e.EventDate.IsZero() {
		e.EventDate = e.EventTime.Truncate(24 * time.Hour)
	}
	return nil
}

func (e *InstallEvent) IsSuccess() bool {
	return e.InstallResult == InstallSuccess
}

func (e *InstallEvent) GetEventDate() time.Time {
	return e.EventTime.Truncate(24 * time.Hour)
}

// 请求和响应结构体
type CreateInstallEventRequest struct {
	AppID            string            `json:"app_id" validate:"required"`
	AppName          string            `json:"app_name" validate:"required,max=100"`
	AppVersion       string            `json:"app_version" validate:"required"`
	AppType          AppType           `json:"app_type" validate:"required,min=1,max=4"`
	EventID          string            `json:"event_id" validate:"required"`
	EventTime        time.Time         `json:"event_time" validate:"required"`
	DeviceID         string            `json:"device_id" validate:"required"`
	ChannelID        string            `json:"channel_id" validate:"required"`
	InstallIP        string            `json:"install_ip" validate:"required,ip"`
	InstallType      InstallType       `json:"install_type" validate:"required,min=1,max=2"`
	InstallResult    InstallResult     `json:"install_result" validate:"required,min=0,max=1"`
	OSLanguage       string            `json:"os_language" validate:"required"`
	OSTimezone       string            `json:"os_timezone" validate:"required"`
	OSName           string            `json:"os_name" validate:"required"`
	OSVersion        string            `json:"os_version" validate:"required"`
	OSBuild          string            `json:"os_build" validate:"required"`
	OSFamily         string            `json:"os_family" validate:"required"`
	SignatureStatus  uint8             `json:"signature_status"`
	SignatureVersion string            `json:"signature_version"`
	SignatureParams  map[string]string `json:"signature_params"`
}

type InstallEventListRequest struct {
	PageRequest
	AppID         string         `form:"app_id,omitempty"`
	AppType       *AppType       `form:"app_type,omitempty"`
	DeviceID      string         `form:"device_id,omitempty"`
	ChannelID     string         `form:"channel_id,omitempty"`
	InstallType   *InstallType   `form:"install_type,omitempty"`
	InstallResult *InstallResult `form:"install_result,omitempty"`
	StartTime     *time.Time     `form:"start_time,omitempty"`
	EndTime       *time.Time     `form:"end_time,omitempty"`
}

type InstallStatsRequest struct {
	AppID     string     `form:"app_id,omitempty"`
	AppType   *AppType   `form:"app_type,omitempty"`
	ChannelID string     `form:"channel_id,omitempty"`
	StartTime *time.Time `form:"start_time,omitempty"`
	EndTime   *time.Time `form:"end_time,omitempty"`
}

type InstallStatsResponse struct {
	TotalEvents      int64                `json:"total_events"`
	SuccessEvents    int64                `json:"success_events"`
	FailedEvents     int64                `json:"failed_events"`
	SuccessRate      float64              `json:"success_rate"`
	FirstInstalls    int64                `json:"first_installs"`
	RepeatInstalls   int64                `json:"repeat_installs"`
	UniqueDevices    int64                `json:"unique_devices"`
	TopChannels      []ChannelInstallStats `json:"top_channels"`
	AppTypeBreakdown []AppTypeInstallStats `json:"app_type_breakdown"`
}

type ChannelInstallStats struct {
	ChannelID string `json:"channel_id"`
	Count     int64  `json:"count"`
}

type AppTypeInstallStats struct {
	AppType AppType `json:"app_type"`
	Count   int64   `json:"count"`
}