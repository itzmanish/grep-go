### Write a command line program that implements Unix command `grep` like functionality. The program should have following features.

- [x] Ability to search for a string in a file. Feel free to assume case-sensitive and exact word matches if required for simplicity’s sake.

```
$ ./mygrep "search_string" filename.txt
I found the search_string in the file.
```

- [x] Ability to search for a string from standard input

```
$ ./mygrep foo
bar
barbazfoo
Foobar
food
^D
```

should produce

```
barbazfoo
food
```

- [x] Ability to write output to a file instead of a standard out.

```
$ ./mygrep lorem loreipsum.txt -o out.txt
```

should create an out.txt file with the output from `mygrep`. for example,

```
$ cat out.txt
lorem ipsum
a dummy text usually contains lorem ipsum
```

- [x] Ability to search for a string recursively in any of the files in a given directory. When searching in multiple files, the output should indicate the file name and all the output from one file should be grouped together in the final output. (in other words, output from two files shouldn't be interleaved in the final output being printed)

```
$ ./mygrep "test" tests
tests/test1.txt:this is a test file
tests/test1.txt:one can test a program by running test cases
tests/inner/test2.txt:this file contains a test line
```

Feel free to make suitable assumptions if needed, ensure to document them in README.md

Criteria for evaluation:

- Use Golang to write your code
- Add unit tests for functions
- Clean and readable code
- Adherance to coding standards and guidelines
