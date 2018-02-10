package sessions_test

import (
	"fmt"
	"testing"

	"github.com/ld100/goblet/test/util"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"
)

var Cleaner = dbcleaner.New()

type ExampleSuite struct {
	suite.Suite
}

func (suite *ExampleSuite) SetupSuite() {
	// Init and set mysql cleanup engine
	dsn := fmt.Sprintf("host=%v user=%v dbname=%v sslmode=disable password=%v port=%v", "postgres", "postgres","goblet_development", "", 5432)
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