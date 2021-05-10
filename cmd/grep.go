package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// Run start searching of string in given input file or os.Stdin in case of input file path not provided.
// It returns result as array of string and error
func Run(search_text, input_file string) ([]string, error) {
	exp := regexp.MustCompile("(?i)" + search_text)
	matched_string := []string{}

	if len(input_file) != 0 {
		exist, isDir, err := Exists(input_file)
		if err != nil {
			return []string{}, err
		}
		if exist && isDir {
			filepaths := []string{}
			err = filepath.Walk(input_file, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					filepaths = append(filepaths, path)
				}
				return nil
			})
			if err != nil {
				return []string{}, err
			}
			for _, path := range filepaths {
				out, err := OpenAndFind(path, exp, true)
				if err != nil {
					return []string{}, err
				}
				matched_string = append(matched_string, out...)
			}
		} else if exist && !isDir {
			matched_string, err = OpenAndFind(input_file, exp, false)
			if err != nil {
				return []string{}, err
			}
		}
	} else {
		matched_string = Find(os.Stdin, exp, "", false)
	}

	return matched_string, nil
}

// Exists returns whether the given file or directory exists
func Exists(path string) (bool, bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return true, info.IsDir(), nil
	}
	return false, false, err
}

// OpenAndFind opens a file by given path and find the expression in that file
// It returns result as array of string and error
// You can get filename in result array by using verbose as true.
func OpenAndFind(path string, exp *regexp.Regexp, verbose bool) ([]string, error) {
	out := []string{}
	f, err := os.Open(path)
	if err != nil {
		return out, err
	}
	defer f.Close()
	out = Find(f, exp, path, verbose)
	return out, nil
}

// Find searches for expression in given Reader interface.
// It returns result as array of string.
// You can get filename in result array by using verbose as true.
func Find(r io.Reader, exp *regexp.Regexp, path string, verbose bool) []string {
	found_strings := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if FindExp(exp, scanner.Text()) {
			matched_line := scanner.Text()
			if verbose {
				matched_line = fmt.Sprintf("%s:%s", path, scanner.Text())
			}
			found_strings = append(found_strings, matched_line)
		}
	}
	return found_strings
}

// FindExp searches for expression in given input string and return whether it exist on input string or not.
func FindExp(exp *regexp.Regexp, input string) bool {
	found_string := exp.FindString(input)
	return len(found_string) != 0
}

// Write writes given array of strings to out file or stdout in case of out filepath not provided.
func Write(out string, lines []string) (int, error) {
	writer := os.Stdout
	var writtenSize int
	if len(out) != 0 {
		f, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		writer = f
	}
	for i := len(lines) - 1; i >= 0; i-- {
		n, err := writer.Write([]byte(lines[i] + "\n"))
		if err != nil {
			return 0, err
		}
		writtenSize += n
	}
	return writtenSize, nil
}
