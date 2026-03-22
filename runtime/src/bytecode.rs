// Set bytecode.go for better documentation of bytecode instructions
#[repr(u8)]
enum Status {
    OpLdcI4 = 0x20,
    OpAdd = 0x58,
    OpPrint = 0x02,
}
