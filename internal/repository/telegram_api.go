package repository

import (
	"context"
	"errors"
	"tservice-checker/internal/config"
	"tservice-checker/internal/core"

	"github.com/gotd/td/telegram"
)

//TClientAPI
type TClientAPI struct {
	*config.Cfg
}

//NewAuthPostgres конструктор
func NewTClientAPI(cfg *config.Cfg) *TClientAPI {
	return &TClientAPI{cfg}
}

//GetTClient создаю клиента для работы с telegram api
func (c *TClientAPI) initClient(sess *core.Session) *telegram.Client {
	return telegram.NewClient(c.AppID, c.AppHash, telegram.Options{
		SessionStorage: sess,
	})
}

//ValidateTSession валидация сессии
func (c *TClientAPI) ValidateTSession(tSession *core.Session) error {
	// st := &session.FileStorage{
	// 	Path: "/home/user/Documents/local/Projects/telegram-tdata-parser/session.json",
	// }
	tClient := c.initClient(tSession)
	api := tClient.API()

	if err := tClient.Run(context.Background(), func(ctx context.Context) error {
		// flow := auth.NewFlow(
		// 	auth.CodeOnly("+79937001034", codeAuthenticatorFunc),
		// 	auth.SendCodeOptions{},
		// )
		// if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		// 	panic(err)
		// }
		status, err := tClient.Auth().Status(ctx)
		if err != nil {
			return err
		}
		if !status.Authorized {
			return errors.New("client not authentificated. drop this session")
		}

		// * доп проверка валидности сессии
		_, err = api.AccountGetPassword(ctx)
		if err != nil {
			return err
		}

		// if err := s.Test(ctx, api); err != nil {
		// 	return err
		// }
		return nil

	}); err != nil {
		return err
	}
	return nil
}

//AccountInfo Получение информации об аккаунте на основе сессии
func (c *TClientAPI) GetAccountInfo(tSession *core.Session) (*core.TelegramAccount, error) {
	tAccount := core.NewTelegramAccount("siesta")
	tClient := c.initClient(tSession)

	if err := tClient.Run(context.Background(), func(ctx context.Context) error {
		status, err := tClient.Auth().Status(ctx)
		if err != nil {
			return err
		}
		if !status.Authorized {
			return errors.New("unable get account info. Session not authorized.")
		}

		// * создаю новый телеграм акк

		tAccount.AccountID = status.User.ID
		tAccount.Phone = status.User.Phone
		tAccount.UserName = status.User.Username
		tAccount.FirstName = status.User.FirstName
		tAccount.LastName = status.User.LastName
		tAccount.Bot = status.User.Bot
		tAccount.Fake = status.User.Fake
		tAccount.Scam = status.User.Scam
		tAccount.Premium = status.User.Premium
		tAccount.Support = status.User.Support
		tAccount.Verified = status.User.Verified
		tAccount.Sessions = append(tAccount.Sessions, core.Session{Data: tSession.Data})

		return nil

	}); err != nil {
		return nil, err
	}
	return tAccount, nil
}
