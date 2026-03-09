package mappers

import (
	"github.com/pkg/errors"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

func ProtoSecretTypeToDomain(t pb.SecretType) (types.SecretType, error) {
	switch t {
	case pb.SecretType_LoginPassword:
		return types.SecretTypeLoginPassword, nil

	case pb.SecretType_Text:
		return types.SecretTypeText, nil

	case pb.SecretType_Binary:
		return types.SecretTypeBinary, nil

	case pb.SecretType_BankCard:
		return types.SecretTypeBankCard, nil

	default:
		return types.SecretTypeUnknown, errors.New("неизвестный тип секрета")
	}
}

func TranslateProtoSecretType(t pb.SecretType) (types.SecretTypeRu, error) {
	switch t {
	case pb.SecretType_LoginPassword:
		return types.SecretTypeLoginPasswordRu, nil

	case pb.SecretType_Text:
		return types.SecretTypeTextRu, nil

	case pb.SecretType_Binary:
		return types.SecretTypeBinaryRu, nil

	case pb.SecretType_BankCard:
		return types.SecretTypeBankCardRu, nil

	default:
		return types.SecretTypeUnknownRu, errors.New("неизвестный тип секрета")
	}
}

func DomainSecretTypeToProto(t types.SecretType) (pb.SecretType, error) {
	switch t {
	case types.SecretTypeLoginPassword:
		return pb.SecretType_LoginPassword, nil

	case types.SecretTypeText:
		return pb.SecretType_Text, nil

	case types.SecretTypeBinary:
		return pb.SecretType_Binary, nil

	case types.SecretTypeBankCard:
		return pb.SecretType_BankCard, nil

	default:
		return pb.SecretType_Unspecified, errors.New("неизвестный тип секрета")
	}
}
