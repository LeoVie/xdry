package normalize

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"testing"
	"x-dry-go/src/_mocks/cli"
	"x-dry-go/src/internal/config"
)

func TestNormalizeErrorsWhenNoNormalizerFoundForFileExtension(t *testing.T) {
	g := NewGomegaWithT(t)

	want := fmt.Errorf("no normalizer found for file extension '.txt'")

	err, _ := Normalize("foo.txt", make(map[string]config.Normalizer), cli.NewMockCommandExecutor(gomock.NewController(t)))

	g.Expect(err).Should(Equal(want))
}

func TestNormalizeErrorsWhenNormalizerErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	want := fmt.Errorf("error")

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	commandExecutor := cli.NewMockCommandExecutor(ctrl)

	commandExecutor.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("", fmt.Errorf("error"))

	normalizers := make(map[string]config.Normalizer)
	normalizers[".txt"] = config.Normalizer{
		Level:     1,
		Extension: ".txt",
		Command:   "pwd",
		Args:      []string{},
	}
	err, _ := Normalize("foo.txt", normalizers, commandExecutor)

	g.Expect(err).Should(Equal(want))
}

func TestNormalize(t *testing.T) {
	g := NewGomegaWithT(t)

	want := "output of the normalizer command"

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	commandExecutor := cli.NewMockCommandExecutor(ctrl)

	commandExecutor.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("output of the normalizer command", nil)

	normalizers := make(map[string]config.Normalizer)
	normalizers[".txt"] = config.Normalizer{
		Level:     1,
		Extension: ".txt",
		Command:   "pwd",
		Args:      []string{},
	}
	err, output := Normalize("foo.txt", normalizers, commandExecutor)

	g.Expect(err).Should(BeNil())
	g.Expect(output).Should(Equal(want))
}
