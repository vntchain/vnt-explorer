package controllers

import (
	"fmt"
	"github.com/vntchain/vnt-explorer/models"
	"time"
)

type SubscribeController struct {
	BaseController
}

func (this *SubscribeController) Subscribe() {
	email := this.GetString("email")
	now := time.Now()

	subscription := &models.Subscription{
		Email:     email,
		TimeStamp: uint64(now.Unix()),
	}

	subscription.TimeStamp = uint64(now.Unix())
	if insertErr := subscription.Insert(); insertErr != nil {
		info, err := subscription.Get(email)
		if err == nil {
			// Duplicate entry, return success
			if info != nil {
				this.ReturnData(info, nil)
				return
			} else {
				this.ReturnErrorMsg("Insert into db err: ", insertErr.Error(), "")
				return
			}
		} else {
			msg := fmt.Sprintf("Insert into db err: %s, and get from db err: %s", insertErr, err)
			this.ReturnErrorMsg(msg, "", "")
			return
		}
	}
	this.ReturnData(subscription, nil)
}
