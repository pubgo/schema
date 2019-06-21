package schema_test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/pubgo/errors"
	"log"
	"strings"

	"github.com/pubgo/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testParams struct {
	DriverName string
	ConnStr    string

	CreateDDL []string
	DropDDL   []string
	DropFn    func()

	TableExpRes []string
	ViewExpRes  []string

	TableNameExpRes string
	ViewNameExpRes  string
}

func SchemaTestRunner(params *testParams) {
	defer errors.Debug()

	setup := func() (*sql.DB, func()) {
		db, err := sql.Open(params.DriverName, params.ConnStr)
		if err != nil {
			log.Fatalf("sql.Open error %v", err)
		}

		for _, ddl := range params.CreateDDL {
			_, err = db.Exec(ddl)
			if err != nil {
				// log.Fatalf("db.Exec (create) error %v", err)
				log.Printf("db.Exec (create) error %v exec %s", err, ddl)
			}
		}

		doneFn := func() {
			for _, ddl := range params.DropDDL {
				_, err = db.Exec(ddl)
				if err != nil {
					// log.Fatalf("db.Exec (drop) error %v", err)
					log.Printf("db.Exec (drop) error %v exec %s", err, ddl)
				}
			}
			err = db.Close()
			if err != nil {
				log.Printf("db.Close error %v", err)
			}
			if params.DropFn != nil {
				params.DropFn()
			}
		}

		return db, doneFn
	}

	Describe("Table", func() {
		defer errors.Debug()

		It("should return the column type info for an existing table", func() {
			db, done := setup()
			defer done()
			ci:= schema.Table(db, "web_resource")
			Expect(ci).To(HaveLen(10))
			var list []string
			for _, c := range ci {
				list = append(list, c.Name()+" "+c.DatabaseTypeName())
			}
			Expect(list).To(Equal(params.TableExpRes))
		})
		It("should return an error for a non-existing table", func() {
			db, done := setup()
			defer done()
			schema.Table(db, "XXX-NO-SUCH-TABLE-XXX")
		})
	})

	Describe("TableNames", func() {
		defer errors.Debug()

		It("should return the table names", func() {
			db, done := setup()
			defer done()

			// err := oraDump(db)
			// Expect(err).To(BeNil())

			sn := schema.TableNames(db)
			Expect(sn).To(HaveLen(1))
			Expect(sn).To(Equal([]string{params.TableNameExpRes}))
		})
	})

	Describe("Tables", func() {
		defer errors.Debug()

		It("should return the column type info for all tables", func() {
			db, done := setup()
			defer done()
			sc:= schema.Tables(db)
			Expect(sc).To(HaveLen(1))
			ci, ok := sc[params.TableNameExpRes]
			Expect(ok).To(BeTrue())
			Expect(ci).To(HaveLen(10))
		})
	})

	Describe("View", func() {
		defer errors.Debug()

		It("should return the column type info for the view", func() {
			db, done := setup()
			defer done()
			ci:= schema.View(db, "web_resource_view")
			Expect(ci).To(HaveLen(2))
			var list []string
			for _, c := range ci {
				list = append(list, c.Name()+" "+c.DatabaseTypeName())
			}
			Expect(list).To(Equal(params.ViewExpRes))
		})
	})

	Describe("ViewNames", func() {
		defer errors.Debug()

		It("should return the view names", func() {
			db, done := setup()
			defer done()
			sn:= schema.ViewNames(db)
			Expect(sn).To(HaveLen(1))
			Expect(sn).To(Equal([]string{params.ViewNameExpRes}))
		})
	})

	Describe("Views", func() {
		defer errors.Debug()

		It("should return the column type info for all views", func() {
			db, done := setup()
			defer done()
			sc:= schema.Views(db)
			Expect(sc).To(HaveLen(1))
			ci, ok := sc[params.ViewNameExpRes]
			Expect(ok).To(BeTrue())
			Expect(ci).To(HaveLen(2))
		})
	})

}

var _ = Describe("schema", func() {
	Context("using a fake db driver", func() {
		defer errors.Debug()

		sql.Register("fakedb", fakeDb{})
		db, _ := sql.Open("fakedb", "")

		It("should return nils for every method", func() {
			ci:= schema.Table(db, "web_resource")
			Expect(ci).To(BeNil())

			tn := schema.TableNames(db)
			Expect(tn).To(BeNil())

			ta:= schema.Tables(db)
			Expect(ta).To(BeNil())

			ci = schema.View(db, "web_resource")
			Expect(ci).To(BeNil())

			vn := schema.ViewNames(db)
			Expect(vn).To(BeNil())

			vw:= schema.Views(db)
			Expect(vw).To(BeNil())
		})
	})
})

type fakeDb struct{}

func (_ fakeDb) Open(name string) (driver.Conn, error) {
	return nil, nil
}

// pack a string, normalising its whitespace.
func pack(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
