package step

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	stepId = "save-dart-cache"

	// Cache key template
	// OS + Arch: guarantees unique cache per stack
	// The cached files are in the home folder, so absolute paths are not portable between different stacks.
	key = `{{ .OS }}-{{ .Arch }}-dart-cache-{{ checksum "pubspec.lock" }}`

	// Cached path
	path = "~/.pub-cache"
)

type Input struct {
	Verbose          bool `env:"verbose,required"`
	CompressionLevel int  `env:"compression_level,range[1..19]"`
}

type SaveCacheStep struct {
	logger       log.Logger
	inputParser  stepconf.InputParser
	pathChecker  pathutil.PathChecker
	pathProvider pathutil.PathProvider
	pathModifier pathutil.PathModifier
	envRepo      env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	pathChecker pathutil.PathChecker,
	pathProvider pathutil.PathProvider,
	pathModifier pathutil.PathModifier,
	envRepo env.Repository,
) SaveCacheStep {
	return SaveCacheStep{
		logger:       logger,
		inputParser:  inputParser,
		pathChecker:  pathChecker,
		pathProvider: pathProvider,
		pathModifier: pathModifier,
		envRepo:      envRepo,
	}
}

func (step SaveCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache key: %s", key)
	step.logger.Printf("Cache paths:")
	step.logger.Printf(path)
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker)
	return saver.Save(cache.SaveCacheInput{
		StepId:           stepId,
		Verbose:          input.Verbose,
		Key:              key,
		Paths:            []string{path},
		IsKeyUnique:      true,
		CompressionLevel: input.CompressionLevel,
	})
}
