package sha224sum

import (
	"github.com/guonaihong/coreutils/hashcore"
)

func Main(argv []string) {
	hashcore.Main(argv, hashcore.Sha224)
}
