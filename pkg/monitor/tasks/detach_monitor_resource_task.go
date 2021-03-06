package tasks

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/monitor/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type DetachMonitorResourceJointTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(&DetachMonitorResourceJointTask{})
}

func (self *DetachMonitorResourceJointTask) OnInit(ctx context.Context, obj db.IStandaloneModel, body jsonutils.JSONObject) {
	alert := obj.(*models.SCommonAlert)
	err := alert.DetachMonitorResourceJoint(ctx, self.GetUserCred())
	if err != nil {
		msg := jsonutils.NewString(fmt.Sprintf("alert:%s DetachMonitorResourceJoint err:%v", alert.Name, err))
		self.taskFail(ctx, alert, msg)
		return
	}
	logclient.AddActionLogWithStartable(self, alert, logclient.ACT_DETACH_MONITOR_RESOURCE_JOINT, nil, self.UserCred, true)
	self.SetStageComplete(ctx, nil)
}

func (self *DetachMonitorResourceJointTask) taskFail(ctx context.Context, alert *models.SCommonAlert, msg jsonutils.JSONObject) {
	db.OpsLog.LogEvent(alert, db.ACT_DETACH_MONITOR_RESOURCE_JOINT, msg, self.GetUserCred())
	logclient.AddActionLogWithStartable(self, alert, logclient.ACT_DETACH_MONITOR_RESOURCE_JOINT, msg, self.UserCred, false)
	self.SetStageFailed(ctx, msg)
	return
}
