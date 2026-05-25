
# Grail

Grail is a terminal-based Gmail manager — a fast, keyboard-friendly TUI (text user interface) for interacting with Gmail from your terminal. It focuses on providing a productive, low-latency email experience without leaving the keyboard: reading, searching, composing, labeling, and managing Gmail accounts for power users and people who live in terminals.

![Grail logo - replace later](assets/logo.png)

Long description
Grail brings essential Gmail workflows into a compact, efficient terminal UI. It supports multiple accounts, threaded conversations, search operators, keyboard-driven navigation, offline caching, and background synchronization so your inbox feels snappy even on slower networks. Grail prioritizes privacy and minimalism: configuration is file-based, credentials are stored using your OS secret store (or a local encrypted token store), and network access is limited to Gmail APIs/IMAP as configured.

Grail is designed for:
- Developers, sysadmins, and terminal power users who prefer keyboard UX over web interfaces.
- Remote work on servers or low-resource systems where a graphical browser is not available.
- Integrations with terminal workflows and scriptable automation.

Example setup video (replace with actual file or hosted link):

[Example setup video placeholder - replace later](assets/setup-example.mp4)

Setup
This section explains how to get Grail installed and running. Choose the option that best matches your environment.

1) Installation

- Download Prebuilt Binaries (recommended)
  - Releases: Visit the GitHub Releases page and download the binary for your platform. Extract and place the `grail` binary on your PATH (for example: `/usr/local/bin/` on Unix).

- Linux
  - Debian / Ubuntu (example):
    - Option A: Use the distribution package if available. Otherwise download the tarball from Releases, extract and move the binary to `/usr/local/bin`.
    - Option B: Install via Snap (if snap package exists): `sudo snap install grail`
    - Common dependencies: a modern terminal (xterm, Alacritty, kitty), and a system secret store if you want token storage integrated.

- Arch / Manjaro
  - Check the AUR for community packages or use the Releases tarball: extract and install the `grail` binary.

- Fedora / CentOS
  - Download release tarball or RPM if provided. Otherwise build from source (see below).

- macOS
  - Homebrew (if a tap exists): `brew install grail` or `brew tap <owner>/tap && brew install grail`.
  - Or download the macOS tarball from Releases and move the `grail` binary to `/usr/local/bin` or `/opt/homebrew/bin`.

- Windows
  - WSL (recommended): install the Linux binary inside your WSL distribution and run from your shell.

- Build from source
  - If you want to build the project locally, clone the repository and follow the project's build instructions. A generic flow:

```sh
git clone https://github.com/<owner>/Grail.git
cd Grail
# If the project uses a Makefile
make build
# or follow the project's language-specific build steps (cargo build --release, go build, python setup, etc.)
sudo make install
```

2) Usage

- Starting Grail
  - Run `grail` from your terminal. The first run will walk you through account setup and authentication.

- Account setup
  - You will be prompted to authorize access to Gmail. Grail supports OAuth flows and local-token fallback. Follow the on-screen instructions to visit an authorization URL and paste the token back into the TUI if required.

- Common keyboard shortcuts
  - j / k: move up/down
  - Enter: open message/thread
  - c: compose
  - r: reply
  - a: archive
  - l: label
  - /: search
  - g to jump, ? for help (these are examples — see the in-app help for a complete list)

- Scripting & Integration
  - Grail offers a command-mode and a short JSON-based API for basic operations (send, fetch, mark-read) suitable for simple shell scripts and hooks. See docs or `grail --help`.

3) Troubleshooting

- Authentication errors
  - Symptom: You can't authorize or tokens are rejected. Solution: ensure your system clock is correct, follow the full OAuth flow, and check that any firewall/HTTP-proxy isn't intercepting requests. If tokens are stored locally, try removing the token store and re-authenticating.

- Rendering issues / corrupted UI
  - Symptom: Boxes, lines, or characters don't look right. Solution: Use a modern terminal emulator and set the locale/encoding to UTF-8 (e.g., export LANG=en_US.UTF-8). Choose a monospaced font that supports box-drawing characters.

- IMAP / API connectivity
  - Symptom: Sync fails or messages do not load. Solution: Verify network connectivity to Gmail, ensure the account has IMAP/API enabled, and check that rate-limiting or account security settings are not blocking access.

- Missing dependencies or failed build
  - Symptom: Build or runtime errors about missing tools. Solution: install basic developer tools (gcc/clang, make, Rust toolchain or Go toolchain or Python 3.x depending on how Grail is implemented). See the project's CONTRIBUTING.md or BUILD.md for exact requirements.

- Performance issues
  - Symptom: Slow sync, laggy UI. Solution: enable offline caching, reduce sync frequency in configuration, and ensure you are not running heavy background tasks. Check logs for repeated API retries.

System architecture preview
Grail is split into focused components so it remains lightweight and testable. The diagram below shows the simplified runtime architecture.

ASCII overview

  +--------------------+      +------------------+
  |    Terminal UI     | <--> |  Input / Keymap  |
  | (ncurses/Textual)  |      +------------------+
  +--------------------+
             |
             v
  +--------------------+      +------------------+      +------------------+
  |  Core Application  | <--> |  Sync Engine     | <--> |  Storage (cache) |
  |  (commands, state) |      |  (IMAP / Gmail)  |      |  (sqlite/ldb)    |
  +--------------------+      +------------------+      +------------------+
             |
             v
  +--------------------+      +------------------+
  |  Auth Manager      | <--> |  OS Secret Store |
  +--------------------+      +------------------+
             |
             v
     Network / Gmail APIs (OAuth / IMAP / SMTP)

Notes:
- The UI renders local state maintained by the core application and listens to events from the sync engine.
- The sync engine handles network interactions, rate limiting, incremental sync and background fetching.
- Storage is an append-friendly local cache (e.g., sqlite, sled, or plain JSON) allowing offline reads.
- Auth Manager isolates credential handling and integrates with OS keyrings when possible.

License
Grail is open source. You can choose a permissive license such as the MIT License or Apache-2.0. The recommended (default) license for this project is MIT.

MIT License (summary)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions: the above copyright
notice and this permission notice shall be included in all copies or
substantial portions of the Software.

For the full license text, see the LICENSE file in this repository.

Contributing
Contributions are welcome. Please open tickets for bugs and feature requests, and submit pull requests with tests and a clear description of changes.

More
For more detailed developer notes, configuration options, and keyboard reference, see the in-repo docs directory and the help built into the application (`grail --help`).

