package pool

import (
	"github.com/astaxie/beego"
	"fmt"
)

type BasicTask struct {
	Name	string
	Pool	*WorkPool
}

func (this *BasicTask) QueuedWork() int32 {
	return this.Pool.QueuedWork()
}

func (this *BasicTask) PreDoWork(workRoutine int) {
	qw := this.Pool.QueuedWork()
	ar := this.Pool.ActiveRoutines()
	beego.Info(fmt.Sprintf("*******> Task: %s WR: %d QW: %d AR: %d Total: %d\n",
		this.Name,
		workRoutine,
		qw,
		ar,
		qw + ar))
}