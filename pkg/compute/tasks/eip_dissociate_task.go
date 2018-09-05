package tasks


import (
	"fmt"
	"context"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
)

type EipDissociateTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(EipDissociateTask{})
}

func (self *EipDissociateTask) OnInit(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	eip := obj.(*models.SElasticip)

	extEip, err := eip.GetIEip()
	if err != nil {
		msg := fmt.Sprintf("fail to find iEIP for eip %s", err)
		eip.SetStatus(self.UserCred, models.EIP_STATUS_DISSOCIATE_FAIL, msg)
		self.SetStageFailed(ctx, msg)
		return
	}

	if len(extEip.GetAssociationExternalId()) > 0 {
		err = extEip.Dissociate()
		if err != nil {
			msg := fmt.Sprintf("fail to remote dissociate eip %s", err)
			eip.SetStatus(self.UserCred, models.EIP_STATUS_DISSOCIATE_FAIL, msg)
			self.SetStageFailed(ctx, msg)
			return
		}
	}

	err = eip.Dissociate(ctx, self.UserCred)
	if err != nil {
		msg := fmt.Sprintf("fail to local dissociate eip %s", err)
		eip.SetStatus(self.UserCred, models.EIP_STATUS_DISSOCIATE_FAIL, msg)
		self.SetStageFailed(ctx, msg)
		return
	}

	eip.SetStatus(self.UserCred, models.EIP_STATUS_READY, "dissociate")

	if eip.AutoDellocate.IsTrue() {
		self.SetStage("on_auto_dellocate_complete", nil)
		eip.StartEipDeallocateTask(ctx, self.UserCred, self.Id)
	} else {
		self.SetStageComplete(ctx, nil)
	}
}

func (self *EipDissociateTask) OnAutoDellocateComplete(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	self.SetStageComplete(ctx, nil)
}

