

package common


import "math"
import "sync/atomic"




type StatMetric struct {
	
	MetricSource *uint64
	DividerSource *uint64
	
	ValueDelta bool
	SpeedDelta bool
	
	ValueThreshold float64
	SpeedThreshold float64
	
	MetricScale float64
	DividerScale float64
	ValueScale float64
	SpeedScale float64
	
	Changed bool
	Invalid bool
	Touched bool
	
	metricLast_0 uint64
	dividerLast_0 uint64
	
	MetricLast float64
	DividerLast float64
	ValueLast float64
	speed0Last float64
	
	TimeDelta uint64
	TimeLast uint64
	TimeFirst uint64
	TimeChanged uint64
	UpdateCount uint64
	ChangedCount uint64
	
	Speed1Last float64
	Speed1Window float64
	Speed1pLast float64
	Speed1paLast float64
	Speed1prLast float64
	Speed1pWindow float64
	Speed1paWindow float64
	Speed1prWindow float64
	
	Speed2Last float64
	Speed2Window float64
	Speed2pLast float64
	Speed2paLast float64
	Speed2prLast float64
	Speed2pWindow float64
	Speed2paWindow float64
	Speed2prWindow float64
	
	WindowSize uint64
}




func (_stat *StatMetric) Update2 (_timeNanoseconds uint64, _changed *bool, _invalid *bool) () {
	_stat.Update (_timeNanoseconds)
	*_changed = *_changed || _stat.Changed
	*_invalid = *_invalid || _stat.Invalid
}




