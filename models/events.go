package models

import (
	"code.cloudfoundry.org/bbs/format"
	"github.com/gogo/protobuf/proto"
)

type Event interface {
	EventType() string
	Key() string
	proto.Message
}

const (
	EventTypeInvalid = ""

	EventTypeDesiredLRPCreated = "desired_lrp_created"
	EventTypeDesiredLRPChanged = "desired_lrp_changed"
	EventTypeDesiredLRPRemoved = "desired_lrp_removed"

	EventTypeActualLRPCreated = "actual_lrp_created" // DEPRECATED
	EventTypeActualLRPChanged = "actual_lrp_changed" // DEPRECATED
	EventTypeActualLRPRemoved = "actual_lrp_removed" // DEPRECATED
	EventTypeActualLRPCrashed = "actual_lrp_crashed"

	EventTypeActualLRPInstanceCreated = "actual_lrp_instance_created"
	EventTypeActualLRPInstanceChanged = "actual_lrp_instance_changed"
	EventTypeActualLRPInstanceRemoved = "actual_lrp_instance_removed"

	EventTypeTaskCreated = "task_created"
	EventTypeTaskChanged = "task_changed"
	EventTypeTaskRemoved = "task_removed"
)

// Downgrade the DesiredLRPEvent payload (i.e. DesiredLRP(s)) to the given
// target version
func VersionDesiredLRPsTo(event Event, target format.Version) Event {
	switch event := event.(type) {
	case *DesiredLRPCreatedEvent:
		return NewDesiredLRPCreatedEvent(event.DesiredLrp.VersionDownTo(target))
	case *DesiredLRPRemovedEvent:
		return NewDesiredLRPRemovedEvent(event.DesiredLrp.VersionDownTo(target))
	case *DesiredLRPChangedEvent:
		return NewDesiredLRPChangedEvent(
			event.Before.VersionDownTo(target),
			event.After.VersionDownTo(target),
		)
	default:
		return event
	}
}

func VersionMiniDesiredLRPsTo(event Event, target format.Version) Event {
	switch event := event.(type) {
	case *MiniDesiredLRPCreatedEvent:
		return NewMiniDesiredLRPCreatedEvent(event.MinidesiredLrp.VersionDownTov2(target))
	case *MiniDesiredLRPRemovedEvent:
		return NewMiniDesiredLRPRemovedEvent(event.MinidesiredLrp.VersionDownTov2(target))
	case *MiniDesiredLRPChangedEvent:
		return NewMiniDesiredLRPChangedEvent(
			event.Before.VersionDownTov2(target),
			event.After.VersionDownTov2(target),
		)
	default:
		return event
	}
}

// Downgrade the TaskEvent payload (i.e. Task(s)) to the given target version
func VersionTaskDefinitionsTo(event Event, target format.Version) Event {
	switch event := event.(type) {
	case *TaskCreatedEvent:
		return NewTaskCreatedEvent(event.Task.VersionDownTo(target))
	case *TaskRemovedEvent:
		return NewTaskRemovedEvent(event.Task.VersionDownTo(target))
	case *TaskChangedEvent:
		return NewTaskChangedEvent(event.Before.VersionDownTo(target), event.After.VersionDownTo(target))
	default:
		return event
	}
}

func NewDesiredLRPCreatedEvent(desiredLRP *DesiredLRP) *DesiredLRPCreatedEvent {
	return &DesiredLRPCreatedEvent{
		DesiredLrp: desiredLRP,
	}
}

func (event *DesiredLRPCreatedEvent) EventType() string {
	return EventTypeDesiredLRPCreated
}

func (event *DesiredLRPCreatedEvent) Key() string {
	return event.DesiredLrp.GetProcessGuid()
}

func NewMiniDesiredLRPCreatedEvent(minidesiredLRP *MiniDesiredLRP) *MiniDesiredLRPCreatedEvent {
	return &MiniDesiredLRPCreatedEvent{
		MinidesiredLrp: minidesiredLRP,
	}
}

func (event *MiniDesiredLRPCreatedEvent) EventType() string {
	return EventTypeDesiredLRPCreated
}

