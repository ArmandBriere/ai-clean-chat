from datetime import datetime

import torch
from transformers import (
    AutoModelForSequenceClassification,
    AutoTokenizer,
    BertModel,
)

from const import DATA_FOLDER, MODERN_BERT_MODEL


class BERTClassifier(torch.nn.Module):
    def __init__(self, bert_model_name, num_classes):
        super(BERTClassifier, self).__init__()
        self.bert: BertModel = BertModel.from_pretrained(bert_model_name)
        self.dropout = torch.nn.Dropout(0.1)
        self.fc = torch.nn.Linear(self.bert.config.hidden_size, num_classes)

    def forward(self, input_ids, attention_mask):
        outputs = self.bert(input_ids=input_ids, attention_mask=attention_mask)
        pooled_output = outputs.pooler_output
        x = self.dropout(pooled_output)
        logits = self.fc(x)
        return logits


class ProfanityModel:
    def __init__(self):
        self.model = None
        self.tokenizer = None

    def load_model(self):
        """Load the BERT model and tokenizer."""
        selected_model = MODERN_BERT_MODEL
        num_classes = 2

        tokenizer = AutoTokenizer.from_pretrained(selected_model)
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        model: AutoModelForSequenceClassification = (
            AutoModelForSequenceClassification.from_pretrained(
                selected_model,
                num_labels=num_classes,
            ).to(device)
        )
        model.load_state_dict(
            torch.load(
                f"./{DATA_FOLDER}/bert_classifier.pth",
                map_location=torch.device("cuda" if torch.cuda.is_available() else "cpu"),
            )
        )

        model.eval()

        self.model = model
        self.tokenizer = tokenizer

    def predict(self, text, max_length=128):
        """Predict the sentiment of the text."""
        if self.model is None:
            self.load_model()

        start = datetime.now()

        encoding = self.tokenizer(
            text,
            return_tensors="pt",
            max_length=max_length,
            padding="max_length",
            truncation=True,
        )
        input_ids = encoding["input_ids"]
        attention_mask = encoding["attention_mask"]

        with torch.no_grad():
            outputs = self.model(input_ids=input_ids, attention_mask=attention_mask)
            # Get the logits from the outputs
            logits = outputs.logits
            # Apply softmax to get probabilities
            probabilities = torch.softmax(logits, dim=1)
            # Get the confidence score for the positive class (index 1)
            confidence_score = probabilities[0][1].item()

        end = datetime.now()
        print(f"Total inference time = {end - start}")
        print(f"Confidence score: {confidence_score}")

        return confidence_score
