package main

import (
	"context"
	"fmt"
	"os"
)

type aKey string

func searchKey(ctx context.Context, k aKey) {
	v := ctx.Value(k)
	if v != nil {
		fmt.Fprintln(os.Stdout, v)
		return
	} else {
		fmt.Fprintln(os.Stdout, k)
	}
}
