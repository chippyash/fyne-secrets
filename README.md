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
 - Work in progress for Android platform
 - Secure for Windows platform
 - Secure for MacOS (Darwin) platform
 - Not secure for iOS platform - TBA

This secrets library **does not** provide encryption of the data.  That is left to the application using the library. Some real
world experience has shown that different platforms can treat the value being stored differently.  For safety you should
put your encrypted data in an envelope, e.g. urlEncode(base64Encode(data)) and unwrap it before using it.  A crypt implementation
is included in this library. You can use it as a starting point for your own implementation or use it as is.

```go
import "github.com/chippyash/fyne-secrets/crypt"
```

See `crypt/crypt.go` for interface details.

## How

```go
import "github.com/chippyash/fyne-secrets/secrets"

func Test_Functionality(t *testing.T) {
    sut, err := secrets.NewSecretStore("app-id", "app description")
    assert.NoError(t, err)
    err = sut.Store("test", []byte("test"))
    assert.NoError(t, err)
    b, err := sut.Exists("test")
    assert.NoError(t, err)
    assert.True(t, b)
    v, err := sut.Load("test")
    assert.NoError(t, err)
    assert.Equal(t, "test", string(v))
    err = sut.Delete("test")
    assert.NoError(t, err)
    b, err = sut.Exists("test")
    assert.NoError(t, err)
    assert.False(t, b)
}
```

See `secrets/secretive.go` for interface details.

For development, 

 - fork the repository
 - clone it locally
 - `go mod tidy`
 - make your changes on a branch
 - push your branch to your repo
 - do a pull request back to this repo.

`make help` to see available Make commands

### Linux platforms
The Linux code looks to see if the `gnome-keyring` daemon is running. If so, it will be used.

If you have `gnome-keyring` installed, you may need to run `sudo gnome-keyring-daemon --start` to start the daemon.

If `gnome-keyring` is not installed, i.e. you may be running the KDE desktop, then you will need to install it.
There is plenty of information on the internet to help you with this for your platform. Note, For Kubuntu, gnome-keyring is preinstalled

You may want to install the [secret-tool](https://man.archlinux.org/man/core/libsecret/secret-tool.1.en) cli command tool if not already installed.

 - Debian/Ubuntu: `sudo apt install libsecret-tools`
 - Fedora: `sudo dnf install libsecret`
 - Arch: `sudo pacman -S libsecret`
 - Suse: `sudo zypper install libsecret`

Verify installation: `secret-tool --help`

You may want to install the [Seahorse GUI](https://wiki.gnome.org/Apps/Seahorse) tool if not already installed.

 - Debian/Ubuntu: `sudo apt install seahorse`
 - Fedora: `sudo dnf install seahorse`
 - Arch: `sudo pacman -S seahorse`
 - Suse: `sudo zypper install seahorse`

If you don't have `gnome-keyring` but `keyctl` is installed, it will be used instead.

If neither `gnome-keyring` nor `keyctl` is installed, the library will fall back to using device file storage.

The preferred method is to use `gnome-keyring`. `keyctl` has a limitation in that the persistent keyring for a user has a timeout
and if the user logs out the keyring is lost after that period. You can use 
`cat /proc/sys/kernel/keys/persistent_keyring_expiry` to see the timeout period in seconds, usually 72 hours.

## License
BSD-3-Clause. See [LICENSE](./LICENSE)

<a href="https://www.flaticon.com/free-icons/confidential" title="confidential icons">Application icon created by surang - Flaticon</a>

## Roadmap

 - Add support for Android
 - Add support for iOS
 - Add support for other Linux flavours

## References

 - [Github](https://github.com/chippyash/fyne-secrets)
 - [Linux Gnome Secrets](https://grahamwatts.co.uk/gnome-secrets/)
 - [vupdate - Semantic Version Updater](https://github.com/chippyash/semantic-version-updater/tree/master)
 - [keyutils](https://man7.org/linux/man-pages/man7/keyutils.7.html)