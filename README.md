# wg-manager

## Description

`wg-manager` is a simple command-line tool for managing WireGuard VPN configurations on Arch Linux. It allows you to list available WireGuard configuration files, bring interfaces up or down using a specified configuration file, and check the status of WireGuard connections with formatted output.

## Features

- **List Configurations:** Discover all available WireGuard configuration files in `/etc/wireguard/`.
- **Manage Configurations:** Bring specific WireGuard interfaces up or down using the `wg-quick` utility.
- **Check Status:** Run the `wg` command and display its output with a visually appealing format.
- **Styled Output:** Use `lipgloss` for consistent and colorful command-line output.

## Requirements

- **Arch Linux:** The tool has been built and tested on Arch Linux, but it should work on any system with WireGuard installed.
- **Go-lang:** Ensure that Go is installed on your system to build the tool from source. You can install it via `pacman`:
  ```bash
  sudo pacman -S go
  ```
- **WireGuard:** Install WireGuard and `wg-quick` using:
  ```bash
  sudo pacman -S wireguard-tools
  ```

## Installation

### From Source

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/yourusername/wg-manager.git
   cd wg-manager
   ```

2. **Build and Install:**
   ```bash
   go build -o wg-manager main.go wireguard.go utils.go style.go
   sudo mv wg-manager /usr/local/bin/
   ```

### From Precompiled Binary

- Download the precompiled binary from the [releases page](https://github.com/yourusername/wg-manager/releases) and place it in your `PATH`.

## Usage

### List Available Configurations
```bash
wg-manager --list
```

**Output Example:**
```
Available WireGuard configurations:
  wg0
  wg1
```

### Bring a Configuration Up or Down
```bash
wg-manager --config <configuration_name> --state up|down
```

**Example:**
```bash
wg-manager --config wg0 --state up
wg-manager --config wg1 --state down
```

### Check WireGuard Status
```bash
wg-manager --check
```

**Output Example:**
```
interface: wg0
  public key: abcdef1234567890abcdef1234567890abcdef1234
  private key: (hidden)
  listening port: 51820

peer: abcdef1234567890abcdef1234567890abcdef1234
  endpoint: 203.0.113.1:51820
  allowed ips: 10.0.0.1/32
  latest handshake: 2023-10-05T14:39:06Z
  transfer: 2.5 MB received, 1.8 MB sent
```

### Help Information
```bash
wg-manager --help
```

**Output:**
```
Usage: wg-manager
--config <config_name> --state <up|down>
--list: list available WireGuard configuration files
--check: Run the wg command and display its output with formatting
```

## Contributing

Contributions are welcome! Feel free to fork the repository, make changes, and submit pull requests. For any issues or feature requests, please open an issue on the [GitHub repository](https://github.com/yourusername/wg-manager).

## License

This project is licensed under the BSD 2 Clause License. See the [LICENSE](LICENSE) file for more details.

## Acknowledgments

- **WireGuard:** A fast, modern, and secure VPN tunneling protocol.
- **lipgloss:** A styling library for the command-line in Go, used for creating visually appealing output.

---