func (event *MiniDesiredLRPCreatedEvent) Key() string {
	return event.MinidesiredLrp.GetProcessGuid()
}

func NewDesiredLRPChangedEvent(before, after *DesiredLRP) *DesiredLRPChangedEvent {
	return &DesiredLRPChangedEvent{
		Before: before,
		After:  after,
	}
}

func (event *DesiredLRPChangedEvent) EventType() string {
	return EventTypeDesiredLRPChanged
}

func (event *DesiredLRPChangedEvent) Key() string {
	return event.Before.GetProcessGuid()
}

func NewMiniDesiredLRPChangedEvent(before, after *MiniDesiredLRP) *MiniDesiredLRPChangedEvent {
	return &MiniDesiredLRPChangedEvent{
		Before: before,
		After:  after,
	}
}

func (event *MiniDesiredLRPChangedEvent) EventType() string {
	return EventTypeDesiredLRPChanged
}

func (event *MiniDesiredLRPChangedEvent) Key() string {
	return event.Before.GetProcessGuid()
}

func NewMiniDesiredLRPRemovedEvent(minidesiredLRP *MiniDesiredLRP) *MiniDesiredLRPRemovedEvent {
	return &MiniDesiredLRPRemovedEvent{
		MinidesiredLrp: minidesiredLRP,
	}
}

func (event *MiniDesiredLRPRemovedEvent) EventType() string {
	return EventTypeDesiredLRPRemoved
}

func (event MiniDesiredLRPRemovedEvent) Key() string {
	return event.MinidesiredLrp.GetProcessGuid()
}

func NewDesiredLRPRemovedEvent(desiredLRP *DesiredLRP) *DesiredLRPRemovedEvent {
	return &DesiredLRPRemovedEvent{
		DesiredLrp: desiredLRP,
	}
}

func (event *DesiredLRPRemovedEvent) EventType() string {
	return EventTypeDesiredLRPRemoved
}

func (event DesiredLRPRemovedEvent) Key() string {
	return event.DesiredLrp.GetProcessGuid()
}

// FIXME: change the signature
func NewActualLRPInstanceChangedEvent(before, after *ActualLRP) *ActualLRPInstanceChangedEvent {
	var (
		actualLRPKey         ActualLRPKey
		actualLRPInstanceKey ActualLRPInstanceKey
	)

	if (before != nil && before.ActualLRPKey != ActualLRPKey{}) {
		actualLRPKey = before.ActualLRPKey
	}
	if (after != nil && after.ActualLRPKey != ActualLRPKey{}) {
		actualLRPKey = after.ActualLRPKey
	}

	if (before != nil && before.ActualLRPInstanceKey != ActualLRPInstanceKey{}) {
		actualLRPInstanceKey = before.ActualLRPInstanceKey
	}
	if (after != nil && after.ActualLRPInstanceKey != ActualLRPInstanceKey{}) {
		actualLRPInstanceKey = after.ActualLRPInstanceKey
	}

	return &ActualLRPInstanceChangedEvent{
		ActualLRPKey:         actualLRPKey,
		ActualLRPInstanceKey: actualLRPInstanceKey,
		Before:               before.ToActualLRPInfo(),
		After:                after.ToActualLRPInfo(),
	}
}

func (event *ActualLRPInstanceChangedEvent) EventType() string {
	return EventTypeActualLRPInstanceChanged
}

func (event *ActualLRPInstanceChangedEvent) Key() string {
	return event.GetInstanceGuid()
}

// DEPRECATED
func NewActualLRPChangedEvent(before, after *ActualLRPGroup) *ActualLRPChangedEvent {
	return &ActualLRPChangedEvent{
		Before: before,
		After:  after,
	}
}

// DEPRECATED
func (event *ActualLRPChangedEvent) EventType() string {
	return EventTypeActualLRPChanged
}

// DEPRECATED
func (event *ActualLRPChangedEvent) Key() string {
	actualLRP, _, resolveError := event.Before.Resolve()
	if resolveError != nil {
		return ""
	}
	return actualLRP.GetInstanceGuid()
}

