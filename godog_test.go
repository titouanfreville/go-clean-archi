// +build test_func

package skeleton_test

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	googlePubSub "cloud.google.com/go/pubsub"
	"github.com/elmagician/godog"
	"github.com/elmagician/godog/colors"
	"github.com/elmagician/kactus/features/definitions"
	"github.com/elmagician/kactus/features/interfaces/api"
	kactusDB "github.com/elmagician/kactus/features/interfaces/database"
	"github.com/elmagician/kactus/features/interfaces/fixtures"
	"github.com/elmagician/kactus/features/interfaces/picker"
	"github.com/elmagician/kactus/features/interfaces/pubsub"
	"google.golang.org/api/option"
)

const emulatorEnv = "PUBSUB_EMULATOR_HOST"

var (
	pickerV2 *picker.Picker
	pgDb     *kactusDB.Postgres
	fix      *fixtures.Fixtures
	gcp      *pubsub.Google
	client   *api.Client
)

var PGStore *database.PGStore

func FeatureContext(s *godog.ScenarioContext) {
	// V3 step installation
	definitions.InstallDebug(s)
	definitions.InstallPicker(s, pickerV2)
	definitions.InstallPostgres(s, pgDb)
	definitions.InstallFixtures(s, fix)
	definitions.InstallGooglePubsub(s, gcp)
	definitions.InstallAPI(s, client)

	// V2 requiring V3 during transition
	// extendApi.Install(s, nil, pickerV2)

	// Auto reset instance before scenario
	s.BeforeScenario(func(_ *godog.Scenario) {
		pgDb.Reset()
		fix.Reset()
		gcp.Reset()
		client.Reset()
	})
}

func suiteContext(s *godog.TestSuiteContext) {
	// retrieving config to initialize Pubsub && Postgres instance for tests
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// initializing dependencies instance
	PGStore = database.NewPGStore(&conf.External.Database)

	cli, err := initPubsubClient(context.Background(), conf)
	if err != nil {
		log.Fatal(err)
	}

	// initializing Kactus services
	pickerV2 = picker.New()
	pickerV2.Reset() // reset picker to setup log

	pgDb = kactusDB.NewPostgres(pickerV2, kactusDB.PostgresInfo{Key: "skeleton", DB: PGStore.Db.DB})
	fix = fixtures.New(pickerV2).WithBasePath("test")
	gcp = pubsub.NewGoogle(pickerV2, pubsub.GoogleInfo{Client: cli, Key: "client"})
	client, err = api.New(pickerV2.This(), true)
	if err != nil {
		log.Fatal(err)
	}

	// reset instance to ensure logger is set up
	pickerV2.Reset()
	pgDb.Reset()
	fix.Reset()
	gcp.Reset()
	client.Reset()
	client.Debug()

	s.AfterSuite(func() {
		PGStore.ClosePool()
		_ = cli.Close()
	})
}

var opts = godog.Options{
	Output:      colors.Colored(os.Stdout),
	Format:      "pretty", // can define default values
	Concurrency: 0,
	Randomize:   -1,
	Strict:      true,
	Paths:       []string{"test"},
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()

	if len(flag.Args()) > 0 {
		opts.Paths = flag.Args()
	}

	st := m.Run()

	status := godog.TestSuite{
		Name:                 "skeleton",
		TestSuiteInitializer: suiteContext,
		ScenarioInitializer:  FeatureContext,
		Options:              &opts,
	}.Run()

	if st > status {
		status = st
	}

	os.Exit(status)
}

func initPubsubClient(ctx context.Context, config *config.Config, opts ...option.ClientOption) (*googlePubSub.Client, error) {
	if config.External.PubSub.EmulatorHost != "" {
		if err := os.Setenv(emulatorEnv, config.External.PubSub.EmulatorHost); err != nil {
			return nil, err
		}
	}

	if config.External.PubSub.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(config.External.PubSub.CredentialsPath))
	}

	return googlePubSub.NewClient(ctx, config.External.PubSub.ProjectID, opts...)
}
