"""Train the model"""

import torch
from datasets import Dataset, load_dataset
from modal import App, Image, Volume
from sklearn.metrics import accuracy_score, classification_report
from sklearn.model_selection import train_test_split
from torch.optim import AdamW
from torch.utils.data import DataLoader
from tqdm import tqdm
from transformers import (
    AutoModelForSequenceClassification,
    AutoTokenizer,
    get_linear_schedule_with_warmup,
)

from const import DATA_FOLDER, MODERN_BERT_MODEL

image = (
    Image.from_registry("nvidia/cuda:12.8.0-cudnn-devel-ubuntu22.04", add_python="3.12")
    .workdir("/workspace/")
    .pip_install_from_requirements("./requirements.txt", force_build=False)
    .add_local_python_source("class_model", "const")
)

app = App(image=image)
vol = Volume.from_name("data")


def get_datasets():
    """Get the datasets."""
    dataset: Dataset = load_dataset(
        "csv",
        data_files="/workspace/data.csv",
        split="train",
    )

    # Make the is_offensive column the label column
    dataset = dataset.rename_column("is_offensive", "label")

    df = dataset.to_pandas()

    text = df["text"].tolist()
    labels = df["label"].tolist()

    return text, labels


class TextClassificationDataset(torch.utils.data.Dataset):
    def __init__(self, texts, labels, tokenizer, max_length):
        self.texts = texts
        self.labels = labels
        self.tokenizer = tokenizer
        self.max_length = max_length

    def __len__(self):
        return len(self.texts)

    def __getitem__(self, idx):
        # Handle both single index and list of indices
        if isinstance(idx, list):
            text = [self.texts[i] for i in idx]
            label = [self.labels[i] for i in idx]
        else:
            text = self.texts[idx]
            label = self.labels[idx]

        if text is None:
            text = ""
            label = 0

        # Tokenize the text(s)
        encoding = self.tokenizer(
            text,
            return_tensors="pt",
            max_length=self.max_length,
            padding="max_length",
            truncation=True,
        )

        # Convert label(s) to tensor
        if isinstance(idx, list):
            label_tensor = torch.tensor(label)
        else:
            label_tensor = torch.tensor(label)

        return {
            "input_ids": encoding["input_ids"].squeeze(),
            "attention_mask": encoding["attention_mask"].squeeze(),
            "label": label_tensor,
        }


def train(
    model: AutoModelForSequenceClassification,
    data_loader: DataLoader,
    optimizer: AdamW,
    scheduler,
    device,
    index,
):
    model.train()

    for batch in tqdm(data_loader):
        optimizer.zero_grad()

        input_ids = batch["input_ids"].to(device)
        attention_mask = batch["attention_mask"].to(device)
        labels = batch["label"].to(device)

        outputs = model(
            input_ids=input_ids, attention_mask=attention_mask, labels=labels
        )

        loss = outputs.loss
        loss.backward()

        optimizer.step()
        scheduler.step()

    torch.save(model.state_dict(), f"/data/model_checkpoint_{index}.pth")
    vol.commit()


def evaluate(model, data_loader, device):
    model.eval()
    predictions = []
    actual_labels = []

    with torch.no_grad():
        for batch in data_loader:
            input_ids = batch["input_ids"].to(device)
            attention_mask = batch["attention_mask"].to(device)
            labels = batch["label"].to(device)

            outputs = model(input_ids=input_ids, attention_mask=attention_mask)

            logits = outputs.logits
            _, preds = torch.max(logits, dim=1)
            predictions.extend(preds.cpu().tolist())
            actual_labels.extend(labels.cpu().tolist())

    return accuracy_score(actual_labels, predictions), classification_report(
        actual_labels, predictions
    )


@app.function(
    gpu="H100",
    image=image.add_local_dir(
        f"{DATA_FOLDER}/", remote_path="/workspace/"
    ).add_local_file(f"{DATA_FOLDER}/data.csv", remote_path="/workspace/data.csv"),
    timeout=60 * 60 * 1,
    volumes={"/data": vol},
)
def train_on_modal():
    import torch

    torch.cuda.empty_cache()

    print("Training model")
    print(f"Memory : {torch.cuda.get_device_properties(0).total_memory}")

    num_classes = 2
    max_length = 128
    batch_size = 16
    num_epochs = 4
    learning_rate = 2e-5
    texts, labels = get_datasets()

    train_texts, val_texts, train_labels, val_labels = train_test_split(
        texts, labels, test_size=0.2, random_state=42
    )

    selected_model = MODERN_BERT_MODEL

    tokenizer = AutoTokenizer.from_pretrained(selected_model)

    train_dataset = TextClassificationDataset(
        train_texts, train_labels, tokenizer, max_length
    )
    val_dataset = TextClassificationDataset(
        val_texts, val_labels, tokenizer, max_length
    )
    train_dataloader: DataLoader = DataLoader(
        train_dataset, batch_size=batch_size, shuffle=True
    )
    val_dataloader: DataLoader = DataLoader(val_dataset, batch_size=batch_size)

    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    model: AutoModelForSequenceClassification = (
        AutoModelForSequenceClassification.from_pretrained(
            selected_model,
            num_labels=num_classes,
        ).to(device)
    )

    optimizer: AdamW = AdamW(model.parameters(), lr=learning_rate)
    total_steps = len(train_dataloader) * num_epochs
    scheduler = get_linear_schedule_with_warmup(
        optimizer, num_warmup_steps=0, num_training_steps=total_steps
    )

    for index, epoch in enumerate(tqdm(range(num_epochs))):
        print(f"Epoch {epoch + 1}/{num_epochs}")
        train(model, train_dataloader, optimizer, scheduler, device, index)

        accuracy, report = evaluate(model, val_dataloader, device)
        print(f"Validation Accuracy: {accuracy:.4f}")
        print(report)

    torch.save(model.state_dict(), "/data/bert_classifier.pth")
    vol.commit()
