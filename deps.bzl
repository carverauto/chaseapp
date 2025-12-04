"""Shared external dependency declarations for Bazel."""

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")


def deps():
    """Register common workspace dependencies."""
    bazel_skylib_workspace()
    gazelle_dependencies()
