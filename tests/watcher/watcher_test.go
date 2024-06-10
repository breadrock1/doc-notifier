package watcher

import (
	"doc-notifier/internal/pkg/reader"
	watcher2 "doc-notifier/internal/watcher"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const TestcaseDirPath = "../testcases/"
const IndexerDirPath = "../../indexer/"

func TestWatcherManager(t *testing.T) {
	watch := watcher2.New(&watcher2.Options{

		WatcherServiceAddress: "0.0.0.0:2893",
		WatchedDirectories:    []string{IndexerDirPath},

		OcrServiceAddress: "http://localhost:8004",
		OcrServiceMode:    "read-raw-file",

		DocSearchAddress: "http://localhost:2892",

		TokenizerServiceAddress: "http://localhost:8001",
		TokenizerServiceMode:    "none",
		TokenizerChunkSize:      0,
		TokenizerChunkOverlap:   0,
		TokenizerReturnChunks:   false,
		TokenizerChunkBySelf:    false,
		TokenizerTimeout:        10,
	})

	t.Run("Append directory to watch", func(t *testing.T) {
		err := watch.AppendDirectories([]string{TestcaseDirPath})
		assert.NoError(t, err, "Failed while appending dir to watch")

		dirs := watch.GetWatchedDirectories()
		assert.Equal(t, len(dirs), 1, "Not equal appended dirs")

		err = watch.RemoveDirectories([]string{TestcaseDirPath})
		assert.NoError(t, err, "Failed while detach dir to watch")
	})

	t.Run("Append multiple dirs to watch", func(t *testing.T) {
		dirs := []string{TestcaseDirPath, IndexerDirPath}
		err := watch.AppendDirectories(dirs)
		assert.NoError(t, err, "Failed while appending dir to watch")

		attached := watch.GetWatchedDirectories()
		assert.Equal(t, len(dirs), len(attached), "Not equal appended dirs")

		err = watch.RemoveDirectories(dirs)
		assert.NoError(t, err, "Failed while detach dir to watch")
	})

	t.Run("Caught error while attach and detach non existing firs", func(t *testing.T) {
		err := watch.AppendDirectories([]string{TestcaseDirPath + "any"})
		assert.Error(t, err, "Failed while catching error to append")

		err = watch.RemoveDirectories([]string{TestcaseDirPath + "any"})
		assert.Error(t, err, "Failed while catching error to detach")
	})

	t.Run("Parse complex structure with OcrMetadata", func(t *testing.T) {
		file, _ := os.ReadFile(TestcaseDirPath + "ocr_result.json")
		var previews []reader.DocumentPreview
		if err := json.Unmarshal(file, &previews); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v", previews)
	})
}
