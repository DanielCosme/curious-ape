#!/usr/bin/env fish

set root_dir (pwd)
set web_ui_root "$root_dir"/web/ui

rm -r "$root_dir"/web/dist
npm install --prefix $web_ui_root; or exit
npm run build --prefix $web_ui_root
mv "$web_ui_root"/dist "$root_dir"/web/