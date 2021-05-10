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

func FindExp(exp *regexp.Regexp, input string) bool {
	found_string := exp.FindString(input)
	return len(found_string) != 0
}

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
