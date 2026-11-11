package main

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

func (a *app) aboutLines(c *ui.Context, t *ui.Theme, text string) {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ui.Text(c, line).FontSize(a.fs(11)).TextColor(t.TextMuted)
	}
}

func (a *app) ensureAboutCommands() {
	L := a.L()
	if a.brewCmd == "" {
		a.brewCmd = "brew install " + aria2NextBrewTap
	}
	if a.dockerCmd == "" {
		a.dockerCmd = L.InstallDockerCmd
	}
}

func (a *app) copyableCommand(c *ui.Context, cmd *string) {
	L := a.L()
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.TextInput(c, cmd).Grow(1).MinWidth(0)
		if ui.Button(c, L.Copy).Padding(4, 10).Clicked() {
			c.WriteClipboard(*cmd)
			a.flash(L.Copied)
		}
	})
}

func (a *app) aboutBody(c *ui.Context, t *ui.Theme) {
	L := a.L()
	a.ensureAboutCommands()

	ui.Column(c).Grow(1).Children(func() {
		ui.Row(c).Padding(10, 14).Gap(10).AlignItems(ui.Center).Background(t.Surface).Children(func() {
			ui.Text(c, L.About).FontSize(a.fs(15)).FontWeight(600).Grow(1)
		})
		ui.Box(c).Height(1).Background(t.Border)

		ui.Scroll(c).Grow(1).Children(func() {
			ui.Column(c).Padding(12, 16).Gap(14).Children(func() {
				a.uiPrefRow(c, t, L.AboutApp, func() {
					ui.Text(c, "aria2-next-mygo").FontSize(a.fs(14)).FontWeight(600)
					ui.Textf(c, "%s  %s", L.AppVersion, appVersion).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Textf(c, "%s  %s", L.HostPlatform, hostPlatformLabel()).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Text(c, L.AboutAppDesc).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Link(c, appRepoURL, appRepoURL).FontSize(a.fs(12))
					ui.Text(c, L.AboutAppFramework).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Link(c, mygoRepoURL, mygoRepoURL).FontSize(a.fs(12))
				})

				a.uiPrefRow(c, t, L.AboutEngine, func() {
					ui.Text(c, "aria2-next").FontSize(a.fs(14)).FontWeight(600)
					ui.Text(c, L.AboutEngineDesc).FontSize(a.fs(12)).TextColor(t.TextMuted)
					ui.Link(c, aria2NextLatestURL, aria2NextLatestURL).FontSize(a.fs(12))
				})

				a.uiPrefRow(c, t, L.InstallEngine, func() {
					ui.Text(c, L.InstallHint).FontSize(a.fs(12)).TextColor(t.TextMuted)

					ui.Text(c, L.InstallBrew).FontSize(a.fs(12)).FontWeight(600)
					a.copyableCommand(c, &a.brewCmd)
					ui.Text(c, L.InstallBrewNote).FontSize(a.fs(11)).TextColor(t.TextMuted)

					ui.Box(c).Height(6)
					ui.Text(c, L.InstallDocker).FontSize(a.fs(12)).FontWeight(600)
					a.copyableCommand(c, &a.dockerCmd)
				})

				a.uiPrefRow(c, t, L.AfterInstall, func() {
					a.aboutLines(c, t, L.AfterInstallSteps)
				})
			})
		})
	})
}
