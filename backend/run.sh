#!/usr/bin/env bash

set -ex

if [ ! -f ./sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18/HLG.fst ]; then
  curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18.tar.bz2
  tar xvf sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18.tar.bz2
  rm sherpa-onnx-streaming-zipformer-ctc-small-2024-03-18.tar.bz2
fi


if [ ! -f ./sherpa-onnx-streaming-zipformer-en-20M-2023-02-17/tokens.txt ]; then
  curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
  tar xvf sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
  rm sherpa-onnx-streaming-zipformer-en-20M-2023-02-17.tar.bz2
fi

go mod tidy
go build

PORT=8083 ./profanity.com
