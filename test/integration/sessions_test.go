package sessions_test

import (
	"github.com/ld100/goblet/test/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var (
		//longBook  Book
		//shortBook Book
	)

	BeforeEach(func() {
		//longBook = Book{
		//	Title:  "Les Miserables",
		//	Author: "Victor Hugo",
		//	Pages:  1488,
		//}
		//
		//shortBook = Book{
		//	Title:  "Fox In Socks",
		//	Author: "Dr. Seuss",
		//	Pages:  24,
		//}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				//Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
				util.Get("http://httpbin.org/get")
				Expect(1).To(Equal(1))
			})
		})

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				//Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
				Expect(1).To(Equal(1))
			})
		})
	})
})