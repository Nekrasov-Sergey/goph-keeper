package client

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

func selectMenu(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func promptString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if input == "" {
				return errors.New("поле не может быть пустым")
			}
			return nil
		},
	}
	return prompt.Run()
}

func promptPassword(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
		Validate: func(input string) error {
			if input == "" {
				return errors.New("поле не может быть пустым")
			}
			return nil
		},
	}
	return prompt.Run()
}

func promptCredentials() (types.LoginPassword, error) {
	login, err := promptString("Логин")
	if err != nil {
		return types.LoginPassword{}, err
	}

	password, err := promptPassword("Пароль")
	if err != nil {
		return types.LoginPassword{}, err
	}

	return types.LoginPassword{
		Login:    login,
		Password: password,
	}, nil
}

func promptSecretID() (int64, error) {
	var parsed int64

	prompt := promptui.Prompt{
		Label: "ID секрета",
		Validate: func(input string) error {
			id, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return errors.New("некорректный ID")
			}
			parsed = id
			return nil
		},
	}

	if _, err := prompt.Run(); err != nil {
		return 0, err
	}

	return parsed, nil
}

func promptCreateSecret() (*types.SecretPayload, error) {
	items := []struct {
		label string
		typ   types.SecretType
	}{
		{string(types.SecretTypeLoginPasswordRu), types.SecretTypeLoginPassword},
		{string(types.SecretTypeTextRu), types.SecretTypeText},
		{string(types.SecretTypeBinaryRu), types.SecretTypeBinary},
		{string(types.SecretTypeBankCardRu), types.SecretTypeBankCard},
		{"Отмена", ""},
	}

	labels := make([]string, len(items))
	for i := range items {
		labels[i] = items[i].label
	}

	prompt := promptui.Select{
		Label: "Тип секрета",
		Items: labels,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	secretType := items[index].typ
	if secretType == "" {
		return nil, errors.New("операция отменена")
	}

	secretPayload := &types.SecretPayload{
		Type: secretType,
	}

	switch secretType {
	case types.SecretTypeLoginPassword:
		secretPayload.Data, err = promptCreateLoginPassword()

	case types.SecretTypeText:
		secretPayload.Data, err = promptCreateText()

	case types.SecretTypeBinary:
		secretPayload.Data, err = promptCreateBinary()

	case types.SecretTypeBankCard:
		secretPayload.Data, err = promptCreateBankCard()
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка подготовки данных")
	}

	name, err := promptString("Название секрета")
	if err != nil {
		return nil, err
	}

	secretPayload.Name = name

	metaChoice, err := selectMenu(
		"Добавить дополнительную информацию?",
		[]string{"Да", "Нет"},
	)
	if err != nil {
		return nil, err
	}

	if metaChoice == "Да" {
		meta, err := promptString("Дополнительная информация")
		if err != nil {
			return nil, err
		}
		secretPayload.Metadata = utils.Ptr(meta)
	}

	return secretPayload, nil
}

func promptUpdateSecret(secretType types.SecretType) (*types.UpdatedSecret, error) {
	result, err := selectMenu(
		"Что обновить?",
		[]string{
			"Название",
			"Основные данные",
			"Дополнительные данные",
			"Отмена",
		},
	)
	if err != nil {
		return nil, err
	}

	updated := &types.UpdatedSecret{}

	switch result {
	case "Название":
		name, err := promptString("Название")
		if err != nil {
			return nil, err
		}
		updated.Name = utils.Ptr(name)

	case "Основные данные":
		switch secretType {
		case types.SecretTypeLoginPassword:
			updated.Data, err = promptCreateLoginPassword()

		case types.SecretTypeText:
			updated.Data, err = promptCreateText()

		case types.SecretTypeBinary:
			updated.Data, err = promptCreateBinary()

		case types.SecretTypeBankCard:
			updated.Data, err = promptCreateBankCard()
		}

		if err != nil {
			return nil, err
		}

	case "Дополнительные данные":
		meta, err := promptString("Дополнительные данные")
		if err != nil {
			return nil, err
		}
		updated.Metadata = utils.Ptr(meta)

	case "Отмена":
		return nil, errors.New("операция отменена")
	}

	return updated, nil
}

func promptCreateLoginPassword() ([]byte, error) {
	creds, err := promptCredentials()
	if err != nil {
		return nil, err
	}
	return json.Marshal(creds)
}

func promptCreateText() ([]byte, error) {
	text, err := promptString("Текст")
	if err != nil {
		return nil, err
	}
	return json.Marshal(text)
}

func promptCreateBinary() ([]byte, error) {
	path, err := promptString("Путь к файлу")
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения информации о файле")
	}

	if info.Size() > 10*1024*1024 {
		return nil, errors.New("файл слишком большой (максимум 10MB)")
	}

	return os.ReadFile(path)
}

func promptCreateBankCard() ([]byte, error) {
	number, err := promptString("Номер карты")
	if err != nil {
		return nil, err
	}

	holder, err := promptString("Владелец карты")
	if err != nil {
		return nil, err
	}

	expiry, err := promptString("Срок действия (MM/YY)")
	if err != nil {
		return nil, err
	}

	cvv, err := promptPassword("CVV")
	if err != nil {
		return nil, err
	}

	card := types.BankCard{
		Number: number,
		Holder: holder,
		Expiry: expiry,
		CVV:    cvv,
	}

	return json.Marshal(card)
}
