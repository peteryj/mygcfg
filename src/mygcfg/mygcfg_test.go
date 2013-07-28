//
// Test functions for mygcfg
//

package mygcfg_test

import (
    "testing"
    "path/filepath"
    "os"
    "mygcfg"
)

var testFileList []string

func parseEachFile(path string, info os.FileInfo, err error) error {

    if err != nil {
        return err
    }

    fi, err := os.Stat(path)
    if err != nil {
        return err
    }

    if fi.IsDir() {
        return err
    }

    testFileList = append(testFileList, path)
    return nil
}

func TestIssue(t *testing.T) {
    testFileList = make([]string, 0, 0);
    err := filepath.Walk("../../test_data/", parseEachFile);
    if err != nil {
        t.Logf("%v", err)
        t.Fail()
    }

    t.Logf("%v\n", testFileList)

    for _, fpath := range testFileList {
        var testParser mygcfg.Parser
        err := testParser.ParseFile(fpath)
        if err != nil {
            t.Logf("f(%s): %v\n", fpath, err)
            t.Fail()
        }
    }
}


