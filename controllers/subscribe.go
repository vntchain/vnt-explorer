package controllers

import (
	"github.com/vntchain/vnt-explorer/models"
	"time"
)

type SubscribeController struct {
	BaseController
}

func (this *SubscribeController) Subscribe() {
	email := this.GetString("email")
	subscription := &models.Subscription{
		Email: email,
	}
	info, err := subscription.Get(email)
	if err == nil && info != nil {
		this.ReturnData(info, nil)
		return
	}

	now := time.Now()
	subscription.TimeStamp = uint64(now.Unix())
	if err = subscription.Insert(); err != nil {
		this.ReturnErrorMsg("Insert into db err ", err.Error(), "")
		return
	}

	this.ReturnData(subscription, nil)
}
