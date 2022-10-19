package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"tservice-checker/internal/core"
	"tservice-checker/internal/repository"
	"tservice-checker/pkg"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
)

// SessionService логика работы с сессиями телеграм
type SessionService struct {
	repo repository.Session
}

// NewSessionService конструктор
func NewSessionService(repo repository.Session) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Extract(files []*multipart.FileHeader) (*core.ExtractResult, error) {
	var extR = core.ExtractResult{}
	extR.TotalFiles = len(files)

	for _, file := range files {
		// * сохранить полученный архив
		filePath, err := s.saveZip(file)
		if err != nil {
			pkg.ErrPrint("service", "internal error while save zip", err, file.Filename)
			continue
		}
		// * разархивирировать из архива директорию tdata
		tdataPath, err := s.unZip(filePath)
		if err != nil {
			pkg.ErrPrint("service", "internal error while unzip", err, file.Filename)
			continue
		}
		// * вытащить из директории tdata сессию
		untrustSessions, err := s.extractSession(tdataPath)
		if err != nil {
			pkg.ErrPrint("service", "internal error while extract session", err, file.Filename)
			continue
		}
		extR.TotalExtractedSessions += len(untrustSessions)
		for _, untSess := range untrustSessions {
			// * проверяю жива ли сессия.
			_, err := s.validateSession(&untSess)
			if err != nil {
				pkg.ErrPrint("service", "internal error while validate session: ", untSess.SessionID, err, file.Filename)
				// if err := s.SaveSession(&sess); err != nil {
				// 	pkg.ErrPrint("service", "internal error while validate session: ", sess.Id, err, file.Filename)
				// }

				continue

			}
			extR.TotalValidSessions++
		}
	}

	return &extR, nil
}

// SaveZip сохраняю полученный архив
func (s *SessionService) saveZip(file *multipart.FileHeader) (string, error) {

	// Source
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err := os.MkdirAll("files/zip/", os.ModePerm); err != nil {
		return "", err
	}

	filePath := "files/zip/" + file.Filename
	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

// Unzip разархивирую из архива в директорию tdata
func (s *SessionService) unZip(src string) (string, error) {
	tdataPath, err := pkg.UnzipSource(src, "files/unzip/")
	if err != nil {
		return "", err
	}
	return tdataPath, nil
}

var (
	errPassCode = errors.New("wrong passcode")
)

// ExtractSession вытаскивую из директории tdata сессию
func (s *SessionService) extractSession(src string) ([]core.UntrustSession, error) {
	var sessionArray = []core.UntrustSession{}
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
		return nil, errPassCode
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

		sessionArray = append(sessionArray, core.UntrustSession{Data: data})
	}

	return sessionArray, nil
}

// ValidateSession проверяю жива ли сессия
func (s *SessionService) validateSession(sess *core.UntrustSession) (*core.TelegramAccount, error) {

	// st := &session.FileStorage{
	// 	Path: "/home/user/Documents/local/Projects/telegram-tdata-parser/session.json",
	// }
	dispatcher := tg.NewUpdateDispatcher()
	tClient := telegram.NewClient(9652426, "c7e1cd3c382656c433835e638965b334", telegram.Options{
		Logger:         nil,
		SessionStorage: sess,
		UpdateHandler:  dispatcher,
	})
	api := tClient.API()
	tAccount := core.NewTelegramAccount("siesta")

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

		ses, err := sess.LoadSession(ctx)
		if err != nil {
			return err
		}
		// * создаю новый телеграм акк
		tAccount.SetAttr(
			status.User.ID,
			status.User.Phone,
			status.User.Username,
			status.User.FirstName,
			status.User.LastName,
			ses,
		)
		tAccount.SetAddAttr(
			status.User.Bot,
			status.User.Fake,
			status.User.Scam,
			status.User.Premium,
			status.User.Support,
			status.User.Verified,
		)

		// if err := s.Test(ctx, api); err != nil {
		// 	return err
		// }
		return err

	}); err != nil {
		return nil, err
	}

	return tAccount, nil
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

func (s *SessionService) Test(ctx context.Context, api *tg.Client) error {

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
