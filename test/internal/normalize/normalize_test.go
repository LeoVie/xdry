package normalize

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"x-dry-go/internal/config"
	"x-dry-go/internal/normalize"
	"x-dry-go/test/_mocks/cli"
)

func TestNormalizeErrorsWhenNoNormalizerFoundForFileExtension(t *testing.T) {
	want := fmt.Errorf("no normalizer found for file extension '.txt'")

	err, _ := normalize.Normalize("foo.txt", []config.Normalizer{}, cli.NewMockCommandExecutor(gomock.NewController(t)))

	assert.Equal(t, want, err)
}

func TestNormalizeErrorsWhenNormalizerErrors(t *testing.T) {
	want := fmt.Errorf("error")

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	commandExecutor := cli.NewMockCommandExecutor(ctrl)

	commandExecutor.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("", fmt.Errorf("error"))

	normalizers := []config.Normalizer{
		{
			Level:     1,
			Extension: ".txt",
			Command:   "pwd",
			Args:      []string{},
		},
	}
	err, _ := normalize.Normalize("foo.txt", normalizers, commandExecutor)

	assert.Equal(t, want, err)
}

func TestNormalize(t *testing.T) {
	want := "output of the normalizer command"

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	commandExecutor := cli.NewMockCommandExecutor(ctrl)

	commandExecutor.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("output of the normalizer command", nil)

	normalizers := []config.Normalizer{
		{
			Level:     1,
			Extension: ".txt",
			Command:   "pwd",
			Args:      []string{},
		},
	}
	err, output := normalize.Normalize("foo.txt", normalizers, commandExecutor)

	assert.Nil(t, err)
	assert.Equal(t, want, output)
}
