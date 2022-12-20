package alarmx

// AlarmsContext 链路
type AlarmsContext interface {
	Alarm(name string)
}

func (ctx *AlarmContext) Alarm(name string) {}

// AlarmContext alarm
type AlarmContext struct {
}
