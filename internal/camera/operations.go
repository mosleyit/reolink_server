package camera

import (
	"context"
	"fmt"
	"time"

	reolink "github.com/mosleyit/reolink_api_wrapper"
)

// ============================================================================
// System API Methods
// ============================================================================

// Reboot reboots the camera
func (c *CameraClient) Reboot(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.Reboot(ctx)
}

// GetTime gets the camera's time configuration
func (c *CameraClient) GetTime(ctx context.Context) (*reolink.TimeConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetTime(ctx)
}

// SetTime sets the camera's time configuration
func (c *CameraClient) SetTime(ctx context.Context, timeConfig *reolink.TimeConfig) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.SetTime(ctx, timeConfig)
}

// GetHddInfo gets HDD/SD card information
func (c *CameraClient) GetHddInfo(ctx context.Context) ([]reolink.HddInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetHddInfo(ctx)
}

// GetChannelStatus gets the channel status
func (c *CameraClient) GetChannelStatus(ctx context.Context) (*reolink.ChannelStatusValue, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetChannelStatus(ctx)
}

// GetAbility gets the camera's capabilities
func (c *CameraClient) GetAbility(ctx context.Context) (*reolink.Ability, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetAbility(ctx)
}

// GetDeviceName gets the camera's device name
func (c *CameraClient) GetDeviceName(ctx context.Context) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return "", fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetDeviceName(ctx)
}

// SetDeviceName sets the camera's device name
func (c *CameraClient) SetDeviceName(ctx context.Context, name string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.SetDeviceName(ctx, name)
}

// GetAutoMaint gets auto maintenance configuration
func (c *CameraClient) GetAutoMaint(ctx context.Context) (*reolink.AutoMaint, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetAutoMaint(ctx)
}

// SetAutoMaint sets auto maintenance configuration
func (c *CameraClient) SetAutoMaint(ctx context.Context, config reolink.AutoMaint) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.SetAutoMaint(ctx, config)
}

// GetAutoUpgrade gets auto upgrade configuration
func (c *CameraClient) GetAutoUpgrade(ctx context.Context) (*reolink.AutoUpgrade, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetAutoUpgrade(ctx)
}

// SetAutoUpgrade sets auto upgrade configuration
func (c *CameraClient) SetAutoUpgrade(ctx context.Context, enable bool) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.SetAutoUpgrade(ctx, enable)
}

// CheckFirmware checks for firmware updates
func (c *CameraClient) CheckFirmware(ctx context.Context) (*reolink.FirmwareCheck, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.CheckFirmware(ctx)
}

// Upgrade upgrades the camera firmware from a file
func (c *CameraClient) Upgrade(ctx context.Context, firmware []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.Upgrade(ctx, firmware)
}

// UpgradeOnline upgrades the camera firmware from online source
func (c *CameraClient) UpgradeOnline(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.UpgradeOnline(ctx)
}

// UpgradePrepare prepares for firmware upgrade
func (c *CameraClient) UpgradePrepare(ctx context.Context, restoreCfg bool, fileName string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.UpgradePrepare(ctx, restoreCfg, fileName)
}

// UpgradeStatus gets the firmware upgrade status
func (c *CameraClient) UpgradeStatus(ctx context.Context) (*reolink.UpgradeStatusInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.UpgradeStatus(ctx)
}

// Format formats the HDD/SD card
func (c *CameraClient) Format(ctx context.Context, hddID int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.Format(ctx, hddID)
}

// Restore performs a factory reset
func (c *CameraClient) Restore(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.Restore(ctx)
}

// GetSysCfg gets system configuration
func (c *CameraClient) GetSysCfg(ctx context.Context) (*reolink.SysCfg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.GetSysCfg(ctx)
}

// SetSysCfg sets system configuration
func (c *CameraClient) SetSysCfg(ctx context.Context, cfg reolink.SysCfg) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.System.SetSysCfg(ctx, cfg)
}

// ============================================================================
// Encoding API Methods
// ============================================================================

// GetSnapshot takes a snapshot from the camera
func (c *CameraClient) GetSnapshot(ctx context.Context, channel int) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Encoding.Snap(ctx, channel)
}

