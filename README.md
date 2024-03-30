# Shell

[![Build](https://github.com/beysed/shell/actions/workflows/build.yml/badge.svg)](https://github.com/beysed/shell/actions/workflows/build.yml)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=beysed_shell&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=beysed_shell)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=beysed_shell&metric=coverage)](https://sonarcloud.io/summary/new_code?id=beysed_shell)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=beysed_shell&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=beysed_shell)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=beysed_shell&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=beysed_shell)

Shell related functions and stuff

## Execute

The module provides convenient way to execute external commands and deal with its `stdout`/`stderr` as well as with `stdin` (see unit tests for examples)

```
command := MakeCommand("sed", "-e", "s/a/A/g")
execution, _ := Execute(command)

execution.Stdin <- []byte("aaa")
close(execution.Stdin)

for run := true; run; {
    select {
    case out := <-execution.Stdout:
        fmt.Print(string(out))
    case err := <-execution.Stderr:
        fmt.Println(string(err))
    case <-execution.Exit:
        run = false
    case <-time.After(time.Second * 3):
        t.Error("process killed by timeout")
        execution.Kill()
        run = false
    }
}
```
