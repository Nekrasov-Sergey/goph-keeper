package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/utils"
)

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
		Validate: func(input string) error {
			if input == "" {
				return errors.New("поле не может быть пустым")
			}
			return nil
		},
		Mask: '*',
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
	prompt := promptui.Prompt{
		Label: "ID секрета",
		Validate: func(input string) error {
			if _, err := strconv.ParseInt(input, 10, 64); err != nil {
				return errors.New("некорректный ID")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	id, _ := strconv.ParseInt(result, 10, 64)
	return id, nil
}

func promptCreateSecret() (*types.SecretInput, error) {
	items := []struct {
		Label string
		Type  types.SecretType
	}{
		{string(types.SecretTypeLoginPasswordRu), types.SecretTypeLoginPassword},
		{string(types.SecretTypeTextRu), types.SecretTypeText},
		{string(types.SecretTypeBinaryRu), types.SecretTypeBinary},
		{string(types.SecretTypeBankCardRu), types.SecretTypeBankCard},
		{"Отмена", ""},
	}

	labels := make([]string, len(items))
	for i := range items {
		labels[i] = items[i].Label
	}

	selectPrompt := promptui.Select{
		Label: "Тип секрета",
		Items: labels,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		return nil, err
	}

	if items[index].Type == "" {
		return nil, fmt.Errorf("операция отменена")
	}

	secret := &types.SecretInput{
		Type: items[index].Type,
	}

	switch secret.Type {
	case types.SecretTypeLoginPassword:
		secret.Data, err = promptCreateLoginPassword()

	case types.SecretTypeText:
		secret.Data, err = promptCreateText()

	case types.SecretTypeBinary:
		secret.Data, err = promptCreateBinary()

	case types.SecretTypeBankCard:
		secret.Data, err = promptCreateBankCard()
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации: %w", err)
	}

	name, err := promptString("Название секрета")
	if err != nil {
		return nil, err
	}

	secret.Name = name

	metaSelect := promptui.Select{
		Label: "Добавить дополнительную информацию?",
		Items: []string{"Да", "Нет"},
	}

	_, metaChoice, _ := metaSelect.Run()

	if metaChoice == "Да" {
		meta, _ := promptString("Дополнительная информация")
		secret.Metadata = utils.Ptr(meta)
	}

	return secret, nil
}

func promptUpdateSecret(secretType types.SecretType) (*types.UpdatedSecret, error) {
	items := []string{
		"Название",
		"Основные данные",
		"Дополнительные данные",
		"Отмена",
	}

	prompt := promptui.Select{
		Label: "Что обновить?",
		Items: items,
	}

	_, result, err := prompt.Run()
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
			data, err := promptCreateLoginPassword()
			if err != nil {
				return nil, err
			}
			updated.Data = data

		case types.SecretTypeText:
			data, err := promptCreateText()
			if err != nil {
				return nil, err
			}
			updated.Data = data

		case types.SecretTypeBinary:
			data, err := promptCreateBinary()
			if err != nil {
				return nil, err
			}
			updated.Data = data

		case types.SecretTypeBankCard:
			data, err := promptCreateBankCard()
			if err != nil {
				return nil, err
			}
			updated.Data = data
		}

	case "Дополнительные данные":

		meta, err := promptString("Дополнительные данные")
		if err != nil {
			return nil, err
		}

		updated.Metadata = utils.Ptr(meta)

	case "Отмена":
		return nil, fmt.Errorf("операция отменена")
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
		return nil, fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	if info.Size() > 10*1024*1024 {
		return nil, fmt.Errorf("файл слишком большой (максимум 10MB)")
	}

	return os.ReadFile(path)
}

func promptCreateBankCard() ([]byte, error) {
	number, _ := promptString("Номер карты")
	holder, _ := promptString("Владелец карты")
	expiry, _ := promptString("Срок действия (MM/YY)")
	cvv, _ := promptPassword("CVV")

	card := types.BankCard{
		Number: number,
		Holder: holder,
		Expiry: expiry,
		CVV:    cvv,
	}

	return json.Marshal(card)
}
