import * as c from 'ansi-colors'

function sortWithIntervals (lst): string[] {
  const out: string[] = []
  const arr: number[] = lst.map((el) => el.number).sort((a, b) => a - b)
  arr.reduce((c, v, x, p) => {
    const n = p[x + 1]
    if (n !== v + 1) {
      out.push(((c !== v) ? ('#' + c.toString() + '-') : '') + ('#' + v.toString()))
      c = n
    }
    return c
  }, arr[0])
  return out
}

export async function summary (
  data: any
): Promise<any> {
  console.log(c.magenta('-- Summary --'))

  if (data.ok.length !== 0) {
    console.log(c.green('success:'), sortWithIntervals(data.ok).join(', '))
  }

  if (data.fail.length !== 0) {
    console.log(c.red('failure:'), sortWithIntervals(data.fail).join(', '))
  }

  if (data.err.length !== 0) {
    console.log(c.yellow('errored:'), sortWithIntervals(data.err).join(', '))
  }

  if (data.skip.length !== 0) {
    console.log(c.gray('skipped:'), sortWithIntervals(data.skip).join(', '))
  }

  console.log(c.magenta('-- Actions --'))

  data.ok.forEach((el) => {
    console.log(el.number, el.labels.map((item) => item.name))
    // was 'fixed?' | 'triage' | 'reproducible' -> remove 'triage' add 'fixed?'
    // TODO Add a comment to the issue
  })

  data.fail.forEach((el) => {
    console.log(el.number, el.labels.map((item) => item.name))
    // was 'fixed?' | 'triage' | 'reproducible' -> remove 'triage' add 'reproducible'
    // TODO Add a comment to the issue
  })

  data.err.forEach((el) => {
    console.log(el.number, el.labels.map((item) => item.name))
    // TODO Add a comment to the issue
    // add comment with the type of error
  })
}
