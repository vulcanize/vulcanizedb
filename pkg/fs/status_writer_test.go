package fs_test

import (
	"io/ioutil"
	"os"

	"github.com/makerdao/vulcanizedb/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("statusWriter", func() {
	var (
		testFileName     = "test-file"
		testFilePath     = "/tmp/" + testFileName
		testFileContents = []byte("test contents\n")
		writer           fs.StatusWriter
	)
	BeforeEach(func() {
		writer = fs.NewStatusWriter(testFilePath, testFileContents)
		err := writer.Write()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(testFilePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("It writes the contents to the given file path", func() {
		info, err := os.Stat(testFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.IsDir()).To(BeFalse())
		Expect(info.Name()).To(Equal(testFileName))

		contents, err := ioutil.ReadFile(testFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(Equal(testFileContents))
	})

	It("appends the contents to the end of the given file", func() {
		newFileContents := []byte("new file contents")
		writer2 := fs.NewStatusWriter(testFilePath, newFileContents)
		err2 := writer2.Write()
		Expect(err2).NotTo(HaveOccurred())

		contents, err := ioutil.ReadFile(testFilePath)
		expectedContents := []byte(testFileContents)
		expectedContents = append(expectedContents, []byte(newFileContents)...)

		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(Equal(expectedContents))
	})
})
