package fantuner

import (
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/internal/proto"
	"github.com/kmou424/syncfans/internal/toolkit/sysfskit"
	"math"
)

// CalcAvgTemp updates temperature sliding window and returns smoothed value
// 更新滑动窗口并返回平滑后的温度值
func (fi *FanInfo) CalcAvgTemp(current float64) float64 {
	// Shift window and append new reading
	copy(fi.tempHist[:], fi.tempHist[1:])
	fi.tempHist[len(fi.tempHist)-1] = current

	// Calculate moving average
	var sum float64
	for _, t := range fi.tempHist {
		sum += t
	}
	return sum / float64(len(fi.tempHist))
}

// CalcTargetSpeed calculates target s based on thermal conditions
// 确定目标转速：根据温度条件计算理论目标转速
func (fi *FanInfo) CalcTargetSpeed(avgTemp float64) (speed int, force bool) {
	curveType := fi.curveType
	curveFactor := fi.curveFactor
	criticalMargin := fi.criMargin

	// Calculate effective temperature thresholds
	// 计算温度阈值
	baseLower, baseUpper := fi.criTemps[0], fi.criTemps[1]
	thresholdLower := baseLower - criticalMargin
	thresholdUpper := baseUpper + criticalMargin

	// Immediate boundary handling
	// 边界处理
	if avgTemp == 0 { // 首次启动
		return fi.MaxSpeed, true
	}
	if avgTemp <= thresholdLower {
		return fi.MinSpeed, true
	}
	if avgTemp >= thresholdUpper {
		return fi.MaxSpeed, true
	}

	// Normalize temperature position
	// 线性处理
	tempRange := thresholdUpper - thresholdLower
	normalized := (avgTemp - thresholdLower) / tempRange
	normalized = math.Max(0, math.Min(normalized, 1))

	// Apply curve algorithm
	// 使用曲线算法
	switch curveType {
	case "s-curve":
		normalized = normalized * normalized * (3 - 2*normalized)
	case "exponential":
		normalized = math.Pow(normalized, curveFactor)
	case "aggressive":
		normalized = 0.5 * (math.Tanh((normalized-0.5)*curveFactor) + 1)
	}

	// Calculate target speed
	// 计算目标转速
	speedRange := float64(fi.MaxSpeed - fi.MinSpeed)
	return fi.MinSpeed + int(math.Round(normalized*speedRange)), true
}

// CalcSmoothedSpeed performs exponential smoothing with momentum reset
// 应用平滑算法：执行指数平滑处理，检测突变时重置动量
func (fi *FanInfo) CalcSmoothedSpeed(target int) int {
	// Detect abrupt changes (>20% speed range)
	// 检测突变：如果目标转速变化幅度大于20%，则重置平滑速度
	speedRange := fi.MaxSpeed - fi.MinSpeed
	if math.Abs(float64(target-fi.lastTarget)) > float64(speedRange)*0.2 {
		fi.smoothedSpeed = float64(target)
	}

	// Apply EMA smoothing
	// 执行指数平滑处理
	alpha := conf.GetServerConfig().Config.SmoothingFactor
	fi.smoothedSpeed = fi.smoothedSpeed*(1-alpha) + float64(target)*alpha

	// Clamp to valid range
	// 限制有效范围
	return clamp(int(math.Round(fi.smoothedSpeed)), fi.MinSpeed, fi.MaxSpeed)
}

// ApplySpeed writes final speed to hardware if changed
// 应用转速调整：将最终的目标转速写入硬件
func (fi *FanInfo) ApplySpeed(finalSpeed int) error {
	if finalSpeed == fi.lastTarget {
		return nil
	}

	slog.Infof("Applying speed %d to %s", finalSpeed, fi.Path)

	if err := sysfskit.WriteInt(fi.Path, finalSpeed); err != nil {
		return caused.FileSystemError(ero.Wrap(err, "failed to write speed to %s", fi.Path))
	}

	slog.Debugf("Applied speed %d to %s", finalSpeed, fi.Path)

	// Update last target
	// 更新上一次的目标转速
	fi.lastTarget = finalSpeed
	return nil
}

func (r *fansManager) Tuning(name string, report *proto.ReportSysInfo) {
	fi, ok := r.fanInfoMap.Get(name)
	if !ok {
		slog.Warnf("Received report for unregistered fan: %s", name)
		return
	}

	// lock fan info to prevent concurrent access
	// 锁定 FanInfo 实例，防止并发访问
	fi.lk.Lock()
	defer fi.lk.Unlock()

	slog.Debugf("Received report from %s: %v", name, report)
	if !fi.limiter.Allow() {
		return
	}

	// Update temperature history, calculate history-smoothed temperature
	// 更新温度历史, 计算历史平滑温度
	avgTemp := fi.CalcAvgTemp(report.Temperature)

	// Calculate target speed
	// 计算目标转速
	target, force := fi.CalcTargetSpeed(avgTemp)

	// When force update, ignore dead zone handling
	// 强制更新时, 忽略死区处理
	if !force {
		// 死区处理：当目标转速变化幅度小于死区阈值时，使用上一次的目标转速
		speedRange := fi.MaxSpeed - fi.MinSpeed
		deadZoneThreshold := float64(speedRange) * fi.deadZoneRatio
		// If speed change is less than dead zone threshold, use last target
		// 当温度变化小于死区范围时保持当前转速
		if math.Abs(float64(target-fi.lastTarget)) < deadZoneThreshold {
			target = fi.lastTarget
		}
	}

	// Smoothing filter: prevent sudden changes
	// 平滑处理: 防突变
	final := fi.CalcSmoothedSpeed(target)

	// Apply final speed
	// 应用转速调整
	if err := fi.ApplySpeed(final); err != nil {
		slog.Warnf(ero.AllTrace(err))
	}
}

// clamp restricts value within min/max boundaries
// 数值范围限制：将输入值限制在最小/最大范围内
func clamp(value, min, max int) int {
	return int(math.Max(float64(min), math.Min(float64(value), float64(max))))
}
