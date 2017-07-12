package alarm

import (
	"encoding/json"
	"fmt"

	"github.com/Cepave/open-falcon-backend/modules/alarm/model/event"
	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	alm "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/alarm"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
)

type InputExternalAlertsToAlarmInputs struct {
	event.ExternalEvent
}

func (mine InputExternalAlertsToAlarmInputs) CheckAlarmName() error {
	db := config.Con()
	atype := alm.AlarmTypes{}
	dt := db.Alarm.Model(atype).Where("name = ?", mine.AlarmType).Scan(&atype)
	if dt.Error != nil && dt.Error.Error() != "record not found" {
		return dt.Error
	}
	if atype.ID == 0 || dt.Error.Error() == "record not found" {
		return fmt.Errorf("alarm type: %v not found", mine.AlarmType)
	}
	return nil
}

func InputExternalAlertsToAlarm(c *gin.Context) {
	var inputs InputExternalAlertsToAlarmInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, fmt.Sprintf("binding input got error: %v", err.Error()))
		return
	}
	if err := inputs.CheckFormating(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckAlarmName(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	b, _ := json.Marshal(inputs)
	err := config.PutAlarmToRedis(string(b))
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	h.JSONR(c, fmt.Sprintf("insert 1 %v alarm succeeded", inputs.AlarmType))
	return
}
