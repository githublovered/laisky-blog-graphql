package telegram

import (
	"fmt"
	"strings"

	"github.com/Laisky/go-utils"
	"github.com/Laisky/zap"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Telegram) monitorHandler() {
	b.bot.Handle("/monitor", func(c *tb.Message) {
		b.userStats.Store(c.Sender.ID, &userStat{
			user:  c.Sender,
			state: userWaitChooseMonitorCmd,
			lastT: utils.Clock.GetUTCNow(),
		})

		if _, err := b.bot.Send(c.Sender, `
Reply number:

	1 - new monitor's name  # reply "1 - demo"
	2 - list all joint alerts  # reply "2"
`); err != nil {
			utils.Logger.Error("reply msg", zap.Error(err))
		}
	})
}

func (b *Telegram) chooseMonitor(us *userStat, msg *tb.Message) {
	utils.Logger.Debug("choose monitor",
		zap.String("user", us.user.Username),
		zap.String("msg", msg.Text))
	defer b.userStats.Delete(us.user.ID)
	var (
		err error
		ans = []string{msg.Text}
	)
	if strings.Contains(msg.Text, " - ") {
		ans = strings.SplitN(msg.Text, " - ", 2)
		if len(ans) < 2 {
			b.PleaseRetry(us.user, msg.Text)
			return
		}
	}

	switch ans[0] {
	case "1": // create new monitor
		if err = b.createNewMonitor(us, ans[1]); err != nil {
			utils.Logger.Error("createNewMonitor", zap.Error(err))
			b.bot.Send(us.user, "[Error] "+err.Error())
		}
	case "2":
		if err = b.listAllMonitorAlerts(us); err != nil {
			utils.Logger.Error("listAllMonitorAlerts", zap.Error(err))
			b.bot.Send(us.user, "[Error] "+err.Error())
		}
	default:
		b.PleaseRetry(us.user, msg.Text)
	}
}

func (b *Telegram) listAllMonitorAlerts(us *userStat) (err error) {
	u, err := b.db.LoadUserByUID(us.user.ID)
	if err != nil {
		return err
	}
	alerts, err := b.db.LoadAlertTypesByUser(u)
	if err != nil {
		return err
	}

	msg := ""
	for _, alert := range alerts {
		msg += "--------------------------------\n"
		msg += "alert_type: " + alert.Name + "\n"
		msg += "push_token: " + alert.PushToken + "\n"
		msg += "join_key: " + alert.JoinKey + "\n"
	}
	msg += "--------------------------------"

	return b.SendMsgToUser(u.UID, msg)
}

func (b *Telegram) createNewMonitor(us *userStat, alertName string) (err error) {
	u, err := b.db.CreateOrGetUser(us.user)
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	a, err := b.db.CreateAlertType(alertName)
	if err != nil {
		return errors.Wrap(err, "create alert_type")
	}

	_, err = b.db.CreateUserAlertRelations(u, a)
	if err != nil {
		return errors.Wrap(err, "create user_alert_relation")
	}

	if _, err = b.bot.Send(us.user, fmt.Sprintf(`
create user & alert_type & user_alert_relations successed!
user: %v
alert_type: %v
join_key: %v
push_token: %v
	`, u.Name,
		a.Name,
		a.JoinKey,
		a.PushToken)); err != nil {
		return errors.Wrap(err, "send msg")
	}

	return nil
}
