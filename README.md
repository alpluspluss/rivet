# Rivet

A fast C/C++ build system with cross-compilation support, workspace management, and incremental builds. 
Rivet supports parallel compilation and caching by default out of the box with workspace support 
for multi-project builds.

## Quick Start

```shell
rivet init --name project # --workspace
rivet build # --release --target [ ... ]
rivet run # --target [ ... ]
rivet clean # --target [ ... ]
```

### Project Configuration

Simply invoke the CLI or create a `rivet.toml` in your project root:

```toml
[project]
name = "root_name"
version = "0.1.0"
type = "binary" # binary, dynamic, static, shared, object

[paths]
src = ["src/"]
include = ["include/"]
output = "build"

[toolchain]
c = { compiler = "gcc", std = "c11" }
cxx = { compiler = "clang++", std = "c++20" }
asm = { compiler = "nasm", fmt = "elf64" }

[profiles]
default = "debug"

[profiles.debug]
optimization = 0
debug_info = true
defines = ["DEBUG"]
compiler_flags = { c = ["-Wall"], cxx = ["-Wall", "-Wextra", "-march=native"], asm = [] }

[cross] # cross compilation support; default to off
enabled = false
target = "aarch64-unknown-linux-gnu"
toolchain_path = "/opt/cross"
sysroot = "/opt/sysroot"
```

### Workspace Support

To create a workspace for multiple projects, call the CLI:

```shell
rivet init --workspace
```

then add these sections:

```toml
# workspace configuration
[workspace]
members = ["app", "lib"]
default_member = "app"  # which one to build if not specified

# shared toolchain configs that members can inherit
[workspace.toolchain]
c = { compiler = "gcc", standard = "c11" }
cxx = { compiler = "clang++", standard = "c++20" }

# shared profiles that members can extend
[workspace.profiles.base]
optimization = 2
debug_info = false
```

Each member can have its own config or inherit from workspace. Members can override workspace 
settings in their own config files. Workspace can have sub-workspaces.

## Building from Scratch

### Requirements

- Go 1.23.3 or compatible
- Compilers that you will be using (i.e. gcc, clang++, MSVC)

### Building

To build, run:
```shell
go run
```

## License

The project is licensed under MIT. See [LICENSE](LICENSE) for more details.

## Contributing

I am no good in Rust so pull requests are always welcome. For major changes, please open an issue first to discuss 
what you would like to change.