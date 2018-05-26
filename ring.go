package ring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	apiUri = "https://api.ring.com"
	//apiUri     = "http://localhost:8080"
	apiVersion = "9"
	hardWareId = "e6b664f0-606d-11e8-bbc0-db91a8e49ced" // This is just a made up uuid
)

// Structure keys
const (
	apiVersionName = "api_version"
)

// Url Paths
const (
	startSessionPath = "/clients_api/session"
)

var loginPostForm = url.Values{
	apiVersionName:                           {apiVersion},
	"device[hardware_id]":                    {hardWareId},
	"device[os]":                             {"android"},
	"device[app_brand]":                      {"ring"},
	"device[metadata][device_model]":         {"KVM"},
	"device[metadata][device_name]":          {"Python"},
	"device[metadata][resolution]":           {"600x800"},
	"device[metadata][app_version]":          {"1.3.806"},
	"device[metadata][app_instalation_date]": {""},
	"device[metadata][manufacturer]":         {"Qemu"},
	"device[metadata][device_type]":          {"desktop"},
	"device[metadata][architecture]":         {"desktop"},
	"device[metadata][language]":             {"en"},
}

// LoginResponse is the JSON structure returned when logging in
type LoginResponse struct {
	Profile Profile `json:"profile"`
}

// Profile is the JSON profile. It is a part of LoginResponse
type Profile struct {
	Id                   int      `json:"id"`
	Email                string   `json:"email"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	PhoneNumber          string   `json:"phone_number"`
	AuthenticationToken  string   `json:"authentication_token"`
	HardwareId           string   `json:"hardware_id"`
	ExplorerProgramTerms string   `json:"explorer_program_terms"`
	UserFlow             string   `json:"user_flow"`
	AppBrand             string   `json:"app_brand"`
	Features             Features `json:"features"`
}

// Features is the doorbell enabled features. It is a part of LoginResponse
type Features struct {
	RemoteLoggingFormatStoring               bool   `json:"remote_logging_format_storing"`
	RemoteLoggingLevel                       int    `json:"remote_logging_level"`
	SubscriptionsEnabled                     bool   `json:"subscriptions_enabled"`
	StickupcamSetupEnabled                   bool   `json:"stickupcam_setup_enabled"`
	VodEnabled                               bool   `json:"vod_enabled"`
	RingplusEnabled                          bool   `json:"ringplus_enabled"`
	LpdEnabled                               bool   `json:"lpd_enabled"`
	ReactiveSnoozingEnabled                  bool   `json:"reactive_snoozing_enabled"`
	ProactiveSnoozingEnabled                 bool   `json:"proactive_snoozing_enabled"`
	OwnerProactiveSnoozingEnabled            bool   `json:"owner_proactive_snoozing_enabled"`
	LiveViewSettingsEnabled                  bool   `json:"live_view_settings_enabled"`
	DeleteAllSettingsEnabled                 bool   `json:"delete_all_settings_enabled"`
	PowerCableEnabled                        bool   `json:"power_cable_enabled"`
	DeviceHealthAlertsEnabled                bool   `json:"device_health_alerts_enabled"`
	ChimeProEnabled                          bool   `json:"chime_pro_enabled"`
	MultipleCallsEnabled                     bool   `json:"multiple_calls_enabled"`
	UjetEnabled                              bool   `json:"ujet_enabled"`
	MultipleDeleteEnabled                    bool   `json:"multiple_delete_enabled"`
	DeleteAllEnabled                         bool   `json:"delete_all_enabled"`
	LpdMotionAnnouncementEnabled             bool   `json:"lpd_motion_announcement_enabled"`
	StarredEventsEnabled                     bool   `json:"starred_events_enabled"`
	ChimeDndEnabled                          bool   `json:"chime_dnd_enabled"`
	VideoSearchEnabled                       bool   `json:"video_search_enabled"`
	FloodlightCamEnabled                     bool   `json:"floodlight_cam_enabled"`
	RingCamBatteryEnabled                    bool   `json:"ring_cam_battery_enabled"`
	EliteCamEnabled                          bool   `json:"elite_cam_enabled"`
	DoorbellV2Enabled                        bool   `json:"doorbell_v2_enabled"`
	SpotlightBatteryDashboardControlsEnabled bool   `json:"spotlight_battery_dashboard_controls_enabled"`
	BypassAccountVerification                bool   `json:"bypass_account_verification"`
	LegacyCvrRetentionEnabled                bool   `json:"legacy_cvr_retention_enabled"`
	RingCamEnabled                           bool   `json:"ring_cam_enabled"`
	RingSearchEnabled                        bool   `json:"ring_search_enabled"`
	RingCamMountEnabled                      bool   `json:"ring_cam_mount_enabled"`
	RingAlarmEnabled                         bool   `json:"ring_alarm_enabled"`
	InAppCallNotifications                   bool   `json:"in_app_call_notifications"`
	RingCashEligibleEnabled                  bool   `json:"ring_cash_eligible_enabled"`
	AppAlertTonesEnabled                     bool   `json:"app_alert_tones_enabled"`
	MotionSnoozingEnabled                    bool   `json:"motion_snoozing_enabled"`
	HistoryClassificationEnabled             bool   `json:"history_classification_enabled"`
	TileDashboardEnabled                     bool   `json:"tile_dashboard_enabled"`
	TileDashboardMode                        string `json:"tile_dashboard_mode"`
	ScrubberAutoLiveEnabled                  bool   `json:"scrubber_auto_live_enabled"`
	ScrubberEnabled                          bool   `json:"scrubber_enabled"`
	NwEnabled                                bool   `json:"nw_enabled"`
	NwV2Enabled                              bool   `json:"nw_v2_enabled"`
	NwFeedTypesEnabled                       bool   `json:"nw_feed_types_enabled"`
	NwLargerAreaEnabled                      bool   `json:"nw_larger_area_enabled"`
	NwUserActivated                          bool   `json:"nw_user_activated"`
	NwNotificationTypesEnabled               bool   `json:"nw_notification_types_enabled"`
	NwNotificationRadiusEnabled              bool   `json:"nw_notification_radius_enabled"`
	NwMapViewFeatureEnabled                  bool   `json:"nw_map_view_feature_enabled"`
}

// Ring implements the non-public Ring Video Doorbell API
type Ring struct {
	client    *http.Client
	config    *Config
	authToken string
}

// Config basic configuration for creating a new Ring
type Config struct {
	Username string
	Password string
}

func New(c *Config) *Ring {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// cookiejar.New only returns an error for legacy reasons, this should never happen
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	return &Ring{
		client:    client,
		config:    c,
		authToken: "",
	}
}

func (r *Ring) Login() (*LoginResponse, error) {
	url := apiUri + startSessionPath

	request, err := http.NewRequest("POST", url, strings.NewReader(loginPostForm.Encode()))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(r.config.Username, r.config.Password)

	query := request.URL.Query()
	query.Add(apiVersionName, apiVersion)

	request.URL.RawQuery = query.Encode()

	resp, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Errorf("Unable to login %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	loginResponse := &LoginResponse{}
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(loginResponse)
	if err != nil {
		return nil, err
	}

	// Save the auth token so it can be used on future requests
	r.authToken = loginResponse.Profile.AuthenticationToken

	return loginResponse, nil
}
