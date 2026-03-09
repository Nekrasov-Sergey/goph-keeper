package grpc

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func PrintError(err error) {
	st, ok := status.FromError(err)
	if !ok {
		fmt.Println(err)
		return
	}
	fmt.Println(st.Message())
}
