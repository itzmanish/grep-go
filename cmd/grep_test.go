package cmd

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"testing"
)

var exp *regexp.Regexp = regexp.MustCompile("(?i)dummy")

func TestWrite(t *testing.T) {
	testcases := []struct {
		name         string
		failExpected bool
		out          string
	}{
		{
			name:         "Success",
			failExpected: false,
			out:          "test_out.txt",
		},
	}
	for _, test := range testcases {
		lines := []string{
			"sdf klsdf skl dummy new",
			"sdjflks",
			"sadfslklf sdfiowejsf dummysdfkl",
		}
		t.Run(test.name, func(t *testing.T) {
			writtenSize, err := Write(test.out, lines)
			if err != nil {
				t.Error(err)
			}
			f, err := os.Stat(test.out)
			if err != nil {
				t.Error(err)
			}
			if f.Size() != int64(writtenSize) {
				t.Error("File not written")
			}
			os.Remove(test.out)

		})
	}
}

func TestFindExp(t *testing.T) {
	testcases := []struct {
		name    string
		input   string
		success bool
	}{
		{
			name:    "match found",
			input:   "sdfo sdf orhte jdDumMysfk",
			success: true,
		},
		{
			name:    "Smatch failed",
			input:   "sdfo sdf orhte dfasdf sdf",
			success: false,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			found := FindExp(exp, test.input)
			if found != test.success {
				t.Errorf("Expected match: %v but got match %v", test.success, found)
			}
		})
	}
}

func TestFind(t *testing.T) {
	testcases := []struct {
		name           string
		sampleInput    []string
		expectedResult []string
		path           string
		verbose        bool
	}{
		{
			name: "Find text with success",
			sampleInput: []string{
				"sdf klsdf skl dummy new",
				"sdjflks",
				"sadfslklf sdfiowejsf dummysdfkl",
			},
			expectedResult: []string{
				"sdf klsdf skl dummy new",
				"sadfslklf sdfiowejsf dummysdfkl",
			},
		},
		{
			name: "Find text with success",
			sampleInput: []string{
				"sdf klsdf skl dummy new",
				"sdjflks",
				"sadfslklf sdfiowejsf dummysdfkl",
			},
			expectedResult: []string{
				"test.txt:sdf klsdf skl dummy new",
				"test.txt:sadfslklf sdfiowejsf dummysdfkl",
			},
			path:    "test.txt",
			verbose: true,
		},
		{
			name: "Not found text",
			sampleInput: []string{
				"sdf klsdf skl new",
				"sdjflks",
				"sadfslklf sdfiowejsf dfkl",
			},
			expectedResult: []string{},
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {

			var stdin bytes.Buffer
			for _, line := range test.sampleInput {
				stdin.Write([]byte(line + "\n"))
			}

			out := Find(&stdin, exp, "", false)
			if len(out) != len(test.expectedResult) {
				t.Errorf("Expected length of out:%v but got %v", len(test.expectedResult), len(out))
				return
			}
			if len(out) > 0 {
				for i := len(out); i <= 0; i++ {
					if out[i] != test.expectedResult[i] {
						t.Errorf("Expected string: [%v] but got [%v]", test.expectedResult[i], out[i])
					}
				}
			}
		})
	}
}

func TestOpenAndFind(t *testing.T) {
	testcases := []struct {
		name        string
		path        string
		expected    []string
		ErrExpected bool
	}{
		{
			name: "Success with right path",
			path: "../test_input.txt",
			expected: []string{
				"a dummy text usually contains lorem ipsum",
			},
			ErrExpected: false,
		},
		{
			name:        "Fail with wrong path",
			path:        "test.txt",
			expected:    []string{},
			ErrExpected: true,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			out, err := OpenAndFind(test.path, exp, false)
			if test.ErrExpected && err == nil {
				t.Error("Expected error but got no err")
			} else if err != nil && !test.ErrExpected {
				t.Errorf("Error not expected but got %v", err)
			}
			if len(out) != len(test.expected) {
				t.Errorf("Expected length of out:%v but got %v", len(test.expected), len(out))
				return
			}
			if len(out) > 0 {
				for i := len(out); i <= 0; i++ {
					if out[i] != test.expected[i] {
						t.Errorf("Expected string: [%v] but got [%v]", test.expected[i], out[i])
					}
				}
			}
		})
	}
}

func TestRun(t *testing.T) {
	testcases := []struct {
		name        string
		searchText  string
		inputFile   string
		outputFile  string
		ErrExpected bool
	}{
		{
			name:        "Success without error",
			searchText:  "dummy",
			inputFile:   "../test_input.txt",
			ErrExpected: false,
		},
		{
			name:        "Success without error",
			searchText:  "dummy",
			inputFile:   "../sample_dir",
			ErrExpected: false,
		},
		{
			name:        "Success with input as stdin",
			searchText:  "dummy",
			ErrExpected: false,
		},
		{
			name:        "Fail with wrong path error",
			searchText:  "dummy",
			inputFile:   "lorem.txt",
			ErrExpected: true,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }() // Restore original Stdin

			if len(test.inputFile) == 0 {
				content := []byte("sadfslklf sdfiowejsf dummysdfkl\ndjflks\nsdf klsdf skl dummy new")
				tmpfile, err := os.CreateTemp("", "")
				if err != nil {
					t.Error(err)
				}

				defer os.Remove(tmpfile.Name()) // clean up

				if _, err := tmpfile.Write(content); err != nil {
					log.Fatal(err)
				}

				if _, err := tmpfile.Seek(0, 0); err != nil {
					log.Fatal(err)
				}
				os.Stdin = tmpfile
			}
			result, err := Run(test.searchText, test.inputFile)
			if err == nil && test.ErrExpected {
				t.Error("Error expected but got no error")
				return
			}
			if err != nil && !test.ErrExpected {
				t.Errorf("Error not expected but got error: %v", err)
				return
			}
			for res := range result.Result {
				if len(res) == 0 && !test.ErrExpected {
					t.Error("Expected some result but got nothing!")
				}
			}
		})
	}
}

func BenchmarkRunSingleFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := Run("dummy", "../test_input.txt")
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkRunSingleFileParallel(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := Run("dummy", "../test_input.txt")
			if err != nil {
				b.Error(err)
			}
		}

	})
}

func BenchmarkRunFolder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := Run("dummy", "../sample_dir")
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkRunFolderParallel(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := Run("dummy", "../sample_dir")
			if err != nil {
				b.Error(err)
			}
		}

	})
}
