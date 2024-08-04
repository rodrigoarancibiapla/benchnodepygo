from numba import jit
import numpy as np
from collections import namedtuple
from flask import Flask, jsonify
import json  # Importar el módulo json

# Definir un namedtuple para la persona
Person = namedtuple('Person', ['email', 'age'])

@jit(nopython=True)
def generate_email():
    return f"user{np.random.randint(1, 1001)}@example.com"

@jit(nopython=True)
def generate_person():
    email = generate_email()
    age = np.random.randint(1, 101)
    return Person(email, age)

def get_people():
    people = [generate_person() for _ in range(5000)]
    return people

app = Flask(__name__)

@app.route('/people')
def get_people_route():
    people = get_people()
    first_20_people = people[:20]
    people_dicts = [{'email': person.email, 'age': person.age} for person in first_20_people]
    response = app.response_class(
        response=json.dumps(people_dicts, indent=4),
        status=200,
        mimetype='application/json'
    )
    return response

if __name__ == '__main__':
    app.run(debug=True)