// GetEnc gets encoding configuration
func (c *CameraClient) GetEnc(ctx context.Context, channel int) (*reolink.EncConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Encoding.GetEnc(ctx, channel)
}

// SetEnc sets encoding configuration
func (c *CameraClient) SetEnc(ctx context.Context, config reolink.EncConfig) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Encoding.SetEnc(ctx, config)
}

// ============================================================================
// PTZ API Methods
// ============================================================================

// PTZMove moves the camera PTZ
func (c *CameraClient) PTZMove(ctx context.Context, operation string, speed int, channel int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	// Use PtzCtrl with PtzCtrlParam
	param := reolink.PtzCtrlParam{
		Channel: channel,
		Op:      operation,
		Speed:   speed,
	}
	return c.Client.PTZ.PtzCtrl(ctx, param)
}

// PTZStop stops PTZ movement
func (c *CameraClient) PTZStop(ctx context.Context, channel int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	// Use PtzCtrl with "Stop" operation
	param := reolink.PtzCtrlParam{
		Channel: channel,
		Op:      "Stop",
	}
	return c.Client.PTZ.PtzCtrl(ctx, param)
}

// PTZGotoPreset moves to a PTZ preset
func (c *CameraClient) PTZGotoPreset(ctx context.Context, channel int, presetID int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	// Use PtzCtrl with "ToPos" operation and preset ID
	param := reolink.PtzCtrlParam{
		Channel: channel,
		Op:      "ToPos",
		ID:      presetID,
	}
	return c.Client.PTZ.PtzCtrl(ctx, param)
}

// GetPtzPreset gets PTZ presets
func (c *CameraClient) GetPtzPreset(ctx context.Context, channel int) ([]reolink.PtzPreset, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetPtzPreset(ctx, channel)
}

// SetPtzPreset sets a PTZ preset
func (c *CameraClient) SetPtzPreset(ctx context.Context, preset reolink.PtzPreset) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.SetPtzPreset(ctx, preset)
}

// GetPtzPatrol gets PTZ patrol configuration
func (c *CameraClient) GetPtzPatrol(ctx context.Context, channel int) (*reolink.PtzPatrol, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetPtzPatrol(ctx, channel)
}

// SetPtzPatrol sets PTZ patrol configuration
func (c *CameraClient) SetPtzPatrol(ctx context.Context, patrol reolink.PtzPatrol) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.SetPtzPatrol(ctx, patrol)
}

// GetPtzGuard gets PTZ guard configuration
func (c *CameraClient) GetPtzGuard(ctx context.Context, channel int) (*reolink.PtzGuard, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetPtzGuard(ctx, channel)
}

// SetPtzGuard sets PTZ guard configuration
func (c *CameraClient) SetPtzGuard(ctx context.Context, guard reolink.PtzGuard) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.SetPtzGuard(ctx, guard)
}

// GetAutoFocus gets auto focus configuration
func (c *CameraClient) GetAutoFocus(ctx context.Context, channel int) (*reolink.AutoFocus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetAutoFocus(ctx, channel)
}

// SetAutoFocus sets auto focus configuration
func (c *CameraClient) SetAutoFocus(ctx context.Context, autoFocus reolink.AutoFocus) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.SetAutoFocus(ctx, autoFocus)
}

// GetZoomFocus gets zoom focus configuration
func (c *CameraClient) GetZoomFocus(ctx context.Context, channel int) (*reolink.ZoomFocus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetZoomFocus(ctx, channel)
}

// StartZoomFocus starts zoom/focus operation
func (c *CameraClient) StartZoomFocus(ctx context.Context, channel int, op string, pos int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.StartZoomFocus(ctx, channel, op, pos)
}

// GetPtzCheckState gets PTZ check state
func (c *CameraClient) GetPtzCheckState(ctx context.Context, channel int) (*reolink.PtzCheckState, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.GetPtzCheckState(ctx, channel)
}

// PtzCheck performs PTZ check
func (c *CameraClient) PtzCheck(ctx context.Context, channel int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.PTZ.PtzCheck(ctx, channel)
}

// ============================================================================
// LED API Methods
// ============================================================================

