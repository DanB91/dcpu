#!/bin/bash

# This is a deploy script for a static website.
# It compresses html, js, css and some image formats.
#
# Dependencies on external tools:
#
# htmlcompressor: http://code.google.com/p/htmlcompressor/
# yuicompressor: http://developer.yahoo.com/yui/compressor/
# pngcrush: http://pmt.sourceforge.net/pngcrush/
#

SRC="data";
DST="deploy";

rm -rf "$DST";
cp -r "$SRC" "$DST";

# Minify HTML files.
echo "[*] Minify HTML files...";
htmlcompressor "$SRC" -o "$DST/" -r --remove-intertag-spaces;

# Minify Javascript files.
echo "[*] Minify Javascript files...";
for FILE in `find "$DST" -type f -name "*.js"`; do
	yuicompressor -o "$FILE" "$FILE";
done

# Minify stylesheets.
echo "[*] Minify CSS files...";
for FILE in `find "$DST" -type f -name "*.css"`; do
	yuicompressor -o "$FILE" "$FILE";
done

# Crush PNG files.
echo "[*] Crushing PNG files...";
for FILE in `find "$DST" -type f -name "*.png"`; do
	pngcrush -q -z 3 -reduce "$FILE" "$FILE.new" 1>/dev/null;

	if [ $? -ne 0 ]; then
		exit $?;
	fi

	cp -f "$FILE.new" "$FILE";
	rm -f "$FILE.new";
done

exit $?;