func NewActualLRPCrashedEvent(before, after *ActualLRP) *ActualLRPCrashedEvent {
	return &ActualLRPCrashedEvent{
		ActualLRPKey:         after.ActualLRPKey,
		ActualLRPInstanceKey: before.ActualLRPInstanceKey,
		CrashCount:           after.CrashCount,
		CrashReason:          after.CrashReason,
		Since:                after.Since,
	}
}

func (event *ActualLRPCrashedEvent) EventType() string {
	return EventTypeActualLRPCrashed
}

func (event *ActualLRPCrashedEvent) Key() string {
	return event.ActualLRPInstanceKey.InstanceGuid
}

// DEPRECATED
func NewActualLRPRemovedEvent(actualLRPGroup *ActualLRPGroup) *ActualLRPRemovedEvent {
	return &ActualLRPRemovedEvent{
		ActualLrpGroup: actualLRPGroup,
	}
}

// DEPRECATED
func (event *ActualLRPRemovedEvent) EventType() string {
	return EventTypeActualLRPRemoved
}

// DEPRECATED
func (event *ActualLRPRemovedEvent) Key() string {
	actualLRP, _, resolveError := event.ActualLrpGroup.Resolve()
	if resolveError != nil {
		return ""
	}
	return actualLRP.GetInstanceGuid()
}

func NewActualLRPInstanceRemovedEvent(actualLrp *ActualLRP) *ActualLRPInstanceRemovedEvent {
	return &ActualLRPInstanceRemovedEvent{
		ActualLrp: actualLrp,
	}
}

func (event *ActualLRPInstanceRemovedEvent) EventType() string {
	return EventTypeActualLRPInstanceRemoved
}

func (event *ActualLRPInstanceRemovedEvent) Key() string {
	if event.ActualLrp == nil {
		return ""
	}
	return event.ActualLrp.GetInstanceGuid()
}

// DEPRECATED
func NewActualLRPCreatedEvent(actualLRPGroup *ActualLRPGroup) *ActualLRPCreatedEvent {
	return &ActualLRPCreatedEvent{
		ActualLrpGroup: actualLRPGroup,
	}
}

// DEPRECATED
func (event *ActualLRPCreatedEvent) EventType() string {
	return EventTypeActualLRPCreated
}

// DEPRECATED
func (event *ActualLRPCreatedEvent) Key() string {
	actualLRP, _, resolveError := event.ActualLrpGroup.Resolve()
	if resolveError != nil {
		return ""
	}
	return actualLRP.GetInstanceGuid()
}

func NewActualLRPInstanceCreatedEvent(actualLrp *ActualLRP) *ActualLRPInstanceCreatedEvent {
	return &ActualLRPInstanceCreatedEvent{
		ActualLrp: actualLrp,
	}
}

func (event *ActualLRPInstanceCreatedEvent) EventType() string {
	return EventTypeActualLRPInstanceCreated
}

func (event *ActualLRPInstanceCreatedEvent) Key() string {
	if event.ActualLrp == nil {
		return ""
	}
	return event.ActualLrp.GetInstanceGuid()
}

func (request *EventsByCellId) Validate() error {
	return nil
}

func NewTaskCreatedEvent(task *Task) *TaskCreatedEvent {
	return &TaskCreatedEvent{
		Task: task,
	}
}

func (event *TaskCreatedEvent) EventType() string {
	return EventTypeTaskCreated
}

func (event *TaskCreatedEvent) Key() string {
	return event.Task.GetTaskGuid()
}

func NewTaskChangedEvent(before, after *Task) *TaskChangedEvent {
	return &TaskChangedEvent{
		Before: before,
		After:  after,
	}
}

func (event *TaskChangedEvent) EventType() string {
	return EventTypeTaskChanged
}

func (event *TaskChangedEvent) Key() string {
	return event.Before.GetTaskGuid()
}

func NewTaskRemovedEvent(task *Task) *TaskRemovedEvent {
	return &TaskRemovedEvent{
		Task: task,
	}
}

func (event *TaskRemovedEvent) EventType() string {
	return EventTypeTaskRemoved
}

func (event TaskRemovedEvent) Key() string {
	return event.Task.GetTaskGuid()
}
