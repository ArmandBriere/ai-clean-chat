import csv
import json
import re
import uuid
from collections import OrderedDict

import matplotlib.pyplot as plt
import pytest

from bert.class_model import ProfanityModel


@pytest.fixture()
def test_data():
    test_data_path = "tests/test_data.json"
    with open(test_data_path, "r") as file:
        data = json.load(file)["dataset"]

    return data


@pytest.fixture()
def model():
    """Load the model for testing."""

    model = ProfanityModel()
    model.load_model()
    return model


# def test_predict(test_data, model):
#     """Test the predict function of the model."""

#     for data in test_data:
#         text = data["sentence"].replace("***", "uck")
#         category = data["category"]

#         result = model.predict(text)
#         print(f"Input: {text}, Prediction: {result}, Category: {category}")
#         assert result >= 0.9 if category == "profane" else result < 0.9


def generate_ordered_word_map(text):
    words = re.findall(r"\b\w+\b", text.lower())

    word_map = OrderedDict()
    for word in words:
        word_map[str(uuid.uuid4())] = {
            "word": word,
            "profanity": 0,
        }

    return word_map


def create_sliding_window(word_map, window_size=8):
    ids = list(word_map.keys())
    words = []
    for value in word_map.values():
        words.append(value.get("word"))
    for i in range(len(words) - window_size + 1):
        window = {
            "input_ids": ids[i : i + window_size],
            "text": " ".join(words[i : i + window_size]),
        }
        yield window


def test_generate_profanity_map(model):
    """Test the generation of the profanity map."""

    data = """
Hey everyone, thanks for joining the call.
oh fuck off, I forgot my slides.
So, before we dive into the main agenda, I just wanted to quickly check—can everyone hear me okay? Great.
Did anyone took time to read the agenda? Why is everybody acting like a donkey today?
Alright, so as you know, we've been making progress on the backend integration, but we hit a bit of a roadblock with the authentication flow.
Oh my god, you are such a fucking idiot.
I spoke with the dev team earlier, and we're considering switching to a different approach using token-based validation instead of session-based.
It should help with scalability, but I want to make sure we're all aligned before we commit to the change.
Please Patricia, stop talking to us, you are so annoying and extremely stupid.
Does anyone have any concerns, or should we move forward with testing it?
"""
    word_map = generate_ordered_word_map(data)
    windows = create_sliding_window(word_map)

    for window in windows:
        profanity_score = model.predict(window.get("text"))

        print(f"{window}: {profanity_score}")
        for id in window.get("input_ids"):
            rounded_profanity_score = float(f"{round(profanity_score, 2):.2f}")
            print(f"{id}: {rounded_profanity_score}")

            word_map[id]["profanity"] = (
                word_map[id].get("profanity", 0) + rounded_profanity_score
            )

    for score in word_map.values():
        with open("profanity_scores.csv", "w", newline="") as csvfile:
            fieldnames = ["word", "profanity"]
            writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

            writer.writeheader()
            for score in word_map.values():
                writer.writerow(
                    {"word": score["word"], "profanity": score["profanity"]}
                )

    words = [value["word"] for value in word_map.values()]
    profanity_scores = [value["profanity"] for value in word_map.values()]
    x_pos = list(range(len(words)))

    plt.figure(figsize=(25, 10))
    plt.bar(x_pos, profanity_scores, color="orange")
    plt.xticks(x_pos, words, rotation=80)


    plt.xlabel("Words")
    plt.ylabel("Profanity Score")
    plt.title("Profanity Score of Words")
    
    plt.savefig("profanity_scores.png", bbox_inches="tight")
