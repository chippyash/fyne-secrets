# Fyne Secrets
## git@github.com:chippyash/fyne-secrets.git

## What
Provides a device dependent secure storage mechanism for [Fyne apps](https://fyne.io/) using the 
device keychain if available.  Will fall back to using device file storage if the keychain is not available. 

## Why
Fyne does not provide a secure storage mechanism at the present time.  My hope is that this code will be incorporated 
into the Fyne project at some point in the future.

## Caveats

 - Secure for Linux platforms that use `gnome-keyring`
 - Secure for Linux platforms that use `keyctl`
 - Work in progress for Android platforms
 - Not secure for Windows platforms - TBA
 - Not secure for MacOS platforms - TBA
 - Not secure for iOS platforms - TBA

This secrets library **does not** provide encryption of the data.  That is left to the application using the library. Some real
world experience has shown that different platforms can treat the value being stored differently.  For safety you should
put your encrypted data in an envelope, e.g. urlEncode(base64Encode(data)) and unwrap it before using it.  A crypt implementation
is included in this library. You can use it as a starting point for your own implementation or use it as is.

```go
import "github.com/chippyash/fyne-secrets/crypt"
```

## How

Import the package for use in your application:
`import "github.com/chippyash/fyne-secrets/secrets"`

For development, 

 - fork the repository
 - clone it locally
 - `go mod tidy`
 - make your changes on a branch
 - push your branch to your repo
 - do a pull request back to this repo.

`make help` to see available Make commands

### Linux platforms
The Linux code looks to see if the `gnome-keyring` daemon is running and that `secret-tool` is installed and uses it if it is.

If you are running Linux and have `gnome-keyring` installed, you may need to run `sudo gnome-keyring-daemon --start` to start the daemon.
In addition, you may need to install, for Ubuntu/Debian `libsecret-tools` to get the `secret-tool` command.

If you are running Linux and `keyctl` installed, no additional steps are required.

The preferred method is to use gnome-keyring. keyctl has a limitation in that the persistent keyring for a user has a timeout
and if the user logs out the keyring is lost after that period. You can use 
`cat /proc/sys/kernel/keys/persistent_keyring_expiry` to see the timeout period in seconds, usually 72 hours.

## License
BSD-3-Clause. See [LICENSE](./LICENSE)

<a href="https://www.flaticon.com/free-icons/confidential" title="confidential icons">Application icon created by surang - Flaticon</a>

## Roadmap

 - Add support for Android
 - Add support for MacOS
 - Add support for iOS
 - Add support for Windows
 - Add support for other Linux flavours

## References

 - [Github](https://github.com/chippyash/fyne-secrets)
 - [Linux Gnome Secrets](https://grahamwatts.co.uk/gnome-secrets/)
 - [Why we need a better Windoze solution](https://grahamwatts.co.uk/windows-secrets/)
 - [vupdate - Semantic Version Updater](https://github.com/chippyash/semantic-version-updater/tree/master)
 - [keyutils](https://man7.org/linux/man-pages/man7/keyutils.7.html)