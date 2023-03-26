all: clean \
	package-win-amd64 \
	package-linux-amd64 \
	package-linux-arm \
	package-darwin-amd64



windows_amd64_target=dist/roundabout_windows_amd64
build-win-amd64:
	mkdir $(windows_amd64_target) ; true
	GOOS=windows GOARCH=amd64 go build -o $(windows_amd64_target)/roundabout.exe  cmd/main/roundabout.go

package-win-amd64: build-win-amd64
	zip -r dist/roundabout_windows_amd64.zip $(windows_amd64_target)

linux_amd64_target=dist/roundabout_linux_amd64
build-linux-amd64:
	mkdir $(linux_amd64_target) ; true
	GOOS=linux GOARCH=amd64 go build -o $(linux_amd64_target)/roundabout  cmd/main/roundabout.go

package-linux-amd64: build-linux-amd64
	tar -zcvf $(linux_amd64_target).tar.gz $(linux_amd64_target)

linux_arm_target=dist/roundabout_linux_arm
build-linux-arm:
	mkdir $(linux_arm_target) ; true
	GOOS=linux GOARCH=arm go build -o $(linux_arm_target)/roundabout  cmd/main/roundabout.go

package-linux-arm: build-linux-arm
	tar -zcvf $(linux_arm_target).tar.gz $(linux_arm_target)

darwin_amd64_target=dist/roundabout_darwin_amd64
build-darwin-amd64:
	mkdir $(darwin_amd64_target) ; true
	GOOS=darwin GOARCH=amd64 go build -o $(darwin_amd64_target)/roundabout  cmd/main/roundabout.go

package-darwin-amd64: build-darwin-amd64
	tar -zcvf $(darwin_amd64_target).tar.gz $(darwin_amd64_target)

clean:
		rm -rf dist/*
