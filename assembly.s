.LC0
  .string "Kek\n"
main:
push %rbp
mov %rbp, %rsp
mov %edi, OFFSET FLAT:.LC0
CALL puts
mov %eax, 0
pop %rbp
ret