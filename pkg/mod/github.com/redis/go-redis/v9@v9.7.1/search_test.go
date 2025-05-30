package redis_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/redis/go-redis/v9"
)

func WaitForIndexing(c *redis.Client, index string) {
	for {
		res, err := c.FTInfo(context.Background(), index).Result()
		Expect(err).NotTo(HaveOccurred())
		if c.Options().Protocol == 2 {
			if res.Indexing == 0 {
				return
			}
			time.Sleep(100 * time.Millisecond)
		} else {
			return
		}
	}
}

var _ = Describe("RediSearch commands Resp 2", Label("search"), func() {
	ctx := context.TODO()
	var client *redis.Client

	BeforeEach(func() {
		client = redis.NewClient(&redis.Options{Addr: ":6379", Protocol: 2})
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("should FTCreate and FTSearch WithScores", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		client.HSet(ctx, "doc1", "txt", "foo baz")
		client.HSet(ctx, "doc2", "txt", "foo bar")
		res, err := client.FTSearchWithArgs(ctx, "txt", "foo ~bar", &redis.FTSearchOptions{WithScores: true}).Result()

		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(int64(2)))
		for _, doc := range res.Docs {
			Expect(*doc.Score).To(BeNumerically(">", 0))
			Expect(doc.ID).To(Or(Equal("doc1"), Equal("doc2")))
		}
	})

	It("should FTCreate and FTSearch stopwords", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{StopWords: []interface{}{"foo", "bar", "baz"}}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		client.HSet(ctx, "doc1", "txt", "foo baz")
		client.HSet(ctx, "doc2", "txt", "hello world")
		res1, err := client.FTSearchWithArgs(ctx, "txt", "foo bar", &redis.FTSearchOptions{NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(0)))
		res2, err := client.FTSearchWithArgs(ctx, "txt", "foo bar hello world", &redis.FTSearchOptions{NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(1)))
	})

	It("should FTCreate and FTSearch filters", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}, &redis.FieldSchema{FieldName: "num", FieldType: redis.SearchFieldTypeNumeric}, &redis.FieldSchema{FieldName: "loc", FieldType: redis.SearchFieldTypeGeo}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		client.HSet(ctx, "doc1", "txt", "foo bar", "num", 3.141, "loc", "-0.441,51.458")
		client.HSet(ctx, "doc2", "txt", "foo baz", "num", 2, "loc", "-0.1,51.2")
		res1, err := client.FTSearchWithArgs(ctx, "txt", "foo", &redis.FTSearchOptions{Filters: []redis.FTSearchFilter{{FieldName: "num", Min: 0, Max: 2}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(1)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("doc2"))
		res2, err := client.FTSearchWithArgs(ctx, "txt", "foo", &redis.FTSearchOptions{Filters: []redis.FTSearchFilter{{FieldName: "num", Min: 0, Max: "+inf"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(2)))
		Expect(res2.Docs[0].ID).To(BeEquivalentTo("doc1"))
		// Test Geo filter
		geoFilter1 := redis.FTSearchGeoFilter{FieldName: "loc", Longitude: -0.44, Latitude: 51.45, Radius: 10, Unit: "km"}
		geoFilter2 := redis.FTSearchGeoFilter{FieldName: "loc", Longitude: -0.44, Latitude: 51.45, Radius: 100, Unit: "km"}
		res3, err := client.FTSearchWithArgs(ctx, "txt", "foo", &redis.FTSearchOptions{GeoFilter: []redis.FTSearchGeoFilter{geoFilter1}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res3.Total).To(BeEquivalentTo(int64(1)))
		Expect(res3.Docs[0].ID).To(BeEquivalentTo("doc1"))
		res4, err := client.FTSearchWithArgs(ctx, "txt", "foo", &redis.FTSearchOptions{GeoFilter: []redis.FTSearchGeoFilter{geoFilter2}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res4.Total).To(BeEquivalentTo(int64(2)))
		docs := []interface{}{res4.Docs[0].ID, res4.Docs[1].ID}
		Expect(docs).To(ContainElement("doc1"))
		Expect(docs).To(ContainElement("doc2"))

	})

	It("should FTCreate and FTSearch sortby", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "num", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}, &redis.FieldSchema{FieldName: "num", FieldType: redis.SearchFieldTypeNumeric, Sortable: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "num")
		client.HSet(ctx, "doc1", "txt", "foo bar", "num", 1)
		client.HSet(ctx, "doc2", "txt", "foo baz", "num", 2)
		client.HSet(ctx, "doc3", "txt", "foo qux", "num", 3)

		sortBy1 := redis.FTSearchSortBy{FieldName: "num", Asc: true}
		sortBy2 := redis.FTSearchSortBy{FieldName: "num", Desc: true}
		res1, err := client.FTSearchWithArgs(ctx, "num", "foo", &redis.FTSearchOptions{NoContent: true, SortBy: []redis.FTSearchSortBy{sortBy1}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(3)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("doc1"))
		Expect(res1.Docs[1].ID).To(BeEquivalentTo("doc2"))
		Expect(res1.Docs[2].ID).To(BeEquivalentTo("doc3"))

		res2, err := client.FTSearchWithArgs(ctx, "num", "foo", &redis.FTSearchOptions{NoContent: true, SortBy: []redis.FTSearchSortBy{sortBy2}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(3)))
		Expect(res2.Docs[2].ID).To(BeEquivalentTo("doc1"))
		Expect(res2.Docs[1].ID).To(BeEquivalentTo("doc2"))
		Expect(res2.Docs[0].ID).To(BeEquivalentTo("doc3"))

		res3, err := client.FTSearchWithArgs(ctx, "num", "foo", &redis.FTSearchOptions{NoContent: true, SortBy: []redis.FTSearchSortBy{sortBy2}, SortByWithCount: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res3.Total).To(BeEquivalentTo(int64(3)))

		res4, err := client.FTSearchWithArgs(ctx, "num", "notpresentf00", &redis.FTSearchOptions{NoContent: true, SortBy: []redis.FTSearchSortBy{sortBy2}, SortByWithCount: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res4.Total).To(BeEquivalentTo(int64(0)))
	})

	It("should FTCreate and FTSearch example", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText, Weight: 5}, &redis.FieldSchema{FieldName: "body", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		client.HSet(ctx, "doc1", "title", "RediSearch", "body", "Redisearch implements a search engine on top of redis")
		res1, err := client.FTSearchWithArgs(ctx, "txt", "search engine", &redis.FTSearchOptions{NoContent: true, Verbatim: true, LimitOffset: 0, Limit: 5}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(1)))

	})

	It("should FTCreate NoIndex", Label("search", "ftcreate", "ftsearch"), func() {
		text1 := &redis.FieldSchema{FieldName: "field", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "text", FieldType: redis.SearchFieldTypeText, NoIndex: true, Sortable: true}
		num := &redis.FieldSchema{FieldName: "numeric", FieldType: redis.SearchFieldTypeNumeric, NoIndex: true, Sortable: true}
		geo := &redis.FieldSchema{FieldName: "geo", FieldType: redis.SearchFieldTypeGeo, NoIndex: true, Sortable: true}
		tag := &redis.FieldSchema{FieldName: "tag", FieldType: redis.SearchFieldTypeTag, NoIndex: true, Sortable: true}
		val, err := client.FTCreate(ctx, "idx", &redis.FTCreateOptions{}, text1, text2, num, geo, tag).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx")
		client.HSet(ctx, "doc1", "field", "aaa", "text", "1", "numeric", 1, "geo", "1,1", "tag", "1")
		client.HSet(ctx, "doc2", "field", "aab", "text", "2", "numeric", 2, "geo", "2,2", "tag", "2")
		res1, err := client.FTSearch(ctx, "idx", "@text:aa*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(0)))
		res2, err := client.FTSearch(ctx, "idx", "@field:aa*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(2)))
		res3, err := client.FTSearchWithArgs(ctx, "idx", "*", &redis.FTSearchOptions{SortBy: []redis.FTSearchSortBy{{FieldName: "text", Desc: true}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res3.Total).To(BeEquivalentTo(int64(2)))
		Expect(res3.Docs[0].ID).To(BeEquivalentTo("doc2"))
		res4, err := client.FTSearchWithArgs(ctx, "idx", "*", &redis.FTSearchOptions{SortBy: []redis.FTSearchSortBy{{FieldName: "text", Asc: true}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res4.Total).To(BeEquivalentTo(int64(2)))
		Expect(res4.Docs[0].ID).To(BeEquivalentTo("doc1"))
		res5, err := client.FTSearchWithArgs(ctx, "idx", "*", &redis.FTSearchOptions{SortBy: []redis.FTSearchSortBy{{FieldName: "numeric", Asc: true}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res5.Docs[0].ID).To(BeEquivalentTo("doc1"))
		res6, err := client.FTSearchWithArgs(ctx, "idx", "*", &redis.FTSearchOptions{SortBy: []redis.FTSearchSortBy{{FieldName: "geo", Asc: true}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res6.Docs[0].ID).To(BeEquivalentTo("doc1"))
		res7, err := client.FTSearchWithArgs(ctx, "idx", "*", &redis.FTSearchOptions{SortBy: []redis.FTSearchSortBy{{FieldName: "tag", Asc: true}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res7.Docs[0].ID).To(BeEquivalentTo("doc1"))

	})

	It("should FTExplain", Label("search", "ftexplain"), func() {
		text1 := &redis.FieldSchema{FieldName: "f1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "f2", FieldType: redis.SearchFieldTypeText}
		text3 := &redis.FieldSchema{FieldName: "f3", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{}, text1, text2, text3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		res1, err := client.FTExplain(ctx, "txt", "@f3:f3_val @f2:f2_val @f1:f1_val").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1).ToNot(BeEmpty())

	})

	It("should FTAlias", Label("search", "ftexplain"), func() {
		text1 := &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText}
		val1, err := client.FTCreate(ctx, "testAlias", &redis.FTCreateOptions{Prefix: []interface{}{"index1:"}}, text1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val1).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "testAlias")
		val2, err := client.FTCreate(ctx, "testAlias2", &redis.FTCreateOptions{Prefix: []interface{}{"index2:"}}, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val2).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "testAlias2")

		client.HSet(ctx, "index1:lonestar", "name", "lonestar")
		client.HSet(ctx, "index2:yogurt", "name", "yogurt")

		res1, err := client.FTSearch(ctx, "testAlias", "*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("index1:lonestar"))

		aliasAddRes, err := client.FTAliasAdd(ctx, "testAlias", "mj23").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(aliasAddRes).To(BeEquivalentTo("OK"))

		res1, err = client.FTSearch(ctx, "mj23", "*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("index1:lonestar"))

		aliasUpdateRes, err := client.FTAliasUpdate(ctx, "testAlias2", "kb24").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(aliasUpdateRes).To(BeEquivalentTo("OK"))

		res3, err := client.FTSearch(ctx, "kb24", "*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res3.Docs[0].ID).To(BeEquivalentTo("index2:yogurt"))

		aliasDelRes, err := client.FTAliasDel(ctx, "mj23").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(aliasDelRes).To(BeEquivalentTo("OK"))

	})

	It("should FTCreate and FTSearch textfield, sortable and nostem ", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText, Sortable: true, NoStem: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resInfo, err := client.FTInfo(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resInfo.Attributes[0].Sortable).To(BeTrue())
		Expect(resInfo.Attributes[0].NoStem).To(BeTrue())

	})

	It("should FTAlter", Label("search", "ftcreate", "ftsearch", "ftalter"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resAlter, err := client.FTAlter(ctx, "idx1", false, []interface{}{"body", redis.SearchFieldTypeText.String()}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resAlter).To(BeEquivalentTo("OK"))

		client.HSet(ctx, "doc1", "title", "MyTitle", "body", "Some content only in the body")
		res1, err := client.FTSearch(ctx, "idx1", "only in the body").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(1)))

		_, err = client.FTSearch(ctx, "idx_not_exist", "only in the body").Result()
		Expect(err).To(HaveOccurred())
	})

	It("should FTSpellCheck", Label("search", "ftcreate", "ftsearch", "ftspellcheck"), func() {
		text1 := &redis.FieldSchema{FieldName: "f1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "f2", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "f1", "some valid content", "f2", "this is sample text")
		client.HSet(ctx, "doc2", "f1", "very important", "f2", "lorem ipsum")

		resSpellCheck, err := client.FTSpellCheck(ctx, "idx1", "impornant").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSpellCheck[0].Suggestions[0].Suggestion).To(BeEquivalentTo("important"))

		resSpellCheck2, err := client.FTSpellCheck(ctx, "idx1", "contnt").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSpellCheck2[0].Suggestions[0].Suggestion).To(BeEquivalentTo("content"))

		// test spellcheck with Levenshtein distance
		resSpellCheck3, err := client.FTSpellCheck(ctx, "idx1", "vlis").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSpellCheck3[0].Term).To(BeEquivalentTo("vlis"))

		resSpellCheck4, err := client.FTSpellCheckWithArgs(ctx, "idx1", "vlis", &redis.FTSpellCheckOptions{Distance: 2}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSpellCheck4[0].Suggestions[0].Suggestion).To(BeEquivalentTo("valid"))

		// test spellcheck include
		resDictAdd, err := client.FTDictAdd(ctx, "dict", "lore", "lorem", "lorm").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDictAdd).To(BeEquivalentTo(3))
		terms := &redis.FTSpellCheckTerms{Inclusion: "INCLUDE", Dictionary: "dict"}
		resSpellCheck5, err := client.FTSpellCheckWithArgs(ctx, "idx1", "lorm", &redis.FTSpellCheckOptions{Terms: terms}).Result()
		Expect(err).NotTo(HaveOccurred())
		lorm := resSpellCheck5[0].Suggestions
		Expect(len(lorm)).To(BeEquivalentTo(3))
		Expect(lorm[0].Score).To(BeEquivalentTo(0.5))
		Expect(lorm[1].Score).To(BeEquivalentTo(0))
		Expect(lorm[2].Score).To(BeEquivalentTo(0))

		terms2 := &redis.FTSpellCheckTerms{Inclusion: "EXCLUDE", Dictionary: "dict"}
		resSpellCheck6, err := client.FTSpellCheckWithArgs(ctx, "idx1", "lorm", &redis.FTSpellCheckOptions{Terms: terms2}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSpellCheck6).To(BeEmpty())
	})

	It("should FTDict opreations", Label("search", "ftdictdump", "ftdictdel", "ftdictadd"), func() {
		text1 := &redis.FieldSchema{FieldName: "f1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "f2", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resDictAdd, err := client.FTDictAdd(ctx, "custom_dict", "item1", "item2", "item3").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDictAdd).To(BeEquivalentTo(3))

		resDictDel, err := client.FTDictDel(ctx, "custom_dict", "item2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDictDel).To(BeEquivalentTo(1))

		resDictDump, err := client.FTDictDump(ctx, "custom_dict").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDictDump).To(BeEquivalentTo([]string{"item1", "item3"}))

		resDictDel2, err := client.FTDictDel(ctx, "custom_dict", "item1", "item3").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDictDel2).To(BeEquivalentTo(2))
	})

	It("should FTSearch phonetic matcher", Label("search", "ftsearch"), func() {
		text1 := &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "name", "Jon")
		client.HSet(ctx, "doc2", "name", "John")

		res1, err := client.FTSearch(ctx, "idx1", "Jon").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(1)))
		Expect(res1.Docs[0].Fields["name"]).To(BeEquivalentTo("Jon"))

		client.FlushDB(ctx)
		text2 := &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText, PhoneticMatcher: "dm:en"}
		val2, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val2).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "name", "Jon")
		client.HSet(ctx, "doc2", "name", "John")

		res2, err := client.FTSearch(ctx, "idx1", "Jon").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(2)))
		names := []interface{}{res2.Docs[0].Fields["name"], res2.Docs[1].Fields["name"]}
		Expect(names).To(ContainElement("Jon"))
		Expect(names).To(ContainElement("John"))
	})

	It("should FTSearch WithScores", Label("search", "ftsearch"), func() {
		text1 := &redis.FieldSchema{FieldName: "description", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "description", "The quick brown fox jumps over the lazy dog")
		client.HSet(ctx, "doc2", "description", "Quick alice was beginning to get very tired of sitting by her quick sister on the bank, and of having nothing to do.")

		res, err := client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(float64(1)))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "TFIDF"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(float64(1)))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "TFIDF.DOCNORM"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(0.14285714285714285))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "BM25"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeNumerically("<=", 0.22471909420069797))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "DISMAX"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(float64(2)))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "DOCSCORE"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(float64(1)))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "quick", &redis.FTSearchOptions{WithScores: true, Scorer: "HAMMING"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(*res.Docs[0].Score).To(BeEquivalentTo(float64(0)))
	})

	It("should FTConfigSet and FTConfigGet ", Label("search", "ftconfigget", "ftconfigset", "NonRedisEnterprise"), func() {
		val, err := client.FTConfigSet(ctx, "TIMEOUT", "100").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))

		res, err := client.FTConfigGet(ctx, "*").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res["TIMEOUT"]).To(BeEquivalentTo("100"))

		res, err = client.FTConfigGet(ctx, "TIMEOUT").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeEquivalentTo(map[string]interface{}{"TIMEOUT": "100"}))

	})

	It("should FTAggregate GroupBy ", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "body", FieldType: redis.SearchFieldTypeText}
		text3 := &redis.FieldSchema{FieldName: "parent", FieldType: redis.SearchFieldTypeText}
		num := &redis.FieldSchema{FieldName: "random_num", FieldType: redis.SearchFieldTypeNumeric}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2, text3, num).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "search", "title", "RediSearch",
			"body", "Redisearch implements a search engine on top of redis",
			"parent", "redis",
			"random_num", 10)
		client.HSet(ctx, "ai", "title", "RedisAI",
			"body", "RedisAI executes Deep Learning/Machine Learning models and managing their data.",
			"parent", "redis",
			"random_num", 3)
		client.HSet(ctx, "json", "title", "RedisJson",
			"body", "RedisJSON implements ECMA-404 The JSON Data Interchange Standard as a native data type.",
			"parent", "redis",
			"random_num", 8)

		reducer := redis.FTAggregateReducer{Reducer: redis.SearchCount}
		options := &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err := client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliascount"]).To(BeEquivalentTo("3"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchCountDistinct, Args: []interface{}{"@title"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliascount_distincttitle"]).To(BeEquivalentTo("3"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchSum, Args: []interface{}{"@random_num"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliassumrandom_num"]).To(BeEquivalentTo("21"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchMin, Args: []interface{}{"@random_num"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliasminrandom_num"]).To(BeEquivalentTo("3"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchMax, Args: []interface{}{"@random_num"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliasmaxrandom_num"]).To(BeEquivalentTo("10"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchAvg, Args: []interface{}{"@random_num"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliasavgrandom_num"]).To(BeEquivalentTo("7"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchStdDev, Args: []interface{}{"@random_num"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliasstddevrandom_num"]).To(BeEquivalentTo("3.60555127546"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchQuantile, Args: []interface{}{"@random_num", 0.5}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliasquantilerandom_num,0.5"]).To(BeEquivalentTo("8"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchToList, Args: []interface{}{"@title"}}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["__generated_aliastolisttitle"]).To(ContainElements("RediSearch", "RedisAI", "RedisJson"))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchFirstValue, Args: []interface{}{"@title"}, As: "first"}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["first"]).To(Or(BeEquivalentTo("RediSearch"), BeEquivalentTo("RedisAI"), BeEquivalentTo("RedisJson")))

		reducer = redis.FTAggregateReducer{Reducer: redis.SearchRandomSample, Args: []interface{}{"@title", 2}, As: "random"}
		options = &redis.FTAggregateOptions{GroupBy: []redis.FTAggregateGroupBy{{Fields: []interface{}{"@parent"}, Reduce: []redis.FTAggregateReducer{reducer}}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "redis", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["parent"]).To(BeEquivalentTo("redis"))
		Expect(res.Rows[0].Fields["random"]).To(Or(
			ContainElement("RediSearch"),
			ContainElement("RedisAI"),
			ContainElement("RedisJson"),
		))

	})

	It("should FTAggregate sort and limit", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "t1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "t2", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "t1", "a", "t2", "b")
		client.HSet(ctx, "doc2", "t1", "b", "t2", "a")

		options := &redis.FTAggregateOptions{SortBy: []redis.FTAggregateSortBy{{FieldName: "@t2", Asc: true}, {FieldName: "@t1", Desc: true}}}
		res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("b"))
		Expect(res.Rows[1].Fields["t1"]).To(BeEquivalentTo("a"))
		Expect(res.Rows[0].Fields["t2"]).To(BeEquivalentTo("a"))
		Expect(res.Rows[1].Fields["t2"]).To(BeEquivalentTo("b"))

		options = &redis.FTAggregateOptions{SortBy: []redis.FTAggregateSortBy{{FieldName: "@t1"}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("a"))
		Expect(res.Rows[1].Fields["t1"]).To(BeEquivalentTo("b"))

		options = &redis.FTAggregateOptions{SortBy: []redis.FTAggregateSortBy{{FieldName: "@t1"}}, SortByMax: 1}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("a"))

		options = &redis.FTAggregateOptions{SortBy: []redis.FTAggregateSortBy{{FieldName: "@t1"}}, Limit: 1, LimitOffset: 1}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("b"))

		options = &redis.FTAggregateOptions{SortBy: []redis.FTAggregateSortBy{{FieldName: "@t1"}}, Limit: 1, LimitOffset: 0}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("a"))
	})

	It("should FTAggregate load ", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "t1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "t2", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "t1", "hello", "t2", "world")

		options := &redis.FTAggregateOptions{Load: []redis.FTAggregateLoad{{Field: "t1"}}}
		res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("hello"))

		options = &redis.FTAggregateOptions{Load: []redis.FTAggregateLoad{{Field: "t2"}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t2"]).To(BeEquivalentTo("world"))

		options = &redis.FTAggregateOptions{Load: []redis.FTAggregateLoad{{Field: "t2", As: "t2alias"}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t2alias"]).To(BeEquivalentTo("world"))

		options = &redis.FTAggregateOptions{Load: []redis.FTAggregateLoad{{Field: "t1"}, {Field: "t2", As: "t2alias"}}}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("hello"))
		Expect(res.Rows[0].Fields["t2alias"]).To(BeEquivalentTo("world"))

		options = &redis.FTAggregateOptions{LoadAll: true}
		res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["t1"]).To(BeEquivalentTo("hello"))
		Expect(res.Rows[0].Fields["t2"]).To(BeEquivalentTo("world"))

		_, err = client.FTAggregateWithArgs(ctx, "idx_not_exist", "*", &redis.FTAggregateOptions{}).Result()
		Expect(err).To(HaveOccurred())
	})

	It("should FTAggregate with scorer and addscores", Label("search", "ftaggregate", "NonRedisEnterprise"), func() {
		title := &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText, Sortable: false}
		description := &redis.FieldSchema{FieldName: "description", FieldType: redis.SearchFieldTypeText, Sortable: false}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnHash: true, Prefix: []interface{}{"product:"}}, title, description).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "product:1", "title", "New Gaming Laptop", "description", "this is not a desktop")
		client.HSet(ctx, "product:2", "title", "Super Old Not Gaming Laptop", "description", "this laptop is not a new laptop but it is a laptop")
		client.HSet(ctx, "product:3", "title", "Office PC", "description", "office desktop pc")

		options := &redis.FTAggregateOptions{
			AddScores: true,
			Scorer:    "BM25",
			SortBy: []redis.FTAggregateSortBy{{
				FieldName: "@__score",
				Desc:      true,
			}},
		}

		res, err := client.FTAggregateWithArgs(ctx, "idx1", "laptop", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).ToNot(BeNil())
		Expect(len(res.Rows)).To(BeEquivalentTo(2))
		score1, err := strconv.ParseFloat(fmt.Sprintf("%s", res.Rows[0].Fields["__score"]), 64)
		Expect(err).NotTo(HaveOccurred())
		score2, err := strconv.ParseFloat(fmt.Sprintf("%s", res.Rows[1].Fields["__score"]), 64)
		Expect(err).NotTo(HaveOccurred())
		Expect(score1).To(BeNumerically(">", score2))

		optionsDM := &redis.FTAggregateOptions{
			AddScores: true,
			Scorer:    "DISMAX",
			SortBy: []redis.FTAggregateSortBy{{
				FieldName: "@__score",
				Desc:      true,
			}},
		}

		resDM, err := client.FTAggregateWithArgs(ctx, "idx1", "laptop", optionsDM).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDM).ToNot(BeNil())
		Expect(len(resDM.Rows)).To(BeEquivalentTo(2))
		score1DM, err := strconv.ParseFloat(fmt.Sprintf("%s", resDM.Rows[0].Fields["__score"]), 64)
		Expect(err).NotTo(HaveOccurred())
		score2DM, err := strconv.ParseFloat(fmt.Sprintf("%s", resDM.Rows[1].Fields["__score"]), 64)
		Expect(err).NotTo(HaveOccurred())
		Expect(score1DM).To(BeNumerically(">", score2DM))

		Expect(score1DM).To(BeEquivalentTo(float64(4)))
		Expect(score2DM).To(BeEquivalentTo(float64(1)))
		Expect(score1).NotTo(BeEquivalentTo(score1DM))
		Expect(score2).NotTo(BeEquivalentTo(score2DM))
	})

	It("should FTAggregate apply and groupby", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "PrimaryKey", FieldType: redis.SearchFieldTypeText, Sortable: true}
		num1 := &redis.FieldSchema{FieldName: "CreatedDateTimeUTC", FieldType: redis.SearchFieldTypeNumeric, Sortable: true}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		// 6 feb
		client.HSet(ctx, "doc1", "PrimaryKey", "9::362330", "CreatedDateTimeUTC", "1738823999")

		// 12 feb
		client.HSet(ctx, "doc2", "PrimaryKey", "9::362329", "CreatedDateTimeUTC", "1739342399")
		client.HSet(ctx, "doc3", "PrimaryKey", "9::362329", "CreatedDateTimeUTC", "1739353199")

		reducer := redis.FTAggregateReducer{Reducer: redis.SearchCount, As: "perDay"}

		options := &redis.FTAggregateOptions{
			Apply: []redis.FTAggregateApply{{Field: "floor(@CreatedDateTimeUTC /(60*60*24))", As: "TimestampAsDay"}},
			GroupBy: []redis.FTAggregateGroupBy{{
				Fields: []interface{}{"@TimestampAsDay"},
				Reduce: []redis.FTAggregateReducer{reducer},
			}},
			SortBy: []redis.FTAggregateSortBy{{
				FieldName: "@perDay",
				Desc:      true,
			}},
		}

		res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).ToNot(BeNil())
		Expect(len(res.Rows)).To(BeEquivalentTo(2))
		Expect(res.Rows[0].Fields["perDay"]).To(BeEquivalentTo("2"))
		Expect(res.Rows[1].Fields["perDay"]).To(BeEquivalentTo("1"))
	})

	It("should FTAggregate apply", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "PrimaryKey", FieldType: redis.SearchFieldTypeText, Sortable: true}
		num1 := &redis.FieldSchema{FieldName: "CreatedDateTimeUTC", FieldType: redis.SearchFieldTypeNumeric, Sortable: true}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "PrimaryKey", "9::362330", "CreatedDateTimeUTC", "637387878524969984")
		client.HSet(ctx, "doc2", "PrimaryKey", "9::362329", "CreatedDateTimeUTC", "637387875859270016")

		options := &redis.FTAggregateOptions{Apply: []redis.FTAggregateApply{{Field: "@CreatedDateTimeUTC * 10", As: "CreatedDateTimeUTC"}}}
		res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Rows[0].Fields["CreatedDateTimeUTC"]).To(Or(BeEquivalentTo("6373878785249699840"), BeEquivalentTo("6373878758592700416")))
		Expect(res.Rows[1].Fields["CreatedDateTimeUTC"]).To(Or(BeEquivalentTo("6373878785249699840"), BeEquivalentTo("6373878758592700416")))

	})

	It("should FTAggregate filter", Label("search", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText, Sortable: true}
		num1 := &redis.FieldSchema{FieldName: "age", FieldType: redis.SearchFieldTypeNumeric, Sortable: true}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "name", "bar", "age", "25")
		client.HSet(ctx, "doc2", "name", "foo", "age", "19")

		for _, dlc := range []int{1, 2} {
			options := &redis.FTAggregateOptions{Filter: "@name=='foo' && @age < 20", DialectVersion: dlc}
			res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Total).To(Or(BeEquivalentTo(2), BeEquivalentTo(1)))
			Expect(res.Rows[0].Fields["name"]).To(BeEquivalentTo("foo"))

			options = &redis.FTAggregateOptions{Filter: "@age > 15", DialectVersion: dlc, SortBy: []redis.FTAggregateSortBy{{FieldName: "@age"}}}
			res, err = client.FTAggregateWithArgs(ctx, "idx1", "*", options).Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Total).To(BeEquivalentTo(2))
			Expect(res.Rows[0].Fields["age"]).To(BeEquivalentTo("19"))
			Expect(res.Rows[1].Fields["age"]).To(BeEquivalentTo("25"))
		}
	})

	It("should FTSearch SkipInitialScan", Label("search", "ftsearch"), func() {
		client.HSet(ctx, "doc1", "foo", "bar")

		text1 := &redis.FieldSchema{FieldName: "foo", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{SkipInitialScan: true}, text1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err := client.FTSearch(ctx, "idx1", "@foo:bar").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(int64(0)))
	})

	It("should FTCreate json", Label("search", "ftcreate"), func() {

		text1 := &redis.FieldSchema{FieldName: "$.name", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true, Prefix: []interface{}{"king:"}}, text1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.JSONSet(ctx, "king:1", "$", `{"name": "henry"}`)
		client.JSONSet(ctx, "king:2", "$", `{"name": "james"}`)

		res, err := client.FTSearch(ctx, "idx1", "henry").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("king:1"))
		Expect(res.Docs[0].Fields["$"]).To(BeEquivalentTo(`{"name":"henry"}`))
	})

	It("should FTCreate json fields as names", Label("search", "ftcreate"), func() {

		text1 := &redis.FieldSchema{FieldName: "$.name", FieldType: redis.SearchFieldTypeText, As: "name"}
		num1 := &redis.FieldSchema{FieldName: "$.age", FieldType: redis.SearchFieldTypeNumeric, As: "just_a_number"}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.JSONSet(ctx, "doc:1", "$", `{"name": "Jon", "age": 25}`)

		res, err := client.FTSearchWithArgs(ctx, "idx1", "Jon", &redis.FTSearchOptions{Return: []redis.FTSearchReturn{{FieldName: "name"}, {FieldName: "just_a_number"}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc:1"))
		Expect(res.Docs[0].Fields["name"]).To(BeEquivalentTo("Jon"))
		Expect(res.Docs[0].Fields["just_a_number"]).To(BeEquivalentTo("25"))
	})

	It("should FTCreate CaseSensitive", Label("search", "ftcreate"), func() {

		tag1 := &redis.FieldSchema{FieldName: "t", FieldType: redis.SearchFieldTypeTag, CaseSensitive: false}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, tag1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "1", "t", "HELLO")
		client.HSet(ctx, "2", "t", "hello")

		res, err := client.FTSearch(ctx, "idx1", "@t:{HELLO}").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(2))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("1"))
		Expect(res.Docs[1].ID).To(BeEquivalentTo("2"))

		resDrop, err := client.FTDropIndex(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDrop).To(BeEquivalentTo("OK"))

		tag2 := &redis.FieldSchema{FieldName: "t", FieldType: redis.SearchFieldTypeTag, CaseSensitive: true}
		val, err = client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, tag2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err = client.FTSearch(ctx, "idx1", "@t:{HELLO}").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("1"))

	})

	It("should FTSearch ReturnFields", Label("search", "ftsearch"), func() {
		resJson, err := client.JSONSet(ctx, "doc:1", "$", `{"t": "riceratops","t2": "telmatosaurus", "n": 9072, "flt": 97.2}`).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resJson).To(BeEquivalentTo("OK"))

		text1 := &redis.FieldSchema{FieldName: "$.t", FieldType: redis.SearchFieldTypeText}
		num1 := &redis.FieldSchema{FieldName: "$.flt", FieldType: redis.SearchFieldTypeNumeric}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err := client.FTSearchWithArgs(ctx, "idx1", "*", &redis.FTSearchOptions{Return: []redis.FTSearchReturn{{FieldName: "$.t", As: "txt"}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc:1"))
		Expect(res.Docs[0].Fields["txt"]).To(BeEquivalentTo("riceratops"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "*", &redis.FTSearchOptions{Return: []redis.FTSearchReturn{{FieldName: "$.t2", As: "txt"}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc:1"))
		Expect(res.Docs[0].Fields["txt"]).To(BeEquivalentTo("telmatosaurus"))
	})

	It("should FTSynUpdate", Label("search", "ftsynupdate"), func() {

		text1 := &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "body", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnHash: true}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resSynUpdate, err := client.FTSynUpdateWithArgs(ctx, "idx1", "id1", &redis.FTSynUpdateOptions{SkipInitialScan: true}, []interface{}{"boy", "child", "offspring"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))
		client.HSet(ctx, "doc1", "title", "he is a baby", "body", "this is a test")

		resSynUpdate, err = client.FTSynUpdateWithArgs(ctx, "idx1", "id1", &redis.FTSynUpdateOptions{SkipInitialScan: true}, []interface{}{"baby"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))
		client.HSet(ctx, "doc2", "title", "he is another baby", "body", "another test")

		res, err := client.FTSearchWithArgs(ctx, "idx1", "child", &redis.FTSearchOptions{Expander: "SYNONYM"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc2"))
		Expect(res.Docs[0].Fields["title"]).To(BeEquivalentTo("he is another baby"))
		Expect(res.Docs[0].Fields["body"]).To(BeEquivalentTo("another test"))
	})

	It("should FTSynDump", Label("search", "ftsyndump"), func() {

		text1 := &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "body", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnHash: true}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resSynUpdate, err := client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"boy", "child", "offspring"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynUpdate, err = client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"baby", "child"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynUpdate, err = client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"tree", "wood"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynDump, err := client.FTSynDump(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynDump[0].Term).To(BeEquivalentTo("baby"))
		Expect(resSynDump[0].Synonyms).To(BeEquivalentTo([]string{"id1"}))
		Expect(resSynDump[1].Term).To(BeEquivalentTo("wood"))
		Expect(resSynDump[1].Synonyms).To(BeEquivalentTo([]string{"id1"}))
		Expect(resSynDump[2].Term).To(BeEquivalentTo("boy"))
		Expect(resSynDump[2].Synonyms).To(BeEquivalentTo([]string{"id1"}))
		Expect(resSynDump[3].Term).To(BeEquivalentTo("tree"))
		Expect(resSynDump[3].Synonyms).To(BeEquivalentTo([]string{"id1"}))
		Expect(resSynDump[4].Term).To(BeEquivalentTo("child"))
		Expect(resSynDump[4].Synonyms).To(Or(BeEquivalentTo([]string{"id1"}), BeEquivalentTo([]string{"id1", "id1"})))
		Expect(resSynDump[5].Term).To(BeEquivalentTo("offspring"))
		Expect(resSynDump[5].Synonyms).To(BeEquivalentTo([]string{"id1"}))

	})

	It("should FTCreate json with alias", Label("search", "ftcreate"), func() {

		text1 := &redis.FieldSchema{FieldName: "$.name", FieldType: redis.SearchFieldTypeText, As: "name"}
		num1 := &redis.FieldSchema{FieldName: "$.num", FieldType: redis.SearchFieldTypeNumeric, As: "num"}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true, Prefix: []interface{}{"king:"}}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.JSONSet(ctx, "king:1", "$", `{"name": "henry", "num": 42}`)
		client.JSONSet(ctx, "king:2", "$", `{"name": "james", "num": 3.14}`)

		res, err := client.FTSearch(ctx, "idx1", "@name:henry").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("king:1"))
		Expect(res.Docs[0].Fields["$"]).To(BeEquivalentTo(`{"name":"henry","num":42}`))

		res, err = client.FTSearch(ctx, "idx1", "@num:[0 10]").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("king:2"))
		Expect(res.Docs[0].Fields["$"]).To(BeEquivalentTo(`{"name":"james","num":3.14}`))
	})

	It("should FTCreate json with multipath", Label("search", "ftcreate"), func() {

		tag1 := &redis.FieldSchema{FieldName: "$..name", FieldType: redis.SearchFieldTypeTag, As: "name"}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true, Prefix: []interface{}{"king:"}}, tag1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.JSONSet(ctx, "king:1", "$", `{"name": "henry", "country": {"name": "england"}}`)

		res, err := client.FTSearch(ctx, "idx1", "@name:{england}").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("king:1"))
		Expect(res.Docs[0].Fields["$"]).To(BeEquivalentTo(`{"name":"henry","country":{"name":"england"}}`))
	})

	It("should FTCreate json with jsonpath", Label("search", "ftcreate"), func() {

		text1 := &redis.FieldSchema{FieldName: `$["prod:name"]`, FieldType: redis.SearchFieldTypeText, As: "name"}
		text2 := &redis.FieldSchema{FieldName: `$.prod:name`, FieldType: redis.SearchFieldTypeText, As: "name_unsupported"}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnJSON: true}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.JSONSet(ctx, "doc:1", "$", `{"prod:name": "RediSearch"}`)

		res, err := client.FTSearch(ctx, "idx1", "@name:RediSearch").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc:1"))
		Expect(res.Docs[0].Fields["$"]).To(BeEquivalentTo(`{"prod:name":"RediSearch"}`))

		res, err = client.FTSearch(ctx, "idx1", "@name_unsupported:RediSearch").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "@name:RediSearch", &redis.FTSearchOptions{Return: []redis.FTSearchReturn{{FieldName: "name"}}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(1))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("doc:1"))
		Expect(res.Docs[0].Fields["name"]).To(BeEquivalentTo("RediSearch"))

	})

	It("should FTCreate VECTOR", Label("search", "ftcreate"), func() {
		hnswOptions := &redis.FTHNSWOptions{Type: "FLOAT32", Dim: 2, DistanceMetric: "L2"}
		val, err := client.FTCreate(ctx, "idx1",
			&redis.FTCreateOptions{},
			&redis.FieldSchema{FieldName: "v", FieldType: redis.SearchFieldTypeVector, VectorArgs: &redis.FTVectorArgs{HNSWOptions: hnswOptions}}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "a", "v", "aaaaaaaa")
		client.HSet(ctx, "b", "v", "aaaabaaa")
		client.HSet(ctx, "c", "v", "aaaaabaa")

		searchOptions := &redis.FTSearchOptions{
			Return:         []redis.FTSearchReturn{{FieldName: "__v_score"}},
			SortBy:         []redis.FTSearchSortBy{{FieldName: "__v_score", Asc: true}},
			DialectVersion: 2,
			Params:         map[string]interface{}{"vec": "aaaaaaaa"},
		}
		res, err := client.FTSearchWithArgs(ctx, "idx1", "*=>[KNN 2 @v $vec]", searchOptions).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("a"))
		Expect(res.Docs[0].Fields["__v_score"]).To(BeEquivalentTo("0"))
	})

	It("should FTCreate and FTSearch text params", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "name", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "name", "Alice")
		client.HSet(ctx, "doc2", "name", "Bob")
		client.HSet(ctx, "doc3", "name", "Carol")

		res1, err := client.FTSearchWithArgs(ctx, "idx1", "@name:($name1 | $name2 )", &redis.FTSearchOptions{Params: map[string]interface{}{"name1": "Alice", "name2": "Bob"}, DialectVersion: 2}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(2)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("doc1"))
		Expect(res1.Docs[1].ID).To(BeEquivalentTo("doc2"))

	})

	It("should FTCreate and FTSearch numeric params", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "numval", FieldType: redis.SearchFieldTypeNumeric}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "numval", 101)
		client.HSet(ctx, "doc2", "numval", 102)
		client.HSet(ctx, "doc3", "numval", 103)

		res1, err := client.FTSearchWithArgs(ctx, "idx1", "@numval:[$min $max]", &redis.FTSearchOptions{Params: map[string]interface{}{"min": 101, "max": 102}, DialectVersion: 2}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(2)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("doc1"))
		Expect(res1.Docs[1].ID).To(BeEquivalentTo("doc2"))

	})

	It("should FTCreate and FTSearch geo params", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "g", FieldType: redis.SearchFieldTypeGeo}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "g", "29.69465, 34.95126")
		client.HSet(ctx, "doc2", "g", "29.69350, 34.94737")
		client.HSet(ctx, "doc3", "g", "29.68746, 34.94882")

		res1, err := client.FTSearchWithArgs(ctx, "idx1", "@g:[$lon $lat $radius $units]", &redis.FTSearchOptions{Params: map[string]interface{}{"lat": "34.95126", "lon": "29.69465", "radius": 1000, "units": "km"}, DialectVersion: 2}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(3)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("doc1"))
		Expect(res1.Docs[1].ID).To(BeEquivalentTo("doc2"))
		Expect(res1.Docs[2].ID).To(BeEquivalentTo("doc3"))

	})

	It("should FTConfigSet and FTConfigGet dialect", Label("search", "ftconfigget", "ftconfigset", "NonRedisEnterprise"), func() {
		res, err := client.FTConfigSet(ctx, "DEFAULT_DIALECT", "1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeEquivalentTo("OK"))

		defDialect, err := client.FTConfigGet(ctx, "DEFAULT_DIALECT").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(defDialect).To(BeEquivalentTo(map[string]interface{}{"DEFAULT_DIALECT": "1"}))

		res, err = client.FTConfigSet(ctx, "DEFAULT_DIALECT", "2").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeEquivalentTo("OK"))

		defDialect, err = client.FTConfigGet(ctx, "DEFAULT_DIALECT").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(defDialect).To(BeEquivalentTo(map[string]interface{}{"DEFAULT_DIALECT": "2"}))
	})

	It("should FTCreate WithSuffixtrie", Label("search", "ftcreate", "ftinfo"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err := client.FTInfo(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Attributes[0].Attribute).To(BeEquivalentTo("txt"))

		resDrop, err := client.FTDropIndex(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDrop).To(BeEquivalentTo("OK"))

		// create withsuffixtrie index - text field
		val, err = client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText, WithSuffixtrie: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err = client.FTInfo(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Attributes[0].WithSuffixtrie).To(BeTrue())

		resDrop, err = client.FTDropIndex(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resDrop).To(BeEquivalentTo("OK"))

		// create withsuffixtrie index - tag field
		val, err = client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "t", FieldType: redis.SearchFieldTypeTag, WithSuffixtrie: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		res, err = client.FTInfo(ctx, "idx1").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Attributes[0].WithSuffixtrie).To(BeTrue())
	})

	It("should test dialect 4", Label("search", "ftcreate", "ftsearch", "NonRedisEnterprise"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{
			Prefix: []interface{}{"resource:"},
		}, &redis.FieldSchema{
			FieldName: "uuid",
			FieldType: redis.SearchFieldTypeTag,
		}, &redis.FieldSchema{
			FieldName: "tags",
			FieldType: redis.SearchFieldTypeTag,
		}, &redis.FieldSchema{
			FieldName: "description",
			FieldType: redis.SearchFieldTypeText,
		}, &redis.FieldSchema{
			FieldName: "rating",
			FieldType: redis.SearchFieldTypeNumeric,
		}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))

		client.HSet(ctx, "resource:1", map[string]interface{}{
			"uuid":        "123e4567-e89b-12d3-a456-426614174000",
			"tags":        "finance|crypto|$btc|blockchain",
			"description": "Analysis of blockchain technologies & Bitcoin's potential.",
			"rating":      5,
		})
		client.HSet(ctx, "resource:2", map[string]interface{}{
			"uuid":        "987e6543-e21c-12d3-a456-426614174999",
			"tags":        "health|well-being|fitness|new-year's-resolutions",
			"description": "Health trends for the new year, including fitness regimes.",
			"rating":      4,
		})

		res, err := client.FTSearchWithArgs(ctx, "idx1", "@uuid:{$uuid}",
			&redis.FTSearchOptions{
				DialectVersion: 2,
				Params:         map[string]interface{}{"uuid": "123e4567-e89b-12d3-a456-426614174000"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(int64(1)))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("resource:1"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "@uuid:{$uuid}",
			&redis.FTSearchOptions{
				DialectVersion: 4,
				Params:         map[string]interface{}{"uuid": "123e4567-e89b-12d3-a456-426614174000"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Total).To(BeEquivalentTo(int64(1)))
		Expect(res.Docs[0].ID).To(BeEquivalentTo("resource:1"))

		client.HSet(ctx, "test:1", map[string]interface{}{
			"uuid":  "3d3586fe-0416-4572-8ce",
			"email": "adriano@acme.com.ie",
			"num":   5,
		})

		// Create the index
		ftCreateOptions := &redis.FTCreateOptions{
			Prefix: []interface{}{"test:"},
		}
		schema := []*redis.FieldSchema{
			{
				FieldName: "uuid",
				FieldType: redis.SearchFieldTypeTag,
			},
			{
				FieldName: "email",
				FieldType: redis.SearchFieldTypeTag,
			},
			{
				FieldName: "num",
				FieldType: redis.SearchFieldTypeNumeric,
			},
		}

		val, err = client.FTCreate(ctx, "idx_hash", ftCreateOptions, schema...).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("OK"))
		WaitForIndexing(client, "idx_hash")

		ftSearchOptions := &redis.FTSearchOptions{
			DialectVersion: 4,
			Params: map[string]interface{}{
				"uuid":  "3d3586fe-0416-4572-8ce",
				"email": "adriano@acme.com.ie",
			},
		}

		res, err = client.FTSearchWithArgs(ctx, "idx_hash", "@uuid:{$uuid}", ftSearchOptions).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("test:1"))
		Expect(res.Docs[0].Fields["uuid"]).To(BeEquivalentTo("3d3586fe-0416-4572-8ce"))

		res, err = client.FTSearchWithArgs(ctx, "idx_hash", "@email:{$email}", ftSearchOptions).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("test:1"))
		Expect(res.Docs[0].Fields["email"]).To(BeEquivalentTo("adriano@acme.com.ie"))

		ftSearchOptions.Params = map[string]interface{}{"num": 5}
		res, err = client.FTSearchWithArgs(ctx, "idx_hash", "@num:[5]", ftSearchOptions).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("test:1"))
		Expect(res.Docs[0].Fields["num"]).To(BeEquivalentTo("5"))
	})

	It("should FTCreate GeoShape", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "geom", FieldType: redis.SearchFieldTypeGeoShape, GeoShapeFieldType: "FLAT"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "small", "geom", "POLYGON((1 1, 1 100, 100 100, 100 1, 1 1))")
		client.HSet(ctx, "large", "geom", "POLYGON((1 1, 1 200, 200 200, 200 1, 1 1))")

		res1, err := client.FTSearchWithArgs(ctx, "idx1", "@geom:[WITHIN $poly]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"poly": "POLYGON((0 0, 0 150, 150 150, 150 0, 0 0))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1.Total).To(BeEquivalentTo(int64(1)))
		Expect(res1.Docs[0].ID).To(BeEquivalentTo("small"))

		res2, err := client.FTSearchWithArgs(ctx, "idx1", "@geom:[CONTAINS $poly]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"poly": "POLYGON((2 2, 2 50, 50 50, 50 2, 2 2))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res2.Total).To(BeEquivalentTo(int64(2)))
	})

	It("should create search index with FLOAT16 and BFLOAT16 vectors", Label("search", "ftcreate", "NonRedisEnterprise"), func() {
		val, err := client.FTCreate(ctx, "index", &redis.FTCreateOptions{},
			&redis.FieldSchema{FieldName: "float16", FieldType: redis.SearchFieldTypeVector, VectorArgs: &redis.FTVectorArgs{FlatOptions: &redis.FTFlatOptions{Type: "FLOAT16", Dim: 768, DistanceMetric: "COSINE"}}},
			&redis.FieldSchema{FieldName: "bfloat16", FieldType: redis.SearchFieldTypeVector, VectorArgs: &redis.FTVectorArgs{FlatOptions: &redis.FTFlatOptions{Type: "BFLOAT16", Dim: 768, DistanceMetric: "COSINE"}}},
		).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "index")
	})

	It("should test geoshapes query intersects and disjoint", Label("NonRedisEnterprise"), func() {
		_, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{
			FieldName:         "g",
			FieldType:         redis.SearchFieldTypeGeoShape,
			GeoShapeFieldType: "FLAT",
		}).Result()
		Expect(err).NotTo(HaveOccurred())

		client.HSet(ctx, "doc_point1", "g", "POINT (10 10)")
		client.HSet(ctx, "doc_point2", "g", "POINT (50 50)")
		client.HSet(ctx, "doc_polygon1", "g", "POLYGON ((20 20, 25 35, 35 25, 20 20))")
		client.HSet(ctx, "doc_polygon2", "g", "POLYGON ((60 60, 65 75, 70 70, 65 55, 60 60))")

		intersection, err := client.FTSearchWithArgs(ctx, "idx1", "@g:[intersects $shape]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"shape": "POLYGON((15 15, 75 15, 50 70, 20 40, 15 15))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		_assert_geosearch_result(&intersection, []string{"doc_point2", "doc_polygon1"})

		disjunction, err := client.FTSearchWithArgs(ctx, "idx1", "@g:[disjoint $shape]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"shape": "POLYGON((15 15, 75 15, 50 70, 20 40, 15 15))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		_assert_geosearch_result(&disjunction, []string{"doc_point1", "doc_polygon2"})
	})

	It("should test geoshapes query contains and within", func() {
		_, err := client.FTCreate(ctx, "idx2", &redis.FTCreateOptions{}, &redis.FieldSchema{
			FieldName:         "g",
			FieldType:         redis.SearchFieldTypeGeoShape,
			GeoShapeFieldType: "FLAT",
		}).Result()
		Expect(err).NotTo(HaveOccurred())

		client.HSet(ctx, "doc_point1", "g", "POINT (10 10)")
		client.HSet(ctx, "doc_point2", "g", "POINT (50 50)")
		client.HSet(ctx, "doc_polygon1", "g", "POLYGON ((20 20, 25 35, 35 25, 20 20))")
		client.HSet(ctx, "doc_polygon2", "g", "POLYGON ((60 60, 65 75, 70 70, 65 55, 60 60))")

		containsA, err := client.FTSearchWithArgs(ctx, "idx2", "@g:[contains $shape]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"shape": "POINT(25 25)"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		_assert_geosearch_result(&containsA, []string{"doc_polygon1"})

		containsB, err := client.FTSearchWithArgs(ctx, "idx2", "@g:[contains $shape]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"shape": "POLYGON((24 24, 24 26, 25 25, 24 24))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		_assert_geosearch_result(&containsB, []string{"doc_polygon1"})

		within, err := client.FTSearchWithArgs(ctx, "idx2", "@g:[within $shape]",
			&redis.FTSearchOptions{
				DialectVersion: 3,
				Params:         map[string]interface{}{"shape": "POLYGON((15 15, 75 15, 50 70, 20 40, 15 15))"},
			}).Result()
		Expect(err).NotTo(HaveOccurred())
		_assert_geosearch_result(&within, []string{"doc_point2", "doc_polygon1"})
	})

	It("should search missing fields", Label("search", "ftcreate", "ftsearch", "NonRedisEnterprise"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{Prefix: []interface{}{"property:"}},
			&redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText, Sortable: true},
			&redis.FieldSchema{FieldName: "features", FieldType: redis.SearchFieldTypeTag, IndexMissing: true},
			&redis.FieldSchema{FieldName: "description", FieldType: redis.SearchFieldTypeText, IndexMissing: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "property:1", map[string]interface{}{
			"title":       "Luxury Villa in Malibu",
			"features":    "pool,sea view,modern",
			"description": "A stunning modern villa overlooking the Pacific Ocean.",
		})

		client.HSet(ctx, "property:2", map[string]interface{}{
			"title":       "Downtown Flat",
			"description": "Modern flat in central Paris with easy access to metro.",
		})

		client.HSet(ctx, "property:3", map[string]interface{}{
			"title":    "Beachfront Bungalow",
			"features": "beachfront,sun deck",
		})

		res, err := client.FTSearchWithArgs(ctx, "idx1", "ismissing(@features)", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:2"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "-ismissing(@features)", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:1"))
		Expect(res.Docs[1].ID).To(BeEquivalentTo("property:3"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "ismissing(@description)", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:3"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "-ismissing(@description)", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:1"))
		Expect(res.Docs[1].ID).To(BeEquivalentTo("property:2"))
	})

	It("should search empty fields", Label("search", "ftcreate", "ftsearch", "NonRedisEnterprise"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{Prefix: []interface{}{"property:"}},
			&redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText, Sortable: true},
			&redis.FieldSchema{FieldName: "features", FieldType: redis.SearchFieldTypeTag, IndexEmpty: true},
			&redis.FieldSchema{FieldName: "description", FieldType: redis.SearchFieldTypeText, IndexEmpty: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "property:1", map[string]interface{}{
			"title":       "Luxury Villa in Malibu",
			"features":    "pool,sea view,modern",
			"description": "A stunning modern villa overlooking the Pacific Ocean.",
		})

		client.HSet(ctx, "property:2", map[string]interface{}{
			"title":       "Downtown Flat",
			"features":    "",
			"description": "Modern flat in central Paris with easy access to metro.",
		})

		client.HSet(ctx, "property:3", map[string]interface{}{
			"title":       "Beachfront Bungalow",
			"features":    "beachfront,sun deck",
			"description": "",
		})

		res, err := client.FTSearchWithArgs(ctx, "idx1", "@features:{\"\"}", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:2"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "-@features:{\"\"}", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:1"))
		Expect(res.Docs[1].ID).To(BeEquivalentTo("property:3"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "@description:''", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:3"))

		res, err = client.FTSearchWithArgs(ctx, "idx1", "-@description:''", &redis.FTSearchOptions{DialectVersion: 4, Return: []redis.FTSearchReturn{{FieldName: "id"}}, NoContent: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Docs[0].ID).To(BeEquivalentTo("property:1"))
		Expect(res.Docs[1].ID).To(BeEquivalentTo("property:2"))
	})
})

func _assert_geosearch_result(result *redis.FTSearchResult, expectedDocIDs []string) {
	ids := make([]string, len(result.Docs))
	for i, doc := range result.Docs {
		ids[i] = doc.ID
	}
	Expect(ids).To(ConsistOf(expectedDocIDs))
	Expect(result.Total).To(BeEquivalentTo(len(expectedDocIDs)))
}

// It("should FTProfile Search and Aggregate", Label("search", "ftprofile"), func() {
// 	val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "t", FieldType: redis.SearchFieldTypeText}).Result()
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(val).To(BeEquivalentTo("OK"))
// 	WaitForIndexing(client, "idx1")

// 	client.HSet(ctx, "1", "t", "hello")
// 	client.HSet(ctx, "2", "t", "world")

// 	// FTProfile Search
// 	query := redis.FTSearchQuery("hello|world", &redis.FTSearchOptions{NoContent: true})
// 	res1, err := client.FTProfile(ctx, "idx1", false, query).Result()
// 	Expect(err).NotTo(HaveOccurred())
// 	panic(res1)
// Expect(len(res1["results"].([]interface{}))).To(BeEquivalentTo(3))
// resProfile := res1["profile"].(map[interface{}]interface{})
// Expect(resProfile["Parsing time"].(float64) < 0.5).To(BeTrue())
// iterProfile0 := resProfile["Iterators profile"].([]interface{})[0].(map[interface{}]interface{})
// Expect(iterProfile0["Counter"]).To(BeEquivalentTo(2.0))
// Expect(iterProfile0["Type"]).To(BeEquivalentTo("UNION"))

// // FTProfile Aggregate
// aggQuery := redis.FTAggregateQuery("*", &redis.FTAggregateOptions{
// 	Load:  []redis.FTAggregateLoad{{Field: "t"}},
// 	Apply: []redis.FTAggregateApply{{Field: "startswith(@t, 'hel')", As: "prefix"}}})
// res2, err := client.FTProfile(ctx, "idx1", false, aggQuery).Result()
// Expect(err).NotTo(HaveOccurred())
// Expect(len(res2["results"].([]interface{}))).To(BeEquivalentTo(2))
// resProfile = res2["profile"].(map[interface{}]interface{})
// iterProfile0 = resProfile["Iterators profile"].([]interface{})[0].(map[interface{}]interface{})
// Expect(iterProfile0["Counter"]).To(BeEquivalentTo(2))
// Expect(iterProfile0["Type"]).To(BeEquivalentTo("WILDCARD"))
// })

// 	It("should FTProfile Search Limited", Label("search", "ftprofile"), func() {
// 		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "t", FieldType: redis.SearchFieldTypeText}).Result()
// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(val).To(BeEquivalentTo("OK"))
// 		WaitForIndexing(client, "idx1")

// 		client.HSet(ctx, "1", "t", "hello")
// 		client.HSet(ctx, "2", "t", "hell")
// 		client.HSet(ctx, "3", "t", "help")
// 		client.HSet(ctx, "4", "t", "helowa")

// 		// FTProfile Search
// 		query := redis.FTSearchQuery("%hell% hel*", &redis.FTSearchOptions{})
// 		res1, err := client.FTProfile(ctx, "idx1", true, query).Result()
// 		Expect(err).NotTo(HaveOccurred())
// 		resProfile := res1["profile"].(map[interface{}]interface{})
// 		iterProfile0 := resProfile["Iterators profile"].([]interface{})[0].(map[interface{}]interface{})
// 		Expect(iterProfile0["Type"]).To(BeEquivalentTo("INTERSECT"))
// 		Expect(len(res1["results"].([]interface{}))).To(BeEquivalentTo(3))
// 		Expect(iterProfile0["Child iterators"].([]interface{})[0].(map[interface{}]interface{})["Child iterators"]).To(BeEquivalentTo("The number of iterators in the union is 3"))
// 		Expect(iterProfile0["Child iterators"].([]interface{})[1].(map[interface{}]interface{})["Child iterators"]).To(BeEquivalentTo("The number of iterators in the union is 4"))
// 	})

// 	It("should FTProfile Search query params", Label("search", "ftprofile"), func() {
// 		hnswOptions := &redis.FTHNSWOptions{Type: "FLOAT32", Dim: 2, DistanceMetric: "L2"}
// 		val, err := client.FTCreate(ctx, "idx1",
// 			&redis.FTCreateOptions{},
// 			&redis.FieldSchema{FieldName: "v", FieldType: redis.SearchFieldTypeVector, VectorArgs: &redis.FTVectorArgs{HNSWOptions: hnswOptions}}).Result()
// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(val).To(BeEquivalentTo("OK"))
// 		WaitForIndexing(client, "idx1")

// 		client.HSet(ctx, "a", "v", "aaaaaaaa")
// 		client.HSet(ctx, "b", "v", "aaaabaaa")
// 		client.HSet(ctx, "c", "v", "aaaaabaa")

// 		// FTProfile Search
// 		searchOptions := &redis.FTSearchOptions{
// 			Return:         []redis.FTSearchReturn{{FieldName: "__v_score"}},
// 			SortBy:         []redis.FTSearchSortBy{{FieldName: "__v_score", Asc: true}},
// 			DialectVersion: 2,
// 			Params:         map[string]interface{}{"vec": "aaaaaaaa"},
// 		}
// 		query := redis.FTSearchQuery("*=>[KNN 2 @v $vec]", searchOptions)
// 		res1, err := client.FTProfile(ctx, "idx1", false, query).Result()
// 		Expect(err).NotTo(HaveOccurred())
// 		resProfile := res1["profile"].(map[interface{}]interface{})
// 		iterProfile0 := resProfile["Iterators profile"].([]interface{})[0].(map[interface{}]interface{})
// 		Expect(iterProfile0["Counter"]).To(BeEquivalentTo(2))
// 		Expect(iterProfile0["Type"]).To(BeEquivalentTo(redis.SearchFieldTypeVector.String()))
// 		Expect(res1["total_results"]).To(BeEquivalentTo(2))
// 		results0 := res1["results"].([]interface{})[0].(map[interface{}]interface{})
// 		Expect(results0["id"]).To(BeEquivalentTo("a"))
// 		Expect(results0["extra_attributes"].(map[interface{}]interface{})["__v_score"]).To(BeEquivalentTo("0"))
// 	})

var _ = Describe("RediSearch commands Resp 3", Label("search"), func() {
	ctx := context.TODO()
	var client *redis.Client
	var client2 *redis.Client

	BeforeEach(func() {
		client = redis.NewClient(&redis.Options{Addr: ":6379", Protocol: 3, UnstableResp3: true})
		client2 = redis.NewClient(&redis.Options{Addr: ":6379", Protocol: 3})
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("should handle FTAggregate with Unstable RESP3 Search Module and without stability", Label("search", "ftcreate", "ftaggregate"), func() {
		text1 := &redis.FieldSchema{FieldName: "PrimaryKey", FieldType: redis.SearchFieldTypeText, Sortable: true}
		num1 := &redis.FieldSchema{FieldName: "CreatedDateTimeUTC", FieldType: redis.SearchFieldTypeNumeric, Sortable: true}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, num1).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "PrimaryKey", "9::362330", "CreatedDateTimeUTC", "637387878524969984")
		client.HSet(ctx, "doc2", "PrimaryKey", "9::362329", "CreatedDateTimeUTC", "637387875859270016")

		options := &redis.FTAggregateOptions{Apply: []redis.FTAggregateApply{{Field: "@CreatedDateTimeUTC * 10", As: "CreatedDateTimeUTC"}}}
		res, err := client.FTAggregateWithArgs(ctx, "idx1", "*", options).RawResult()
		results := res.(map[interface{}]interface{})["results"].([]interface{})
		Expect(results[0].(map[interface{}]interface{})["extra_attributes"].(map[interface{}]interface{})["CreatedDateTimeUTC"]).
			To(Or(BeEquivalentTo("6373878785249699840"), BeEquivalentTo("6373878758592700416")))
		Expect(results[1].(map[interface{}]interface{})["extra_attributes"].(map[interface{}]interface{})["CreatedDateTimeUTC"]).
			To(Or(BeEquivalentTo("6373878785249699840"), BeEquivalentTo("6373878758592700416")))

		rawVal := client.FTAggregateWithArgs(ctx, "idx1", "*", options).RawVal()
		rawValResults := rawVal.(map[interface{}]interface{})["results"].([]interface{})
		Expect(err).NotTo(HaveOccurred())
		Expect(rawValResults[0]).To(Or(BeEquivalentTo(results[0]), BeEquivalentTo(results[1])))
		Expect(rawValResults[1]).To(Or(BeEquivalentTo(results[0]), BeEquivalentTo(results[1])))

		// Test with UnstableResp3 false
		Expect(func() {
			options = &redis.FTAggregateOptions{Apply: []redis.FTAggregateApply{{Field: "@CreatedDateTimeUTC * 10", As: "CreatedDateTimeUTC"}}}
			rawRes, _ := client2.FTAggregateWithArgs(ctx, "idx1", "*", options).RawResult()
			rawVal = client2.FTAggregateWithArgs(ctx, "idx1", "*", options).RawVal()
			Expect(rawRes).To(BeNil())
			Expect(rawVal).To(BeNil())
		}).Should(Panic())

	})

	It("should handle FTInfo with Unstable RESP3 Search Module and without stability", Label("search", "ftcreate", "ftinfo"), func() {
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText, Sortable: true, NoStem: true}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resInfo, err := client.FTInfo(ctx, "idx1").RawResult()
		Expect(err).NotTo(HaveOccurred())
		attributes := resInfo.(map[interface{}]interface{})["attributes"].([]interface{})
		flags := attributes[0].(map[interface{}]interface{})["flags"].([]interface{})
		Expect(flags).To(ConsistOf("SORTABLE", "NOSTEM"))

		valInfo := client.FTInfo(ctx, "idx1").RawVal()
		attributes = valInfo.(map[interface{}]interface{})["attributes"].([]interface{})
		flags = attributes[0].(map[interface{}]interface{})["flags"].([]interface{})
		Expect(flags).To(ConsistOf("SORTABLE", "NOSTEM"))

		// Test with UnstableResp3 false
		Expect(func() {
			rawResInfo, _ := client2.FTInfo(ctx, "idx1").RawResult()
			rawValInfo := client2.FTInfo(ctx, "idx1").RawVal()
			Expect(rawResInfo).To(BeNil())
			Expect(rawValInfo).To(BeNil())
		}).Should(Panic())
	})

	It("should handle FTSpellCheck with Unstable RESP3 Search Module and without stability", Label("search", "ftcreate", "ftspellcheck"), func() {
		text1 := &redis.FieldSchema{FieldName: "f1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "f2", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		client.HSet(ctx, "doc1", "f1", "some valid content", "f2", "this is sample text")
		client.HSet(ctx, "doc2", "f1", "very important", "f2", "lorem ipsum")

		resSpellCheck, err := client.FTSpellCheck(ctx, "idx1", "impornant").RawResult()
		valSpellCheck := client.FTSpellCheck(ctx, "idx1", "impornant").RawVal()
		Expect(err).NotTo(HaveOccurred())
		Expect(valSpellCheck).To(BeEquivalentTo(resSpellCheck))
		results := resSpellCheck.(map[interface{}]interface{})["results"].(map[interface{}]interface{})
		Expect(results["impornant"].([]interface{})[0].(map[interface{}]interface{})["important"]).To(BeEquivalentTo(0.5))

		// Test with UnstableResp3 false
		Expect(func() {
			rawResSpellCheck, _ := client2.FTSpellCheck(ctx, "idx1", "impornant").RawResult()
			rawValSpellCheck := client2.FTSpellCheck(ctx, "idx1", "impornant").RawVal()
			Expect(rawResSpellCheck).To(BeNil())
			Expect(rawValSpellCheck).To(BeNil())
		}).Should(Panic())
	})

	It("should handle FTSearch with Unstable RESP3 Search Module and without stability", Label("search", "ftcreate", "ftsearch"), func() {
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{StopWords: []interface{}{"foo", "bar", "baz"}}, &redis.FieldSchema{FieldName: "txt", FieldType: redis.SearchFieldTypeText}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		client.HSet(ctx, "doc1", "txt", "foo baz")
		client.HSet(ctx, "doc2", "txt", "hello world")
		res1, err := client.FTSearchWithArgs(ctx, "txt", "foo bar", &redis.FTSearchOptions{NoContent: true}).RawResult()
		val1 := client.FTSearchWithArgs(ctx, "txt", "foo bar", &redis.FTSearchOptions{NoContent: true}).RawVal()
		Expect(err).NotTo(HaveOccurred())
		Expect(val1).To(BeEquivalentTo(res1))
		totalResults := res1.(map[interface{}]interface{})["total_results"]
		Expect(totalResults).To(BeEquivalentTo(int64(0)))
		res2, err := client.FTSearchWithArgs(ctx, "txt", "foo bar hello world", &redis.FTSearchOptions{NoContent: true}).RawResult()
		Expect(err).NotTo(HaveOccurred())
		totalResults2 := res2.(map[interface{}]interface{})["total_results"]
		Expect(totalResults2).To(BeEquivalentTo(int64(1)))

		// Test with UnstableResp3 false
		Expect(func() {
			rawRes2, _ := client2.FTSearchWithArgs(ctx, "txt", "foo bar hello world", &redis.FTSearchOptions{NoContent: true}).RawResult()
			rawVal2 := client2.FTSearchWithArgs(ctx, "txt", "foo bar hello world", &redis.FTSearchOptions{NoContent: true}).RawVal()
			Expect(rawRes2).To(BeNil())
			Expect(rawVal2).To(BeNil())
		}).Should(Panic())
	})
	It("should handle FTSynDump with Unstable RESP3 Search Module and without stability", Label("search", "ftsyndump"), func() {
		text1 := &redis.FieldSchema{FieldName: "title", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "body", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "idx1", &redis.FTCreateOptions{OnHash: true}, text1, text2).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "idx1")

		resSynUpdate, err := client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"boy", "child", "offspring"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynUpdate, err = client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"baby", "child"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynUpdate, err = client.FTSynUpdate(ctx, "idx1", "id1", []interface{}{"tree", "wood"}).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(resSynUpdate).To(BeEquivalentTo("OK"))

		resSynDump, err := client.FTSynDump(ctx, "idx1").RawResult()
		valSynDump := client.FTSynDump(ctx, "idx1").RawVal()
		Expect(err).NotTo(HaveOccurred())
		Expect(valSynDump).To(BeEquivalentTo(resSynDump))
		Expect(resSynDump.(map[interface{}]interface{})["baby"]).To(BeEquivalentTo([]interface{}{"id1"}))

		// Test with UnstableResp3 false
		Expect(func() {
			rawResSynDump, _ := client2.FTSynDump(ctx, "idx1").RawResult()
			rawValSynDump := client2.FTSynDump(ctx, "idx1").RawVal()
			Expect(rawResSynDump).To(BeNil())
			Expect(rawValSynDump).To(BeNil())
		}).Should(Panic())
	})

	It("should test not affected Resp 3 Search method - FTExplain", Label("search", "ftexplain"), func() {
		text1 := &redis.FieldSchema{FieldName: "f1", FieldType: redis.SearchFieldTypeText}
		text2 := &redis.FieldSchema{FieldName: "f2", FieldType: redis.SearchFieldTypeText}
		text3 := &redis.FieldSchema{FieldName: "f3", FieldType: redis.SearchFieldTypeText}
		val, err := client.FTCreate(ctx, "txt", &redis.FTCreateOptions{}, text1, text2, text3).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeEquivalentTo("OK"))
		WaitForIndexing(client, "txt")
		res1, err := client.FTExplain(ctx, "txt", "@f3:f3_val @f2:f2_val @f1:f1_val").Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(res1).ToNot(BeEmpty())

		// Test with UnstableResp3 false
		Expect(func() {
			res2, err := client2.FTExplain(ctx, "txt", "@f3:f3_val @f2:f2_val @f1:f1_val").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(res2).ToNot(BeEmpty())
		}).ShouldNot(Panic())
	})
})
