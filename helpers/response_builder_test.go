package helpers

import (
	"reflect"
	"testing"
)

type testCase struct {
    // Inputs
    message string
    data    any
    err     any

    // Output
    response Response
}

func TestBuildSuccessfulResponse(t *testing.T) {
    testCases := []testCase{
        {
            message: "Example Message",
            data:    100,
            response: Response{
                Message: "Example Message",
                Data:    100,
            },
        },
    }

    for _, tc := range testCases {
        result := BuildSuccessfulResponse(tc.message, tc.data)
        
        if !reflect.DeepEqual(result, tc.response) {
            t.Errorf("BuildSuccessfulResponse(%s, %v) = %v, expected %v", tc.message, tc.data, result, tc.response)
        }
    }
}

func TestBuildFailedResponse(t *testing.T) {
    testCases := []testCase{
        {
            message: "Error Message",
            data:    nil,
            err:     "an error occurred",
            response: Response{
                Message: "Error Message",
                Error:   "an error occurred",
                Data:    nil,
            },
        },
    }

    for _, tc := range testCases {
        result := BuildFailedResponse(tc.message, tc.err, tc.data)
        
        if !reflect.DeepEqual(result, tc.response) {
            t.Errorf("BuildFailedResponse(%s, %v, %v) = %v, expected %v", tc.message, tc.err, tc.data, result, tc.response)
        }
    }
}
