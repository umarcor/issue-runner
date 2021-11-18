import { writeFileSync } from 'fs'
import { endGroup, startGroup } from '@actions/core'
import { exec } from '@actions/exec'
import * as c from 'ansi-colors'

const exitCodes = {
  empty: 1,
  exec: 2,
  format: 3,
  docker: 4,
  fail: 5,
  default: 6
}

export async function runner (
  body: string | undefined,
  dir: string | undefined,
  allowHost?: boolean
): Promise<any> {
  const out = {
    result: -1,
    stdout: '',
    stderr: ''
  }
  if (body) {
    const fname = dir ? dir + '.md' : 'body.md'
    // TODO: can we pipe the content to the tool (through stdin), instead of writing the body to disk?
    // FIXME: is each line still preprended with '\n'? Is it coming from GitHub's payload?
    writeFileSync(fname, body.replace(/\r\n/g, '\n'), { mode: 0o644 })

    out.result = await exec('issue-runner', [
      allowHost ? '-y' : '--no-host',
      '-c',
      '-d',
      dir || 'tmp-dir',
      fname
    ], {
      ignoreReturnCode: true,
      silent: true,
      listeners: {
        stdline: (data: string): void => { out.stdout += data + '\n' },
        errline: (data: string): void => { out.stderr += data + '\n' }
      }
    })
  }
  return out
}

export async function execute (
  issues: any[]
): Promise<any> {
  console.log(c.magenta('-- Run --'))

  const results: any[] = []
  await Promise.all(issues.map(async (el) => {
    console.log(`#${el.number}...`)
    const out = await runner(el.body, el.number.toString(10))
    out.number = el.number
    out.labels = el.labels
    results.push(out)
  }))

  const sum = {
    skip: [] as any[],
    err: [] as any[],
    fail: [] as any[],
    ok: [] as any[]
  }

  console.log(c.magenta('-- Logs --'))

  results.forEach((el) => {
    startGroup(`#${el.number}`)
    if (el.stdout.length > 0) {
      console.log(c.magenta('stdout'))
      console.log(el.stdout)
    }
    if (el.stderr.length > 0) {
      console.log(c.magenta('stderr'))
      console.log(el.stderr)
    }
    endGroup()
    switch (el.result) {
      case 0:
        sum.ok.push(el)
        break
      case exitCodes.empty:
        sum.skip.push(el)
        break
      case exitCodes.fail:
        sum.fail.push(el)
        break
      case exitCodes.exec:
        el.errortype = 'Exec'
        sum.err.push(el)
        break
      case exitCodes.format:
        el.errortype = 'Format'
        sum.err.push(el)
        break
      case exitCodes.docker:
        el.errortype = 'Docker'
        sum.err.push(el)
        break
      default:
        sum.err.push(el)
        break
    }
  })

  return sum
}
