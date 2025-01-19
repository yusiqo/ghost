# Ghost

![](https://raw.githubusercontent.com/yusiqo/ghost/refs/heads/main/banner.png)

Ghost is a lightweight package manager designed for Debian and Arch Linux users. It allows you to fetch and install packages effortlessly while also offering compatibility with other package managers like `yay` and `apt` for missing dependencies.

---

## Features
- **Custom Package Management**: Fetch and install custom packages from the Ghost repository.
- **Dependency Handling**: Automatically checks and installs missing dependencies.
- **Update Functionality**: Ensures you have the latest version of Ghost with a simple update command.
- **Fallback Support**: Tries popular package managers (`yay`, `apt`) when the desired package is unavailable in the Ghost repository.

---

## Installation
To install Ghost, run the following command:

```bash
sudo curl -L https://github.com/yusiqo/ghost/releases/latest/download/ghost -o /usr/local/bin/ghost && sudo chmod +x /usr/local/bin/ghost
```

---

## Usage

### Check for Updates
Update Ghost to the latest version:

```bash
ghost update
```

### Install a Package
Install a package by specifying its name:

```bash
ghost install <package-name>
```

Ghost will:
1. Search for the package in the Ghost repository.
2. Install any required dependencies.
3. Attempt installation using alternative package managers (`yay`, `apt`) if the package is not found.

---

## Reporting Missing Packages
If a package is missing from the Ghost repository, it will be reported automatically to the Ghost server. You can also manually report missing packages if needed.

---

## Contribution
Contributions are welcome! If you have any ideas or improvements, feel free to submit an issue or pull request on the [GitHub repository](https://github.com/yusiqo/ghost).

---

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/yusiqo/ghost/blob/main/LICENSE) file for details.

---

## Acknowledgments
Special thanks to the open-source community for making tools like this possible!

