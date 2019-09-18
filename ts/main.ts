import * as core from '@actions/core';
import * as github from '@actions/github';
import * as exec from '@actions/exec';

import { writeFileSync } from 'fs';
//import { join } from 'path';

async function runner(issue) {
  // TODO: can we pipe the content to the tool (through stdin), instead of writing the body to disk?
  core.startGroup('write body.md');
  // FIXME: is each line still preprended with '\n'? Is it coming from GitHub's payload?
  writeFileSync('body.md', issue.body.replace(/\r\n/g, "\n"), { mode: 0o644 });
  core.endGroup()

  let cli_args = ''
  const allow_host = core.getInput('token', {required: true});
  if (allow_host) {
    cli_args = '-y '
  }

  writeFileSync('task.sh', `#!/usr/bin/env sh
curl -fsSL -o issue-runner https://github.com/1138-4EB/issue-runner/releases/download/tip/issue-runner_lin_amd64
chmod +x issue-runner
./issue-runner ` + cli_args + `body.md
`, { mode: 0o755 });

  await exec.exec(`cat`, ['task.sh']);

  // FIXME: it should be possible to execute 'task.sh' directly
  await exec.exec(`sh`, ['-c', './task.sh']);
}

export async function run() {
  try {
    const ctx = github.context;

    //const issue: { owner: string; repo: string; number: number; } = ctx.issue;

    if (ctx.eventName != 'issues' || !ctx.payload.issue ) {
      console.log('not an issue, skipping');
      return
    }

    const act = ctx.payload.action;
    if (act != 'opened' && act != 'edited') {
      console.log("issue neither 'opened' nor 'edited', skipping");
      return;
    }

    core.startGroup('Print issue labels');
    console.log(ctx.payload.issue.labels);
    core.endGroup()

    runner(ctx.payload.issue);

    //const repoToken: string = core.getInput('token', {required: true});

/*
    const octokit: github.GitHub = new github.GitHub(repoToken);
    const { data: comment } = await octokit.issues.createComment({
      owner: issue.owner,
      repo: issue.repo,
      issue_number: issue.number,
      body: 'Welcome message!'
    });
*/
  }

  catch (error) {
    core.setFailed(error.message);
    throw error;
  }
}

run();

/*
pull_request:      issues:       issue_comment:

opened             opened         created
edited             edited         edited
closed             closed
                   deleted        deleted
assigned           assigned
unassigned         unassigned
labeled            labeled
unlabeled          unlabeled
reopened           reopened
locked             locked
unlocked           unlocked
synchronize
ready_for_review
                   pinned
                   unpinned
                   milestoned
                   demilestoned
                   transferred
*/

//https://github.com/actions/labeler
//https://github.com/actions/first-interaction/blob/master/src/main.ts

//https://octokit.github.io/rest.js/