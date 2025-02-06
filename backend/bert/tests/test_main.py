import json

import pytest

from bert.class_model import ProfanityModel


@pytest.fixture()
def test_data():
    test_data_path = "tests/test_data.json"
    with open(test_data_path, "r") as file:
        data = json.load(file)['dataset']

    return data


@pytest.fixture()
def model():
    """Load the model for testing."""

    model = ProfanityModel()
    model.load_model()
    return model



def test_predict(test_data, model):
    """Test the predict function of the model."""

    for data in test_data:
        text = data['sentence'].replace("***", "uck")
        category = data['category']

        result = model.predict(text)
        print(f"Input: {text}, Prediction: {result}, Category: {category}")
        assert result >= 0.9 if category == "profane" else result < 0.9
