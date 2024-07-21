from flask import Flask, jsonify,Response, json
import random
import string

app = Flask(__name__)

def generate_email():
    username = ''.join(random.choices(string.ascii_lowercase, k=7))
    domain = ''.join(random.choices(string.ascii_lowercase, k=5))
    return f"{username}@{domain}.com"


def generate_person():
    email = generate_email()
    age = random.randint(1, 100)
    return {'email': email, 'age': age}

@app.route('/people', methods=['GET'])
def get_people():
    people = [generate_person() for _ in range(5000)]
    response = json.dumps(people[:20], indent=4)
    return Response(response, mimetype='application/json')



if __name__ == '__main__':
    app.run(debug=False)
