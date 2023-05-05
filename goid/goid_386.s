#include "go_asm.h"
#include "textflag.h"

TEXT ·getg(SB), NOSPLIT, $0-4
    MOVL TLS, CX
    MOVL 0(CX)(TLS*1), AX
    MOVL AX, ret+0(FP)
    RET
