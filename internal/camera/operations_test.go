package camera

import (
	"context"
	"testing"
	"time"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test camera client with circuit breaker open
func createTestCameraClientWithCircuitOpen() *CameraClient {
	return &CameraClient{
		Camera: &models.Camera{
			ID:       "test-camera",
			Name:     "Test Camera",
			Host:     "192.168.1.100",
			Port:     80,
			Username: "admin",
			Password: "password",
		},
		Client:       nil, // No real client needed for circuit breaker tests
		LastHealthy:  time.Now().Add(-1 * time.Hour),
		FailureCount: 5,
		CircuitOpen:  true,
	}
}

func TestCameraClient_Reboot_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.Reboot(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_GetSnapshot_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	data, err := client.GetSnapshot(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_PTZMove_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.PTZMove(ctx, "Up", 32, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_PTZStop_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.PTZStop(ctx, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_PTZGotoPreset_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.PTZGotoPreset(ctx, 0, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_SetIRLights_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetIRLights(ctx, 0, "Auto")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_SetWhiteLED_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	config := &reolink.WhiteLed{
		Channel: 0,
		Mode:    1,
	}

	err := client.SetWhiteLED(ctx, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_TriggerSiren_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.TriggerSiren(ctx, 0, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_GetMotionState_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	state, err := client.GetMotionState(ctx, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, state)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

func TestCameraClient_GetAIState_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	state, err := client.GetAIState(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, state)
	assert.Contains(t, err.Error(), "circuit open")
	assert.Contains(t, err.Error(), "test-camera")
}

// ============================================================================
// System API Tests
// ============================================================================

func TestCameraClient_GetTime_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	timeConfig, err := client.GetTime(ctx)
	assert.Error(t, err)
	assert.Nil(t, timeConfig)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetTime_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetTime(ctx, &reolink.TimeConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetHddInfo_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	hddInfo, err := client.GetHddInfo(ctx)
	assert.Error(t, err)
	assert.Nil(t, hddInfo)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetChannelStatus_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	status, err := client.GetChannelStatus(ctx)
	assert.Error(t, err)
	assert.Nil(t, status)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAbility_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	ability, err := client.GetAbility(ctx)
	assert.Error(t, err)
	assert.Nil(t, ability)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetDeviceName_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	name, err := client.GetDeviceName(ctx)
	assert.Error(t, err)
	assert.Empty(t, name)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetDeviceName_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetDeviceName(ctx, "New Name")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAutoMaint_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	autoMaint, err := client.GetAutoMaint(ctx)
	assert.Error(t, err)
	assert.Nil(t, autoMaint)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAutoMaint_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAutoMaint(ctx, reolink.AutoMaint{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAutoUpgrade_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	autoUpgrade, err := client.GetAutoUpgrade(ctx)
	assert.Error(t, err)
	assert.Nil(t, autoUpgrade)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAutoUpgrade_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAutoUpgrade(ctx, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_CheckFirmware_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	firmwareCheck, err := client.CheckFirmware(ctx)
	assert.Error(t, err)
	assert.Nil(t, firmwareCheck)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_Upgrade_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.Upgrade(ctx, []byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_UpgradeOnline_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.UpgradeOnline(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_Format_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.Format(ctx, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_Restore_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.Restore(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetSysCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	sysCfg, err := client.GetSysCfg(ctx)
	assert.Error(t, err)
	assert.Nil(t, sysCfg)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetSysCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetSysCfg(ctx, reolink.SysCfg{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Encoding API Tests
// ============================================================================

func TestCameraClient_GetEnc_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	encConfig, err := client.GetEnc(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, encConfig)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetEnc_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetEnc(ctx, reolink.EncConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// PTZ API Tests
// ============================================================================

func TestCameraClient_GetPtzPreset_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	presets, err := client.GetPtzPreset(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, presets)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPtzPreset_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPtzPreset(ctx, reolink.PtzPreset{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPtzPatrol_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	patrol, err := client.GetPtzPatrol(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, patrol)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPtzPatrol_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPtzPatrol(ctx, reolink.PtzPatrol{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPtzGuard_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	guard, err := client.GetPtzGuard(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, guard)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPtzGuard_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPtzGuard(ctx, reolink.PtzGuard{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAutoFocus_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	autoFocus, err := client.GetAutoFocus(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, autoFocus)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAutoFocus_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAutoFocus(ctx, reolink.AutoFocus{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetZoomFocus_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	zoomFocus, err := client.GetZoomFocus(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, zoomFocus)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_StartZoomFocus_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.StartZoomFocus(ctx, 0, "ZoomInc", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPtzCheckState_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	checkState, err := client.GetPtzCheckState(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, checkState)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_PtzCheck_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.PtzCheck(ctx, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// LED API Tests
// ============================================================================

func TestCameraClient_GetIrLights_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	irLights, err := client.GetIrLights(ctx)
	assert.Error(t, err)
	assert.Nil(t, irLights)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetWhiteLed_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	whiteLed, err := client.GetWhiteLed(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, whiteLed)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPowerLed_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	powerLed, err := client.GetPowerLed(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, powerLed)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPowerLed_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPowerLed(ctx, 0, "Auto")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAlarmArea_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAlarmArea(ctx, map[string]interface{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAiAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	aiAlarm, err := client.GetAiAlarm(ctx, 0, "people")
	assert.Error(t, err)
	assert.Nil(t, aiAlarm)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAiAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAiAlarm(ctx, 0, reolink.AiAlarm{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Alarm API Tests
// ============================================================================

func TestCameraClient_GetMdAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	mdAlarm, err := client.GetMdAlarm(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, mdAlarm)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetMdAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetMdAlarm(ctx, reolink.MdAlarm{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	alarm, err := client.GetAlarm(ctx, 0, "md")
	assert.Error(t, err)
	assert.Nil(t, alarm)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAlarm(ctx, reolink.Alarm{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetAudioAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	audioAlarm, err := client.GetAudioAlarm(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, audioAlarm)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAudioAlarm_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAudioAlarm(ctx, reolink.AudioAlarm{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetBuzzerAlarmV20_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	buzzerAlarm, err := client.GetBuzzerAlarmV20(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, buzzerAlarm)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetBuzzerAlarmV20_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetBuzzerAlarmV20(ctx, reolink.BuzzerAlarm{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// AI API Tests
// ============================================================================

func TestCameraClient_GetAiCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	aiCfg, err := client.GetAiCfg(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, aiCfg)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetAiCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetAiCfg(ctx, reolink.AiCfg{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Recording API Tests
// ============================================================================

func TestCameraClient_GetRec_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	rec, err := client.GetRec(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, rec)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetRec_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetRec(ctx, reolink.Rec{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetRecV20_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	rec, err := client.GetRecV20(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, rec)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetRecV20_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetRecV20(ctx, reolink.Rec{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_Search_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	results, err := client.Search(ctx, 0, time.Now(), time.Now(), "main")
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_NvrDownload_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.NvrDownload(ctx, map[string]interface{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Video API Tests
// ============================================================================

func TestCameraClient_GetOsd_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	osd, err := client.GetOsd(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, osd)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetOsd_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetOsd(ctx, reolink.Osd{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetImage_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	image, err := client.GetImage(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, image)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetImage_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetImage(ctx, reolink.Image{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetIsp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	isp, err := client.GetIsp(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, isp)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetIsp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetIsp(ctx, reolink.Isp{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetMask_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	mask, err := client.GetMask(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, mask)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetMask_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetMask(ctx, reolink.Mask{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetCrop_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	crop, err := client.GetCrop(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, crop)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetCrop_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetCrop(ctx, reolink.Crop{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetStitch_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	stitch, err := client.GetStitch(ctx)
	assert.Error(t, err)
	assert.Nil(t, stitch)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetStitch_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetStitch(ctx, reolink.Stitch{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Network API Tests
// ============================================================================

func TestCameraClient_GetNetPort_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	netPort, err := client.GetNetPort(ctx)
	assert.Error(t, err)
	assert.Nil(t, netPort)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetNetPort_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetNetPort(ctx, reolink.NetPort{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetLocalLink_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	localLink, err := client.GetLocalLink(ctx)
	assert.Error(t, err)
	assert.Nil(t, localLink)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetLocalLink_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetLocalLink(ctx, reolink.LocalLink{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetNtp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	ntp, err := client.GetNtp(ctx)
	assert.Error(t, err)
	assert.Nil(t, ntp)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetNtp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetNtp(ctx, reolink.Ntp{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetWifi_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	wifi, err := client.GetWifi(ctx)
	assert.Error(t, err)
	assert.Nil(t, wifi)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetWifi_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetWifi(ctx, reolink.Wifi{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_ScanWifi_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	networks, err := client.ScanWifi(ctx)
	assert.Error(t, err)
	assert.Nil(t, networks)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetWifiSignal_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	signal, err := client.GetWifiSignal(ctx)
	assert.Error(t, err)
	assert.Nil(t, signal)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetDdns_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	ddns, err := client.GetDdns(ctx)
	assert.Error(t, err)
	assert.Nil(t, ddns)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetDdns_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetDdns(ctx, reolink.Ddns{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetEmail_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	email, err := client.GetEmail(ctx)
	assert.Error(t, err)
	assert.Nil(t, email)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetEmail_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetEmail(ctx, reolink.Email{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetFtp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	ftp, err := client.GetFtp(ctx)
	assert.Error(t, err)
	assert.Nil(t, ftp)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetFtp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetFtp(ctx, reolink.Ftp{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPush_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	push, err := client.GetPush(ctx)
	assert.Error(t, err)
	assert.Nil(t, push)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPush_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPush(ctx, reolink.Push{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetPushCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	pushCfg, err := client.GetPushCfg(ctx)
	assert.Error(t, err)
	assert.Nil(t, pushCfg)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetPushCfg_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetPushCfg(ctx, reolink.PushCfg{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetP2p_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	p2p, err := client.GetP2p(ctx)
	assert.Error(t, err)
	assert.Nil(t, p2p)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetP2p_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetP2p(ctx, reolink.P2p{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetUpnp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	upnp, err := client.GetUpnp(ctx)
	assert.Error(t, err)
	assert.Nil(t, upnp)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_SetUpnp_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.SetUpnp(ctx, reolink.Upnp{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetRtspUrl_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	rtspUrl, err := client.GetRtspUrl(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, rtspUrl)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Security API Tests
// ============================================================================

func TestCameraClient_GetUsers_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	users, err := client.GetUsers(ctx)
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_AddUser_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.AddUser(ctx, reolink.User{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_ModifyUser_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.ModifyUser(ctx, reolink.User{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_DeleteUser_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.DeleteUser(ctx, "testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetOnlineUsers_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	users, err := client.GetOnlineUsers(ctx)
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_DisconnectUser_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.DisconnectUser(ctx, "testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_GetCertificateInfo_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	certInfo, err := client.GetCertificateInfo(ctx)
	assert.Error(t, err)
	assert.Nil(t, certInfo)
	assert.Contains(t, err.Error(), "circuit open")
}

func TestCameraClient_CertificateClear_CircuitOpen(t *testing.T) {
	client := createTestCameraClientWithCircuitOpen()
	ctx := context.Background()

	err := client.CertificateClear(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit open")
}

// ============================================================================
// Streaming API Tests
// ============================================================================

func TestCameraClient_GetRTSPURL(t *testing.T) {
	// Create a real client to test URL generation
	camera := &models.Camera{
		ID:       "test-camera",
		Name:     "Test Camera",
		Host:     "192.168.1.100",
		Port:     80,
		Username: "admin",
		Password: "password",
	}

	client := reolink.NewClient(camera.Host,
		reolink.WithCredentials(camera.Username, camera.Password),
	)

	cameraClient := &CameraClient{
		Camera:       camera,
		Client:       client,
		LastHealthy:  time.Now(),
		FailureCount: 0,
		CircuitOpen:  false,
	}

	// Test with main stream
	url := cameraClient.GetRTSPURL(reolink.StreamMain, 0)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "rtsp://")
	assert.Contains(t, url, camera.Host)

	// Test with sub stream
	url = cameraClient.GetRTSPURL(reolink.StreamSub, 0)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "rtsp://")
	assert.Contains(t, url, camera.Host)
}

func TestCameraClient_GetFLVURL(t *testing.T) {
	camera := &models.Camera{
		ID:       "test-camera",
		Name:     "Test Camera",
		Host:     "192.168.1.100",
		Port:     80,
		Username: "admin",
		Password: "password",
	}

	client := reolink.NewClient(camera.Host,
		reolink.WithCredentials(camera.Username, camera.Password),
	)

	cameraClient := &CameraClient{
		Camera:       camera,
		Client:       client,
		LastHealthy:  time.Now(),
		FailureCount: 0,
		CircuitOpen:  false,
	}

	url := cameraClient.GetFLVURL(reolink.StreamMain, 0)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, camera.Host)
}

func TestCameraClient_GetRTMPURL(t *testing.T) {
	camera := &models.Camera{
		ID:       "test-camera",
		Name:     "Test Camera",
		Host:     "192.168.1.100",
		Port:     80,
		Username: "admin",
		Password: "password",
	}

	client := reolink.NewClient(camera.Host,
		reolink.WithCredentials(camera.Username, camera.Password),
	)

	cameraClient := &CameraClient{
		Camera:       camera,
		Client:       client,
		LastHealthy:  time.Now(),
		FailureCount: 0,
		CircuitOpen:  false,
	}

	url := cameraClient.GetRTMPURL(reolink.StreamMain, 0)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, camera.Host)
}

func TestCameraClient_PTZOperations_Parameters(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		speed     int
		channel   int
	}{
		{
			name:      "Move Up",
			operation: "Up",
			speed:     32,
			channel:   0,
		},
		{
			name:      "Move Down",
			operation: "Down",
			speed:     16,
			channel:   0,
		},
		{
			name:      "Move Left",
			operation: "Left",
			speed:     8,
			channel:   1,
		},
		{
			name:      "Move Right",
			operation: "Right",
			speed:     64,
			channel:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with circuit open to verify parameters are passed correctly
			client := createTestCameraClientWithCircuitOpen()
			ctx := context.Background()

			err := client.PTZMove(ctx, tt.operation, tt.speed, tt.channel)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "circuit open")
		})
	}
}

func TestCameraClient_PTZGotoPreset_Parameters(t *testing.T) {
	tests := []struct {
		name     string
		channel  int
		presetID int
	}{
		{
			name:     "Preset 1",
			channel:  0,
			presetID: 1,
		},
		{
			name:     "Preset 5",
			channel:  0,
			presetID: 5,
		},
		{
			name:     "Preset 10 on channel 1",
			channel:  1,
			presetID: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := createTestCameraClientWithCircuitOpen()
			ctx := context.Background()

			err := client.PTZGotoPreset(ctx, tt.channel, tt.presetID)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "circuit open")
		})
	}
}

// Benchmark tests
func BenchmarkCameraClient_GetRTSPURL(b *testing.B) {
	camera := &models.Camera{
		ID:       "test-camera",
		Host:     "192.168.1.100",
		Port:     80,
		Username: "admin",
		Password: "password",
	}

	client := reolink.NewClient(camera.Host,
		reolink.WithCredentials(camera.Username, camera.Password),
	)
	cameraClient := &CameraClient{
		Camera:      camera,
		Client:      client,
		CircuitOpen: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cameraClient.GetRTSPURL(reolink.StreamMain, 0)
	}
}
