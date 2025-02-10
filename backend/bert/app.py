import logging

from flask import Flask, jsonify, request
from flask_cors import cross_origin

from class_model import ProfanityModel

# Configure logging to show INFO messages
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

app = Flask(__name__)

model = ProfanityModel()
logging.info("Start loading the model...")
model.load_model()


@app.route("/profanity", methods=["POST"])
@cross_origin()
def post_profanity():
    data = request.get_json()

    if "text" not in data:
        return 'Error: Invalid request. Please provide a "text" parameter.', 400

    text = data["text"].strip()

    profanity_score = model.predict(text)

    response = jsonify({"profanity_score": profanity_score})
    return response


@app.route("/health", methods=["GET"])
@cross_origin()
def health():
    return {"status": "healthy"}



if __name__ == "__main__":
    app.run()
