#!/usr/bin/env bash

set -ex

if [ ! -f ./sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/tokens.txt ]; then
  curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
  tar xvf sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
  rm sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
  mv sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/decoder-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/decoder.onnx
  mv sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/encoder-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/encoder.onnx
  mv sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/joiner-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/joiner.onnx
fi

if [ ! -f ./sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/tokens.txt ]; then
  curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20.tar.bz2
  tar xvf sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20.tar.bz2
  rm sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20.tar.bz2
  mv sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/decoder-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/decoder.onnx
  mv sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/encoder-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/encoder.onnx
  mv sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/joiner-epoch-99-avg-1.int8.onnx sherpa-onnx-streaming-zipformer-bilingual-zh-en-2023-02-20/joiner.onnx
fi

if [ ! -f ./sherpa-onnx-streaming-zipformer-en-2023-06-26/tokens.txt ]; then
  curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2
  tar xvf sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2
  rm sherpa-onnx-streaming-zipformer-en-2023-06-26.tar.bz2
  mv sherpa-onnx-streaming-zipformer-en-2023-06-26/decoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx sherpa-onnx-streaming-zipformer-en-2023-06-26/decoder.onnx
  mv sherpa-onnx-streaming-zipformer-en-2023-06-26/encoder-epoch-99-avg-1-chunk-16-left-128.int8.onnx sherpa-onnx-streaming-zipformer-en-2023-06-26/encoder.onnx
  mv sherpa-onnx-streaming-zipformer-en-2023-06-26/joiner-epoch-99-avg-1-chunk-16-left-128.int8.onnx sherpa-onnx-streaming-zipformer-en-2023-06-26/joiner.onnx
fi
