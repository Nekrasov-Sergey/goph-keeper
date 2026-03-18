// Package grpc реализует gRPC-клиент для взаимодействия с сервером.
package grpc

import (
	"fmt"

	"google.golang.org/grpc/status"
)

// PrintError выводит ошибку gRPC в читаемом формате.
func PrintError(err error) {
	fmt.Println()
	if st, ok := status.FromError(err); ok {
		fmt.Println(st.Message())
		return
	}
	fmt.Println(err)
}
