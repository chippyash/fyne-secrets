# Fyne Secrets
## git@github.com:chippyash/fyne-secrets.git

## What
Provides a device dependent secure storage mechanism for [Fyne apps](https://fyne.io/) using the 
device keychain if available.  Will fall back to using device file storage if the keychain is not available. 

## Why
Fyne does not provide a secure storage mechanism at the present time.  My hope is that this code will be incorporated 
into the Fyne project at some point in the future.

## Caveats

 - Secure for Gnome based Linux platforms
 - Work in progress for Android platforms
 - Not secure for Windows platforms - TBA
 - Not secure for MacOS platforms - TBA
 - Not secure for iOS platforms - TBA

## How



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