func (_stat *StatMetric) Update (_timeNanoseconds uint64) () {
	
	_invalid := false
	
	_timeNow := _timeNanoseconds
	
	_metricNow_0 := atomic.LoadUint64 (_stat.MetricSource)
	_metricNow := float64 (_metricNow_0)
	if _stat.MetricScale != 0 {
		_metricNow = _metricNow / _stat.MetricScale
	}
	
	_dividerNow_0 := uint64 (0)
	if _stat.DividerSource != nil {
		_dividerNow_0 = atomic.LoadUint64 (_stat.DividerSource)
	}
	_dividerNow := float64 (_dividerNow_0)
	if _stat.DividerScale != 0 {
		_dividerNow = _dividerNow / _stat.DividerScale
	}
	
	_timeDelta := int64 (_timeNow) - int64 (_stat.TimeLast)
	_timeDeltaSec := float64 (_timeDelta) / 1000000000
	
	if _timeDelta <= 0 {
		_invalid = true
	}
	
	_valueNow := _metricNow
	_speed0Now := _metricNow
	
	if _stat.ValueDelta {
		_valueNow = _valueNow - _stat.MetricLast
	}
	if _stat.SpeedDelta {
		_speed0Now = _speed0Now - _stat.MetricLast
	}
	
	if _stat.DividerSource != nil {
		
		if _stat.ValueDelta {
			_dividerDelta := _dividerNow - _stat.DividerLast
			if _dividerDelta > 0 {
				_valueNow = _valueNow / _dividerDelta
			} else {
				_invalid = true
			}
		} else {
			if _dividerNow > 0 {
				_valueNow = _valueNow / _dividerNow
			} else {
				_invalid = true
			}
		}
		
		if _stat.SpeedDelta {
			_dividerDelta := _dividerNow - _stat.DividerLast
			if _dividerDelta > 0 {
				_speed0Now = _speed0Now / _dividerDelta
			} else {
				_invalid = true
			}
		} else {
			if _dividerNow > 0 {
				_speed0Now = _speed0Now / _dividerNow
			} else {
				_invalid = true
			}
		}
	}
	
	if _stat.ValueScale != 0 {
		_valueNow = _valueNow / _stat.ValueScale
	}
	if _stat.SpeedScale != 0 {
		_speed0Now = _speed0Now / _stat.SpeedScale
	}
	
	if _stat.UpdateCount == 0 {
		_stat.TimeFirst = _timeNow
		_invalid = true
	}
	
	_stat.Changed = false
	_thresholdUsed := false
	_thresholdMatched := false
	
	if _invalid {
		goto _return
	}
	
	_stat.Changed = (_stat.metricLast_0 != _metricNow_0) || (_stat.dividerLast_0 != _dividerNow_0)
	
	if _stat.Changed {
		
		_speed1Now := (_speed0Now - _stat.speed0Last) / _timeDeltaSec
		_speed1pNow := _speed1Now / _stat.Speed1Last - 1
		
		_speed2Now := (_speed1Now - _stat.Speed1Last) / _timeDeltaSec
		_speed2pNow := _speed2Now / _stat.Speed2Last - 1
		
		_stat.Speed1Last = _speed1Now
		_stat.Speed2Last = _speed2Now
		
		_windowNew := float64 (0)
		const _windowSizeNormal = 12
		if _stat.WindowSize >= _windowSizeNormal {
			_windowNew = 1 / float64 (_windowSizeNormal + 1)
		} else {
			_windowNew = 1 / float64 (_stat.WindowSize + 1)
		}
		_windowOld := 1 - _windowNew
		
		if ((_timeNow - _stat.TimeChanged) / 1000000000 > 6) || (_stat.TimeChanged == 0) || (_stat.WindowSize == 0) {
			_stat.Speed1pLast = 0
			_stat.Speed1Window = _speed1Now
			_stat.Speed1pWindow = 0
			_stat.Speed2pLast = 0
			_stat.Speed2Window = _speed2Now
			_stat.Speed2pWindow = 0
		} else {
			_stat.Speed1pLast = _speed1pNow
			_stat.Speed1Window = _stat.Speed1Window * _windowOld + _speed1Now * _windowNew
			_stat.Speed1pWindow = _stat.Speed1pWindow * _windowOld + _speed1pNow * _windowNew
			_stat.Speed2pLast = _speed2pNow
			_stat.Speed2Window = _stat.Speed2Window * _windowOld + _speed2Now * _windowNew
			_stat.Speed2pWindow = _stat.Speed2pWindow * _windowOld + _speed2pNow * _windowNew
		}
		
		_stat.TimeChanged = _timeNow
		_stat.ChangedCount += 1
		_stat.WindowSize++
		
	} else {
		
		if ((_timeNow - _stat.TimeChanged) / 1000000000 > 6) || (_stat.TimeChanged == 0) {
			_stat.Speed1Last = 0.0
			_stat.Speed1Window = 0.0
			_stat.Speed1pWindow = 0.0
			_stat.Speed1pLast = 0.0
			_stat.Speed1paLast = 0.0
			_stat.Speed2Last = 0.0
			_stat.Speed2Window = 0.0
			_stat.Speed2pWindow = 0.0
			_stat.Speed2pLast = 0.0
			_stat.Speed2paLast = 0.0
			_stat.WindowSize = 0
		}
	}
	
	if _stat.ValueThreshold != 0 {
		_thresholdUsed = true
		_thresholdMatched = _thresholdMatched || (math.Abs (_valueNow) >= _stat.ValueThreshold)
	}
	if _stat.SpeedThreshold != 0 {
		_thresholdUsed = true
		_thresholdMatched = _thresholdMatched || (math.Abs (_stat.Speed1Last) >= _stat.SpeedThreshold)
	}
	if _thresholdUsed && !_thresholdMatched {
		_stat.Changed = false
	}
	
	_return :
	
	_stat.metricLast_0 = _metricNow_0
	_stat.dividerLast_0 = _dividerNow_0
	
	_stat.MetricLast = _metricNow
	_stat.DividerLast = _dividerNow
	
	_stat.ValueLast = _valueNow
	_stat.speed0Last = _speed0Now
	
	_stat.TimeLast = _timeNow
	_stat.TimeDelta = uint64 (_timeDelta)
	
	_stat.UpdateCount += 1
	
	_stat.Invalid = _invalid
	_stat.Touched = (_stat.ChangedCount > 0)
	
	
	_stat.Speed1paLast = math.Abs (_stat.Speed1pLast)
	_stat.Speed1paWindow = math.Abs (_stat.Speed1pWindow)
	_stat.Speed2paLast = math.Abs (_stat.Speed2pLast)
	_stat.Speed1paWindow = math.Abs (_stat.Speed1paWindow)
	
	_stat.Speed1prLast = math.Round (_stat.Speed1pLast * 100 * 100) / 100
	_stat.Speed1prWindow = math.Round (_stat.Speed1pWindow * 100 * 100) / 100
	_stat.Speed2prLast = math.Round (_stat.Speed2pLast * 100 * 100) / 100
	_stat.Speed1prWindow = math.Round (_stat.Speed1paWindow * 100 * 100) / 100
	
	_infinite := float64 (1.0)
	_infinite = _infinite / 0
	if _stat.Speed1paLast > 1000 {
		_stat.Speed1prLast = math.Copysign (_infinite, _stat.Speed1prLast)
	}
	if _stat.Speed1paWindow > 1000 {
		_stat.Speed1prWindow = math.Copysign (_infinite, _stat.Speed1prWindow)
	}
	if _stat.Speed2paLast > 1000 {
		_stat.Speed2prLast = math.Copysign (_infinite, _stat.Speed2prLast)
	}
	if _stat.Speed2paWindow > 1000 {
		_stat.Speed2prWindow = math.Copysign (_infinite, _stat.Speed2prWindow)
	}
}

