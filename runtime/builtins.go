package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

func (env *env) evalPrint(args []int) (int, error) {
	var sb strings.Builder

	for i, arg := range args {
		sb.WriteString(strconv.Itoa(arg))
		if i != len(args)-1 {
			sb.WriteString(" ")
		}
	}

	fmt.Fprintln(env.out, sb.String())

	// TODO print should not return anything
	return 0, nil
}