// SetIRLights controls the IR lights
func (c *CameraClient) SetIRLights(ctx context.Context, channel int, state string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.SetIrLights(ctx, channel, state)
}

// SetWhiteLED controls the white LED
func (c *CameraClient) SetWhiteLED(ctx context.Context, config *reolink.WhiteLed) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.SetWhiteLed(ctx, *config)
}

// GetIrLights gets IR lights configuration
func (c *CameraClient) GetIrLights(ctx context.Context) (*reolink.IrLights, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.GetIrLights(ctx)
}

// GetWhiteLed gets white LED configuration
func (c *CameraClient) GetWhiteLed(ctx context.Context, channel int) (*reolink.WhiteLed, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.GetWhiteLed(ctx, channel)
}

// GetPowerLed gets power LED configuration
func (c *CameraClient) GetPowerLed(ctx context.Context, channel int) (*reolink.PowerLed, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.GetPowerLed(ctx, channel)
}

// SetPowerLed sets power LED configuration
func (c *CameraClient) SetPowerLed(ctx context.Context, channel int, state string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.SetPowerLed(ctx, channel, state)
}

// SetAlarmArea sets alarm detection area/zone
func (c *CameraClient) SetAlarmArea(ctx context.Context, params map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.SetAlarmArea(ctx, params)
}

// GetAiAlarm gets AI alarm configuration
func (c *CameraClient) GetAiAlarm(ctx context.Context, channel int, aiType string) (*reolink.AiAlarm, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.GetAiAlarm(ctx, channel, aiType)
}

// SetAiAlarm sets AI alarm configuration
func (c *CameraClient) SetAiAlarm(ctx context.Context, channel int, alarm reolink.AiAlarm) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.LED.SetAiAlarm(ctx, channel, alarm)
}

// ============================================================================
// Alarm API Methods
// ============================================================================

// TriggerSiren triggers the camera siren
func (c *CameraClient) TriggerSiren(ctx context.Context, channel int, duration int) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	// AudioAlarmPlayParam fields: Channel, AlarmMode, ManualSwitch, Times
	param := reolink.AudioAlarmPlayParam{
		Channel:      channel,
		AlarmMode:    "manul", // Manual mode (note: API uses "manul" typo)
		ManualSwitch: 1,       // Enable
		Times:        duration,
	}

	return c.Client.Alarm.AudioAlarmPlay(ctx, param)
}

// GetMotionState gets the current motion detection state
func (c *CameraClient) GetMotionState(ctx context.Context, channel int) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return 0, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	// GetMdState returns (int, error) not (*MdStateValue, error)
	return c.Client.Alarm.GetMdState(ctx, channel)
}

// GetMdAlarm gets motion detection alarm configuration
func (c *CameraClient) GetMdAlarm(ctx context.Context, channel int) (*reolink.MdAlarm, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.GetMdAlarm(ctx, channel)
}

// SetMdAlarm sets motion detection alarm configuration
func (c *CameraClient) SetMdAlarm(ctx context.Context, config reolink.MdAlarm) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.SetMdAlarm(ctx, config)
}

// GetAlarm gets alarm configuration
func (c *CameraClient) GetAlarm(ctx context.Context, channel int, alarmType string) (*reolink.Alarm, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.GetAlarm(ctx, channel, alarmType)
}

// SetAlarm sets alarm configuration
func (c *CameraClient) SetAlarm(ctx context.Context, alarm reolink.Alarm) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.SetAlarm(ctx, alarm)
}

// GetAudioAlarm gets audio alarm configuration
func (c *CameraClient) GetAudioAlarm(ctx context.Context, channel int) (*reolink.AudioAlarm, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.GetAudioAlarm(ctx, channel)
}

// SetAudioAlarm sets audio alarm configuration
func (c *CameraClient) SetAudioAlarm(ctx context.Context, audioAlarm reolink.AudioAlarm) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.SetAudioAlarm(ctx, audioAlarm)
}

// GetBuzzerAlarmV20 gets buzzer alarm configuration (V20 API)
func (c *CameraClient) GetBuzzerAlarmV20(ctx context.Context, channel int) (*reolink.BuzzerAlarm, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.GetBuzzerAlarmV20(ctx, channel)
}

