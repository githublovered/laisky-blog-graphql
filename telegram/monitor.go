package telegram

import (
	"strings"

	"github.com/Laisky/go-utils"
	"github.com/Laisky/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Telegram) monitorHandler() {
	b.bot.Handle("/monitor", func(c *tb.Message) {
		b.userStats.Store(c.Sender, &userStat{
			user:  c.Sender,
			state: userWaitChooseMonitorCmd,
			lastT: utils.Clock.GetUTCNow(),
		})

		if _, err := b.bot.Send(c.Sender, `
Reply number:

	1 - new monitor's name
`); err != nil {
			utils.Logger.Error("reply msg", zap.Error(err))
		}
	})
}

func (b *Telegram) chooseMonitor(us *userStat, msg *tb.Message) {
	ans := strings.SplitN(msg.Text, " - ", 1)
	if len(ans) < 2 {
		b.PleaseRetry(us.user, msg.Text)
	}

	switch ans[0] {
	case "1": // create new monitor
		b.createNewMonitor(us)
	default:
		b.PleaseRetry(us.user, msg.Text)
	}
}

func (b *Telegram) createNewMonitor(us *userStat) {
	u, err := b.db.CreateOrGetUser(us.user)
	if err != nil {
		utils.Logger.Error("create or get user", zap.Error(err), zap.Int("uid", us.user.ID))
	}
}
