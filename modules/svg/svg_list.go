package svg

func generateSVGIconMap() map[string]string {
	const cloudfrontBaseURL = "https://s21gfi7kzrpyzs.cloudfront.net/"

	return map[string]string{
		"telegram":                        "public/assets/img/telegram.png",
        "gogs":                        "public/assets/img/emoji/gogs.png",
        "404":                        "public/assets/img/404.png",
        "avatar_default":                        "public/assets/img/avatar_default.png",
        "apple-touch-icon":                        "public/assets/img/apple-touch-icon.png",
        "failed":                        "public/assets/img/failed.png",
        "packagist":                        "public/assets/img/packagist.png",
        "500":                        "public/assets/img/500.png",
        "gitea":                        "public/assets/img/emoji/gitea.png",
        "checkmark":                        "public/assets/img/checkmark.png",
        "favicon":                        "public/assets/img/favicon.svg",
        "msteams":                        "public/assets/img/msteams.png",
        "wechatwork":                        "public/assets/img/wechatwork.png",
        "logo":                        "public/assets/img/logo.svg",
        "slack":                        "public/assets/img/slack.png",
        "openid-16x16":                        "public/assets/img/openid-16x16.png",
        "discord":                        "public/assets/img/discord.png",
        "feishu":                        "public/assets/img/feishu.png",
        "repo_default":                        "public/assets/img/repo_default.png",
        "loading":                        "public/assets/img/loading.png",
        "dingtalk":                        "public/assets/img/dingtalk.ico",
        "octicon-three-bars":                        "public/assets/img/svg/octicon-three-bars.svg",
        "octicon-cloud":                        "public/assets/img/svg/octicon-cloud.svg",
        "octicon-sort-asc":                        "public/assets/img/svg/octicon-sort-asc.svg",
        "gitea-yandex":                        "public/assets/img/svg/gitea-yandex.svg",
        "octicon-duplicate":                        "public/assets/img/svg/octicon-duplicate.svg",
        "gitea-google":                        "public/assets/img/svg/gitea-google.svg",
        "octicon-alert-fill":                        "public/assets/img/svg/octicon-alert-fill.svg",
        "octicon-package-dependencies":                        "public/assets/img/svg/octicon-package-dependencies.svg",
        "octicon-file":                        "public/assets/img/svg/octicon-file.svg",
        "octicon-globe":                        "public/assets/img/svg/octicon-globe.svg",
        "gitea-azureadv2":                        "public/assets/img/svg/gitea-azureadv2.svg",
        "octicon-issue-tracked-by":                        "public/assets/img/svg/octicon-issue-tracked-by.svg",
        "octicon-file-symlink-file":                        "public/assets/img/svg/octicon-file-symlink-file.svg",
        "octicon-mirror":                        "public/assets/img/svg/octicon-mirror.svg",
        "octicon-sign-out":                        "public/assets/img/svg/octicon-sign-out.svg",
        "octicon-pin-slash":                        "public/assets/img/svg/octicon-pin-slash.svg",
        "octicon-location":                        "public/assets/img/svg/octicon-location.svg",
        "octicon-issue-reopened":                        "public/assets/img/svg/octicon-issue-reopened.svg",
        "gitea-whitespace":                        "public/assets/img/svg/gitea-whitespace.svg",
        "octicon-file-moved":                        "public/assets/img/svg/octicon-file-moved.svg",
        "octicon-sliders":                        "public/assets/img/svg/octicon-sliders.svg",
        "octicon-heading":                        "public/assets/img/svg/octicon-heading.svg",
        "octicon-stack":                        "public/assets/img/svg/octicon-stack.svg",
        "octicon-issue-draft":                        "public/assets/img/svg/octicon-issue-draft.svg",
        "octicon-video":                        "public/assets/img/svg/octicon-video.svg",
        "octicon-discussion-outdated":                        "public/assets/img/svg/octicon-discussion-outdated.svg",
        "gitea-onedev":                        "public/assets/img/svg/gitea-onedev.svg",
        "octicon-triangle-down":                        "public/assets/img/svg/octicon-triangle-down.svg",
        "octicon-devices":                        "public/assets/img/svg/octicon-devices.svg",
        "octicon-arrow-up":                        "public/assets/img/svg/octicon-arrow-up.svg",
        "octicon-unfold":                        "public/assets/img/svg/octicon-unfold.svg",
        "octicon-share":                        "public/assets/img/svg/octicon-share.svg",
        "octicon-rocket":                        "public/assets/img/svg/octicon-rocket.svg",
        "octicon-code-of-conduct":                        "public/assets/img/svg/octicon-code-of-conduct.svg",
        "octicon-list-unordered":                        "public/assets/img/svg/octicon-list-unordered.svg",
        "gitea-facebook":                        "public/assets/img/svg/gitea-facebook.svg",
        "octicon-dot-fill":                        "public/assets/img/svg/octicon-dot-fill.svg",
        "octicon-upload":                        "public/assets/img/svg/octicon-upload.svg",
        "octicon-feed-trophy":                        "public/assets/img/svg/octicon-feed-trophy.svg",
        "octicon-calendar":                        "public/assets/img/svg/octicon-calendar.svg",
        "gitea-unlock":                        "public/assets/img/svg/gitea-unlock.svg",
        "octicon-repo-clone":                        "public/assets/img/svg/octicon-repo-clone.svg",
        "octicon-person-fill":                        "public/assets/img/svg/octicon-person-fill.svg",
        "octicon-briefcase":                        "public/assets/img/svg/octicon-briefcase.svg",
        "octicon-ellipsis":                        "public/assets/img/svg/octicon-ellipsis.svg",
        "octicon-clock":                        "public/assets/img/svg/octicon-clock.svg",
        "octicon-copilot":                        "public/assets/img/svg/octicon-copilot.svg",
        "octicon-tracked-by-closed-completed":                        "public/assets/img/svg/octicon-tracked-by-closed-completed.svg",
        "octicon-clock-fill":                        "public/assets/img/svg/octicon-clock-fill.svg",
        "octicon-repo-template":                        "public/assets/img/svg/octicon-repo-template.svg",
        "octicon-chevron-down":                        "public/assets/img/svg/octicon-chevron-down.svg",
        "octicon-screen-full":                        "public/assets/img/svg/octicon-screen-full.svg",
        "octicon-repo-forked":                        "public/assets/img/svg/octicon-repo-forked.svg",
        "octicon-arrow-both":                        "public/assets/img/svg/octicon-arrow-both.svg",
        "octicon-filter":                        "public/assets/img/svg/octicon-filter.svg",
        "octicon-heart":                        "public/assets/img/svg/octicon-heart.svg",
        "octicon-italic":                        "public/assets/img/svg/octicon-italic.svg",
        "octicon-mortar-board":                        "public/assets/img/svg/octicon-mortar-board.svg",
        "octicon-tab-external":                        "public/assets/img/svg/octicon-tab-external.svg",
        "octicon-bookmark":                        "public/assets/img/svg/octicon-bookmark.svg",
        "octicon-codescan-checkmark":                        "public/assets/img/svg/octicon-codescan-checkmark.svg",
        "octicon-hourglass":                        "public/assets/img/svg/octicon-hourglass.svg",
        "octicon-pulse":                        "public/assets/img/svg/octicon-pulse.svg",
        "octicon-sun":                        "public/assets/img/svg/octicon-sun.svg",
        "octicon-columns":                        "public/assets/img/svg/octicon-columns.svg",
        "gitea-jetbrains":                        "public/assets/img/svg/gitea-jetbrains.svg",
        "octicon-squirrel":                        "public/assets/img/svg/octicon-squirrel.svg",
        "octicon-home":                        "public/assets/img/svg/octicon-home.svg",
        "octicon-graph":                        "public/assets/img/svg/octicon-graph.svg",
        "gitea-codebase":                        "public/assets/img/svg/gitea-codebase.svg",
        "octicon-fold":                        "public/assets/img/svg/octicon-fold.svg",
        "octicon-git-compare":                        "public/assets/img/svg/octicon-git-compare.svg",
        "octicon-feed-repo":                        "public/assets/img/svg/octicon-feed-repo.svg",
        "gitea-rubygems":                        "public/assets/img/svg/gitea-rubygems.svg",
        "octicon-feed-star":                        "public/assets/img/svg/octicon-feed-star.svg",
        "octicon-feed-pull-request-open":                        "public/assets/img/svg/octicon-feed-pull-request-open.svg",
        "octicon-bug":                        "public/assets/img/svg/octicon-bug.svg",
        "octicon-feed-issue-closed":                        "public/assets/img/svg/octicon-feed-issue-closed.svg",
        "octicon-megaphone":                        "public/assets/img/svg/octicon-megaphone.svg",
        "octicon-flame":                        "public/assets/img/svg/octicon-flame.svg",
        "octicon-skip":                        "public/assets/img/svg/octicon-skip.svg",
        "octicon-heart-fill":                        "public/assets/img/svg/octicon-heart-fill.svg",
        "octicon-diff-modified":                        "public/assets/img/svg/octicon-diff-modified.svg",
        "octicon-git-commit":                        "public/assets/img/svg/octicon-git-commit.svg",
        "octicon-person-add":                        "public/assets/img/svg/octicon-person-add.svg",
        "octicon-project-symlink":                        "public/assets/img/svg/octicon-project-symlink.svg",
        "octicon-feed-issue-open":                        "public/assets/img/svg/octicon-feed-issue-open.svg",
        "octicon-codescan":                        "public/assets/img/svg/octicon-codescan.svg",
        "gitea-twitter":                        "public/assets/img/svg/gitea-twitter.svg",
        "octicon-credit-card":                        "public/assets/img/svg/octicon-credit-card.svg",
        "octicon-chevron-right":                        "public/assets/img/svg/octicon-chevron-right.svg",
        "gitea-discord":                        "public/assets/img/svg/gitea-discord.svg",
        "octicon-log":                        "public/assets/img/svg/octicon-log.svg",
        "octicon-plus-circle":                        "public/assets/img/svg/octicon-plus-circle.svg",
        "octicon-bookmark-slash":                        "public/assets/img/svg/octicon-bookmark-slash.svg",
        "octicon-quote":                        "public/assets/img/svg/octicon-quote.svg",
        "octicon-unverified":                        "public/assets/img/svg/octicon-unverified.svg",
        "octicon-bold":                        "public/assets/img/svg/octicon-bold.svg",
        "gitea-vagrant":                        "public/assets/img/svg/gitea-vagrant.svg",
        "octicon-hash":                        "public/assets/img/svg/octicon-hash.svg",
        "gitea-gogs":                        "public/assets/img/svg/gitea-gogs.svg",
        "octicon-feed-issue-reopen":                        "public/assets/img/svg/octicon-feed-issue-reopen.svg",
        "octicon-table":                        "public/assets/img/svg/octicon-table.svg",
        "octicon-unlink":                        "public/assets/img/svg/octicon-unlink.svg",
        "gitea-azuread":                        "public/assets/img/svg/gitea-azuread.svg",
        "octicon-logo-github":                        "public/assets/img/svg/octicon-logo-github.svg",
        "octicon-tracked-by-closed-not-planned":                        "public/assets/img/svg/octicon-tracked-by-closed-not-planned.svg",
        "octicon-plus":                        "public/assets/img/svg/octicon-plus.svg",
        "gitea-double-chevron-right":                        "public/assets/img/svg/gitea-double-chevron-right.svg",
        "octicon-dash":                        "public/assets/img/svg/octicon-dash.svg",
        "octicon-paper-airplane":                        "public/assets/img/svg/octicon-paper-airplane.svg",
        "octicon-device-camera-video":                        "public/assets/img/svg/octicon-device-camera-video.svg",
        "octicon-list-ordered":                        "public/assets/img/svg/octicon-list-ordered.svg",
        "octicon-circle-slash":                        "public/assets/img/svg/octicon-circle-slash.svg",
        "octicon-zap":                        "public/assets/img/svg/octicon-zap.svg",
        "octicon-pencil":                        "public/assets/img/svg/octicon-pencil.svg",
        "octicon-file-directory":                        "public/assets/img/svg/octicon-file-directory.svg",
        "octicon-move-to-end":                        "public/assets/img/svg/octicon-move-to-end.svg",
        "octicon-gear":                        "public/assets/img/svg/octicon-gear.svg",
        "octicon-unlock":                        "public/assets/img/svg/octicon-unlock.svg",
        "octicon-eye":                        "public/assets/img/svg/octicon-eye.svg",
        "octicon-info":                        "public/assets/img/svg/octicon-info.svg",
        "octicon-file-binary":                        "public/assets/img/svg/octicon-file-binary.svg",
        "octicon-project-roadmap":                        "public/assets/img/svg/octicon-project-roadmap.svg",
        "octicon-book":                        "public/assets/img/svg/octicon-book.svg",
        "octicon-x":                        "public/assets/img/svg/octicon-x.svg",
        "octicon-arrow-up-left":                        "public/assets/img/svg/octicon-arrow-up-left.svg",
        "octicon-image":                        "public/assets/img/svg/octicon-image.svg",
        "octicon-zoom-out":                        "public/assets/img/svg/octicon-zoom-out.svg",
        "octicon-bell":                        "public/assets/img/svg/octicon-bell.svg",
        "gitea-gitea":                        "public/assets/img/svg/gitea-gitea.svg",
        "octicon-plug":                        "public/assets/img/svg/octicon-plug.svg",
        "octicon-telescope-fill":                        "public/assets/img/svg/octicon-telescope-fill.svg",
        "octicon-goal":                        "public/assets/img/svg/octicon-goal.svg",
        "octicon-link-external":                        "public/assets/img/svg/octicon-link-external.svg",
        "octicon-strikethrough":                        "public/assets/img/svg/octicon-strikethrough.svg",
        "octicon-feed-issue-draft":                        "public/assets/img/svg/octicon-feed-issue-draft.svg",
        "octicon-kebab-horizontal":                        "public/assets/img/svg/octicon-kebab-horizontal.svg",
        "octicon-move-to-start":                        "public/assets/img/svg/octicon-move-to-start.svg",
        "octicon-code":                        "public/assets/img/svg/octicon-code.svg",
        "octicon-feed-heart":                        "public/assets/img/svg/octicon-feed-heart.svg",
        "octicon-container":                        "public/assets/img/svg/octicon-container.svg",
        "octicon-milestone":                        "public/assets/img/svg/octicon-milestone.svg",
        "octicon-stop":                        "public/assets/img/svg/octicon-stop.svg",
        "octicon-issue-tracks":                        "public/assets/img/svg/octicon-issue-tracks.svg",
        "octicon-workflow":                        "public/assets/img/svg/octicon-workflow.svg",
        "octicon-discussion-duplicate":                        "public/assets/img/svg/octicon-discussion-duplicate.svg",
        "octicon-sort-desc":                        "public/assets/img/svg/octicon-sort-desc.svg",
        "octicon-report":                        "public/assets/img/svg/octicon-report.svg",
        "octicon-square":                        "public/assets/img/svg/octicon-square.svg",
        "octicon-diff":                        "public/assets/img/svg/octicon-diff.svg",
        "octicon-feed-tag":                        "public/assets/img/svg/octicon-feed-tag.svg",
        "octicon-package":                        "public/assets/img/svg/octicon-package.svg",
        "octicon-triangle-right":                        "public/assets/img/svg/octicon-triangle-right.svg",
        "octicon-mail":                        "public/assets/img/svg/octicon-mail.svg",
        "octicon-arrow-switch":                        "public/assets/img/svg/octicon-arrow-switch.svg",
        "gitea-double-chevron-left":                        "public/assets/img/svg/gitea-double-chevron-left.svg",
        "octicon-stopwatch":                        "public/assets/img/svg/octicon-stopwatch.svg",
        "octicon-horizontal-rule":                        "public/assets/img/svg/octicon-horizontal-rule.svg",
        "octicon-light-bulb":                        "public/assets/img/svg/octicon-light-bulb.svg",
        "octicon-infinity":                        "public/assets/img/svg/octicon-infinity.svg",
        "octicon-apps":                        "public/assets/img/svg/octicon-apps.svg",
        "octicon-arrow-up-right":                        "public/assets/img/svg/octicon-arrow-up-right.svg",
        "octicon-codespaces":                        "public/assets/img/svg/octicon-codespaces.svg",
        "octicon-copilot-error":                        "public/assets/img/svg/octicon-copilot-error.svg",
        "octicon-note":                        "public/assets/img/svg/octicon-note.svg",
        "octicon-git-pull-request-closed":                        "public/assets/img/svg/octicon-git-pull-request-closed.svg",
        "octicon-passkey-fill":                        "public/assets/img/svg/octicon-passkey-fill.svg",
        "octicon-check-circle-fill":                        "public/assets/img/svg/octicon-check-circle-fill.svg",
        "gitea-go":                        "public/assets/img/svg/gitea-go.svg",
        "octicon-fold-down":                        "public/assets/img/svg/octicon-fold-down.svg",
        "octicon-lock":                        "public/assets/img/svg/octicon-lock.svg",
        "octicon-arrow-down-right":                        "public/assets/img/svg/octicon-arrow-down-right.svg",
        "octicon-checkbox":                        "public/assets/img/svg/octicon-checkbox.svg",
        "octicon-code-square":                        "public/assets/img/svg/octicon-code-square.svg",
        "octicon-people":                        "public/assets/img/svg/octicon-people.svg",
        "gitea-vscodium":                        "public/assets/img/svg/gitea-vscodium.svg",
        "octicon-paste":                        "public/assets/img/svg/octicon-paste.svg",
        "octicon-trash":                        "public/assets/img/svg/octicon-trash.svg",
        "octicon-feed-public":                        "public/assets/img/svg/octicon-feed-public.svg",
        "octicon-arrow-left":                        "public/assets/img/svg/octicon-arrow-left.svg",
        "octicon-paintbrush":                        "public/assets/img/svg/octicon-paintbrush.svg",
        "fontawesome-openid":                        "public/assets/img/svg/fontawesome-openid.svg",
        "octicon-bell-slash":                        "public/assets/img/svg/octicon-bell-slash.svg",
        "octicon-link":                        "public/assets/img/svg/octicon-link.svg",
        "octicon-feed-person":                        "public/assets/img/svg/octicon-feed-person.svg",
        "octicon-telescope":                        "public/assets/img/svg/octicon-telescope.svg",
        "gitea-debian":                        "public/assets/img/svg/gitea-debian.svg",
        "octicon-ruby":                        "public/assets/img/svg/octicon-ruby.svg",
        "octicon-star-fill":                        "public/assets/img/svg/octicon-star-fill.svg",
        "gitea-composer":                        "public/assets/img/svg/gitea-composer.svg",
        "octicon-feed-forked":                        "public/assets/img/svg/octicon-feed-forked.svg",
        "octicon-archive":                        "public/assets/img/svg/octicon-archive.svg",
        "octicon-sponsor-tiers":                        "public/assets/img/svg/octicon-sponsor-tiers.svg",
        "octicon-repo-deleted":                        "public/assets/img/svg/octicon-repo-deleted.svg",
        "octicon-move-to-bottom":                        "public/assets/img/svg/octicon-move-to-bottom.svg",
        "gitea-gitlab":                        "public/assets/img/svg/gitea-gitlab.svg",
        "gitea-openid":                        "public/assets/img/svg/gitea-openid.svg",
        "octicon-feed-rocket":                        "public/assets/img/svg/octicon-feed-rocket.svg",
        "gitea-maven":                        "public/assets/img/svg/gitea-maven.svg",
        "octicon-question":                        "public/assets/img/svg/octicon-question.svg",
        "octicon-thumbsup":                        "public/assets/img/svg/octicon-thumbsup.svg",
        "octicon-feed-pull-request-draft":                        "public/assets/img/svg/octicon-feed-pull-request-draft.svg",
        "octicon-single-select":                        "public/assets/img/svg/octicon-single-select.svg",
        "octicon-trophy":                        "public/assets/img/svg/octicon-trophy.svg",
        "octicon-device-camera":                        "public/assets/img/svg/octicon-device-camera.svg",
        "octicon-accessibility":                        "public/assets/img/svg/octicon-accessibility.svg",
        "gitea-cran":                        "public/assets/img/svg/gitea-cran.svg",
        "octicon-repo-push":                        "public/assets/img/svg/octicon-repo-push.svg",
        "gitea-chef":                        "public/assets/img/svg/gitea-chef.svg",
        "octicon-hubot":                        "public/assets/img/svg/octicon-hubot.svg",
        "gitea-helm":                        "public/assets/img/svg/gitea-helm.svg",
        "octicon-typography":                        "public/assets/img/svg/octicon-typography.svg",
        "octicon-file-submodule":                        "public/assets/img/svg/octicon-file-submodule.svg",
        "octicon-package-dependents":                        "public/assets/img/svg/octicon-package-dependents.svg",
        "octicon-skip-fill":                        "public/assets/img/svg/octicon-skip-fill.svg",
        "octicon-north-star":                        "public/assets/img/svg/octicon-north-star.svg",
        "octicon-file-directory-symlink":                        "public/assets/img/svg/octicon-file-directory-symlink.svg",
        "octicon-chevron-up":                        "public/assets/img/svg/octicon-chevron-up.svg",
        "octicon-mute":                        "public/assets/img/svg/octicon-mute.svg",
        "fontawesome-windows":                        "public/assets/img/svg/fontawesome-windows.svg",
        "octicon-desktop-download":                        "public/assets/img/svg/octicon-desktop-download.svg",
        "octicon-shield-lock":                        "public/assets/img/svg/octicon-shield-lock.svg",
        "octicon-download":                        "public/assets/img/svg/octicon-download.svg",
        "octicon-browser":                        "public/assets/img/svg/octicon-browser.svg",
        "octicon-checklist":                        "public/assets/img/svg/octicon-checklist.svg",
        "gitea-pub":                        "public/assets/img/svg/gitea-pub.svg",
        "octicon-repo":                        "public/assets/img/svg/octicon-repo.svg",
        "octicon-star":                        "public/assets/img/svg/octicon-star.svg",
        "gitea-vscode":                        "public/assets/img/svg/gitea-vscode.svg",
        "octicon-number":                        "public/assets/img/svg/octicon-number.svg",
        "octicon-shield-x":                        "public/assets/img/svg/octicon-shield-x.svg",
        "octicon-smiley":                        "public/assets/img/svg/octicon-smiley.svg",
        "octicon-play":                        "public/assets/img/svg/octicon-play.svg",
        "octicon-broadcast":                        "public/assets/img/svg/octicon-broadcast.svg",
        "octicon-paperclip":                        "public/assets/img/svg/octicon-paperclip.svg",
        "octicon-thumbsdown":                        "public/assets/img/svg/octicon-thumbsdown.svg",
        "octicon-diff-ignored":                        "public/assets/img/svg/octicon-diff-ignored.svg",
        "octicon-server":                        "public/assets/img/svg/octicon-server.svg",
        "octicon-check":                        "public/assets/img/svg/octicon-check.svg",
        "octicon-arrow-down":                        "public/assets/img/svg/octicon-arrow-down.svg",
        "octicon-tag":                        "public/assets/img/svg/octicon-tag.svg",
        "octicon-diff-removed":                        "public/assets/img/svg/octicon-diff-removed.svg",
        "octicon-shield":                        "public/assets/img/svg/octicon-shield.svg",
        "gitea-conda":                        "public/assets/img/svg/gitea-conda.svg",
        "octicon-cloud-offline":                        "public/assets/img/svg/octicon-cloud-offline.svg",
        "octicon-redo":                        "public/assets/img/svg/octicon-redo.svg",
        "octicon-zoom-in":                        "public/assets/img/svg/octicon-zoom-in.svg",
        "octicon-git-merge-queue":                        "public/assets/img/svg/octicon-git-merge-queue.svg",
        "octicon-markdown":                        "public/assets/img/svg/octicon-markdown.svg",
        "gitea-matrix":                        "public/assets/img/svg/gitea-matrix.svg",
        "octicon-unmute":                        "public/assets/img/svg/octicon-unmute.svg",
        "gitea-npm":                        "public/assets/img/svg/gitea-npm.svg",
        "gitea-split":                        "public/assets/img/svg/gitea-split.svg",
        "octicon-file-added":                        "public/assets/img/svg/octicon-file-added.svg",
        "octicon-organization":                        "public/assets/img/svg/octicon-organization.svg",
        "octicon-copilot-warning":                        "public/assets/img/svg/octicon-copilot-warning.svg",
        "gitea-empty-checkbox":                        "public/assets/img/svg/gitea-empty-checkbox.svg",
        "octicon-git-merge":                        "public/assets/img/svg/octicon-git-merge.svg",
        "octicon-copy":                        "public/assets/img/svg/octicon-copy.svg",
        "octicon-check-circle":                        "public/assets/img/svg/octicon-check-circle.svg",
        "octicon-feed-pull-request-closed":                        "public/assets/img/svg/octicon-feed-pull-request-closed.svg",
        "material-invert-colors":                        "public/assets/img/svg/material-invert-colors.svg",
        "octicon-beaker":                        "public/assets/img/svg/octicon-beaker.svg",
        "octicon-move-to-top":                        "public/assets/img/svg/octicon-move-to-top.svg",
        "octicon-file-directory-open-fill":                        "public/assets/img/svg/octicon-file-directory-open-fill.svg",
        "octicon-git-pull-request-draft":                        "public/assets/img/svg/octicon-git-pull-request-draft.svg",
        "octicon-grabber":                        "public/assets/img/svg/octicon-grabber.svg",
        "octicon-versions":                        "public/assets/img/svg/octicon-versions.svg",
        "octicon-webhook":                        "public/assets/img/svg/octicon-webhook.svg",
        "octicon-chevron-left":                        "public/assets/img/svg/octicon-chevron-left.svg",
        "octicon-pivot-column":                        "public/assets/img/svg/octicon-pivot-column.svg",
        "gitea-rpm":                        "public/assets/img/svg/gitea-rpm.svg",
        "octicon-dot":                        "public/assets/img/svg/octicon-dot.svg",
        "octicon-rel-file-path":                        "public/assets/img/svg/octicon-rel-file-path.svg",
        "octicon-sidebar-expand":                        "public/assets/img/svg/octicon-sidebar-expand.svg",
        "octicon-discussion-closed":                        "public/assets/img/svg/octicon-discussion-closed.svg",
        "gitea-lock-cog":                        "public/assets/img/svg/gitea-lock-cog.svg",
        "octicon-file-directory-fill":                        "public/assets/img/svg/octicon-file-directory-fill.svg",
        "gitea-mastodon":                        "public/assets/img/svg/gitea-mastodon.svg",
        "octicon-arrow-right":                        "public/assets/img/svg/octicon-arrow-right.svg",
        "octicon-mention":                        "public/assets/img/svg/octicon-mention.svg",
        "octicon-square-fill":                        "public/assets/img/svg/octicon-square-fill.svg",
        "octicon-issue-opened":                        "public/assets/img/svg/octicon-issue-opened.svg",
        "octicon-git-branch":                        "public/assets/img/svg/octicon-git-branch.svg",
        "octicon-diff-renamed":                        "public/assets/img/svg/octicon-diff-renamed.svg",
        "octicon-reply":                        "public/assets/img/svg/octicon-reply.svg",
        "octicon-repo-locked":                        "public/assets/img/svg/octicon-repo-locked.svg",
        "gitea-nextcloud":                        "public/assets/img/svg/gitea-nextcloud.svg",
        "octicon-shield-check":                        "public/assets/img/svg/octicon-shield-check.svg",
        "octicon-moon":                        "public/assets/img/svg/octicon-moon.svg",
        "octicon-law":                        "public/assets/img/svg/octicon-law.svg",
        "octicon-share-android":                        "public/assets/img/svg/octicon-share-android.svg",
        "octicon-key":                        "public/assets/img/svg/octicon-key.svg",
        "octicon-file-zip":                        "public/assets/img/svg/octicon-file-zip.svg",
        "octicon-device-desktop":                        "public/assets/img/svg/octicon-device-desktop.svg",
        "octicon-fold-up":                        "public/assets/img/svg/octicon-fold-up.svg",
        "octicon-multi-select":                        "public/assets/img/svg/octicon-multi-select.svg",
        "octicon-device-mobile":                        "public/assets/img/svg/octicon-device-mobile.svg",
        "octicon-tools":                        "public/assets/img/svg/octicon-tools.svg",
        "octicon-triangle-left":                        "public/assets/img/svg/octicon-triangle-left.svg",
        "gitea-join":                        "public/assets/img/svg/gitea-join.svg",
        "octicon-read":                        "public/assets/img/svg/octicon-read.svg",
        "gitea-git":                        "public/assets/img/svg/gitea-git.svg",
        "octicon-file-badge":                        "public/assets/img/svg/octicon-file-badge.svg",
        "octicon-file-diff":                        "public/assets/img/svg/octicon-file-diff.svg",
        "gitea-alpine":                        "public/assets/img/svg/gitea-alpine.svg",
        "octicon-issue-closed":                        "public/assets/img/svg/octicon-issue-closed.svg",
        "octicon-key-asterisk":                        "public/assets/img/svg/octicon-key-asterisk.svg",
        "octicon-feed-merged":                        "public/assets/img/svg/octicon-feed-merged.svg",
        "gitea-bitbucket":                        "public/assets/img/svg/gitea-bitbucket.svg",
        "octicon-meter":                        "public/assets/img/svg/octicon-meter.svg",
        "octicon-logo-gist":                        "public/assets/img/svg/octicon-logo-gist.svg",
        "octicon-id-badge":                        "public/assets/img/svg/octicon-id-badge.svg",
        "material-palette":                        "public/assets/img/svg/material-palette.svg",
        "octicon-x-circle":                        "public/assets/img/svg/octicon-x-circle.svg",
        "octicon-file-removed":                        "public/assets/img/svg/octicon-file-removed.svg",
        "octicon-code-review":                        "public/assets/img/svg/octicon-code-review.svg",
        "octicon-shield-slash":                        "public/assets/img/svg/octicon-shield-slash.svg",
        "octicon-history":                        "public/assets/img/svg/octicon-history.svg",
        "octicon-project":                        "public/assets/img/svg/octicon-project.svg",
        "octicon-iterations":                        "public/assets/img/svg/octicon-iterations.svg",
        "octicon-x-circle-fill":                        "public/assets/img/svg/octicon-x-circle-fill.svg",
        "octicon-sparkle-fill":                        "public/assets/img/svg/octicon-sparkle-fill.svg",
        "octicon-pin":                        "public/assets/img/svg/octicon-pin.svg",
        "octicon-cross-reference":                        "public/assets/img/svg/octicon-cross-reference.svg",
        "octicon-screen-normal":                        "public/assets/img/svg/octicon-screen-normal.svg",
        "octicon-undo":                        "public/assets/img/svg/octicon-undo.svg",
        "gitea-gitbucket":                        "public/assets/img/svg/gitea-gitbucket.svg",
        "octicon-filter-remove":                        "public/assets/img/svg/octicon-filter-remove.svg",
        "fontawesome-save":                        "public/assets/img/svg/fontawesome-save.svg",
        "octicon-dependabot":                        "public/assets/img/svg/octicon-dependabot.svg",
        "gitea-exclamation":                        "public/assets/img/svg/gitea-exclamation.svg",
        "octicon-mark-github":                        "public/assets/img/svg/octicon-mark-github.svg",
        "octicon-diamond":                        "public/assets/img/svg/octicon-diamond.svg",
        "octicon-file-code":                        "public/assets/img/svg/octicon-file-code.svg",
        "octicon-circle":                        "public/assets/img/svg/octicon-circle.svg",
        "octicon-triangle-up":                        "public/assets/img/svg/octicon-triangle-up.svg",
        "octicon-person":                        "public/assets/img/svg/octicon-person.svg",
        "octicon-bell-fill":                        "public/assets/img/svg/octicon-bell-fill.svg",
        "gitea-conan":                        "public/assets/img/svg/gitea-conan.svg",
        "octicon-eye-closed":                        "public/assets/img/svg/octicon-eye-closed.svg",
        "octicon-inbox":                        "public/assets/img/svg/octicon-inbox.svg",
        "octicon-cache":                        "public/assets/img/svg/octicon-cache.svg",
        "octicon-blocked":                        "public/assets/img/svg/octicon-blocked.svg",
        "octicon-unread":                        "public/assets/img/svg/octicon-unread.svg",
        "gitea-cargo":                        "public/assets/img/svg/gitea-cargo.svg",
        "octicon-verified":                        "public/assets/img/svg/octicon-verified.svg",
        "octicon-tasklist":                        "public/assets/img/svg/octicon-tasklist.svg",
        "octicon-database":                        "public/assets/img/svg/octicon-database.svg",
        "octicon-diff-added":                        "public/assets/img/svg/octicon-diff-added.svg",
        "fontawesome-send":                        "public/assets/img/svg/fontawesome-send.svg",
        "octicon-sync":                        "public/assets/img/svg/octicon-sync.svg",
        "octicon-rows":                        "public/assets/img/svg/octicon-rows.svg",
        "octicon-command-palette":                        "public/assets/img/svg/octicon-command-palette.svg",
        "octicon-project-template":                        "public/assets/img/svg/octicon-project-template.svg",
        "octicon-sidebar-collapse":                        "public/assets/img/svg/octicon-sidebar-collapse.svg",
        "gitea-microsoftonline":                        "public/assets/img/svg/gitea-microsoftonline.svg",
        "octicon-feed-plus":                        "public/assets/img/svg/octicon-feed-plus.svg",
        "octicon-search":                        "public/assets/img/svg/octicon-search.svg",
        "octicon-cpu":                        "public/assets/img/svg/octicon-cpu.svg",
        "octicon-sign-in":                        "public/assets/img/svg/octicon-sign-in.svg",
        "octicon-no-entry":                        "public/assets/img/svg/octicon-no-entry.svg",
        "gitea-python":                        "public/assets/img/svg/gitea-python.svg",
        "gitea-lock":                        "public/assets/img/svg/gitea-lock.svg",
        "octicon-comment-discussion":                        "public/assets/img/svg/octicon-comment-discussion.svg",
        "octicon-arrow-down-left":                        "public/assets/img/svg/octicon-arrow-down-left.svg",
        "gitea-nuget":                        "public/assets/img/svg/gitea-nuget.svg",
        "octicon-git-pull-request":                        "public/assets/img/svg/octicon-git-pull-request.svg",
        "octicon-fiscal-host":                        "public/assets/img/svg/octicon-fiscal-host.svg",
        "octicon-terminal":                        "public/assets/img/svg/octicon-terminal.svg",
        "octicon-repo-pull":                        "public/assets/img/svg/octicon-repo-pull.svg",
        "gitea-dropbox":                        "public/assets/img/svg/gitea-dropbox.svg",
        "octicon-feed-discussion":                        "public/assets/img/svg/octicon-feed-discussion.svg",
        "gitea-swift":                        "public/assets/img/svg/gitea-swift.svg",
        "octicon-accessibility-inset":                        "public/assets/img/svg/octicon-accessibility-inset.svg",
        "octicon-alert":                        "public/assets/img/svg/octicon-alert.svg",
        "octicon-gift":                        "public/assets/img/svg/octicon-gift.svg",
        "octicon-comment":                        "public/assets/img/svg/octicon-comment.svg",
        "octicon-rss":                        "public/assets/img/svg/octicon-rss.svg",
        "git":                        "public/assets/img/emoji/git.png",
        "github":                        "public/assets/img/emoji/github.png",
        "gitlab":                        "public/assets/img/emoji/gitlab.png",
        "codeberg":                        "public/assets/img/emoji/codeberg.png",

	}
}