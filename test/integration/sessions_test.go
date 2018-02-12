package sessions_test

import (
	"fmt"
	"testing"

	"github.com/ld100/goblet/test/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

var Cleaner = dbcleaner.New()

type ExampleSuite struct {
	suite.Suite
}

func (suite *ExampleSuite) SetupSuite() {
	// Init and set mysql cleanup engine
	dsn := fmt.Sprintf("host=%v user=%v dbname=%v sslmode=disable password=%v port=%v", "postgres", "postgres", "goblet_development", "", 5432)
	postgres := engine.NewPostgresEngine(dsn)
	Cleaner.SetEngine(postgres)
}

func (suite *ExampleSuite) SetupTest() {
	Cleaner.Acquire("users")
}

func (suite *ExampleSuite) TearDownTest() {
	Cleaner.Clean("users")
}

func (suite *ExampleSuite) TestSomething() {
	// Have some meaningful test
	suite.Equal(true, true)
	fmt.Println(util.VERSION)
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(ExampleSuite))
}
