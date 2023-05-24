package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hitshop/internal/core"
	"hitshop/internal/repository"
	"io/fs"
	"os"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

// SessionService логика работы с сессиями телеграм
type TSessionService struct {
	repo    repository.TSession
	tClient repository.TClient
}

// NewSessionService конструктор
func NewTSessionService(repo repository.TSession, tClient repository.TClient) *TSessionService {
	return &TSessionService{repo: repo, tClient: tClient}
}

// ExtractSession вытаскивую из директории tdata сессию
func (s *TSessionService) extractSession(src string) ([]core.Session, error) {
	var sessionArray = []core.Session{}
	dir, err := os.Stat(src)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, errors.Errorf("unable find tdata direcrory (path: %q)", src)
	case err != nil:
		return nil, err
	case !dir.IsDir():
		return nil, errors.Errorf("%q is not a directory", src)
	}

	accounts, err := tdesktop.Read(src, []byte("")) // todo add passcode
	switch {
	case errors.Is(err, tdesktop.ErrKeyInfoDecrypt):
		return nil, errors.New("wrong passcode")
	case err != nil:
		return nil, err
	}

	for _, account := range accounts {

		session, err := session.TDesktopSession(account)
		if err != nil {
			return nil, errors.Wrap(err, "convert")
		}

		a := core.SessionData{Version: 1, Data: *session}

		data, err := json.Marshal(a)
		if err != nil {
			return nil, err
		}

		sessionArray = append(sessionArray, core.Session{Data: data})
	}

	return sessionArray, nil
}

// ValidateSession проверяю жива ли сессия
func (s *TSessionService) ValidateSession(sess *core.Session) error {
	return s.tClient.ValidateTSession(sess)
}

// var codeAuthenticatorFunc auth.CodeAuthenticatorFunc = func(
// 	ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
// 	fmt.Print("Enter code: ")
// 	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(code), nil
// }

func (s *TSessionService) Test(ctx context.Context, api *tg.Client) error {

	ph, err := api.AccountGetAuthorizations(ctx)
	if err != nil {
		return err
	}

	for i, a := range ph.Authorizations {
		fmt.Printf("%v:  %+v \n\n\n", i, a)
	}

	// st, err := api.AccountResetAuthorization(ctx, 5903999754715481020)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("reset another auth? :  %+v \n\n\n", st)
	q := query.GetDialogs(api)

	collect, err := q.Collect(ctx)
	if err != nil {
		return err
	}
	fmt.Println("len ", len(collect))

	for _, elem := range collect {

		b, _ := json.MarshalIndent(elem, " ", "  ")
		fmt.Println(string(b))
		fmt.Println("*****************************************************")
	}
	// Return to close client connection and free up resources.

	// api.ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
	// 	Channel: c.channel.InputChannel(),
	// 	Filter:  c.filter,
	// 	Offset:  offset,
	// 	Limit:   limit,
	// })
	return nil
}
