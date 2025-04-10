#!/usr/bin/bash
rm -rf gio
git clone https://git.sr.ht/~eliasnaur/gio
cd gio||exit
rm -rf .git
cp -rf internal helpers
rm -rf internal
cp -rf app/internal app/helpers
rm -rf app/internal
cp -rf gpu/internal gpu/helpers
rm -rf gpu/internal
rm go.???
find ./ -type f -exec sed -i 's/\/internal/\/helpers/g' {} \;
find ./ -type f -exec sed -i 's/\"gioui.org/\"gio.mleku.dev\/gio/g' {} \;
find ./ -type f -exec sed -i 's/\"gioui.org\/\"internal/gio.mleku.dev\/gio\/helpers/g' {} \;
find ./ -type f -exec sed -i 's/\"gioui.org\/app\/internal/\"gio.mleku.dev\/gio\/app\/helpers/g' {} \;
find ./ -type f -exec sed -i 's/\"gio.mleku.dev\/gio\/shader/\"gioui.org\/shader/g' {} \;
cd ..
go mod tidy
git add .