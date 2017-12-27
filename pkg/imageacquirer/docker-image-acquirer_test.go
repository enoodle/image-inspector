package imageacquirer

import (
	"io"
	"testing"
)

func Test_decodeDockerResponse(t *testing.T) {
	no_error_input := "{\"Status\": \"fine\"}"
	one_error := "{\"Status\": \"fine\"}{\"Error\": \"Oops\"}{\"Status\": \"fine\"}"
	decode_error := "{}{}what"
	decode_error_message := "Error decoding json: invalid character 'w' looking for beginning of value"
	tests := map[string]struct {
		readerInput    string
		expectedErrors bool
		errorMessage   string
	}{
		"no error":      {readerInput: no_error_input, expectedErrors: false},
		"error":         {readerInput: one_error, expectedErrors: true, errorMessage: "Oops"},
		"decode errror": {readerInput: decode_error, expectedErrors: true, errorMessage: decode_error_message},
	}

	for test_name, test_params := range tests {
		parsedErrors := make(chan error, 100)
		defer func() { close(parsedErrors) }()

		go func() {
			reader, writer := io.Pipe()
			// handle closing the reader/writer in the method that creates them
			defer reader.Close()
			defer writer.Close()
			go decodeDockerResponse(parsedErrors, reader)
			writer.Write([]byte(test_params.readerInput))
		}()

		select {
		case decodedErrors := <-parsedErrors:
			if decodedErrors == nil && test_params.expectedErrors {
				t.Errorf("Expected to parse an error, but non was parsed in test %s", test_name)
			}
			if decodedErrors != nil {
				if !test_params.expectedErrors {
					t.Errorf("Expected not to get errors in test %s but got: %v", test_name, decodedErrors)
				} else {
					if decodedErrors.Error() != test_params.errorMessage {
						t.Errorf("Expected error message is different than expected in test %s. Expected %v received %v",
							test_name, test_params.errorMessage, decodedErrors.Error())
					}
				}
			}
		}
	}
}
