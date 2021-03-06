// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package runner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/drone/drone/core"
)

func systemEnviron(system *core.System) map[string]string {
	return map[string]string{
		"CI":                    "true",
		"DRONE":                 "true",
		"DRONE_SYSTEM_PROTO":    system.Proto,
		"DRONE_SYSTEM_HOST":     system.Host,
		"DRONE_SYSTEM_HOSTNAME": system.Host,
		"DRONE_SYSTEM_VERSION":  fmt.Sprint(system.Version),
	}
}

func agentEnviron(runner *Runner) map[string]string {
	return map[string]string{
		"DRONE_MACHINE":         runner.Machine,
		"DRONE_RUNNER_HOST":     runner.Machine,
		"DRONE_RUNNER_HOSTNAME": runner.Machine,
		"DRONE_RUNNER_PLATFORM": runner.Platform,
	}
}

func repoEnviron(repo *core.Repository) map[string]string {
	return map[string]string{
		"DRONE_REPO":            repo.Slug,
		"DRONE_REPO_SCM":        repo.SCM,
		"DRONE_REPO_OWNER":      repo.Namespace,
		"DRONE_REPO_NAMESPACE":  repo.Namespace,
		"DRONE_REPO_NAME":       repo.Name,
		"DRONE_REPO_LINK":       repo.Link,
		"DRONE_REPO_BRANCH":     repo.Branch,
		"DRONE_REMOTE_URL":      repo.HTTPURL,
		"DRONE_GIT_HTTP_URL":    repo.HTTPURL,
		"DRONE_GIT_SSH_URL":     repo.SSHURL,
		"DRONE_REPO_VISIBILITY": repo.Visibility,
		"DRONE_REPO_PRIVATE":    fmt.Sprint(repo.Private),
	}
}

func stageEnviron(stage *core.Stage) map[string]string {
	return map[string]string{
		"DRONE_STAGE_KIND":       "pipeline",
		"DRONE_STAGE_NAME":       stage.Name,
		"DRONE_STAGE_NUMBER":     fmt.Sprint(stage.Number),
		"DRONE_STAGE_MACHINE":    stage.Machine,
		"DRONE_STAGE_OS":         stage.OS,
		"DRONE_STAGE_ARCH":       stage.Arch,
		"DRONE_STAGE_VARIANT":    stage.Variant,
		"DRONE_STAGE_DEPENDS_ON": strings.Join(stage.DependsOn, ","),
	}
}

func buildEnviron(build *core.Build) map[string]string {
	env := map[string]string{
		"DRONE_BRANCH":               build.Target,
		"DRONE_SOURCE_BRANCH":        build.Source,
		"DRONE_TARGET_BRANCH":        build.Target,
		"DRONE_COMMIT":               build.After,
		"DRONE_COMMIT_SHA":           build.After,
		"DRONE_COMMIT_BEFORE":        build.Before,
		"DRONE_COMMIT_AFTER":         build.After,
		"DRONE_COMMIT_REF":           build.Ref,
		"DRONE_COMMIT_BRANCH":        build.Target,
		"DRONE_COMMIT_LINK":          build.Link,
		"DRONE_COMMIT_MESSAGE":       build.Message,
		"DRONE_COMMIT_AUTHOR":        build.Author,
		"DRONE_COMMIT_AUTHOR_EMAIL":  build.AuthorEmail,
		"DRONE_COMMIT_AUTHOR_AVATAR": build.AuthorAvatar,
		"DRONE_COMMIT_AUTHOR_NAME":   build.AuthorName,
		"DRONE_BUILD_NUMBER":         fmt.Sprint(build.Number),
		"DRONE_BUILD_EVENT":          build.Event,
		"DRONE_BUILD_ACTION":         build.Action,
		"DRONE_BUILD_CREATED":        fmt.Sprint(build.Created),
		"DRONE_BUILD_STARTED":        fmt.Sprint(build.Started),
		"DRONE_BUILD_FINISHED":       fmt.Sprint(build.Finished),
		"DRONE_DEPLOY_TO":            build.Deploy,
	}
	if strings.HasPrefix(build.Ref, "refs/tags/") {
		env["DRONE_TAG"] = strings.TrimPrefix(build.Ref, "refs/tags/")
	}
	if build.Event == core.EventPullRequest {
		env["DRONE_PULL_REQUEST"] = re.FindString(build.Ref)
	}
	return env
}

func linkEnviron(repo *core.Repository, build *core.Build, system *core.System) map[string]string {
	return map[string]string{
		"DRONE_BUILD_LINK": fmt.Sprintf(
			"%s://%s/%s/%d",
			system.Proto,
			system.Host,
			repo.Slug,
			build.Number,
		),
	}
}

// regular expression to extract the pull request number
// from the git ref (e.g. refs/pulls/{d}/head)
var re = regexp.MustCompile("\\d+")

// helper function combines one or more maps of environment
// variables into a single map.
func combineEnviron(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}
