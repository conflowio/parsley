// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

var _ = Describe("File", func() {

	var (
		fs   *parsley.FileSet
		f    *text.File
		data []byte
	)

	JustBeforeEach(func() {
		f = text.NewFile("testfile", data)
		fs = parsley.NewFileSet(f)
	})

	BeforeEach(func() {
		data = []byte("ab\nc")
	})

	It("should implement the parsley.File interface", func() {
		var _ parsley.File = &text.File{}
	})

	Describe("Position()", func() {
		It("should return with the position for an offset", func() {
			Expect(f.Position(0)).To(Equal(text.NewPosition("testfile", 1, 1)))
			Expect(f.Position(1)).To(Equal(text.NewPosition("testfile", 1, 2)))
			Expect(f.Position(2)).To(Equal(text.NewPosition("testfile", 1, 3)))
			Expect(f.Position(3)).To(Equal(text.NewPosition("testfile", 2, 1)))
			Expect(f.Position(4)).To(Equal(text.NewPosition("testfile", 2, 2)))
		})

		It("should return nil position for an invalid offset", func() {
			Expect(f.Position(5)).To(Equal(parsley.NilPosition))
		})
	})

	Describe("Len()", func() {
		Context("when data is empty", func() {
			BeforeEach(func() {
				data = []byte{}
			})

			It("should return zero", func() {
				Expect(f.Len()).To(Equal(0))
			})
		})

		It("should return the data length", func() {
			Expect(f.Len()).To(Equal(4))
		})
	})

	Context("when data contains Windows-style line endings", func() {
		BeforeEach(func() {
			data = []byte("a\r\nb\r\nc\n")
		})

		It("should remove the \r characters", func() {
			Expect(f.Len()).To(Equal(6))
			Expect(f.Position(1).(*text.Position).Line).To(Equal(1))
			Expect(f.Position(2).(*text.Position).Line).To(Equal(2))
			Expect(f.Position(3).(*text.Position).Line).To(Equal(2))
			Expect(f.Position(4).(*text.Position).Line).To(Equal(3))
		})
	})

	Context("ReadFile", func() {
		var (
			readFileErr error
			filename    string
		)

		JustBeforeEach(func() {
			f, readFileErr = text.ReadFile(filename)
		})

		Context("reading an existing file", func() {
			var tmpDir string

			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "parsley-test-")
				Expect(err).ToNot(HaveOccurred())
				filename = filepath.Join(tmpDir, "testfile")
				ioutil.WriteFile(filename, []byte("ab\nc"), 0600)
			})

			AfterEach(func() {
				if tmpDir != "" {
					os.RemoveAll(tmpDir)
				}
			})

			It("should succeed", func() {
				Expect(readFileErr).ToNot(HaveOccurred())
			})

			Describe("Position", func() {
				It("should return the position for an offset", func() {
					Expect(f.Position(0)).To(Equal(text.NewPosition(filename, 1, 1)))
					Expect(f.Position(1)).To(Equal(text.NewPosition(filename, 1, 2)))
					Expect(f.Position(2)).To(Equal(text.NewPosition(filename, 1, 3)))
					Expect(f.Position(3)).To(Equal(text.NewPosition(filename, 2, 1)))
					Expect(f.Position(4)).To(Equal(text.NewPosition(filename, 2, 2)))
				})

				It("should return nil position for an invalid offset", func() {
					Expect(f.Position(5)).To(Equal(parsley.NilPosition))
				})

				It("should contain a readable position string", func() {
					Expect(f.Position(2).String()).To(Equal(filename + ":1:3"))
				})
			})

			Describe("Len", func() {
				It("should return with the filesize", func() {
					Expect(f.Len()).To(Equal(4))
				})
			})
		})

		Context("reading a non-existing file", func() {
			BeforeEach(func() {
				filename = "/tmp/non-existing-filename"
			})

			It("should throw an error", func() {
				Expect(readFileErr).To(MatchError("can not read /tmp/non-existing-filename"))
			})
		})
	})

	Describe("Pos()", func() {
		It("should return with a global position (pos 0)", func() {
			pos := f.Pos(0)
			Expect(pos).To(Equal(parsley.Pos(1)))
			Expect(fs.Position(pos)).To(Equal(text.NewPosition("testfile", 1, 1)))
		})

		It("should return with a global position (at the end)", func() {
			pos := f.Pos(3)
			Expect(pos).To(Equal(parsley.Pos(4)))
			Expect(fs.Position(pos)).To(Equal(text.NewPosition("testfile", 2, 1)))
		})
	})
})
