import * as core from '@actions/core';
import * as github from '@actions/github';
import * as exec from '@actions/exec';

import { existsSync, mkdirSync, writeFileSync } from 'fs';
import { join } from 'path';
import codeBlocks from 'gfm-code-blocks';

async function runner(issue) {
  core.startGroup('Print issue labels');
  console.log(issue.labels);
  core.endGroup()

  //core.startGroup('Print issue body');
  //console.log(issue.body);
  //core.endGroup()

  // TODO: is it worth replacing this library with a regexp?
  // https://coderwall.com/p/r6b4xg/regex-to-match-github-s-markdown-code-blocks
  // var code = s.match(/```([^`]*)```/)[1]
  const blocks = codeBlocks(issue.body);
  const l = blocks.length;

  // FIXME: Check whether any attached tarball/zipfile exists.
  if (l === 0) {
    console.log("no code blocks found in issue body, skipping");
    return;
  }

  const dir = join(__dirname, 'tmp-dir');

  if (!existsSync(dir)){
    mkdirSync(dir);
  }

  var img = 'host';

  // Write each code block to a file
  blocks.forEach(function(d, i) {
    process.stdout.write('Processing file '+(i+1)+'/'+l+'... ');
    const c = d.code.slice(1);
    var fname = c.match(/:file:.*/g);
    if (!fname) {
      fname = c.match(/:image:.*/g);
      if (!fname) {
        console.log("code block does not contain ':file:' or ':image:', skipping");
        return;
      }
      img = fname[0].replace(':image:','').trim();
      fname = 'run';
    } else {
      fname = fname[0].replace(':file:','').trim();
    }
    console.log(fname);
    var mode = 0o644;
    if (fname === 'run') {
      mode = 0o744;
    }
    // FIXME: why is each line preprended with '\n'? Is it coming from GitHub's payload or is it added by codeBlocks?
    writeFileSync(join(dir, fname), c.replace(/\r\n/g, "\n"), { mode: mode });
  });

  if (!existsSync(join(dir, 'run'))) {
    console.log("file 'run' not provided, skipping");
    return;
  }

  if ( img != 'host' ) {
    core.startGroup('Docker pull '+img);
    await exec.exec(`docker`, ['pull', img]);
    core.endGroup()

    core.startGroup('Execute in docker container');
    await exec.exec(`docker`, ['run', '--rm', '-tv', dir+':/src', '-w', '/src', img, `./run`]);
    core.endGroup()
  } else {
    core.startGroup('Execute');
    await exec.exec(`./run`, [], { cwd: dir });
    core.endGroup()
  }

  console.log(issue);

/*
    import * as io from '@actions/io';

    const pythonPath: string = await io.which('python', true);
    await exec.exec(`"${pythonPath}"`, ['main.py']);
*/

/*
let myOutput = '';
let myError = '';

await exec.exec('node', ['index.js', 'foo=bar'], {
  listeners: {
    stdout: (data: Buffer) => { myOutput += data.toString(); },
    stderr: (data: Buffer) => { myError += data.toString(); }
  },
  cwd: './lib'
});
*/
}

export async function run() {
  try {
    //const repoToken: string = core.getInput('token', {required: true});
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

    runner(ctx.payload.issue);

/*
    const client: github.GitHub = new github.GitHub(repoToken);
    await client.issues.createComment({
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
//https://octokit.github.io/rest.js/#octokit-routes-issues-list-comments
//https://octokit.github.io/rest.js/#octokit-routes-issues-update-comment
//https://octokit.github.io/rest.js/#octokit-routes-issues-get-comment