// SetBuzzerAlarmV20 sets buzzer alarm configuration (V20 API)
func (c *CameraClient) SetBuzzerAlarmV20(ctx context.Context, buzzerAlarm reolink.BuzzerAlarm) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Alarm.SetBuzzerAlarmV20(ctx, buzzerAlarm)
}

// ============================================================================
// AI API Methods
// ============================================================================

// GetAIState gets the current AI detection state
func (c *CameraClient) GetAIState(ctx context.Context, channel int) (*reolink.AiState, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.AI.GetAiState(ctx, channel)
}

// GetAiCfg gets AI detection configuration
func (c *CameraClient) GetAiCfg(ctx context.Context, channel int) (*reolink.AiCfg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.AI.GetAiCfg(ctx, channel)
}

// SetAiCfg sets AI detection configuration
func (c *CameraClient) SetAiCfg(ctx context.Context, config reolink.AiCfg) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.AI.SetAiCfg(ctx, config)
}

// ============================================================================
// Streaming API Methods
// ============================================================================

// GetRTSPURL gets the RTSP URL for the camera
func (c *CameraClient) GetRTSPURL(streamType reolink.StreamType, channel int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Client.Streaming.GetRTSPURL(streamType, channel)
}

// GetFLVURL gets the FLV URL for the camera
func (c *CameraClient) GetFLVURL(streamType reolink.StreamType, channelID int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Client.Streaming.GetFLVURL(streamType, channelID)
}

// GetRTMPURL gets the RTMP URL for the camera
func (c *CameraClient) GetRTMPURL(streamType reolink.StreamType, channelID int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Client.Streaming.GetRTMPURL(streamType, channelID)
}

// ============================================================================
// Recording API Methods
// ============================================================================

// GetRec gets recording configuration (v1.0)
func (c *CameraClient) GetRec(ctx context.Context, channel int) (*reolink.Rec, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.GetRec(ctx, channel)
}

// SetRec sets recording configuration (v1.0)
func (c *CameraClient) SetRec(ctx context.Context, rec reolink.Rec) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.SetRec(ctx, rec)
}

// GetRecV20 gets recording configuration (v2.0)
func (c *CameraClient) GetRecV20(ctx context.Context, channel int) (*reolink.Rec, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.GetRecV20(ctx, channel)
}

// SetRecV20 sets recording configuration (v2.0)
func (c *CameraClient) SetRecV20(ctx context.Context, rec reolink.Rec) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.SetRecV20(ctx, rec)
}

// Search searches for recordings within a time range
func (c *CameraClient) Search(ctx context.Context, channel int, startTime, endTime time.Time, streamType string) ([]reolink.SearchResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.Search(ctx, channel, startTime, endTime, streamType)
}

// Download downloads a recording
func (c *CameraClient) Download(source, output string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Client.Recording.Download(source, output)
}

// Playback plays back a recording
func (c *CameraClient) Playback(source, output string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Client.Recording.Playback(source, output)
}

// NvrDownload downloads a recording from NVR
func (c *CameraClient) NvrDownload(ctx context.Context, params map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Recording.NvrDownload(ctx, params)
}

// ============================================================================
// Video API Methods
// ============================================================================

// GetOsd gets OSD (On-Screen Display) configuration
func (c *CameraClient) GetOsd(ctx context.Context, channel int) (*reolink.Osd, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetOsd(ctx, channel)
}

// SetOsd sets OSD (On-Screen Display) configuration
func (c *CameraClient) SetOsd(ctx context.Context, osd reolink.Osd) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetOsd(ctx, osd)
}

// GetImage gets image settings (brightness, contrast, saturation, etc.)
func (c *CameraClient) GetImage(ctx context.Context, channel int) (*reolink.Image, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetImage(ctx, channel)
}

// SetImage sets image settings (brightness, contrast, saturation, etc.)
func (c *CameraClient) SetImage(ctx context.Context, image reolink.Image) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetImage(ctx, image)
}

// GetIsp gets ISP (Image Signal Processing) settings
func (c *CameraClient) GetIsp(ctx context.Context, channel int) (*reolink.Isp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetIsp(ctx, channel)
}

