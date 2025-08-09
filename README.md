# wg-manager

`wg-manager` is a simple management tool for WireGuard VPN configurations on Linux systems. It aims to simplify the creation, modification, and management of WireGuard interfaces.

**Note:** This tool has been primarily tested on Arch Linux. While it may work on other distributions, support is not guaranteed.

## Features

- **Create and Configure WireGuard Interfaces:** Easily set up new VPN interfaces based on user input or predefined templates.
- **Manage Peers:** Add, remove, and modify peers for your WireGuard interfaces.
- **Export Configurations:** Save configurations to files for backup or manual editing.
- **Start and Stop Interfaces:** Manage the lifecycle of your WireGuard interfaces directly through the tool.

## Prerequisites

- **WireGuard Tools:** Ensure you have `wireguard-tools` installed on your system.
  ```bash
  sudo pacman -S wireguard-tools
  ```
- **Python 3:** `wg-manager` is written in Python and requires Python 3.6 or higher.
  ```bash
  sudo pacman -S python
  ```

## Installation

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/absenth/wg-manager.git
   cd wg-manager
   ```

2. **Install Dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

## Usage

### Create a New WireGuard Interface

```bash
python wg-manager.py create --interface <interface_name>
```

### Add a Peer to an Existing Interface

```bash
python wg-manager.py add-peer --interface <interface_name> --peer-name <peer_name>
```

### Start a WireGuard Interface

```bash
python wg-manager.py start --interface <interface_name>
```

### Stop a WireGuard Interface

```bash
python wg-manager.py stop --interface <interface_name>
```

For more detailed usage, see the `--help` flag:

```bash
python wg-manager.py --help
```

## Contributing

Contributions to `wg-manager` are welcome! Feel free to fork the repository, make your changes, and submit a pull request. Please ensure that any new features are tested on Arch Linux before submission.

## License

`wg-manager` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Support

If you encounter any issues or have feature requests, please open an issue on the [GitHub repository](https://github.com/absenth/wg-manager/issues).

---

Enjoy managing your WireGuard VPNs with `wg-manager`!
