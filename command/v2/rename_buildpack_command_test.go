package v2_test

import (
	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/flag"
	. "code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/v2fakes"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("rename buildpack command", func() {
	var (
		cmd             RenameBuildpackCommand
		fakeActor       *v2fakes.FakeRenameBuildpackActor
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		testUI          *ui.UI
		executeErr      error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)
		fakeActor = new(v2fakes.FakeRenameBuildpackActor)

		cmd = RenameBuildpackCommand{
			UI:          testUI,
			Actor:       fakeActor,
			Config:      fakeConfig,
			SharedActor: fakeSharedActor,
		}
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("when checking the target fails", func() {
		var binaryName string

		BeforeEach(func() {
			binaryName = "faceman"
			fakeConfig.BinaryNameReturns(binaryName)
			fakeSharedActor.CheckTargetReturns(actionerror.NotLoggedInError{BinaryName: binaryName})
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError(actionerror.NotLoggedInError{BinaryName: "faceman"}))
			Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
			checkTargetedOrg, checkTargetedSpace := fakeSharedActor.CheckTargetArgsForCall(0)
			Expect(checkTargetedOrg).To(BeFalse())
			Expect(checkTargetedSpace).To(BeFalse())
		})
	})

	Context("when checking the target succeeds", func() {
		var (
			oldName string
			newName string
		)

		BeforeEach(func() {
			oldName = "some-old-name"
			newName = "some-new-name"

			cmd.RequiredArgs = flag.RenameBuildpackArgs{
				OldBuildpackName: oldName,
				NewBuildpackName: newName,
			}
		})

		Context("when renaming to a unique buildpack name", func() {
			BeforeEach(func() {
				fakeActor.RenameBuildpackReturns(v2action.Warnings{"warning1", "warning2"}, nil)
			})

			It("successfully renames the buildpack and displays any warnings", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeActor.RenameBuildpackCallCount()).To(Equal(1))

				oldBuildpackName, newBuildpackName := fakeActor.RenameBuildpackArgsForCall(0)
				Expect(oldBuildpackName).To(Equal(oldName))
				Expect(newBuildpackName).To(Equal(newName))

				Expect(testUI.Out).To(Say("Renaming buildpack %s to %s...", oldName, newName))
				Expect(testUI.Err).To(Say("warning1"))
				Expect(testUI.Err).To(Say("warning2"))
			})
		})

		Context("when the actor errors because a buildpack with the desired name already exists", func() {
			BeforeEach(func() {
				fakeActor.RenameBuildpackReturns(
					v2action.Warnings{"warning1", "warning2"},
					actionerror.BuildpackNameTakenError{Name: newName})
			})

			It("returns an error and prints warnings", func() {
				Expect(executeErr).To(MatchError(actionerror.BuildpackNameTakenError{Name: newName}))
				Expect(testUI.Err).To(Say("warning1"))
				Expect(testUI.Err).To(Say("warning2"))
			})
		})
	})
})
