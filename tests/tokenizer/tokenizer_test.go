package tokenizer

import (
	"context"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/tokenizer"
	"doc-notifier/tests/mocked"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseOtherDirPath = "../../testcases/directory/"
const TestcaseFilePath = TestcaseOtherDirPath + "test_file_1.txt"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestComputeContentTokens(t *testing.T) {
	tokenizerService := tokenizer.New(&tokenizer.Options{
		Mode:         tokenizer.GetModeFromString("assistant"),
		Address:      "http://localhost:3451",
		ChunkSize:    0,
		ChunkedFlag:  false,
		ChunkOverlap: 0,
	})

	t.Run("Compute tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseFilePath)
		document.Content = "test_file_1"
		document.DocumentName = "test_file_1.txt"
		tokens, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.NoError(t, computeErr, "Returned error while storing document")
		assert.Equal(t, tokens.Chunks, 1, "Non correct returned chunks size")
		assert.Equal(t, tokens.ChunkedText[0], "test_file_1", "Non correct returned chunks data")
		assert.Equal(t, tokens.Vectors[0], []float64{0.345, 0.045}, "Non correct returned vectors")

		time.AfterFunc(200*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while computing tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		_, parseErr := reader.ParseFile(TestcaseNonExistingFilePath)
		//_, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.Error(t, parseErr, "Returned error while parsing file")
		//assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		_ = tokenizer.New(&tokenizer.Options{
			Mode:         tokenizer.GetModeFromString("assistant"),
			Address:      "http://localhost:4444",
			ChunkSize:    0,
			ChunkedFlag:  false,
			ChunkOverlap: 0,
		})

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		_, parseErr := reader.ParseFile(TestcaseNonExistingFilePath)
		//_, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.Error(t, parseErr, "Returned error while parsing file")
		//assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}
