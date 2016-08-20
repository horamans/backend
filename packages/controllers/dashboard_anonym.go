package controllers

import (
	"github.com/DayLightProject/go-daylight/packages/utils"
)


type dashboardAnonymPage struct {
	Lang                  map[string]string
	Title                 string
	CountSign             int
	CountSignArr          []int
	BlockExplorer         []map[string]string
	SignData              string
	ShowSignData          bool
}

func (c *Controller) DashboardAnonym() (string, error) {

	blockExplorer,err := c.GetAll("SELECT hash, cb_id, wallet_id, time, tx, id FROM block_chain order by id desc limit 0, 30",-1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for ind := range blockExplorer {
		blockExplorer[ind][`hash`] = string(utils.BinToHex([]byte(blockExplorer[ind][`hash`])))
	}
		
	TemplateStr, err := makeTemplate("dashboard_anonym", "dashboardAnonym", &dashboardAnonymPage{
		CountSignArr:          c.CountSignArr,
		CountSign:             c.CountSign,
		Lang:                  c.Lang,
		BlockExplorer:         blockExplorer,
		Title:                 "Home",
		ShowSignData:          c.ShowSignData,
		SignData:              ""})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
