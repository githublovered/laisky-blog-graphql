package telegram

import (
	"context"
	"time"

	"github.com/Laisky/zap"

	"github.com/Laisky/go-utils"

	"github.com/pkg/errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Telegram client
type Telegram struct {
	bot *tb.Bot
}

// NewTelegram create new telegram client
func NewTelegram(ctx context.Context, token, api string) (*Telegram, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token: token,
		URL:   api,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "new telegram bot")
	}

	// start default handler
	bot.Handle(tb.OnText, func(m *tb.Message) {
		utils.Logger.Debug("got message", zap.String("msg", m.Text))
		if _, err = bot.Send(m.Sender, "NotImplement for "+m.Text); err != nil {
			utils.Logger.Error("send msg", zap.Error(err), zap.String("to", m.Sender.Username))
		}
	})

	tel := &Telegram{
		bot: bot,
	}
	go bot.Start()
	go tel.runDispatcher(ctx)
	go func() {
		<-ctx.Done()
		bot.Stop()
	}()
	return tel, nil
}

func (b *Telegram) runDispatcher(ctx context.Context) {

}
