package grpc

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthInterceptor(jwtSecret []byte) grpc.UnaryServerInterceptor {
	publicMethods := map[string]struct{}{
		pb.Keeper_Register_FullMethodName: {},
		pb.Keeper_Login_FullMethodName:    {},
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "нет metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "нет токена")
		}

		tokenString := strings.TrimPrefix(values[0], "Bearer ")

		token, err := jwt.Parse(tokenString, func(*jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "невалидный токен")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "не удалось прочитать claims")
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "user_id отсутствует в токене")
		}

		userID := int64(userIDFloat)

		ctx = context.WithValue(ctx, UserIDKey, userID)

		return handler(ctx, req)
	}
}

func GetUserID(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, errors.New("пользователь не найден в контексте")
	}
	return id, nil
}
