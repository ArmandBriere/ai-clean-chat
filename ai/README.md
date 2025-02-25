# AI

The application is using two main AI models:

- A speech-to-text model to transcribe the conversation
- A text classification model to detect inappropriate content

Both of the model are only available in English.

## Speech-to-text

The speech-to-text model is Sherpa-onnx, a model developed by [K2-fsa](https://github.com/k2-fsa/sherpa-onnx) that is able to transcribe speech to text.

The code is available in the `backend` folder in the `sherpa-onnx` folder.

## Text classification

The text classification model is a [BERT](https://huggingface.co/docs/transformers/en/model_doc/bert) custom model trained on open source data to detect inappropriate content.
