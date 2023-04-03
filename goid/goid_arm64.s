//go:build arm64

#include "textflag.h"

TEXT ·getg(SB),NOSPLIT,$0-8
	MOVD g, ret+0(FP)
	RET
