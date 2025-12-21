ANDROID_HOME=${HOME}/Android/Sdk
ADB=${ANDROID_HOME}/platform-tools/adb
NDK_HOME=${ANDROID_HOME}/ndk/29.0.14206865
EMULATOR=${ANDROID_HOME}/emulator/emulator
EMULATOR_DEVICE=Samsung_Galaxy_S10
#BUILD=`cat BUILD`
VERSION=`cat VERSION`
APP_NAME=fyne-secrets
APP_ID=biz.zf4.fyne-secrets
GOROOT=${HOME}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.4.linux-amd64

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the APK image
	vupdate -ppatch
	export ANDROID_NDK_HOME=${NDK_HOME} && fyne package -os android --app-version ${VERSION}

.PHONY: emulate
emulate: ## Install APK to running Android Emulator
	${ADB} install ${APP_NAME}.apk

.PHONY: delete
delete: ## Removes the APK from the emulator but leaves app data intact
	${ADB} shell cmd package uninstall -k ${APP_ID}

.PHONY: purge
purge: ## Removes the APK and its data from the emulator
	${ADB} uninstall ${APP_ID}

.PHONY: update
update: delete build emulate ## Rebuild and push APK to emulator

.PHONY: local
local: ## Run app as local window app
	go run .

.PHONY: web
web: ## Run app as a web app on http://localhost:8080
	fyne serve

.PHONY: tag
tag: ## Tag the code for pushing to Github
	vupdate -ppatch
	git add .
	git commit -m"tag code at $(cat VERSION)"
	cat VERSION | xargs git tag

#.PHONY: release
#release: ## Build a release version of the app for public distribution
#	vupdate -pfeature
#	export ANDROID_NDK_HOME=${NDK_HOME} && fyne package -os android --release --app-version ${VERSION} --app-build ${BUILD}

.PHONY: start-adb
start-adb: stop-adb ## Start/Restart ADB server
	${ADB} start-server
	${ADB} devices

.PHONY: stop-adb
stop-adb: ## Stop the ADB server
	${ADB} kill-server

.PHONY: start-emulator
start-emulator:  ## Start the Android Emulator
	${EMULATOR} -avd ${EMULATOR_DEVICE}

.PHONY: shell
shell: ## Go into emulator shell
	${ADB} shell

.PHONY: flush-local
flush-local:  ## Flush locally cached secrets
	rm -f ${HOME}/.config/fyne/${APP_ID}/preferences.json

.PHONY: test-generic
test-generic: ## Test generic functionality
	${GOROOT}/bin/go test --tags generic ./tests

.PHONY: test-linux
test-linux: ## Test linux functionality
	${GOROOT}/bin/go test --tags linux ./tests