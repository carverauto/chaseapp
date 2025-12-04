workspace(name = "chaseapp")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Core utility rules.
http_archive(
    name = "bazel_skylib",
    urls = ["https://github.com/bazelbuild/bazel-skylib/releases/download/1.5.0/bazel-skylib-1.5.0.tar.gz"],
)

http_archive(
    name = "build_bazel_rules_nodejs",
    url = "https://github.com/bazelbuild/rules_nodejs/archive/refs/tags/5.8.0.zip",
    sha256 = "adabe513387911365169a1403ca04f72ad5c4c079489fd5896b15ddb526ce3bd",
    strip_prefix = "rules_nodejs-5.8.0",
)

# Go rules + Gazelle.
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
    sha256 = "6dc2da7ab4cf5d7bfc7c949776b1b7c733f05e56edc4bcd9022bb249d2e2a996",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz"],
)

load("@build_bazel_rules_nodejs//:repositories.bzl", "build_bazel_rules_nodejs_dependencies")
build_bazel_rules_nodejs_dependencies()

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "npm_install")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//tools/build_rules:flutter.bzl", "flutter_repositories")
load("//:deps.bzl", "deps")

# Node toolchain + npm install for the web app.
node_repositories(
    node_version = "18.12.1",
)

npm_install(
    name = "npm",
    package_json = "//web:package.json",
    package_lock_json = "//web:package-lock.json",
    args = ["--legacy-peer-deps"],
    symlink_node_modules = True,
)

# Go + Gazelle.
go_rules_dependencies()
go_register_toolchains(version = "1.19.0")
gazelle_dependencies()

# Flutter placeholder (replace with real toolchain when ready).
flutter_repositories()

# Shared dependency registration.
deps()
