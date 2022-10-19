package service

import (
	"mime/multipart"
	"tservice-checker/internal/core"
	"tservice-checker/internal/repository"
	"tservice-checker/pkg"
)

// SessionService логика работы с сессиями телеграм
type TAccountService struct {
	repo    repository.TAccount
	tClient repository.TClient
}

// NewSessionService конструктор
func NewTAccountService(repo repository.TAccount, tClient repository.TClient) *TAccountService {
	return &TAccountService{repo: repo, tClient: tClient}
}

/* MultipartSave метод извлекает сессии из переданного слайса []*multipart.FileHeader. Работает только с архивами .zip
1. разархивирует архивы;
2. ищет в них tdata;
3. извлекает сессии;
4. валидирует сессии;
5. создает обьект tAccount, добавляет валидные сессии в поле sessions
6. сохраняет все данные обькта tAccount (включая валидные сессии) в базу;
7. возвращает итог*/
func (s *TAccountService) MultipartSave(files []*multipart.FileHeader) (*core.ExtractResult, error) {
	var sessionSVC = NewTSessionService(nil, s.tClient)
	var extR = core.ExtractResult{TotalFiles: len(files)}

	for _, file := range files {
		// * сохранить полученный архив
		filePath, err := pkg.SaveZip(file)
		if err != nil {
			pkg.ErrPrint("service", "internal error while save zip", err, file.Filename)
			continue
		}
		// * разархивирировать из архива директорию tdata
		tdataPath, err := pkg.UnzipSource(filePath, "files/unzip/")
		if err != nil {
			pkg.ErrPrint("service", "internal error while unzip", err, file.Filename)
			continue
		}
		// * вытащить из директории tdata сессии
		untrustSessions, err := sessionSVC.extractSession(tdataPath)
		if err != nil {
			pkg.ErrPrint("service", "internal error while extract session", err, file.Filename)
			continue
		}
		extR.TotalExtractedSessions += len(untrustSessions)
		// * проверяю живы ли сессии.
		for _, tSession := range untrustSessions {
			if err := sessionSVC.ValidateSession(&tSession); err != nil {
				pkg.ErrPrint("service", "internal error while validate session: ", tSession.SessionID, err, file.Filename)
				continue
			}

			// * получаю инфу об аккаунте на основе валидной сессии
			tAccount, err := sessionSVC.tClient.GetAccountInfo(&tSession)
			if err != nil {
				pkg.ErrPrint("service", "internal error while get account info: ", err, file.Filename)
				continue

			}

			//* сохраняю аккаунт, всю доп инфу и валидные сессии в базу
			if err := s.repo.Save(tAccount); err != nil {
				pkg.ErrPrint("service", "internal error while save account in database: ", err, file.Filename)
				continue
			}

			extR.TotalValidSessions++
		}
	}
	return &extR, nil
}

// func (s *TAccountService) Test(ctx context.Context, api *tg.Client) error {

// 	ph, err := api.AccountGetAuthorizations(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	for i, a := range ph.Authorizations {
// 		fmt.Printf("%v:  %+v \n\n\n", i, a)
// 	}

// 	// st, err := api.AccountResetAuthorization(ctx, 5903999754715481020)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// fmt.Printf("reset another auth? :  %+v \n\n\n", st)
// 	q := query.GetDialogs(api)

// 	collect, err := q.Collect(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("len ", len(collect))

// 	for _, elem := range collect {

// 		b, _ := json.MarshalIndent(elem, " ", "  ")
// 		fmt.Println(string(b))
// 		fmt.Println("*****************************************************")
// 	}
// 	// Return to close client connection and free up resources.

// 	// api.ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
// 	// 	Channel: c.channel.InputChannel(),
// 	// 	Filter:  c.filter,
// 	// 	Offset:  offset,
// 	// 	Limit:   limit,
// 	// })
// 	return nil
// }