// SetIsp sets ISP (Image Signal Processing) settings
func (c *CameraClient) SetIsp(ctx context.Context, isp reolink.Isp) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetIsp(ctx, isp)
}

// GetMask gets privacy mask configuration
func (c *CameraClient) GetMask(ctx context.Context, channel int) (*reolink.Mask, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetMask(ctx, channel)
}

// SetMask sets privacy mask configuration
func (c *CameraClient) SetMask(ctx context.Context, mask reolink.Mask) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetMask(ctx, mask)
}

// GetCrop gets video crop configuration
func (c *CameraClient) GetCrop(ctx context.Context, channel int) (*reolink.Crop, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetCrop(ctx, channel)
}

// SetCrop sets video crop configuration
func (c *CameraClient) SetCrop(ctx context.Context, crop reolink.Crop) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetCrop(ctx, crop)
}

// GetStitch gets panoramic stitching configuration
func (c *CameraClient) GetStitch(ctx context.Context) (*reolink.Stitch, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.GetStitch(ctx)
}

// SetStitch sets panoramic stitching configuration
func (c *CameraClient) SetStitch(ctx context.Context, stitch reolink.Stitch) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Video.SetStitch(ctx, stitch)
}

// ============================================================================
// Network API Methods
// ============================================================================

// GetNetPort gets network port configuration (HTTP, RTSP, RTMP, ONVIF)
func (c *CameraClient) GetNetPort(ctx context.Context) (*reolink.NetPort, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetNetPort(ctx)
}

// SetNetPort sets network port configuration
func (c *CameraClient) SetNetPort(ctx context.Context, netPort reolink.NetPort) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetNetPort(ctx, netPort)
}

// GetLocalLink gets local link configuration
func (c *CameraClient) GetLocalLink(ctx context.Context) (*reolink.LocalLink, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetLocalLink(ctx)
}

// SetLocalLink sets local link configuration
func (c *CameraClient) SetLocalLink(ctx context.Context, localLink reolink.LocalLink) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetLocalLink(ctx, localLink)
}

// GetNtp gets NTP configuration
func (c *CameraClient) GetNtp(ctx context.Context) (*reolink.Ntp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetNtp(ctx)
}

// SetNtp sets NTP configuration
func (c *CameraClient) SetNtp(ctx context.Context, ntp reolink.Ntp) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetNtp(ctx, ntp)
}

// GetWifi gets WiFi configuration
func (c *CameraClient) GetWifi(ctx context.Context) (*reolink.Wifi, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetWifi(ctx)
}

// SetWifi sets WiFi configuration
func (c *CameraClient) SetWifi(ctx context.Context, wifi reolink.Wifi) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetWifi(ctx, wifi)
}

// ScanWifi scans for available WiFi networks
func (c *CameraClient) ScanWifi(ctx context.Context) ([]reolink.WifiNetwork, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.ScanWifi(ctx)
}

// GetWifiSignal gets WiFi signal strength
func (c *CameraClient) GetWifiSignal(ctx context.Context) (*reolink.WifiSignal, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetWifiSignal(ctx)
}

// GetDdns gets DDNS configuration
func (c *CameraClient) GetDdns(ctx context.Context) (*reolink.Ddns, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetDdns(ctx)
}

// SetDdns sets DDNS configuration
func (c *CameraClient) SetDdns(ctx context.Context, ddns reolink.Ddns) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetDdns(ctx, ddns)
}

// GetEmail gets email notification configuration
func (c *CameraClient) GetEmail(ctx context.Context) (*reolink.Email, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetEmail(ctx)
}

// SetEmail sets email notification configuration
func (c *CameraClient) SetEmail(ctx context.Context, email reolink.Email) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetEmail(ctx, email)
}

// GetEmailV20 gets email notification configuration (v2.0)
func (c *CameraClient) GetEmailV20(ctx context.Context, channel int) (*reolink.Email, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetEmailV20(ctx, channel)
}

// SetEmailV20 sets email notification configuration (v2.0)
func (c *CameraClient) SetEmailV20(ctx context.Context, channel int, email reolink.Email) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetEmailV20(ctx, channel, email)
}

