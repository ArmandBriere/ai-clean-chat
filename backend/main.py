import os

import requests
from faster_whisper import WhisperModel

model_size = "tiny.en"

# Run on CPU with INT8
model = WhisperModel(model_size, device="cpu", compute_type="int8")

if not os.path.exists("data/jfk.flac"):
    data = requests.get(
        "https://github.com/SYSTRAN/faster-whisper/raw/97a4785fa13d067c300f8b6e40c4381ad0381c02/docker/jfk.flac"
    )
    with open("data/jfk.flac", "wb") as f:
        f.write(data.content)

segments, info = model.transcribe("data/jfk.flac")

for segment in segments:
    print("[%.2fs -> %.2fs] %s" % (segment.start, segment.end, segment.text))
