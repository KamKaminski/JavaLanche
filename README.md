# JavaLanche

JavaLanche is a programming language developed as an assignment for the "Language Design and Implementation" module at the University of Derby. 

The language supports a variety of features including arithmetic operations, boolean logic, string manipulation, global variables, and repetitive logic control structures.

## Features

* **Arithmetic Operations:** Perform calculations using operators like `+`, `-`, `*`, and `/`.
* **Boolean Logic:** Evaluate logical expressions with operators like `&&`, `||`, and `!`.
* **String Manipulation:** Combine Strings.
* **Global Variables:** Define variables that can be used anywhere in your code.
* **Control Structures:** Implement loops and conditional logic to control the flow of your program.

## Usage

JavaLanche can be interacted with via command line arguments or by pipelining. 

For example, you can pass in a JavaLanche script as a command line argument:
.\main.exe fizzBuzz.javalanche


Or you can use input redirection to read a script from a file:

go run ./... < fizzBuzz.javalanche


## Example Code

Example JavaLanche code can be found in the source files. This will give you a feel for the syntax and capabilities of the language.

## License

This project is licensed under the terms of the [MIT License](LICENSE).

## Known Issues
Race condition that causes sometimes errors and wrong outputs. - Working on the fix