// GetFtp gets FTP configuration
func (c *CameraClient) GetFtp(ctx context.Context) (*reolink.Ftp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetFtp(ctx)
}

// SetFtp sets FTP configuration
func (c *CameraClient) SetFtp(ctx context.Context, ftp reolink.Ftp) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetFtp(ctx, ftp)
}

// GetFtpV20 gets FTP configuration (v2.0)
func (c *CameraClient) GetFtpV20(ctx context.Context, channel int) (*reolink.Ftp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetFtpV20(ctx, channel)
}

// SetFtpV20 sets FTP configuration (v2.0)
func (c *CameraClient) SetFtpV20(ctx context.Context, channel int, ftp reolink.Ftp) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetFtpV20(ctx, channel, ftp)
}

// GetPush gets push notification configuration
func (c *CameraClient) GetPush(ctx context.Context) (*reolink.Push, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetPush(ctx)
}

// SetPush sets push notification configuration
func (c *CameraClient) SetPush(ctx context.Context, push reolink.Push) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetPush(ctx, push)
}

// GetPushV20 gets push notification configuration (v2.0)
func (c *CameraClient) GetPushV20(ctx context.Context, channel int) (*reolink.Push, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetPushV20(ctx, channel)
}

// SetPushV20 sets push notification configuration (v2.0)
func (c *CameraClient) SetPushV20(ctx context.Context, channel int, push reolink.Push) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetPushV20(ctx, channel, push)
}

// GetPushCfg gets push configuration
func (c *CameraClient) GetPushCfg(ctx context.Context) (*reolink.PushCfg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetPushCfg(ctx)
}

// SetPushCfg sets push configuration
func (c *CameraClient) SetPushCfg(ctx context.Context, pushCfg reolink.PushCfg) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetPushCfg(ctx, pushCfg)
}

// GetP2p gets P2P configuration
func (c *CameraClient) GetP2p(ctx context.Context) (*reolink.P2p, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetP2p(ctx)
}

// SetP2p sets P2P configuration
func (c *CameraClient) SetP2p(ctx context.Context, p2p reolink.P2p) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetP2p(ctx, p2p)
}

// GetUpnp gets UPnP configuration
func (c *CameraClient) GetUpnp(ctx context.Context) (*reolink.Upnp, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetUpnp(ctx)
}

// SetUpnp sets UPnP configuration
func (c *CameraClient) SetUpnp(ctx context.Context, upnp reolink.Upnp) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.SetUpnp(ctx, upnp)
}

// GetRtspUrl gets RTSP URL configuration
func (c *CameraClient) GetRtspUrl(ctx context.Context, channel int) (*reolink.RtspUrl, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Network.GetRtspUrl(ctx, channel)
}

// ============================================================================
// Security API Methods
// ============================================================================

// GetUsers gets list of users
func (c *CameraClient) GetUsers(ctx context.Context) ([]reolink.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.GetUsers(ctx)
}

// AddUser adds a new user
func (c *CameraClient) AddUser(ctx context.Context, user reolink.User) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.AddUser(ctx, user)
}

// ModifyUser modifies an existing user
func (c *CameraClient) ModifyUser(ctx context.Context, user reolink.User) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.ModifyUser(ctx, user)
}

// DeleteUser deletes a user
func (c *CameraClient) DeleteUser(ctx context.Context, username string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.DeleteUser(ctx, username)
}

// GetOnlineUsers gets list of currently online users
func (c *CameraClient) GetOnlineUsers(ctx context.Context) ([]reolink.OnlineUser, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.GetOnlineUsers(ctx)
}

// DisconnectUser disconnects a user session
func (c *CameraClient) DisconnectUser(ctx context.Context, username string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.DisconnectUser(ctx, username)
}

// GetCertificateInfo gets SSL certificate information
func (c *CameraClient) GetCertificateInfo(ctx context.Context) (*reolink.CertificateInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return nil, fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.GetCertificateInfo(ctx)
}

// CertificateClear clears SSL certificate
func (c *CameraClient) CertificateClear(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.CircuitOpen {
		return fmt.Errorf("circuit open for camera %s", c.Camera.ID)
	}

	return c.Client.Security.CertificateClear(ctx)
}
