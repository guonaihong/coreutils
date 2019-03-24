package utils

func NewArgs(args []string) []string {
	if len(args) == 0 {
		return append(args, "-")
	}
	return args
}
