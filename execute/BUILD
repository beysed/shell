load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "Execute.go",
        "Execution.go",
        "ExecutionError.go",
    ],
    importpath = "github.com/beysed/shell/execute",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["Execution_test.go"],
    embed = [":go_default_library"],
    deps = ["@com_github_stretchr_testify//assert:go_default_library"],
)
