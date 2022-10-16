package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"tservice-checker/internal/core"
	"tservice-checker/internal/repository"
	"tservice-checker/pkg"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
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

// SaveZip сохраняю полученный архив
func (s *SessionService) SaveZip(file *multipart.FileHeader) (string, error) {

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
func (s *SessionService) Unzip(src string) (string, error) {
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
func (s *SessionService) ExtractSession(src string) ([]core.Session, error) {
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

		sessionArray = append(sessionArray, core.Session{Data: data})
	}

	return sessionArray, nil
}

// ValidateSession проверяю жива ли сессия, сохраняю ее в базе
func (s *SessionService) ValidateSession(sess *core.Session) error {

	// st := &session.FileStorage{
	// 	Path: "/home/user/Documents/local/Projects/telegram-tdata-parser/session.json",
	// }

	client := telegram.NewClient(9652426, "c7e1cd3c382656c433835e638965b334", telegram.Options{
		Logger:         nil,
		SessionStorage: sess,
	})

	if err := client.Run(context.Background(), func(ctx context.Context) error {
		// flow := auth.NewFlow(
		// 	auth.CodeOnly("+79937001034", codeAuthenticatorFunc),
		// 	auth.SendCodeOptions{},
		// )
		// if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		// 	panic(err)
		// }
		status, err := client.Auth().Status(ctx)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!", status.User.Phone)

		switch {
		case status.Authorized:
			ses, err := sess.LoadSession(ctx)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("======================================")
			fmt.Println(string(ses))
			fmt.Println("======================================")

			fmt.Println("AUTHENTIFICATED!!!!!!!!!!!!!!!!!!!!!")
			break
		case !status.Authorized:
			fmt.Println("NOT AUTHENTIFICATED!!!!!!!!!!!!!!!!!!!!!")
			break
		}
		api := client.API()
		_, err = api.AccountGetPassword(ctx)
		if err != nil {
			return err
		}

		ph, err := api.AccountGetAuthorizations(ctx)
		if err != nil {
			return err
		}

		for i, a := range ph.Authorizations {

			fmt.Printf("%v:  %+v \n\n\n", i, a)
		}

		st, err := api.AccountResetAuthorization(ctx, 5903999754715481020)
		if err != nil {
			return err
		}
		fmt.Printf("reset another auth? :  %+v \n\n\n", st)

		api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
			Message: "Hi",
		})
		// Return to close client connection and free up resources.
		return nil
	}); err != nil {
		return err
	}

	return nil
}

var codeAuthenticatorFunc auth.CodeAuthenticatorFunc = func(
	ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
