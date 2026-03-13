package grpc

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func PrintError(err error) {
	fmt.Println()
	if st, ok := status.FromError(err); ok {
		fmt.Println(st.Message())
		return
	}
	fmt.Println(err)
}
