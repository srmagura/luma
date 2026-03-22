mod bytecode;

use std::fs;
use std::io;

fn main() -> io::Result<()> {
    let input_path = "../test_runner/helloWorld/helloWorld.out";
    let bytes = fs::read(input_path)?;
    Ok(())
}
