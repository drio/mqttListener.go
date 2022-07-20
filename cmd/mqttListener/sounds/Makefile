all:
	for f in *.wav; do ffmpeg -i "$$f" "$${f/%wav/mp3}"; done
