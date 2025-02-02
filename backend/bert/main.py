from class_model import ProfanityModel
from flask import Flask, jsonify, request
from flask_cors import cross_origin

app = Flask(__name__)


@app.route("/insult", methods=["POST"])
@cross_origin()
def post_insult():
    data = request.get_json()

    if "text" not in data:
        return 'Error: Invalid request. Please provide a "text" parameter.', 400

    text = data["text"].strip()

    profanity_score = model.predict(text)

    response = jsonify({"profanity_score": profanity_score})
    return response


if __name__ == "__main__":
    model = ProfanityModel()
    model.load_model()

    app.run()
